package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/atotto/clipboard"
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

	// Print or run the final built command.
	if a, ok := finalModel.(tui.AppModel); ok {
		if cmd := a.GetFinalCommand(); cmd != "" {
			if a.GetSettings().RunOnEnter {
				shell := os.Getenv("SHELL")
				if shell == "" {
					shell = "/bin/sh"
				}
				c := exec.Command(shell, "-c", cmd)
				c.Stdin = os.Stdin
				c.Stdout = os.Stdout
				c.Stderr = os.Stderr
				if runErr := c.Run(); runErr != nil {
					fmt.Fprintf(os.Stderr, "Error running command: %v\n", runErr)
					os.Exit(1)
				}
			} else {
				fmt.Println(cmd)
				if err := clipboard.WriteAll(cmd); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: could not copy to clipboard: %v\n", err)
				}
			}
		}
	}
}
