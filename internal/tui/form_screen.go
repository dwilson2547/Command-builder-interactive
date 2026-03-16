package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
)

// FormModel is the command-builder form screen.
type FormModel struct {
	cfg    *config.Config
	cmd    *config.Command
	opt    *config.Option
	inputs []textinput.Model
	focus  int // index of focused input

	// flagStates tracks the on/off state of each input with type "flag".
	// Indexed in parallel with opt.Inputs.
	flagStates []bool

	completions    []string
	compIdx        int
	showCompletion bool

	builtCmd    string
	runOnEnter  bool
	width       int
	height      int
	scroll      int // vertical scroll offset for the form
}

// NewFormModel creates a new form for the given option.
func NewFormModel(cfg *config.Config, cmd *config.Command, opt *config.Option, w, h int, runOnEnter bool) FormModel {
	inputs := make([]textinput.Model, len(opt.Inputs))
	for i, inp := range opt.Inputs {
		ti := textinput.New()
		ti.Width = max(30, w-20)
		if inp.Default != "" {
			ti.SetValue(inp.Default)
			ti.Placeholder = inp.Default
		} else {
			ti.Placeholder = inp.Description
		}
		if inp.Type == "string" && strings.Contains(strings.ToLower(inp.Name), "password") {
			ti.EchoMode = textinput.EchoPassword
		}
		inputs[i] = ti
	}

	flagStates := make([]bool, len(opt.Inputs))
	m := FormModel{
		cfg:        cfg,
		cmd:        cmd,
		opt:        opt,
		inputs:     inputs,
		flagStates: flagStates,
		runOnEnter: runOnEnter,
		width:      w,
		height:     h,
	}
	m.rebuildCommand()
	return m
}

// Init satisfies tea.Model – focus the first input.
func (m FormModel) Init() tea.Cmd {
	if len(m.inputs) == 0 {
		return nil
	}
	return m.inputs[0].Focus()
}

// Update satisfies tea.Model.
func (m FormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputW := max(30, msg.Width-20)
		for i := range m.inputs {
			m.inputs[i].Width = inputW
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			return m, func() tea.Msg { return backToSearchMsg{} }

		case tea.KeySpace:
			// Toggle flag-type inputs.
			if m.focus < len(m.opt.Inputs) && m.opt.Inputs[m.focus].Type == "flag" {
				m.flagStates[m.focus] = !m.flagStates[m.focus]
				m.rebuildCommand()
				return m, nil
			}

		case tea.KeyEnter:
			if m.builtCmd != "" && m.allRequiredFilled() {
				cmd := m.builtCmd
				return m, func() tea.Msg { return commandConfirmedMsg{command: cmd} }
			}
			// Move to next input on Enter.
			m = m.moveFocus(1)
			return m, nil

		case tea.KeyTab:
			if m.showCompletion {
				// Cycle through completions.
				if len(m.completions) > 0 {
					m.compIdx = (m.compIdx + 1) % len(m.completions)
					m.inputs[m.focus].SetValue(m.completions[m.compIdx])
					m.rebuildCommand()
				}
				return m, nil
			}
			// If current field is file/dir, compute completions.
			if m.focus < len(m.opt.Inputs) {
				inp := m.opt.Inputs[m.focus]
				if inp.Type == "file" || inp.Type == "dir" {
					prefix := m.inputs[m.focus].Value()
					comps := getPathCompletions(prefix, inp.Type == "dir")
					if len(comps) == 1 {
						m.inputs[m.focus].SetValue(comps[0])
						m.rebuildCommand()
						return m, nil
					} else if len(comps) > 1 {
						// Complete to longest common prefix.
						cp := commonPrefix(comps)
						if cp != prefix {
							m.inputs[m.focus].SetValue(cp)
							m.rebuildCommand()
						}
						m.completions = comps
						m.compIdx = 0
						m.showCompletion = true
						return m, nil
					}
				}
			}
			// Otherwise move to next field.
			m = m.moveFocus(1)
			return m, nil

		case tea.KeyShiftTab:
			m.showCompletion = false
			m = m.moveFocus(-1)
			return m, nil

		case tea.KeyUp:
			if m.showCompletion {
				if m.compIdx > 0 {
					m.compIdx--
					m.inputs[m.focus].SetValue(m.completions[m.compIdx])
					m.rebuildCommand()
				}
				return m, nil
			}
			m = m.moveFocus(-1)
			return m, nil

		case tea.KeyDown:
			if m.showCompletion {
				if m.compIdx < len(m.completions)-1 {
					m.compIdx++
					m.inputs[m.focus].SetValue(m.completions[m.compIdx])
					m.rebuildCommand()
				}
				return m, nil
			}
			m = m.moveFocus(1)
			return m, nil

		default:
			// Any other key hides completions.
			m.showCompletion = false
		}
	}

	// Delegate to focused input (skip flag types — they use Space to toggle).
	if m.focus < len(m.inputs) && m.focus < len(m.opt.Inputs) && m.opt.Inputs[m.focus].Type != "flag" {
		var cmd tea.Cmd
		m.inputs[m.focus], cmd = m.inputs[m.focus].Update(msg)
		m.rebuildCommand()
		// Recompute completions live — but only when the user is NOT already
		// browsing the completion list. If showCompletion is true, the input
		// value has been set to the currently-highlighted completion; recomputing
		// here would replace the original listing with that subdirectory's
		// contents and break navigation.
		if !m.showCompletion {
			inp := m.opt.Inputs[m.focus]
			if inp.Type == "file" || inp.Type == "dir" {
				prefix := m.inputs[m.focus].Value()
				m.completions = getPathCompletions(prefix, inp.Type == "dir")
			}
		}
		return m, cmd
	}
	return m, nil
}

