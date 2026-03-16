package tui

import (
	"fmt"
	"os"
	"path/filepath"
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
	actionImportFile
	actionPull
)

// ConfigScreenModel is the config management screen.
type ConfigScreenModel struct {
	mgr            *config.Manager
	configs        []*config.Config
	selected       int
	action         configAction
	input          textinput.Model // used for new name, export path, import URL/file
	message        string
	width          int
	height         int
	completions    []string // tab-complete candidates for file import
	completionIdx  int      // cycling index (-1 = at longest-common-prefix stage)
	completionBase string   // input value that produced current completions
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
				m.action = actionDelete
				m.input.SetValue("")
				if cfg.FilePath == "" {
					m.input.Placeholder = fmt.Sprintf("Type %q to confirm deletion (built-in — will be hidden on next launch)", cfg.Name)
				} else {
					m.input.Placeholder = fmt.Sprintf("Type %q to confirm deletion", cfg.Name)
				}
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
					// Promote the embedded default to a real file before editing.
					if err := m.mgr.PromoteDefaultConfig(cfg); err != nil {
						m.message = StyleError.Render("Could not save default config to disk: " + err.Error())
						break
					}
					m.message = StyleInfo.Render("Default config saved to disk – opening editor")
				}
				return m, func() tea.Msg { return goToEditMsg{cfg: cfg} }
			}

		case "i", "I":
			m.action = actionImport
			m.input.SetValue("")
			m.input.Placeholder = "Config URL (http:// or https://)…"
			return m, m.input.Focus()

		case "f", "F":
			m.action = actionImportFile
			m.input.SetValue("")
			m.input.Placeholder = "Local file path (Tab to autocomplete)…"
			m.completions = nil
			m.completionIdx = -1
			m.completionBase = ""
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

		case actionImportFile:
			if val == "" {
				m.message = StyleError.Render("File path cannot be empty")
				break
			}
			cfg, err := m.mgr.ImportConfigFromFile(val)
			if err != nil {
				m.message = StyleError.Render("Import failed: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf("Imported config %q", cfg.Name))
				m.configs = m.mgr.ListConfigs()
			}
			m.action = actionList
			m.input.Blur()
			m.completions = nil
			m.completionBase = ""

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

	// Tab-completion for file import.
	if msg.Type == tea.KeyTab && m.action == actionImportFile {
		val := m.input.Value()
		// Recompute completions if the input value changed since last Tab.
		if val != m.completionBase || len(m.completions) == 0 {
			m.completions = pathCompletions(val)
			m.completionBase = val
			m.completionIdx = -1
		}
		if len(m.completions) == 0 {
			return m, nil
		}
		if len(m.completions) == 1 {
			completed := m.completions[0]
			m.input.SetValue(completed)
			m.input.CursorEnd()
			m.completions = nil
			m.completionBase = completed
			return m, nil
		}
		// Multiple matches: first Tab fills the longest common prefix;
		// each subsequent Tab cycles through individual completions.
		prefix := longestCommonPrefix(m.completions)
		if m.completionIdx == -1 && prefix != val {
			m.input.SetValue(prefix)
			m.input.CursorEnd()
			m.completionBase = prefix
		} else {
			if m.completionIdx == -1 {
				m.completionIdx = 0
			} else {
				m.completionIdx = (m.completionIdx + 1) % len(m.completions)
			}
			m.input.SetValue(m.completions[m.completionIdx])
			m.input.CursorEnd()
		}
		return m, nil
	}

	// Delegate to text input.
	prevVal := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	// Reset completions whenever the typed value changes.
	if m.input.Value() != prevVal {
		m.completions = nil
		m.completionIdx = -1
		m.completionBase = ""
	}
	return m, cmd
}

// pathCompletions returns glob-based file path completions for the given partial
// path. A leading ~ is expanded to the user home directory. Directories in the
// result have a trailing slash appended.
func pathCompletions(partial string) []string {
	if partial == "" {
		return nil
	}
	home := ""
	expanded := partial
	if h, err := os.UserHomeDir(); err == nil {
		home = h
		if partial == "~" {
			expanded = home
		} else if strings.HasPrefix(partial, "~/") {
			expanded = filepath.Join(home, partial[2:])
		}
	}

	matches, err := filepath.Glob(expanded + "*")
	if err != nil || len(matches) == 0 {
		return nil
	}

	// Append trailing slash for directories; re-apply ~ prefix.
	for i, m := range matches {
		if info, serr := os.Stat(m); serr == nil && info.IsDir() {
			matches[i] = m + "/"
		}
		if home != "" && strings.HasPrefix(matches[i], home+"/") {
			matches[i] = "~/" + matches[i][len(home)+1:]
		}
	}
	return matches
}

// longestCommonPrefix returns the longest string that is a prefix of every
// element in strs.
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	prefix := strs[0]
	for _, s := range strs[1:] {
		for !strings.HasPrefix(s, prefix) {
			if len(prefix) == 0 {
				return ""
			}
			prefix = prefix[:len(prefix)-1]
		}
	}
	return prefix
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
	title := StyleTitle.Width(w).Render(
		"⚡ " + AppDisplayName + "  " + StyleResultDesc.Render("Config Manager"),
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
			b.WriteString(StyleConfigItemSelected.Width(w).Render(line) + "\n")
		} else {
			b.WriteString(StyleConfigItem.Width(w).Render(line) + "\n")
		}
	}

	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Active action input ────────────────────────────────────────────────
	if m.action != actionList {
		prompt := m.actionPrompt()
		b.WriteString("\n" + StyleInputLabel.Render("  "+prompt) + "\n")
		b.WriteString(StyleInputFocused.Width(w-6).Render(m.input.View()) + "\n")
		// Show tab-complete suggestions for file import.
		if m.action == actionImportFile && len(m.completions) > 1 {
			const maxShow = 5
			shown := m.completions
			if len(shown) > maxShow {
				shown = shown[:maxShow]
			}
			for i, c := range shown {
				if i == m.completionIdx {
					b.WriteString(StyleResultSelected.Width(w - 4).Render("  "+c) + "\n")
				} else {
					b.WriteString(StyleResultNormal.Render("  "+c) + "\n")
				}
			}
			if len(m.completions) > maxShow {
				b.WriteString(StyleResultDesc.Render(fmt.Sprintf("  … and %d more (Tab to cycle)", len(m.completions)-maxShow)) + "\n")
			}
		}
		hint := "  Enter to confirm · Esc to cancel"
		if m.action == actionImportFile {
			hint += " · Tab to autocomplete"
		}
		b.WriteString(StyleResultDesc.Render(hint) + "\n")
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
		StyleStatusKey.Render("  f") + StyleStatus.Render(" import file") +
		StyleStatusKey.Render("  u") + StyleStatus.Render(" pull update") +
		StyleStatusKey.Render("  Esc") + StyleStatus.Render(" back")
	b.WriteString(renderFooter(w, keys, footerVersion()))

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
	case actionImportFile:
		return "Import config from local file:"
	case actionPull:
		if len(m.configs) > 0 {
			cfg := m.configs[m.selected]
			return fmt.Sprintf("Re-pull %q from %s (type \"yes\"):", cfg.Name, cfg.SourceURL)
		}
		return "Re-pull config (type \"yes\"):"
	}
	return ""
}
