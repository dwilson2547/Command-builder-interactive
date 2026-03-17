package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
)

// SubCmdItem represents a single entry returned by a sub-command picker.
type SubCmdItem struct {
	Value  string // inserted into the input when selected
	Detail string // display-only detail (may be empty)
}

// subCmdResultMsg carries the result of an async sub-command execution.
type subCmdResultMsg struct {
	items    []SubCmdItem
	err      error
	inputIdx int // index of the input that triggered the run
}

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

	// subCmd picker state.
	subCmdItems   []SubCmdItem
	subCmdIdx     int
	showSubCmd    bool
	loadingSubCmd bool
	subCmdErr     string

	// star naming state — active while the user is entering a custom name.
	starNaming    bool
	starNameInput textinput.Model

	builtCmd   string
	formMsg    string // transient status message (e.g. star confirmation)
	runOnEnter bool
	width      int
	height     int
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

	nameInput := textinput.New()
	nameInput.Width = max(30, w-20)
	nameInput.Placeholder = cmd.Name + " › " + opt.Name

	m := FormModel{
		cfg:           cfg,
		cmd:           cmd,
		opt:           opt,
		inputs:        inputs,
		flagStates:    flagStates,
		starNameInput: nameInput,
		runOnEnter:    runOnEnter,
		width:         w,
		height:        h,
	}
	m.rebuildCommand()
	return m
}

