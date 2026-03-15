package config

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager holds all loaded configs and provides management operations.
type Manager struct {
	configs []*Config
}

// NewManager creates a Manager, loading the embedded default config and any
// user configs from ~/.config/command-builder/configs/.
func NewManager(defaultData []byte) (*Manager, error) {
	m := &Manager{}

	// Load the embedded default config.
	if len(defaultData) > 0 {
		cfg, err := LoadConfigFromBytes(defaultData)
		if err == nil {
			cfg.FilePath = "" // embedded – no file path
			m.configs = append(m.configs, cfg)
		}
	}

	// Load user configs.
	if err := EnsureConfigDir(); err == nil {
		userCfgs, _ := LoadAllConfigs(GetConfigDir())
		for _, uc := range userCfgs {
			// Don't duplicate a config whose name matches the default.
			if uc.Name == "default" {
				continue
			}
			m.configs = append(m.configs, uc)
		}
	}

	return m, nil
}

// ListConfigs returns all loaded configs.
func (m *Manager) ListConfigs() []*Config {
	return m.configs
}

// GetConfig returns the Config with the given name, or nil.
func (m *Manager) GetConfig(name string) *Config {
	for _, c := range m.configs {
		if c.Name == name {
			return c
		}
	}
	return nil
}

// AddConfig adds a new config and saves it to the config directory.
func (m *Manager) AddConfig(cfg *Config) error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}
	path := filepath.Join(GetConfigDir(), sanitizeName(cfg.Name)+".yaml")
	cfg.FilePath = path
	if err := SaveConfig(cfg, path); err != nil {
		return err
	}
	// Replace existing or append.
	for i, c := range m.configs {
		if c.Name == cfg.Name {
			m.configs[i] = cfg
			return nil
		}
	}
	m.configs = append(m.configs, cfg)
	return nil
}

// UpdateConfig persists in-memory changes made to an already-loaded config.
// Returns an error if the config has no file path (built-in configs).
func (m *Manager) UpdateConfig(cfg *Config) error {
	if cfg.FilePath == "" {
		return fmt.Errorf("config %q is a built-in config and cannot be modified", cfg.Name)
	}
	return SaveConfig(cfg, cfg.FilePath)
}

// DeleteConfig removes a config by name, deleting its file if it has one.
func (m *Manager) DeleteConfig(name string) error {
	for i, c := range m.configs {
		if c.Name == name {
			if c.FilePath != "" {
				if err := os.Remove(c.FilePath); err != nil && !os.IsNotExist(err) {
					return err
				}
			}
			m.configs = append(m.configs[:i], m.configs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("config %q not found", name)
}

// ExportConfig copies a config's YAML to the given destination path.
func (m *Manager) ExportConfig(name, destPath string) error {
	cfg := m.GetConfig(name)
	if cfg == nil {
		return fmt.Errorf("config %q not found", name)
	}
	return SaveConfig(cfg, destPath)
}

// ImportConfigFromURL fetches a YAML config from a URL and adds it.
func (m *Manager) ImportConfigFromURL(rawURL string) (*Config, error) {
	// Validate the URL: only http and https are allowed.
	parsed, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("invalid URL %q: %w", rawURL, err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme %q: only http and https are allowed", parsed.Scheme)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(rawURL) //nolint:noctx // timeout is set on the client
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", rawURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, rawURL)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	cfg, err := LoadConfigFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	cfg.SourceURL = rawURL

	// Resolve name collision.
	base := cfg.Name
	counter := 1
	for m.GetConfig(cfg.Name) != nil {
		cfg.Name = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}

	if err := m.AddConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// ImportConfigFromFile loads a YAML config from a local file path and adds it.
// A leading ~ in the path is expanded to the user's home directory.
// The config is copied into the command-builder config directory, mirroring
// the behaviour of ImportConfigFromURL.
func (m *Manager) ImportConfigFromFile(rawPath string) (*Config, error) {
	// Expand ~ prefix.
	path := rawPath
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		path = filepath.Join(home, path[2:])
	} else if path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("cannot determine home directory: %w", err)
		}
		path = home
	}

	// Resolve to absolute path.
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path %q: %w", rawPath, err)
	}

	cfg, err := LoadConfig(abs)
	if err != nil {
		return nil, fmt.Errorf("load config from file: %w", err)
	}

	// Resolve name collision.
	base := cfg.Name
	counter := 1
	for m.GetConfig(cfg.Name) != nil {
		cfg.Name = fmt.Sprintf("%s-%d", base, counter)
		counter++
	}

	// Copy into the managed config directory.
	if err := m.AddConfig(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// PullConfig re-fetches a config from its SourceURL and replaces its commands
// in place, preserving the config name and file path.
func (m *Manager) PullConfig(name string) (*Config, error) {
	existing := m.GetConfig(name)
	if existing == nil {
		return nil, fmt.Errorf("config %q not found", name)
	}
	if existing.SourceURL == "" {
		return nil, fmt.Errorf("config %q has no source URL", name)
	}

	parsed, err := url.ParseRequestURI(existing.SourceURL)
	if err != nil {
		return nil, fmt.Errorf("invalid source URL: %w", err)
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, fmt.Errorf("unsupported URL scheme %q", parsed.Scheme)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(existing.SourceURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("fetch %s: %w", existing.SourceURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d fetching %s", resp.StatusCode, existing.SourceURL)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	fresh, err := LoadConfigFromBytes(data)
	if err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	// Preserve identity fields; replace content.
	existing.Description = fresh.Description
	existing.Version = fresh.Version
	existing.Commands = fresh.Commands
	// Keep existing.SourceURL and existing.FilePath unchanged.

	if err := m.UpdateConfig(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

// RenameConfig renames a config and its backing file.
func (m *Manager) RenameConfig(oldName, newName string) error {
	cfg := m.GetConfig(oldName)
	if cfg == nil {
		return fmt.Errorf("config %q not found", oldName)
	}
	if m.GetConfig(newName) != nil {
		return fmt.Errorf("config %q already exists", newName)
	}

	oldPath := cfg.FilePath
	cfg.Name = newName

	newPath := filepath.Join(GetConfigDir(), sanitizeName(newName)+".yaml")
	cfg.FilePath = newPath
	if err := SaveConfig(cfg, newPath); err != nil {
		return err
	}
	if oldPath != "" && oldPath != newPath {
		os.Remove(oldPath)
	}
	return nil
}

// sanitizeName converts a config name to a safe filename stem.
func sanitizeName(name string) string {
	replacer := strings.NewReplacer(
		" ", "-",
		"/", "-",
		"\\", "-",
		":", "-",
		"*", "-",
		"?", "-",
		"\"", "-",
		"<", "-",
		">", "-",
		"|", "-",
	)
	return replacer.Replace(name)
}
