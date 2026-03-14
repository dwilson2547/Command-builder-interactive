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
