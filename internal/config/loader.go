package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// LoadConfig loads a YAML config file from the given path.
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	cfg.FilePath = path
	return &cfg, nil
}

// LoadConfigFromBytes loads a Config from raw YAML bytes.
func LoadConfigFromBytes(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig saves a Config to the given path as YAML.
func SaveConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// LoadAllConfigs loads all .yaml/.yml configs from a directory.
func LoadAllConfigs(configDir string) ([]*Config, error) {
	entries, err := os.ReadDir(configDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var configs []*Config
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := filepath.Ext(entry.Name())
		if ext != ".yaml" && ext != ".yml" {
			continue
		}
		path := filepath.Join(configDir, entry.Name())
		cfg, err := LoadConfig(path)
		if err != nil {
			continue // skip invalid configs
		}
		configs = append(configs, cfg)
	}

	return configs, nil
}

// GetConfigDir returns the user's command-builder config directory.
func GetConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "command-builder", "configs")
}

// EnsureConfigDir creates the config directory if it doesn't exist.
func EnsureConfigDir() error {
	return os.MkdirAll(GetConfigDir(), 0o755)
}

// GetDefaultConfigPath returns the path to the default config, searching
// next to the executable and then the current working directory.
func GetDefaultConfigPath() string {
	// Try next to the executable first.
	if exec, err := os.Executable(); err == nil {
		p := filepath.Join(filepath.Dir(exec), "configs", "default.yaml")
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	// Fall back to current working directory.
	p := filepath.Join("configs", "default.yaml")
	if _, err := os.Stat(p); err == nil {
		return p
	}

	return ""
}
