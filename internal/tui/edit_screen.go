package tui

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
)

// goToEditMsg navigates to the config edit screen for the given config.
type goToEditMsg struct{ cfg *config.Config }

// ---- edit screen level / submode ----

// editLevel represents which tier of the config hierarchy is displayed.
type editLevel int

const (
	editLevelCommands editLevel = iota // list of commands in the config
	editLevelOptions                   // list of options inside a command
	editLevelInputs                    // list of inputs inside an option
)

// editSubmode represents the current interactive mode.
type editSubmode int

const (
	editSubBrowse editSubmode = iota // navigating the list
	editSubForm                      // filling a create/edit form
	editSubDelete                    // waiting for y / n delete confirmation
)

// EditScreenModel is the config-editing TUI screen.
// It lets the user browse and mutate Commands → Options → Inputs within a
// single Config, persisting changes to disk on every save.
type EditScreenModel struct {
	mgr *config.Manager
	cfg *config.Config // pointer – mutations go straight to the loaded config

	level  editLevel
	cmdIdx int // selected command index
	optIdx int // selected option within selected command
	inpIdx int // selected input within selected option

	submode   editSubmode
	formIsNew bool // true = creating new item; false = editing existing
	formTitle  string
	formLabels []string
	formFields []textinput.Model
	formFocus  int

	message string
	width   int
	height  int
}

// NewEditScreenModel creates the editor for cfg.
func NewEditScreenModel(mgr *config.Manager, cfg *config.Config, w, h int) EditScreenModel {
	return EditScreenModel{
		mgr:    mgr,
		cfg:    cfg,
		level:  editLevelCommands,
		width:  w,
		height: h,
	}
}

// Init satisfies tea.Model.
func (m EditScreenModel) Init() tea.Cmd { return nil }

// ---- tea.Model: Update -------------------------------------------------------

func (m EditScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		inputW := max(30, msg.Width-20)
		for i := range m.formFields {
			m.formFields[i].Width = inputW
		}
		return m, nil

	case tea.KeyMsg:
		switch m.submode {
		case editSubBrowse:
			return m.updateBrowse(msg)
		case editSubForm:
			return m.updateForm(msg)
		case editSubDelete:
			return m.updateDelete(msg)
		}
	}
	return m, nil
}

// ---- browse submode ----------------------------------------------------------

func (m EditScreenModel) updateBrowse(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	listLen := m.currentListLen()
	idx := m.currentIdx()

	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		if m.level == editLevelCommands {
			return m, func() tea.Msg { return goToConfigMsg{} }
		}
		m.level--
		m.message = ""
		return m, nil

	case tea.KeyUp:
		if idx > 0 {
			m.setIdx(idx - 1)
		}

	case tea.KeyDown:
		if listLen > 0 && idx < listLen-1 {
			m.setIdx(idx + 1)
		}

	case tea.KeyEnter:
		// Drill down one level (commands → options → inputs).
		if listLen > 0 && m.level < editLevelInputs {
			m.level++
			m.setIdx(0)
			m.message = ""
		}

	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "n", "N":
			newM, initCmd := m.openForm(true)
			return newM, initCmd

		case "e", "E":
			if listLen > 0 {
				newM, initCmd := m.openForm(false)
				return newM, initCmd
			}

		case "d", "D":
			if listLen > 0 {
				m.submode = editSubDelete
				m.message = ""
			}

		case "q", "Q":
			return m, func() tea.Msg { return goToConfigMsg{} }
		}
	}
	return m, nil
}

// ---- delete-confirm submode --------------------------------------------------

func (m EditScreenModel) updateDelete(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		m.submode = editSubBrowse
		m.message = ""
		return m, nil

	case tea.KeyRunes:
		switch string(msg.Runes) {
		case "y", "Y":
			newM, err := m.deleteSelected()
			if err != nil {
				m.message = StyleError.Render("Error: " + err.Error())
				m.submode = editSubBrowse
				return m, nil
			}
			newM.message = StyleInfo.Render("Deleted successfully.")
			newM.submode = editSubBrowse
			return newM, nil

		case "n", "N":
			m.submode = editSubBrowse
			m.message = ""
		}
	}
	return m, nil
}

// ---- form submode ------------------------------------------------------------

