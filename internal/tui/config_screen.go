package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
)

// configAction represents the current action in the config screen.
type configAction int

const (
	actionList configAction = iota
	actionNew
	actionDelete
	actionExport
	actionImport
	actionPull
)

// ConfigScreenModel is the config management screen.
type ConfigScreenModel struct {
	mgr      *config.Manager
	configs  []*config.Config
	selected int
	action   configAction
	input    textinput.Model // used for new name, export path, import URL
	message  string
	width    int
	height   int
}

// NewConfigScreenModel creates a new config screen.
func NewConfigScreenModel(mgr *config.Manager, w, h int) ConfigScreenModel {
	ti := textinput.New()
	ti.Width = max(30, w-10)

	m := ConfigScreenModel{
		mgr:    mgr,
		width:  w,
		height: h,
		input:  ti,
		action: actionList,
	}
	m.configs = mgr.ListConfigs()
	return m
}

// Init satisfies tea.Model.
func (m ConfigScreenModel) Init() tea.Cmd { return nil }

// Update satisfies tea.Model.
func (m ConfigScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = max(30, msg.Width-10)
		return m, nil

	case tea.KeyMsg:
		switch m.action {
		case actionList:
			return m.updateList(msg)
		default:
			return m.updateInput(msg)
		}
	}
	return m, nil
}

func (m ConfigScreenModel) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		return m, func() tea.Msg { return backToSearchMsg{} }

	case tea.KeyUp:
		if m.selected > 0 {
			m.selected--
		}

	case tea.KeyDown:
		if m.selected < len(m.configs)-1 {
			m.selected++
		}

	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "n", "N":
			m.action = actionNew
			m.input.SetValue("")
			m.input.Placeholder = "New config name…"
			return m, m.input.Focus()

		case "d", "D":
			if len(m.configs) > 0 {
				cfg := m.configs[m.selected]
				if cfg.Name == "default" {
					m.message = StyleError.Render("Cannot delete the built-in default config")
					break
				}
				m.action = actionDelete
				m.input.SetValue("")
				m.input.Placeholder = fmt.Sprintf("Type %q to confirm deletion", cfg.Name)
				return m, m.input.Focus()
			}

		case "x", "X":
			if len(m.configs) > 0 {
				m.action = actionExport
				m.input.SetValue("")
				m.input.Placeholder = "Export path (e.g. ~/my-config.yaml)…"
				return m, m.input.Focus()
			}

		case "e", "E":
			if len(m.configs) > 0 {
				cfg := m.configs[m.selected]
				if cfg.FilePath == "" {
					m.message = StyleError.Render("Cannot edit the built-in default config")
					break
				}
				return m, func() tea.Msg { return goToEditMsg{cfg: cfg} }
			}

		case "i", "I":
			m.action = actionImport
			m.input.SetValue("")
			m.input.Placeholder = "Config URL…"
			return m, m.input.Focus()

		case "u", "U":
			if len(m.configs) > 0 {
				cfg := m.configs[m.selected]
				if cfg.SourceURL == "" {
					m.message = StyleError.Render("Config has no source URL")
					break
				}
				m.action = actionPull
				m.input.SetValue("")
				m.input.Placeholder = fmt.Sprintf("Re-pull %q from %s ? (yes to confirm)", cfg.Name, cfg.SourceURL)
				return m, m.input.Focus()
			}

		case "q", "Q":
			return m, func() tea.Msg { return backToSearchMsg{} }
		}
	}
	return m, nil
}

