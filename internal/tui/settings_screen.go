package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dwilson2547/command-builder/internal/config"
)

// ---- messages ---------------------------------------------------------------

// goToSettingsMsg navigates to the /settings screen.
type goToSettingsMsg struct{}

// themeChangedMsg signals that ApplyTheme has been called and all screens
// should re-render with the new palette.
type themeChangedMsg struct{}

// ---- colour entry table -----------------------------------------------------

// settingsEntry describes one editable colour slot.
type settingsEntry struct {
	label  string
	desc   string
	getVal func(config.AppSettings) string
	setVal func(*config.AppSettings, string)
	defVal string
}

var settingsEntries = func() []settingsEntry {
	def := config.DefaultSettings()
	return []settingsEntry{
		{
			label:  "Primary",
			desc:   "Commands, borders, title highlights",
			getVal: func(s config.AppSettings) string { return s.ColorPrimary },
			setVal: func(s *config.AppSettings, v string) { s.ColorPrimary = v },
			defVal: def.ColorPrimary,
		},
		{
			label:  "Accent",
			desc:   "Options, focused inputs, required fields",
			getVal: func(s config.AppSettings) string { return s.ColorAccent },
			setVal: func(s *config.AppSettings, v string) { s.ColorAccent = v },
			defVal: def.ColorAccent,
		},
		{
			label:  "Success",
			desc:   "Command previews, success messages",
			getVal: func(s config.AppSettings) string { return s.ColorSuccess },
			setVal: func(s *config.AppSettings, v string) { s.ColorSuccess = v },
			defVal: def.ColorSuccess,
		},
		{
			label:  "Warning",
			desc:   "Completion overlays, warnings",
			getVal: func(s config.AppSettings) string { return s.ColorWarning },
			setVal: func(s *config.AppSettings, v string) { s.ColorWarning = v },
			defVal: def.ColorWarning,
		},
		{
			label:  "Error",
			desc:   "Error messages and banners",
			getVal: func(s config.AppSettings) string { return s.ColorError },
			setVal: func(s *config.AppSettings, v string) { s.ColorError = v },
			defVal: def.ColorError,
		},
		{
			label:  "Muted",
			desc:   "Descriptions, hints, separators",
			getVal: func(s config.AppSettings) string { return s.ColorMuted },
			setVal: func(s *config.AppSettings, v string) { s.ColorMuted = v },
			defVal: def.ColorMuted,
		},
		{
			label:  "Text",
			desc:   "Normal result rows and labels",
			getVal: func(s config.AppSettings) string { return s.ColorText },
			setVal: func(s *config.AppSettings, v string) { s.ColorText = v },
			defVal: def.ColorText,
		},
		{
			label:  "Selected BG",
			desc:   "Background of selected list rows",
			getVal: func(s config.AppSettings) string { return s.ColorSelected },
			setVal: func(s *config.AppSettings, v string) { s.ColorSelected = v },
			defVal: def.ColorSelected,
		},
	}
}()

// ---- model ------------------------------------------------------------------

// SettingsScreenModel is the /settings TUI screen.
type SettingsScreenModel struct {
	settings config.AppSettings
	selected int
	editing  bool
	input    textinput.Model
	message  string
	width    int
	height   int
}

// NewSettingsScreenModel creates a settings screen loaded with the current settings.
func NewSettingsScreenModel(s config.AppSettings, w, h int) SettingsScreenModel {
	ti := textinput.New()
	ti.Width = 22
	ti.Placeholder = "ANSI 0-255 or #rrggbb"
	return SettingsScreenModel{
		settings: s,
		width:    w,
		height:   h,
		input:    ti,
	}
}

// Init satisfies tea.Model.
func (m SettingsScreenModel) Init() tea.Cmd { return nil }

// Update satisfies tea.Model.
func (m SettingsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.editing {
			return m.updateEdit(msg)
		}
		return m.updateBrowse(msg)
	}
	return m, nil
}