func (m EditScreenModel) updateForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		return m, tea.Quit

	case tea.KeyEsc:
		m.submode = editSubBrowse
		m.message = ""
		// Blur all inputs.
		for i := range m.formFields {
			m.formFields[i].Blur()
		}
		return m, nil

	case tea.KeyCtrlS:
		newM, err := m.saveForm()
		if err != nil {
			m.message = StyleError.Render("Error: " + err.Error())
			return m, nil
		}
		newM.message = StyleInfo.Render("Saved.")
		newM.submode = editSubBrowse
		for i := range newM.formFields {
			newM.formFields[i].Blur()
		}
		return newM, nil

	case tea.KeyTab, tea.KeyDown:
		if len(m.formFields) > 0 {
			m.formFields[m.formFocus].Blur()
			m.formFocus = (m.formFocus + 1) % len(m.formFields)
			return m, m.formFields[m.formFocus].Focus()
		}

	case tea.KeyShiftTab, tea.KeyUp:
		if len(m.formFields) > 0 {
			m.formFields[m.formFocus].Blur()
			m.formFocus = (m.formFocus - 1 + len(m.formFields)) % len(m.formFields)
			return m, m.formFields[m.formFocus].Focus()
		}
	}

	// Delegate key to focused text input.
	if m.formFocus < len(m.formFields) {
		var cmd tea.Cmd
		m.formFields[m.formFocus], cmd = m.formFields[m.formFocus].Update(msg)
		return m, cmd
	}
	return m, nil
}

// ---- form helpers ------------------------------------------------------------

// openForm prepares the form fields for the current level.
// isNew == true means creating, false means editing the selected item.
func (m EditScreenModel) openForm(isNew bool) (EditScreenModel, tea.Cmd) {
	m.submode = editSubForm
	m.formIsNew = isNew
	m.formFocus = 0
	inputW := max(30, m.width-20)

	action := "Edit"
	if isNew {
		action = "New"
	}

	switch m.level {
	case editLevelCommands:
		name, desc := "", ""
		if !isNew && m.cmdIdx < len(m.cfg.Commands) {
			name = m.cfg.Commands[m.cmdIdx].Name
			desc = m.cfg.Commands[m.cmdIdx].Description
		}
		m.formTitle = action + " Command"
		m.formLabels = []string{"Name", "Description"}
		m.formFields = makeEditInputs([]string{name, desc}, inputW)

	case editLevelOptions:
		name, desc, tmpl, tags := "", "", "", ""
		if !isNew && m.cmdIdx < len(m.cfg.Commands) {
			cmd := m.cfg.Commands[m.cmdIdx]
			if m.optIdx < len(cmd.Options) {
				name = cmd.Options[m.optIdx].Name
				desc = cmd.Options[m.optIdx].Description
				tmpl = cmd.Options[m.optIdx].Template
				tags = strings.Join(cmd.Options[m.optIdx].Tags, ", ")
			}
		}
		m.formTitle = action + " Option"
		m.formLabels = []string{"Name", "Description", "Template", "Tags  (comma-separated search aliases)"}
		m.formFields = makeEditInputs([]string{name, desc, tmpl, tags}, inputW)

	case editLevelInputs:
		var name, typ, desc, req, def string
		typ = "string"
		req = "false"
		if !isNew && m.cmdIdx < len(m.cfg.Commands) {
			cmd := m.cfg.Commands[m.cmdIdx]
			if m.optIdx < len(cmd.Options) {
				opt := cmd.Options[m.optIdx]
				if m.inpIdx < len(opt.Inputs) {
					inp := opt.Inputs[m.inpIdx]
					name = inp.Name
					typ = inp.Type
					if typ == "" {
						typ = "string"
					}
					desc = inp.Description
					if inp.Required {
						req = "true"
					}
					def = inp.Default
				}
			}
		}
		m.formTitle = action + " Input"
		m.formLabels = []string{
			"Name",
			"Type  (string / file / dir / flag)",
			"Description",
			"Required  (true / false)",
			"Default value",
		}
		m.formFields = makeEditInputs([]string{name, typ, desc, req, def}, inputW)
	}

	return m, m.formFields[0].Focus()
}

