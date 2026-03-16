package tui

import (
	"fmt"
	"os"
	"path/filepath"
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

// settingsEntry describes one editable settings slot.
type settingsEntry struct {
	label       string
	desc        string
	getVal      func(config.AppSettings) string
	setVal      func(*config.AppSettings, string)
	defVal      string
	isColor     bool   // renders colour swatch; false = plain text field
	isToggle    bool   // renders [✓]/[ ] toggle, toggled with Enter (no text input)
	placeholder string // hint shown inside the textinput
}

// firstColorIndex is the index of the first colour entry in settingsEntries.
// All entries before it are general (non-colour) settings.
const firstColorIndex = 2

var settingsEntries = func() []settingsEntry {
	def := config.DefaultSettings()
	return []settingsEntry{
		// ── General ────────────────────────────────────────────────────────────
		{
			label:       "App Name",
			desc:        "Name shown in the header (also used for the shell alias)",
			getVal:      func(s config.AppSettings) string { return s.AppName },
			setVal:      func(s *config.AppSettings, v string) { s.AppName = v },
			defVal:      def.AppName,
			isColor:     false,
			placeholder: "e.g. devtools or My CLI",
		},
		{
			label:    "Run on Enter",
			desc:     "Execute the built command immediately instead of printing it",
			getVal: func(s config.AppSettings) string {
				if s.RunOnEnter {
					return "true"
				}
				return "false"
			},
			setVal: func(s *config.AppSettings, v string) { s.RunOnEnter = v == "true" },
			defVal:   "false",
			isToggle: true,
		},
		// ── Colour Palette ─────────────────────────────────────────────────────
		{
			label:       "Primary",
			desc:        "Commands, borders, title highlights",
			getVal:      func(s config.AppSettings) string { return s.ColorPrimary },
			setVal:      func(s *config.AppSettings, v string) { s.ColorPrimary = v },
			defVal:      def.ColorPrimary,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Accent",
			desc:        "Options, focused inputs, required fields",
			getVal:      func(s config.AppSettings) string { return s.ColorAccent },
			setVal:      func(s *config.AppSettings, v string) { s.ColorAccent = v },
			defVal:      def.ColorAccent,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Success",
			desc:        "Command previews, success messages",
			getVal:      func(s config.AppSettings) string { return s.ColorSuccess },
			setVal:      func(s *config.AppSettings, v string) { s.ColorSuccess = v },
			defVal:      def.ColorSuccess,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Warning",
			desc:        "Completion overlays, warnings",
			getVal:      func(s config.AppSettings) string { return s.ColorWarning },
			setVal:      func(s *config.AppSettings, v string) { s.ColorWarning = v },
			defVal:      def.ColorWarning,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Error",
			desc:        "Error messages and banners",
			getVal:      func(s config.AppSettings) string { return s.ColorError },
			setVal:      func(s *config.AppSettings, v string) { s.ColorError = v },
			defVal:      def.ColorError,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Muted",
			desc:        "Descriptions, hints, separators",
			getVal:      func(s config.AppSettings) string { return s.ColorMuted },
			setVal:      func(s *config.AppSettings, v string) { s.ColorMuted = v },
			defVal:      def.ColorMuted,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Text",
			desc:        "Normal result rows and labels",
			getVal:      func(s config.AppSettings) string { return s.ColorText },
			setVal:      func(s *config.AppSettings, v string) { s.ColorText = v },
			defVal:      def.ColorText,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
		{
			label:       "Selected BG",
			desc:        "Background of selected list rows",
			getVal:      func(s config.AppSettings) string { return s.ColorSelected },
			setVal:      func(s *config.AppSettings, v string) { s.ColorSelected = v },
			defVal:      def.ColorSelected,
			isColor:     true,
			placeholder: "ANSI 0-255 or #rrggbb",
		},
	}
}()

// ---- model ------------------------------------------------------------------

// SettingsScreenModel is the /settings TUI screen.
type SettingsScreenModel struct {
	settings         config.AppSettings
	selected         int
	editing          bool
	confirmingBashrc bool   // waiting for y/n answer about bashrc alias
	pendingAliasLine string // the exact alias line to be written
	pendingName      string // the new app name being confirmed
	input            textinput.Model
	message          string
	width            int
	height            int
}

// NewSettingsScreenModel creates a settings screen loaded with the current settings.
func NewSettingsScreenModel(s config.AppSettings, w, h int) SettingsScreenModel {
	ti := textinput.New()
	ti.Width = 28
	ti.Placeholder = "ANSI 0-255 or #rrggbb"
	return SettingsScreenModel{
		settings: s,
		width:    w,
		height:   h,
		input:    ti,
	}
}

// bashrcAliasName converts a display name to a safe bash alias identifier.
// e.g. "My CLI" → "my-cli", "devtools" → "devtools"
func bashrcAliasName(name string) string {
	n := strings.ToLower(strings.TrimSpace(name))
	n = strings.ReplaceAll(n, " ", "-")
	return n
}

// appendAliasToBashrc appends aliasLine to ~/.bashrc.
func appendAliasToBashrc(aliasLine string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	bashrcPath := filepath.Join(home, ".bashrc")
	f, err := os.OpenFile(bashrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "\n# command-builder alias (added by command-builder /settings)\n%s\n", aliasLine)
	return err
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
	// Intercept key presses while the bashrc confirmation is pending.
	if m.confirmingBashrc {
		return m.updateConfirmBashrc(msg)
	}
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
		if settingsEntries[m.selected].isToggle {
			return m.toggleBool()
		}
		return m.startEditing()
	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "e", "E":
			if settingsEntries[m.selected].isToggle {
				return m.toggleBool()
			}
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
			m.message = StyleInfo.Render("All settings reset to defaults")
			if svErr := config.SaveSettings(m.settings); svErr != nil {
				m.message = StyleError.Render("Save failed: " + svErr.Error())
			}
			ApplyTheme(m.settings)
			return m, func() tea.Msg { return themeChangedMsg{} }
		}
	}
	return m, nil
}