func (m SettingsScreenModel) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
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
		if m.selected < len(settingsEntries)-1 {
			m.selected++
		}
	case tea.KeyEnter:
		return m.startEditing()
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "e", "E":
			return m.startEditing()
		case "r":
			// Reset selected entry to its default value.
			e := settingsEntries[m.selected]
			e.setVal(&m.settings, e.defVal)
			m.message = StyleInfo.Render(fmt.Sprintf("Reset \"%s\" to default (%s)", e.label, e.defVal))
			if svErr := config.SaveSettings(m.settings); svErr != nil {
				m.message = StyleError.Render("Save failed: " + svErr.Error())
			}
			ApplyTheme(m.settings)
			return m, func() tea.Msg { return themeChangedMsg{} }
		case "R":
			// Reset ALL entries to defaults.
			m.settings = config.DefaultSettings()
			m.message = StyleInfo.Render("All colours reset to defaults")
			if svErr := config.SaveSettings(m.settings); svErr != nil {
				m.message = StyleError.Render("Save failed: " + svErr.Error())
			}
			ApplyTheme(m.settings)
			return m, func() tea.Msg { return themeChangedMsg{} }
		}
	}
	return m, nil
}

func (m SettingsScreenModel) startEditing() (tea.Model, tea.Cmd) {
	e := settingsEntries[m.selected]
	m.input.SetValue(e.getVal(m.settings))
	m.input.CursorEnd()
	m.editing = true
	m.message = ""
	return m, m.input.Focus()
}

func (m SettingsScreenModel) updateEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.editing = false
		m.message = ""
		return m, nil
	case tea.KeyEnter:
		val := strings.TrimSpace(m.input.Value())
		if val == "" {
			m.editing = false
			return m, nil
		}
		e := settingsEntries[m.selected]
		e.setVal(&m.settings, val)
		m.editing = false
		m.message = StyleInfo.Render(fmt.Sprintf("Updated \"%s\" → %s", e.label, val))
		if svErr := config.SaveSettings(m.settings); svErr != nil {
			m.message = StyleError.Render("Save failed: " + svErr.Error())
		}
		ApplyTheme(m.settings)
		return m, func() tea.Msg { return themeChangedMsg{} }
	}
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	return m, inputCmd
}

// View satisfies tea.Model.
func (m SettingsScreenModel) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	var b strings.Builder

	// ── Title bar ──────────────────────────────────────────────────────────
	title := StyleTitle.Copy().Width(w).Render(
		"⚡ Command Builder " + StyleTitleVersion.Render(AppVersion) + "  " +
			StyleResultDesc.Render("Settings"),
	)
	b.WriteString(title + "\n")

	// ── Header ─────────────────────────────────────────────────────────────
	b.WriteString(StyleConfigHeader.Render("Colour Palette") + "\n")
	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Colour entry rows ──────────────────────────────────────────────────
	for i, e := range settingsEntries {
		val := e.getVal(m.settings)
		swatch := lipgloss.NewStyle().
			Background(lipgloss.Color(val)).
			Foreground(lipgloss.Color(val)).
			Render("  █  ")
		labelPart := fmt.Sprintf("%-12s", e.label)
		valPart := fmt.Sprintf("%-10s", val)

		if i == m.selected {
			b.WriteString(StyleConfigItemSelected.Copy().Width(w).Render(
				"  "+labelPart+"  "+swatch+"  "+valPart+"  "+e.desc,
			) + "\n")
		} else {
			line := "  " +
				StyleResultCommand.Render(labelPart) +
				"  " + swatch + "  " +
				StyleResultOption.Render(valPart) +
				"  " + StyleResultDesc.Render(e.desc)
			b.WriteString(line + "\n")
		}
	}

	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Edit input (shown only when editing) ───────────────────────────────
	if m.editing {
		e := settingsEntries[m.selected]
		b.WriteString("\n")
		b.WriteString(StyleInputLabel.Render(
			fmt.Sprintf("  Edit \"%s\" (current: %s):", e.label, e.getVal(m.settings)),
		) + "\n")
		b.WriteString("  " + StyleInputFocused.Copy().Width(28).Render(m.input.View()) + "\n")
		b.WriteString(StyleResultDesc.Render("  Enter to confirm · Esc to cancel") + "\n")
	}

	// ── Status / message ───────────────────────────────────────────────────
	if m.message != "" {
		b.WriteString("\n  " + m.message + "\n")
	}

	// ── Pad remaining rows ─────────────────────────────────────────────────
	used := strings.Count(b.String(), "\n") + 2
	for i := used; i < h-1; i++ {
		b.WriteString("\n")
	}

	// ── Status bar ─────────────────────────────────────────────────────────
	keys := StyleStatusKey.Render(" e") + StyleStatus.Render(" edit") +
		StyleStatusKey.Render("  r") + StyleStatus.Render(" reset") +
		StyleStatusKey.Render("  R") + StyleStatus.Render(" reset all") +
		StyleStatusKey.Render("  Esc") + StyleStatus.Render(" back")
	b.WriteString(renderFooter(w, keys, footerVersion()))

	return b.String()
}
