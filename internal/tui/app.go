package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
)

// ---- inter-screen messages ----

// selectOptionMsg is sent when the user selects a search result.
type selectOptionMsg struct {
	cfg *config.Config
	cmd *config.Command
	opt *config.Option
}

// goToConfigMsg navigates to the config management screen.
type goToConfigMsg struct{}

// backToSearchMsg navigates back to the search screen.
type backToSearchMsg struct{}

// commandConfirmedMsg carries the final built command when the user confirms.
type commandConfirmedMsg struct{ command string }

// importURLMsg requests an import-from-URL operation.
type importURLMsg struct{ url string }

// importFileMsg requests an import-from-local-file operation.
type importFileMsg struct{ path string }

// ---- screen identifiers ----

const (
	screenSearch = iota
	screenForm
	screenConfig
	screenEdit
	screenSettings
)

// AppModel is the root bubbletea model. It routes messages to the active screen.
type AppModel struct {
	mgr             *config.Manager
	currentSettings config.AppSettings
	activeScreen    int
	search          SearchModel
	form            FormModel
	cfgScreen       ConfigScreenModel
	editScreen      EditScreenModel
	settingsScreen  SettingsScreenModel
	finalCmd        string
	width           int
	height          int
}

// NewApp creates the root application model.
// settings should be the already-loaded (and already applied via ApplyTheme)
// user settings so the settings screen reflects the correct initial values.
func NewApp(mgr *config.Manager, settings config.AppSettings) AppModel {
	return AppModel{
		mgr:             mgr,
		currentSettings: settings,
		activeScreen:    screenSearch,
		search:          NewSearchModel(mgr, 0, 0),
	}
}

// GetFinalCommand returns the command built by the user (if any).
func (a AppModel) GetFinalCommand() string { return a.finalCmd }

// Init satisfies tea.Model.
func (a AppModel) Init() tea.Cmd {
	return a.search.Init()
}

// Update satisfies tea.Model.
func (a AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		var cmds []tea.Cmd
		{
			s, c := a.search.Update(msg)
			a.search = s.(SearchModel)
			cmds = append(cmds, c)
		}
		if a.activeScreen == screenForm {
			f, c := a.form.Update(msg)
			a.form = f.(FormModel)
			cmds = append(cmds, c)
		}
		if a.activeScreen == screenConfig {
			cs, c := a.cfgScreen.Update(msg)
			a.cfgScreen = cs.(ConfigScreenModel)
			cmds = append(cmds, c)
		}
		if a.activeScreen == screenEdit {
			es, c := a.editScreen.Update(msg)
			a.editScreen = es.(EditScreenModel)
			cmds = append(cmds, c)
		}
		if a.activeScreen == screenSettings {
			ss, c := a.settingsScreen.Update(msg)
			a.settingsScreen = ss.(SettingsScreenModel)
			cmds = append(cmds, c)
		}
		return a, tea.Batch(cmds...)

	case selectOptionMsg:
		a.form = NewFormModel(msg.cfg, msg.cmd, msg.opt, a.width, a.height)
		a.activeScreen = screenForm
		return a, a.form.Init()

	case goToConfigMsg:
		a.cfgScreen = NewConfigScreenModel(a.mgr, a.width, a.height)
		a.activeScreen = screenConfig
		return a, a.cfgScreen.Init()

	case goToEditMsg:
		a.editScreen = NewEditScreenModel(a.mgr, msg.cfg, a.width, a.height)
		a.activeScreen = screenEdit
		return a, a.editScreen.Init()

	case goToSettingsMsg:
		a.settingsScreen = NewSettingsScreenModel(a.currentSettings, a.width, a.height)
		a.activeScreen = screenSettings
		return a, a.settingsScreen.Init()

	case themeChangedMsg:
		// Sync the updated settings from the settings screen back into AppModel
		// so future navigations to /settings preserve the current state.
		if a.activeScreen == screenSettings {
			a.currentSettings = a.settingsScreen.settings
		}
		return a, nil

	case backToSearchMsg:
		a.activeScreen = screenSearch
		// Refresh config list in case configs changed.
		// Pass current terminal dimensions so the input is immediately full-width.
		a.search = NewSearchModel(a.mgr, a.width, a.height)
		return a, a.search.Init()

	case commandConfirmedMsg:
		a.finalCmd = msg.command
		return a, tea.Quit
	}

	// Delegate to active screen.
	var cmd tea.Cmd
	switch a.activeScreen {
	case screenSearch:
		var m tea.Model
		m, cmd = a.search.Update(msg)
		a.search = m.(SearchModel)
	case screenForm:
		var m tea.Model
		m, cmd = a.form.Update(msg)
		a.form = m.(FormModel)
	case screenConfig:
		var m tea.Model
		m, cmd = a.cfgScreen.Update(msg)
		a.cfgScreen = m.(ConfigScreenModel)
	case screenEdit:
		var m tea.Model
		m, cmd = a.editScreen.Update(msg)
		a.editScreen = m.(EditScreenModel)
	case screenSettings:
		var m tea.Model
		m, cmd = a.settingsScreen.Update(msg)
		a.settingsScreen = m.(SettingsScreenModel)
	}
	return a, cmd
}

// View satisfies tea.Model.
func (a AppModel) View() string {
	switch a.activeScreen {
	case screenForm:
		return a.form.View()
	case screenConfig:
		return a.cfgScreen.View()
	case screenEdit:
		return a.editScreen.View()
	case screenSettings:
		return a.settingsScreen.View()
	default:
		return a.search.View()
	}
}
