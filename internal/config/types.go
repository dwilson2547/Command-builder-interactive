package config

// Config represents a command configuration file.
type Config struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Version     string    `yaml:"version"`
	FilePath    string    `yaml:"-"` // set at runtime, not persisted
	Commands    []Command `yaml:"commands"`
}

// Command represents a CLI tool with multiple sub-options.
type Command struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Options     []Option `yaml:"options"`
}

// Option represents a specific use-case for a command.
type Option struct {
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Template    string  `yaml:"template"`
	Inputs      []Input `yaml:"inputs"`
}

// Input represents a user-fillable field in an option template.
type Input struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"` // "file", "dir", "string", "flag"
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
}