// updateConfirmBashrc handles the y/n prompt after an app name change.
func (m SettingsScreenModel) updateConfirmBashrc(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit
	case tea.KeyEsc:
		m.confirmingBashrc = false
		m.message = StyleInfo.Render(fmt.Sprintf("Name updated to \"%s\" (no alias added)", m.pendingName))
		return m, nil
	case tea.KeyRunes:
		switch strings.ToLower(string(msg.Runes)) {
		case "y":
			m.confirmingBashrc = false
			if err := appendAliasToBashrc(m.pendingAliasLine); err != nil {
				m.message = StyleError.Render("Failed to update ~/.bashrc: " + err.Error())
			} else {
				m.message = StyleInfo.Render(fmt.Sprintf(
					"Added to ~/.bashrc: %s  — run 'source ~/.bashrc' to activate",
					m.pendingAliasLine,
				))
			}
		case "n":
			m.confirmingBashrc = false
			m.message = StyleInfo.Render(fmt.Sprintf("Name updated to \"%s\" (no alias added)", m.pendingName))
		}
	}
	return m, nil
}

func (m SettingsScreenModel) startEditing() (tea.Model, tea.Cmd) {
	e := settingsEntries[m.selected]
	if e.isToggle {
		return m.toggleBool()
	}
	m.input.SetValue(e.getVal(m.settings))
	m.input.Placeholder = e.placeholder
	m.input.CursorEnd()
	m.editing = true
	m.message = ""
	return m, m.input.Focus()
}