// saveForm reads the current form values, applies them to cfg and persists.
func (m EditScreenModel) saveForm() (EditScreenModel, error) {
	switch m.level {
	case editLevelCommands:
		name := strings.TrimSpace(m.formFields[0].Value())
		desc := strings.TrimSpace(m.formFields[1].Value())
		if name == "" {
			return m, fmt.Errorf("name cannot be empty")
		}
		if m.formIsNew {
			m.cfg.Commands = append(m.cfg.Commands, config.Command{
				Name:        name,
				Description: desc,
			})
			m.cmdIdx = len(m.cfg.Commands) - 1
		} else if m.cmdIdx < len(m.cfg.Commands) {
			m.cfg.Commands[m.cmdIdx].Name = name
			m.cfg.Commands[m.cmdIdx].Description = desc
		}

	case editLevelOptions:
		name := strings.TrimSpace(m.formFields[0].Value())
		desc := strings.TrimSpace(m.formFields[1].Value())
		tmpl := strings.TrimSpace(m.formFields[2].Value())
		tagsRaw := strings.TrimSpace(m.formFields[3].Value())
		if name == "" {
			return m, fmt.Errorf("name cannot be empty")
		}
		if m.cmdIdx >= len(m.cfg.Commands) {
			return m, fmt.Errorf("no command selected")
		}
		parsedTags := parseTags(tagsRaw)
		cmd := &m.cfg.Commands[m.cmdIdx]
		if m.formIsNew {
			cmd.Options = append(cmd.Options, config.Option{
				Name:        name,
				Description: desc,
				Template:    tmpl,
				Tags:        parsedTags,
			})
			m.optIdx = len(cmd.Options) - 1
		} else if m.optIdx < len(cmd.Options) {
			cmd.Options[m.optIdx].Name = name
			cmd.Options[m.optIdx].Description = desc
			cmd.Options[m.optIdx].Template = tmpl
			cmd.Options[m.optIdx].Tags = parsedTags
		}

		// Auto-detect {{varName}} placeholders in the template and add any
		// missing inputs as optional string inputs.
		if tmpl != "" && m.optIdx < len(cmd.Options) {
			opt := &cmd.Options[m.optIdx]
			varRe := regexp.MustCompile(`\{\{(\w+)\}\}`)
			existing := make(map[string]bool, len(opt.Inputs))
			for _, inp := range opt.Inputs {
				existing[inp.Name] = true
			}
			for _, sub := range varRe.FindAllStringSubmatch(tmpl, -1) {
				varName := sub[1]
				if !existing[varName] {
					opt.Inputs = append(opt.Inputs, config.Input{
						Name:     varName,
						Type:     "string",
						Required: false,
					})
					existing[varName] = true
				}
			}
		}

	case editLevelInputs:
		name := strings.TrimSpace(m.formFields[0].Value())
		typ := strings.TrimSpace(m.formFields[1].Value())
		desc := strings.TrimSpace(m.formFields[2].Value())
		reqStr := strings.TrimSpace(m.formFields[3].Value())
		def := strings.TrimSpace(m.formFields[4].Value())
		if name == "" {
			return m, fmt.Errorf("name cannot be empty")
		}
		if typ == "" {
			typ = "string"
		}
		req := strings.ToLower(reqStr) == "true"
		if m.cmdIdx >= len(m.cfg.Commands) {
			return m, fmt.Errorf("no command selected")
		}
		cmd := &m.cfg.Commands[m.cmdIdx]
		if m.optIdx >= len(cmd.Options) {
			return m, fmt.Errorf("no option selected")
		}
		opt := &cmd.Options[m.optIdx]
		inp := config.Input{
			Name:        name,
			Type:        typ,
			Description: desc,
			Required:    req,
			Default:     def,
		}
		if m.formIsNew {
			opt.Inputs = append(opt.Inputs, inp)
			m.inpIdx = len(opt.Inputs) - 1
		} else if m.inpIdx < len(opt.Inputs) {
			opt.Inputs[m.inpIdx] = inp
		}
	}

	return m, m.mgr.UpdateConfig(m.cfg)
}