// View satisfies tea.Model.
func (m FormModel) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	var b strings.Builder

	// ── Title bar ─────────────────────────────────────────────────────────
	title := StyleTitle.Copy().Width(w).Render(
		"⚡ " + AppDisplayName + StyleResultDesc.Render("  ← Esc to go back"),
	)
	b.WriteString(title + "\n")

	// ── Command header ────────────────────────────────────────────────────
	header := StyleFormHeader.Copy().Width(w - 4).Render(
		StyleResultCommand.Render(m.cmd.Name) +
			"  " + StyleResultOption.Render(m.opt.Name) +
			"\n" + StyleResultDesc.Render(m.opt.Description),
	)
	b.WriteString(header + "\n")

	// ── Inputs ────────────────────────────────────────────────────────────
	for i, inp := range m.opt.Inputs {
		req := ""
		labelStyle := StyleInputLabel
		if inp.Required {
			req = StyleInputLabelRequired.Render(" *")
			labelStyle = StyleInputLabelRequired
		}
		label := labelStyle.Render(inp.Name) + req
		desc := StyleInputDesc.Render("  " + inp.Description)

		if inp.Type == "flag" {
			on := m.flagStates[i]
			focused := i == m.focus

			var checkbox string
			var rowStyle lipgloss.Style
			if focused {
				if on {
					checkbox = "✓ " + inp.Description + "  (Space: toggle)"
				} else {
					checkbox = "  " + inp.Description + "  (Space: toggle)"
				}
				rowStyle = StyleFlagFocused
			} else if on {
				checkbox = "✓ " + inp.Description
				rowStyle = StyleFlagOn
			} else {
				checkbox = "  " + inp.Description
				rowStyle = StyleFlagOff
			}

			b.WriteString("  " + label + "\n")
			b.WriteString(rowStyle.Copy().Width(w-6).Render(checkbox) + "\n\n")
			continue
		}

		var inputView string
		if i == m.focus {
			inputView = StyleInputFocused.Copy().Width(w - 6).Render(m.inputs[i].View())
		} else {
			inputView = StyleInputBlurred.Copy().Width(w - 6).Render(m.inputs[i].View())
		}

		b.WriteString("  " + label + desc + "\n")
		b.WriteString(inputView + "\n")

		// Show completions below focused file/dir field.
		if i == m.focus && m.showCompletion && len(m.completions) > 0 {
			b.WriteString(m.renderCompletions(w) + "\n")
		}
		b.WriteString("\n")
	}

	// ── Tab hint for file/dir fields ──────────────────────────────────────
	if m.focus < len(m.opt.Inputs) {
		inp := m.opt.Inputs[m.focus]
		if inp.Type == "file" || inp.Type == "dir" {
			hint := StyleResultDesc.Render("  Tab: path completion")
			b.WriteString(hint + "\n")
		}
	}

	// ── Command preview ───────────────────────────────────────────────────
	b.WriteString("\n")
	if m.builtCmd != "" {
		previewLabel := StylePreviewLabel.Render("  $ ")
		preview := StylePreviewBox.Copy().Width(w - 4).Render(m.builtCmd)
		b.WriteString(previewLabel + "\n")
		b.WriteString(preview + "\n")
	}

	// ── Status bar ────────────────────────────────────────────────────────
	// Fill remaining height.
	used := strings.Count(b.String(), "\n") + 2
	for i := used; i < h-1; i++ {
		b.WriteString("\n")
	}

	var statusMsg string
	if m.allRequiredFilled() && m.builtCmd != "" {
		enterLabel := "copy command & quit"
		if m.runOnEnter {
			enterLabel = "run command & quit"
		}
		statusMsg = StyleStatusKey.Render(" Enter") + StyleStatus.Render(" "+enterLabel) +
			StyleStatusKey.Render("  Esc") + StyleStatus.Render(" back")
	} else {
		missing := m.missingRequired()
		if len(missing) > 0 {
			statusMsg = StyleStatus.Render(" Required: ") + StyleError.Render(strings.Join(missing, ", "))
		}
	}
	b.WriteString(renderFooter(w, statusMsg, footerVersion()))

	return b.String()
}

