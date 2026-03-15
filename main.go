package main

import (
	_ "embed"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
	"github.com/dwilson2547/command-builder/internal/tui"
)

//go:embed configs/default.yaml
var defaultConfigData []byte

func main() {
	mgr, err := config.NewManager(defaultConfigData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configs: %v\n", err)
		os.Exit(1)
	}

	// Load user settings (falls back to built-in defaults gracefully).
	settings := config.LoadSettings()
	// Apply the palette so all styles reflect user colours from the first render.
	tui.ApplyTheme(settings)

	app := tui.NewApp(mgr, settings)
	p := tea.NewProgram(app, tea.WithAltScreen())

	finalModel, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}

	// Print the final built command so the user can use it.
	if a, ok := finalModel.(tui.AppModel); ok {
		if cmd := a.GetFinalCommand(); cmd != "" {
			fmt.Println(cmd)
		}
	}
}