func (m ConfigScreenModel) updateInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {

	case tea.KeyEsc:
		m.action = actionList
		m.input.Blur()
		m.message = ""
		return m, nil

	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEnter:
		val := strings.TrimSpace(m.input.Value())
		switch m.action {

		case actionNew:
			if val == "" {
				m.message = StyleError.Render("Config name cannot be empty")
				break
			}
			cfg := &config.Config{
				Name:        val,
				Description: "Custom config",
				Version:     "1.0.0",
			}
			if err := m.mgr.AddConfig(cfg); err != nil {
				m.message = StyleError.Render("Error: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf("Created config %q", val))
				m.configs = m.mgr.ListConfigs()
			}
			m.action = actionList
			m.input.Blur()

		case actionDelete:
			if len(m.configs) == 0 {
				break
			}
			target := m.configs[m.selected]
			if val != target.Name {
				m.message = StyleError.Render("Name doesn't match – deletion cancelled")
				m.action = actionList
				m.input.Blur()
				break
			}
			if err := m.mgr.DeleteConfig(target.Name); err != nil {
				m.message = StyleError.Render("Error: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf("Deleted config %q", target.Name))
				m.configs = m.mgr.ListConfigs()
				if m.selected >= len(m.configs) {
					m.selected = max(0, len(m.configs)-1)
				}
			}
			m.action = actionList
			m.input.Blur()

		case actionExport:
			if len(m.configs) == 0 {
				break
			}
			if val == "" {
				m.message = StyleError.Render("Path cannot be empty")
				break
			}
			target := m.configs[m.selected]
			if err := m.mgr.ExportConfig(target.Name, val); err != nil {
				m.message = StyleError.Render("Export failed: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf("Exported %q → %s", target.Name, val))
			}
			m.action = actionList
			m.input.Blur()

		case actionImport:
			if val == "" {
				m.message = StyleError.Render("URL cannot be empty")
				break
			}
			cfg, err := m.mgr.ImportConfigFromURL(val)
			if err != nil {
				m.message = StyleError.Render("Import failed: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf("Imported config %q", cfg.Name))
				m.configs = m.mgr.ListConfigs()
			}
			m.action = actionList
			m.input.Blur()

		case actionPull:
			if strings.ToLower(val) != "yes" {
				m.message = StyleResultDesc.Render("Pull cancelled")
				m.action = actionList
				m.input.Blur()
				break
			}
			if len(m.configs) == 0 {
				break
			}
			target := m.configs[m.selected]
			cfg, err := m.mgr.PullConfig(target.Name)
			if err != nil {
				m.message = StyleError.Render("Pull failed: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf("Updated config %q from %s", cfg.Name, cfg.SourceURL))
				m.configs = m.mgr.ListConfigs()
			}
			m.action = actionList
			m.input.Blur()
		}
		return m, nil
	}

	// Delegate to text input.
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// View satisfies tea.Model.
func (m ConfigScreenModel) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	var b strings.Builder

	// ── Title ──────────────────────────────────────────────────────────────
	title := StyleTitle.Copy().Width(w).Render(
		"⚡ Command Builder " + StyleTitleVersion.Render(AppVersion) + "  " +
			StyleResultDesc.Render("Config Manager"),
	)
	b.WriteString(title + "\n")

	// ── Header ─────────────────────────────────────────────────────────────
	b.WriteString(StyleConfigHeader.Render("Loaded Configurations") + "\n")
	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Config list ────────────────────────────────────────────────────────
	if len(m.configs) == 0 {
		b.WriteString(StyleResultDesc.Padding(1, 2).Render("No configs loaded.") + "\n")
	}

	visRows := h - 10
	if visRows < 1 {
		visRows = 1
	}
	for i, cfg := range m.configs {
		if i >= visRows {
			break
		}
		badge := ""
		if cfg.FilePath == "" {
			badge = " " + StyleResultConfig.Render("[built-in]")
		} else if cfg.SourceURL != "" {
			badge = " " + StyleResultConfig.Render("[url]")
		}
		cmds := fmt.Sprintf("%d cmd(s)", len(cfg.Commands))
		line := fmt.Sprintf("%-20s  %-30s  %s%s", cfg.Name, cfg.Description, cmds, badge)
		if i == m.selected {
			b.WriteString(StyleConfigItemSelected.Copy().Width(w).Render(line) + "\n")
		} else {
			b.WriteString(StyleConfigItem.Copy().Width(w).Render(line) + "\n")
		}
	}

	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Active action input ────────────────────────────────────────────────
	if m.action != actionList {
		prompt := m.actionPrompt()
		b.WriteString("\n" + StyleInputLabel.Render("  "+prompt) + "\n")
		b.WriteString(StyleInputFocused.Copy().Width(w-6).Render(m.input.View()) + "\n")
		b.WriteString(StyleResultDesc.Render("  Enter to confirm · Esc to cancel") + "\n")
	}

	if m.message != "" {
		b.WriteString("\n" + "  " + m.message + "\n")
	}

	// ── Status bar ─────────────────────────────────────────────────────────
	used := strings.Count(b.String(), "\n") + 2
	for i := used; i < h-1; i++ {
		b.WriteString("\n")
	}

	keys := StyleStatusKey.Render(" n") + StyleStatus.Render(" new") +
		StyleStatusKey.Render("  e") + StyleStatus.Render(" edit") +
		StyleStatusKey.Render("  d") + StyleStatus.Render(" delete") +
		StyleStatusKey.Render("  x") + StyleStatus.Render(" export") +
		StyleStatusKey.Render("  i") + StyleStatus.Render(" import URL") +
		StyleStatusKey.Render("  u") + StyleStatus.Render(" pull update") +
		StyleStatusKey.Render("  Esc") + StyleStatus.Render(" back")
	b.WriteString(StyleStatus.Copy().Width(w).Render(keys))

	return b.String()
}

func (m ConfigScreenModel) actionPrompt() string {
	switch m.action {
	case actionNew:
		return "New config name:"
	case actionDelete:
		if len(m.configs) > 0 {
			return fmt.Sprintf("Confirm deletion of %q (type name):", m.configs[m.selected].Name)
		}
		return "Confirm deletion:"
	case actionExport:
		if len(m.configs) > 0 {
			return fmt.Sprintf("Export %q to path:", m.configs[m.selected].Name)
		}
		return "Export path:"
	case actionImport:
		return "Import config from URL:"
	case actionPull:
		if len(m.configs) > 0 {
			cfg := m.configs[m.selected]
			return fmt.Sprintf("Re-pull %q from %s (type \"yes\"):", cfg.Name, cfg.SourceURL)
		}
		return "Re-pull config (type \"yes\"):"
	}
	return ""
}