func (m FormModel) renderCompletions(width int) string {
	maxShow := 6
	items := m.completions
	if len(items) > maxShow {
		items = items[:maxShow]
	}
	var rows []string
	for i, c := range items {
		var s string
		if i == m.compIdx {
			s = StyleCompletionSelected.Render(c)
		} else {
			s = StyleCompletionItem.Render(c)
		}
		rows = append(rows, s)
	}
	if len(m.completions) > maxShow {
		rows = append(rows, StyleResultDesc.Render(fmt.Sprintf("  … %d more", len(m.completions)-maxShow)))
	}
	content := strings.Join(rows, "\n")
	return StyleCompletionBox.Copy().Width(width - 6).Render(content)
}

func (m FormModel) allRequiredFilled() bool {
	for i, inp := range m.opt.Inputs {
		if inp.Required && strings.TrimSpace(m.inputs[i].Value()) == "" {
			return false
		}
	}
	return true
}

func (m FormModel) missingRequired() []string {
	var missing []string
	for i, inp := range m.opt.Inputs {
		if inp.Required && strings.TrimSpace(m.inputs[i].Value()) == "" {
			missing = append(missing, inp.Name)
		}
	}
	return missing
}

// spaceCollapser collapses consecutive whitespace into a single space.
var spaceCollapser = regexp.MustCompile(`\s{2,}`)

func (m *FormModel) rebuildCommand() {
	result := m.opt.Template
	for i, inp := range m.opt.Inputs {
		var val string
		if inp.Type == "flag" {
			// Flags: use the Default string when toggled on, empty when off.
			if m.flagStates[i] {
				val = inp.Default
			}
		} else {
			val = m.inputs[i].Value()
		}

		if val == "" {
			if inp.Required {
				// Required but unfilled: show placeholder.
				result = strings.ReplaceAll(result, "{"+"{" +inp.Name+"}}", "<"+inp.Name+">")
			} else {
				// Optional and empty: omit entirely from the command.
				result = strings.ReplaceAll(result, "{"+"{" +inp.Name+"}}", "")
			}
		} else {
			result = strings.ReplaceAll(result, "{"+"{" +inp.Name+"}}", val)
		}
	}
	// Clean up extra whitespace left by omitted optional arguments.
	result = strings.TrimSpace(spaceCollapser.ReplaceAllString(result, " "))
	m.builtCmd = result
}

func (m FormModel) moveFocus(delta int) FormModel {
	if len(m.inputs) == 0 {
		return m
	}
	m.inputs[m.focus].Blur()
	m.focus = (m.focus + delta + len(m.inputs)) % len(m.inputs)
	m.inputs[m.focus].Focus()
	m.showCompletion = false
	return m
}

// getPathCompletions returns filesystem paths matching the given prefix.
func getPathCompletions(prefix string, dirsOnly bool) []string {
	if prefix == "" {
		entries, err := os.ReadDir(".")
		if err != nil {
			return nil
		}
		var out []string
		for _, e := range entries {
			if dirsOnly && !e.IsDir() {
				continue
			}
			name := e.Name()
			if e.IsDir() {
				name += "/"
			}
			out = append(out, name)
		}
		return out
	}

	dir := filepath.Dir(prefix)
	base := filepath.Base(prefix)

	// If prefix ends with "/" treat the whole thing as the directory.
	if strings.HasSuffix(prefix, "/") {
		dir = prefix
		base = ""
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var out []string
	for _, e := range entries {
		if dirsOnly && !e.IsDir() {
			continue
		}
		if !strings.HasPrefix(e.Name(), base) {
			continue
		}
		full := filepath.Join(dir, e.Name())
		if e.IsDir() {
			full += "/"
		}
		out = append(out, full)
	}
	return out
}

// commonPrefix returns the longest common prefix of all strings.
func commonPrefix(strs []string) string {
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