// NewPrefilledFormModel creates a form pre-filled with the saved values from a
// starred command, so the user can review or tweak inputs before running.
func NewPrefilledFormModel(cfg *config.Config, cmd *config.Command, opt *config.Option, w, h int, runOnEnter bool, values map[string]string, flagStates map[string]bool) FormModel {
	m := NewFormModel(cfg, cmd, opt, w, h, runOnEnter)
	for i, inp := range opt.Inputs {
		if inp.Type == "flag" {
			if on, ok := flagStates[inp.Name]; ok {
				m.flagStates[i] = on
			}
		} else {
			if val, ok := values[inp.Name]; ok && val != "" {
				m.inputs[i].SetValue(val)
			}
		}
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

	case subCmdResultMsg:
		m.loadingSubCmd = false
		if msg.inputIdx != m.focus {
			// User moved away while command was running; discard.
			return m, nil
		}
		if msg.err != nil {
			m.subCmdErr = msg.err.Error()
			m.subCmdItems = nil
		} else {
			m.subCmdErr = ""
			m.subCmdItems = msg.items
		}
		m.subCmdIdx = 0
		m.showSubCmd = true
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputW := max(30, msg.Width-20)
		for i := range m.inputs {
			m.inputs[i].Width = inputW
		}
		m.starNameInput.Width = inputW
		return m, nil

	case tea.KeyMsg:
		// ── Star naming mode: capture the custom name before saving. ─────────
		if m.starNaming {
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEsc:
				m.starNaming = false
				m.starNameInput.SetValue("")
				return m, nil
			case tea.KeyEnter:
				name := strings.TrimSpace(m.starNameInput.Value())
				m.starNaming = false
				m.starNameInput.SetValue("")
				return m.starCurrentCommand(name)
			default:
				var cmd tea.Cmd
				m.starNameInput, cmd = m.starNameInput.Update(msg)
				return m, cmd
			}
		}

		switch msg.Type {

		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEsc:
			if m.showSubCmd {
				m.showSubCmd = false
				return m, nil
			}
			return m, func() tea.Msg { return backToSearchMsg{} }

		case tea.KeySpace:
			// Toggle flag-type inputs.
			if m.focus < len(m.opt.Inputs) && m.opt.Inputs[m.focus].Type == "flag" {
				m.flagStates[m.focus] = !m.flagStates[m.focus]
				m.rebuildCommand()
				return m, nil
			}

		case tea.KeyEnter:
			if m.showSubCmd && len(m.subCmdItems) > 0 {
				item := m.subCmdItems[m.subCmdIdx]
				m.inputs[m.focus].SetValue(item.Value)
				m.showSubCmd = false
				m.rebuildCommand()
				return m, nil
			}
			if m.builtCmd != "" && m.allRequiredFilled() {
				cmd := m.builtCmd
				return m, func() tea.Msg { return commandConfirmedMsg{command: cmd} }
			}
			// Move to next input on Enter.
			m = m.moveFocus(1)
			return m, nil

		case tea.KeyTab:
			if m.showSubCmd {
				// Close the picker and move to the next field.
				m.showSubCmd = false
				m = m.moveFocus(1)
				return m, nil
			}
			if m.showCompletion {
				// Cycle through completions.
				if len(m.completions) > 0 {
					m.compIdx = (m.compIdx + 1) % len(m.completions)
					m.inputs[m.focus].SetValue(m.completions[m.compIdx])
					m.rebuildCommand()
				}
				return m, nil
			}
			// If current field has a sub_command, trigger the picker.
			if m.focus < len(m.opt.Inputs) {
				inp := m.opt.Inputs[m.focus]
				if inp.SubCommand != "" && !m.loadingSubCmd {
					m.loadingSubCmd = true
					m.showCompletion = false
					idx := m.focus
					subCmd := inp.SubCommand
					return m, func() tea.Msg {
						return runSubCommand(subCmd, idx)
					}
				}
				// If current field is file/dir, compute completions.
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
			if m.showSubCmd {
				if m.subCmdIdx > 0 {
					m.subCmdIdx--
				}
				return m, nil
			}
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
			if m.showSubCmd {
				if m.subCmdIdx < len(m.subCmdItems)-1 {
					m.subCmdIdx++
				}
				return m, nil
			}
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
			// Enter star naming mode — prompt for a custom name.
			if msg.Type == tea.KeyRunes && string(msg.Runes) == "*" {
				m.starNaming = true
				m.starNameInput.SetValue("")
				m.starNameInput.Placeholder = m.cmd.Name + " › " + m.opt.Name
				return m, m.starNameInput.Focus()
			}
			m.showCompletion = false
			m.showSubCmd = false
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
	title := StyleTitle.Width(w).Render(
		"⚡ " + AppDisplayName + StyleResultDesc.Render("  ← Esc to go back"),
	)
	b.WriteString(title + "\n")

	// ── Command header ────────────────────────────────────────────────────
	header := StyleFormHeader.Width(w - 4).Render(
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
			b.WriteString(rowStyle.Width(w-6).Render(checkbox) + "\n\n")
			continue
		}

		var inputView string
		if i == m.focus {
			inputView = StyleInputFocused.Width(w - 6).Render(m.inputs[i].View())
		} else {
			inputView = StyleInputBlurred.Width(w - 6).Render(m.inputs[i].View())
		}

		b.WriteString("  " + label + desc + "\n")
		b.WriteString(inputView + "\n")

		// Show completions, loading indicator, or sub-command picker below focused field.
		if i == m.focus {
			if m.loadingSubCmd {
				b.WriteString(StyleSubCmdBox.Width(w-6).Render(
					StyleSubCmdLoading.Render("Loading…"),
				) + "\n")
			} else if m.showSubCmd {
				b.WriteString(m.renderSubCmdPicker(w) + "\n")
			} else if m.showCompletion && len(m.completions) > 0 {
				b.WriteString(m.renderCompletions(w) + "\n")
			}
		}
		b.WriteString("\n")
	}

	// ── Hint for focused field ────────────────────────────────────────────
	if m.focus < len(m.opt.Inputs) {
		inp := m.opt.Inputs[m.focus]
		if inp.SubCommand != "" {
			if m.showSubCmd {
				hint := StyleResultDesc.Render("  ↑↓: navigate  Enter: select  Esc: close")
				b.WriteString(hint + "\n")
			} else if !m.loadingSubCmd {
				hint := StyleResultDesc.Render("  Tab: pick value")
				b.WriteString(hint + "\n")
			}
		} else if inp.Type == "file" || inp.Type == "dir" {
			hint := StyleResultDesc.Render("  Tab: path completion")
			b.WriteString(hint + "\n")
		}
	}

	// ── Command preview ───────────────────────────────────────────────────
	b.WriteString("\n")
	if m.builtCmd != "" {
		previewLabel := StylePreviewLabel.Render("  $ ")
		preview := StylePreviewBox.Width(w - 4).Render(m.builtCmd)
		b.WriteString(previewLabel + "\n")
		b.WriteString(preview + "\n")
	}

	if m.starNaming {
		prompt := StyleInputLabel.Render("  Name this star") +
			StyleResultDesc.Render("  (leave blank to keep default, Esc to cancel)")
		nameView := StyleInputFocused.Width(w - 6).Render(m.starNameInput.View())
		b.WriteString("\n" + prompt + "\n")
		b.WriteString(nameView + "\n")
	} else if m.formMsg != "" {
		b.WriteString("\n  " + m.formMsg + "\n")
	}

	// ── Status bar ────────────────────────────────────────────────────────
	// Fill remaining height.
	used := strings.Count(b.String(), "\n") + 2
	for i := used; i < h-1; i++ {
		b.WriteString("\n")
	}

	var statusMsg string
	if m.starNaming {
		statusMsg = StyleStatusKey.Render(" Enter") + StyleStatus.Render(" confirm name") +
			StyleStatusKey.Render("  Esc") + StyleStatus.Render(" cancel")
	} else if m.allRequiredFilled() && m.builtCmd != "" {
		enterLabel := "copy command & quit"
		if m.runOnEnter {
			enterLabel = "run command & quit"
		}
		statusMsg = StyleStatusKey.Render(" Enter") + StyleStatus.Render(" "+enterLabel) +
			StyleStatusKey.Render("  *") + StyleStatus.Render(" star") +
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
	return StyleCompletionBox.Width(width - 6).Render(content)
}

func (m FormModel) renderSubCmdPicker(width int) string {
	const maxVisible = 8

	if m.subCmdErr != "" {
		return StyleSubCmdBox.Width(width - 6).Render(
			StyleSubCmdLoading.Render("Error: " + m.subCmdErr),
		)
	}
	if len(m.subCmdItems) == 0 {
		return StyleSubCmdBox.Width(width - 6).Render(
			StyleSubCmdLoading.Render("No results"),
		)
	}

	// Compute scroll window around the selected index.
	start := 0
	if m.subCmdIdx >= maxVisible {
		start = m.subCmdIdx - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(m.subCmdItems) {
		end = len(m.subCmdItems)
	}
	visible := m.subCmdItems[start:end]

	// Find the longest value in the visible window for alignment.
	maxVal := 0
	for _, item := range visible {
		if len(item.Value) > maxVal {
			maxVal = len(item.Value)
		}
	}

	innerWidth := max(20, width-10)

	var rows []string
	for i, item := range visible {
		absIdx := start + i
		selected := absIdx == m.subCmdIdx

		var row string
		if item.Detail != "" {
			pad := strings.Repeat(" ", max(1, maxVal-len(item.Value)+2))
			row = item.Value + pad + item.Detail
		} else {
			row = item.Value
		}

		if selected {
			rows = append(rows, StyleSubCmdSelected.Width(innerWidth).Render(row))
		} else {
			rows = append(rows, StyleSubCmdItem.Width(innerWidth).Render(row))
		}
	}

	if len(m.subCmdItems) > maxVisible {
		rows = append(rows, StyleResultDesc.Render(
			fmt.Sprintf("  %d / %d", m.subCmdIdx+1, len(m.subCmdItems)),
		))
	}

	return StyleSubCmdBox.Width(width - 6).Render(strings.Join(rows, "\n"))
}

// runSubCommand executes a shell command and parses its stdout as CSV, returning
// a subCmdResultMsg. Each line is split on the first comma: col 0 = value,
// col 1 (optional) = detail.
func runSubCommand(command string, inputIdx int) subCmdResultMsg {
	out, err := exec.Command("sh", "-c", command).Output()
	if err != nil {
		return subCmdResultMsg{err: err, inputIdx: inputIdx}
	}

	var items []SubCmdItem
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 2)
		item := SubCmdItem{Value: strings.TrimSpace(parts[0])}
		if len(parts) == 2 {
			item.Detail = strings.TrimSpace(parts[1])
		}
		if item.Value != "" {
			items = append(items, item)
		}
	}
	return subCmdResultMsg{items: items, inputIdx: inputIdx}
}

// starCurrentCommand saves the form's current input values as a starred command.
// customName is shown in the stars list; leave empty to use the default label.
func (m FormModel) starCurrentCommand(customName string) (tea.Model, tea.Cmd) {
	values := make(map[string]string)
	flagStates := make(map[string]bool)
	for i, inp := range m.opt.Inputs {
		if inp.Type == "flag" {
			flagStates[inp.Name] = m.flagStates[i]
		} else {
			values[inp.Name] = m.inputs[i].Value()
		}
	}
	star := config.Star{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()),
		ConfigName:  m.cfg.Name,
		CommandName: m.cmd.Name,
		OptionName:  m.opt.Name,
		Label:       m.cmd.Name + " › " + m.opt.Name,
		CustomName:  customName,
		Values:      values,
		FlagStates:  flagStates,
		CreatedAt:   time.Now(),
	}
	if err := config.AddStar(star); err != nil {
		m.formMsg = StyleError.Render("Failed to star: " + err.Error())
	} else {
		m.formMsg = StyleInfo.Render("★ Starred! (view with /s)")
	}
	return m, nil
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
	m.showSubCmd = false
	m.loadingSubCmd = false
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
