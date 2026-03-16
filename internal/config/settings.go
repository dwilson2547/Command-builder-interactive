package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AppSettings holds user-tunable global settings that persist between sessions.
// Color values are strings accepted by lipgloss.Color – either ANSI terminal
// colour codes ("0"–"255") or hex colours ("#rrggbb").
type AppSettings struct {
	AppName         string `json:"app_name"`
	RunOnEnter      bool   `json:"run_on_enter"`
	ColorPrimary    string `json:"color_primary"`
	ColorAccent     string `json:"color_accent"`
	ColorSuccess    string `json:"color_success"`
	ColorWarning    string `json:"color_warning"`
	ColorError      string `json:"color_error"`
	ColorMuted      string `json:"color_muted"`
	ColorText       string `json:"color_text"`
	ColorBackground string `json:"color_background"`
	ColorSelected   string `json:"color_selected"`
}

// DefaultSettings returns a settings instance with the built-in colour palette.
func DefaultSettings() AppSettings {
	return AppSettings{
		AppName:         "Command Builder",
		RunOnEnter:      false,
		ColorPrimary:    "39",
		ColorAccent:     "213",
		ColorSuccess:    "76",
		ColorWarning:    "220",
		ColorError:      "196",
		ColorMuted:      "241",
		ColorText:       "252",
		ColorBackground: "235",
		ColorSelected:   "24",
	}
}

// settingsFilePath returns the path to the persisted settings JSON file.
func settingsFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "command-builder", "settings.json")
}

// LoadSettings loads settings from disk, merging over built-in defaults so
// that any value absent from the file falls back to the default value.
func LoadSettings() AppSettings {
	s := DefaultSettings()
	data, err := os.ReadFile(settingsFilePath())
	if err != nil {
		return s
	}
	// Unmarshal into a temporary struct; only non-empty fields override defaults.
	var tmp AppSettings
	if err := json.Unmarshal(data, &tmp); err != nil {
		return s
	}
	if tmp.ColorPrimary != "" {
		s.ColorPrimary = tmp.ColorPrimary
	}
	if tmp.ColorAccent != "" {
		s.ColorAccent = tmp.ColorAccent
	}
	if tmp.ColorSuccess != "" {
		s.ColorSuccess = tmp.ColorSuccess
	}
	if tmp.ColorWarning != "" {
		s.ColorWarning = tmp.ColorWarning
	}
	if tmp.ColorError != "" {
		s.ColorError = tmp.ColorError
	}
	if tmp.ColorMuted != "" {
		s.ColorMuted = tmp.ColorMuted
	}
	if tmp.ColorText != "" {
		s.ColorText = tmp.ColorText
	}
	if tmp.ColorBackground != "" {
		s.ColorBackground = tmp.ColorBackground
	}
	if tmp.ColorSelected != "" {
		s.ColorSelected = tmp.ColorSelected
	}
	if tmp.AppName != "" {
		s.AppName = tmp.AppName
	}
	// bool: always copy (false is a valid user choice)
	s.RunOnEnter = tmp.RunOnEnter
	return s
}

// SaveSettings writes the settings to disk, creating parent directories as
// needed.
func SaveSettings(s AppSettings) error {
	path := settingsFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
