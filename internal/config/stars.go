package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Star represents a saved command with its input values, allowing users to
// quickly re-run frequently-used commands with pre-filled parameters.
type Star struct {
	ID          string            `json:"id"`
	ConfigName  string            `json:"config_name"`
	CommandName string            `json:"command_name"`
	OptionName  string            `json:"option_name"`
	Label       string            `json:"label"`
	CustomName  string            `json:"custom_name,omitempty"` // optional user-defined display name
	Values      map[string]string `json:"values"`                // input name → value (string/file/dir inputs)
	FlagStates  map[string]bool   `json:"flag_states"`           // input name → bool (flag inputs)
	CreatedAt   time.Time         `json:"created_at"`
}

// DisplayName returns the custom name if one was set, otherwise the default label.
func (s Star) DisplayName() string {
	if s.CustomName != "" {
		return s.CustomName
	}
	return s.Label
}

// starsFilePath returns the path to the persisted stars JSON file.
func starsFilePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "command-builder", "stars.json")
}

// LoadStars loads the saved stars from disk. Returns an empty slice on any
// error so callers can treat a missing file as "no stars yet".
func LoadStars() []Star {
	data, err := os.ReadFile(starsFilePath())
	if err != nil {
		return []Star{}
	}
	var stars []Star
	if err := json.Unmarshal(data, &stars); err != nil {
		return []Star{}
	}
	return stars
}

// SaveStars writes the full stars list to disk.
func SaveStars(stars []Star) error {
	path := starsFilePath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(stars, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// AddStar appends a star to the saved list.
func AddStar(star Star) error {
	if star.ID == "" {
		star.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if star.CreatedAt.IsZero() {
		star.CreatedAt = time.Now()
	}
	stars := LoadStars()
	stars = append(stars, star)
	return SaveStars(stars)
}

// DeleteStar removes the star with the given ID from the saved list.
func DeleteStar(id string) error {
	stars := LoadStars()
	filtered := stars[:0]
	for _, s := range stars {
		if s.ID != id {
			filtered = append(filtered, s)
		}
	}
	return SaveStars(filtered)
}