// toggleBool flips a boolean toggle entry and persists the change.
func (m SettingsScreenModel) toggleBool() (tea.Model, tea.Cmd) {
	e := settingsEntries[m.selected]
	if e.getVal(m.settings) == "true" {
		e.setVal(&m.settings, "false")
	} else {
		e.setVal(&m.settings, "true")
	}
	m.message = StyleInfo.Render(fmt.Sprintf("Updated \"%s\"", e.label))
	if svErr := config.SaveSettings(m.settings); svErr != nil {
		m.message = StyleError.Render("Save failed: " + svErr.Error())
	}
	return m, func() tea.Msg { return themeChangedMsg{} }
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
		if svErr := config.SaveSettings(m.settings); svErr != nil {
			m.message = StyleError.Render("Save failed: " + svErr.Error())
			return m, nil
		}
		ApplyTheme(m.settings)
		if !e.isColor && !e.isToggle {
			// Non-colour, non-toggle entry (e.g. App Name): offer to create a shell alias.
			aliasID := bashrcAliasName(val)
			m.pendingName = val
			m.pendingAliasLine = fmt.Sprintf("alias %s='command-builder'", aliasID)
			m.confirmingBashrc = true
			m.message = ""
			return m, func() tea.Msg { return themeChangedMsg{} }
		}
		m.message = StyleInfo.Render(fmt.Sprintf("Updated \"%s\" → %s", e.label, val))
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
		"⚡ " + AppDisplayName + "  " + StyleResultDesc.Render("Settings"),
	)
	b.WriteString(title + "\n")

	// ── Settings entry rows ────────────────────────────────────────────────
	for i, e := range settingsEntries {
		// Print section headers at section boundaries.
		if i == 0 {
			b.WriteString(StyleConfigHeader.Render("General") + "\n")
			b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")
		} else if i == firstColorIndex {
			b.WriteString(StyleConfigHeader.Render("Colour Palette") + "\n")
			b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")
		}

		val := e.getVal(m.settings)
		labelPart := fmt.Sprintf("%-12s", e.label)
		valPart := fmt.Sprintf("%-22s", val)

		if e.isColor {
			swatch := lipgloss.NewStyle().
				Background(lipgloss.Color(val)).
				Foreground(lipgloss.Color(val)).
				Render("  █  ")
			colorValPart := fmt.Sprintf("%-10s", val)
			if i == m.selected {
				b.WriteString(StyleConfigItemSelected.Copy().Width(w).Render(
					"  "+labelPart+"  "+swatch+"  "+colorValPart+"  "+e.desc,
				) + "\n")
			} else {
				line := "  " +
					StyleResultCommand.Render(labelPart) +
					"  " + swatch + "  " +
					StyleResultOption.Render(colorValPart) +
					"  " + StyleResultDesc.Render(e.desc)
				b.WriteString(line + "\n")
			}
		} else if e.isToggle {
			checkbox := "[ ]"
			if val == "true" {
				checkbox = "[✓]"
			}
			if i == m.selected {
				b.WriteString(StyleConfigItemSelected.Copy().Width(w).Render(
					"  "+labelPart+"  "+checkbox+"  "+e.desc,
				) + "\n")
			} else {
				line := "  " +
					StyleResultCommand.Render(labelPart) +
					"  " + StyleResultOption.Render(checkbox) +
					"  " + StyleResultDesc.Render(e.desc)
				b.WriteString(line + "\n")
			}
		} else {
			if i == m.selected {
				b.WriteString(StyleConfigItemSelected.Copy().Width(w).Render(
					"  "+labelPart+"  "+valPart+"  "+e.desc,
				) + "\n")
			} else {
				line := "  " +
					StyleResultCommand.Render(labelPart) +
					"  " + StyleResultOption.Render(valPart) +
					"  " + StyleResultDesc.Render(e.desc)
				b.WriteString(line + "\n")
			}
		}
	}

	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Bashrc alias confirmation ──────────────────────────────────────────
	if m.confirmingBashrc {
		b.WriteString("\n")
		b.WriteString(StyleInfo.Render(
			fmt.Sprintf("  Add alias to ~/.bashrc for \"%s\"?", m.pendingName),
		) + "\n")
		b.WriteString(StyleResultOption.Render("  "+m.pendingAliasLine) + "\n")
		b.WriteString(StyleResultDesc.Render("  [y] Yes, add alias    [n] No, skip    Esc to skip") + "\n")
	}

	// ── Edit input (shown only when editing) ───────────────────────────────
	if m.editing {
		e := settingsEntries[m.selected]
		b.WriteString("\n")
		b.WriteString(StyleInputLabel.Render(
			fmt.Sprintf("  Edit \"%s\" (current: %s):", e.label, e.getVal(m.settings)),
		) + "\n")
		b.WriteString("  " + StyleInputFocused.Copy().Width(34).Render(m.input.View()) + "\n")
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