// deleteSelected removes the currently highlighted item and saves.
func (m EditScreenModel) deleteSelected() (EditScreenModel, error) {
	switch m.level {
	case editLevelCommands:
		if m.cmdIdx >= len(m.cfg.Commands) {
			return m, nil
		}
		m.cfg.Commands = append(
			m.cfg.Commands[:m.cmdIdx],
			m.cfg.Commands[m.cmdIdx+1:]...,
		)
		if m.cmdIdx > 0 && m.cmdIdx >= len(m.cfg.Commands) {
			m.cmdIdx--
		}

	case editLevelOptions:
		if m.cmdIdx >= len(m.cfg.Commands) {
			return m, nil
		}
		cmd := &m.cfg.Commands[m.cmdIdx]
		if m.optIdx >= len(cmd.Options) {
			return m, nil
		}
		cmd.Options = append(cmd.Options[:m.optIdx], cmd.Options[m.optIdx+1:]...)
		if m.optIdx > 0 && m.optIdx >= len(cmd.Options) {
			m.optIdx--
		}

	case editLevelInputs:
		if m.cmdIdx >= len(m.cfg.Commands) {
			return m, nil
		}
		cmd := &m.cfg.Commands[m.cmdIdx]
		if m.optIdx >= len(cmd.Options) {
			return m, nil
		}
		opt := &cmd.Options[m.optIdx]
		if m.inpIdx >= len(opt.Inputs) {
			return m, nil
		}
		opt.Inputs = append(opt.Inputs[:m.inpIdx], opt.Inputs[m.inpIdx+1:]...)
		if m.inpIdx > 0 && m.inpIdx >= len(opt.Inputs) {
			m.inpIdx--
		}
	}

	return m, m.mgr.UpdateConfig(m.cfg)
}

// ---- index helpers ----------------------------------------------------------

func (m EditScreenModel) currentListLen() int {
	switch m.level {
	case editLevelCommands:
		return len(m.cfg.Commands)
	case editLevelOptions:
		if m.cmdIdx < len(m.cfg.Commands) {
			return len(m.cfg.Commands[m.cmdIdx].Options)
		}
	case editLevelInputs:
		if m.cmdIdx < len(m.cfg.Commands) {
			cmd := m.cfg.Commands[m.cmdIdx]
			if m.optIdx < len(cmd.Options) {
				return len(cmd.Options[m.optIdx].Inputs)
			}
		}
	}
	return 0
}

func (m EditScreenModel) currentIdx() int {
	switch m.level {
	case editLevelCommands:
		return m.cmdIdx
	case editLevelOptions:
		return m.optIdx
	case editLevelInputs:
		return m.inpIdx
	}
	return 0
}

func (m *EditScreenModel) setIdx(i int) {
	switch m.level {
	case editLevelCommands:
		m.cmdIdx = i
	case editLevelOptions:
		m.optIdx = i
	case editLevelInputs:
		m.inpIdx = i
	}
}

// makeEditInputs creates a slice of text inputs with preset values.
func makeEditInputs(values []string, width int) []textinput.Model {
	inputs := make([]textinput.Model, len(values))
	for i, v := range values {
		ti := textinput.New()
		ti.Width = width
		ti.SetValue(v)
		inputs[i] = ti
	}
	return inputs
}

// parseTags splits a comma-separated tag string into a trimmed, non-empty slice.
func parseTags(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var tags []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}

// ---- View -------------------------------------------------------------------

func (m EditScreenModel) View() string {
	w := m.width
	if w <= 0 {
		w = 80
	}
	h := m.height
	if h <= 0 {
		h = 24
	}

	var b strings.Builder

	// Title.
	levelLabel := "Commands"
	switch m.level {
	case editLevelOptions:
		levelLabel = "Options"
	case editLevelInputs:
		levelLabel = "Inputs"
	}
	title := StyleTitle.Copy().Width(w).Render(
		"⚡ " + AppDisplayName + "  " +
			StyleResultDesc.Render("Edit Config: "+m.cfg.Name+" › "+levelLabel),
	)
	b.WriteString(title + "\n")
	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// Breadcrumb context bar.
	if m.level > editLevelCommands {
		ctx := ""
		if m.cmdIdx < len(m.cfg.Commands) {
			ctx = StyleResultCommand.Render(m.cfg.Commands[m.cmdIdx].Name)
		}
		if m.level == editLevelInputs {
			if m.cmdIdx < len(m.cfg.Commands) {
				cmd := m.cfg.Commands[m.cmdIdx]
				if m.optIdx < len(cmd.Options) {
					ctx += StyleResultDesc.Render(" › ") + StyleResultOption.Render(cmd.Options[m.optIdx].Name)
				}
			}
		}
		b.WriteString("  " + ctx + "\n")
		b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")
	}

	// Main content area.
	switch m.submode {
	case editSubForm:
		b.WriteString(m.renderForm(w))
	case editSubDelete:
		b.WriteString(m.renderList(w, h))
		b.WriteString("\n" + StyleError.Render("  Delete selected item? (y / n)") + "\n")
	default:
		b.WriteString(m.renderList(w, h))
	}

	// Message line.
	if m.message != "" {
		b.WriteString("\n  " + m.message + "\n")
	}

	// Fill remaining vertical space then render status bar.
	used := strings.Count(b.String(), "\n") + 2
	for i := used; i < h-1; i++ {
		b.WriteString("\n")
	}

	var keys string
	switch m.submode {
	case editSubForm:
		keys = StyleStatusKey.Render(" Ctrl+S") + StyleStatus.Render(" save") +
			StyleStatusKey.Render("  Tab") + StyleStatus.Render(" next field") +
			StyleStatusKey.Render("  Esc") + StyleStatus.Render(" cancel")
	case editSubDelete:
		keys = StyleStatusKey.Render(" y") + StyleStatus.Render(" confirm delete") +
			StyleStatusKey.Render("  n / Esc") + StyleStatus.Render(" cancel")
	default:
		keys = StyleStatusKey.Render(" n") + StyleStatus.Render(" new") +
			StyleStatusKey.Render("  e") + StyleStatus.Render(" edit") +
			StyleStatusKey.Render("  d") + StyleStatus.Render(" delete")
		if m.level < editLevelInputs {
			keys += StyleStatusKey.Render("  Enter") + StyleStatus.Render(" open")
		}
		keys += StyleStatusKey.Render("  Esc") + StyleStatus.Render(" back")
	}
	b.WriteString(renderFooter(w, keys, footerVersion()))

	return b.String()
}

