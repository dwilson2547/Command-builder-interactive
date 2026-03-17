package config

// Config represents a command configuration file.
type Config struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Version     string    `yaml:"version"`
	SourceURL   string    `yaml:"source_url,omitempty"` // URL this config was fetched from, if any
	FilePath    string    `yaml:"-"`                    // set at runtime, not persisted
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
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Template    string   `yaml:"template"`
	Tags        []string `yaml:"tags,omitempty"` // searchable aliases/keywords
	Inputs      []Input  `yaml:"inputs"`
}

// Input represents a user-fillable field in an option template.
type Input struct {
	Name        string `yaml:"name"`
	Type        string `yaml:"type"` // "file", "dir", "string", "flag"
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Default     string `yaml:"default"`
	// SubCommand is an optional shell command whose stdout is parsed as CSV to
	// populate a dynamic picker. Column 0 is the value inserted into the field;
	// column 1 (optional) is a display detail shown alongside the value.
	SubCommand string `yaml:"sub_command,omitempty"`
}