// renderList renders the browseable list for the current level.
func (m EditScreenModel) renderList(w, h int) string {
	var b strings.Builder

	type row struct{ text string }
	var rows []row

	switch m.level {
	case editLevelCommands:
		b.WriteString(StyleConfigHeader.Render("Commands") + "\n")
		for _, cmd := range m.cfg.Commands {
			rows = append(rows, row{
				fmt.Sprintf("%-22s  %s", cmd.Name, cmd.Description),
			})
		}

	case editLevelOptions:
		if m.cmdIdx < len(m.cfg.Commands) {
			cmd := m.cfg.Commands[m.cmdIdx]
			b.WriteString(StyleConfigHeader.Render("Options in "+cmd.Name) + "\n")
			for _, opt := range cmd.Options {
				tagStr := ""
				if len(opt.Tags) > 0 {
					tagStr = "  [" + strings.Join(opt.Tags, ", ") + "]"
				}
				rows = append(rows, row{
					fmt.Sprintf("%-22s  %-28s  %s", opt.Name, opt.Description, opt.Template) + tagStr,
				})
			}
		}

	case editLevelInputs:
		if m.cmdIdx < len(m.cfg.Commands) {
			cmd := m.cfg.Commands[m.cmdIdx]
			if m.optIdx < len(cmd.Options) {
				opt := cmd.Options[m.optIdx]
				b.WriteString(StyleConfigHeader.Render("Inputs for option: "+opt.Name) + "\n")
				for _, inp := range opt.Inputs {
					req := ""
					if inp.Required {
						req = StyleInputLabelRequired.Render(" *")
					}
					rows = append(rows, row{
						fmt.Sprintf("%-18s  %-8s  %s", inp.Name, inp.Type, inp.Description) + req,
					})
				}
			}
		}
	}

	idx := m.currentIdx()
	visRows := h - 12
	if visRows < 1 {
		visRows = 1
	}

	if len(rows) == 0 {
		b.WriteString(StyleResultDesc.Padding(1, 2).Render("No items yet. Press n to add one.") + "\n")
	} else {
		for i, r := range rows {
			if i >= visRows {
				break
			}
			if i == idx {
				b.WriteString(StyleConfigItemSelected.Copy().Width(w).Render(r.text) + "\n")
			} else {
				b.WriteString(StyleConfigItem.Copy().Width(w).Render(r.text) + "\n")
			}
		}
	}

	return b.String()
}

// renderForm renders the create / edit form.
func (m EditScreenModel) renderForm(w int) string {
	var b strings.Builder

	b.WriteString(StyleConfigHeader.Render(m.formTitle) + "\n\n")

	for i, label := range m.formLabels {
		b.WriteString("  " + StyleInputLabel.Render(label) + "\n")
		var fieldView string
		if i == m.formFocus {
			fieldView = StyleInputFocused.Copy().Width(w - 6).Render(m.formFields[i].View())
		} else {
			fieldView = StyleInputBlurred.Copy().Width(w - 6).Render(m.formFields[i].View())
		}
		b.WriteString(fieldView + "\n\n")
	}

	return b.String()
}
