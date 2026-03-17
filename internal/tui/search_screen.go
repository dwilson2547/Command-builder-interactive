package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dwilson2547/command-builder/internal/config"
	"github.com/dwilson2547/command-builder/internal/search"
)

// SearchModel is the main search screen.
type SearchModel struct {
	mgr           *config.Manager
	input         textinput.Model
	results       []search.SearchResult
	selectedIdx   int
	scrollTop     int
	width         int
	height        int
	message       string // transient status message
	completions   []string
	completionIdx int
	completionBase string

	// star mode — active when the query is "/s" or "/s <term>"
	starMode     bool
	starList     []config.Star
	starSelected int
}

// NewSearchModel creates a new search screen bound to the given manager.
// Pass the current terminal dimensions so the input is sized correctly
// immediately; use 0, 0 when the size is not yet known.
func NewSearchModel(mgr *config.Manager, w, h int) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search commands… (e.g. 'openssl print p12')"
	if w > 0 {
		ti.Width = max(20, w-10)
	} else {
		ti.Width = 60
	}
	ti.Focus()

	m := SearchModel{
		mgr:    mgr,
		input:  ti,
		width:  w,
		height: h,
	}
	m.results = runSearch("", mgr)
	return m
}

func runSearch(query string, mgr *config.Manager) []search.SearchResult {
	filter, _ := search.ParseQuery(query)
	return search.Search(query, mgr.ListConfigs(), filter)
}

// Init satisfies tea.Model – start the cursor blink animation.
func (m SearchModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update satisfies tea.Model.
func (m SearchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = max(20, msg.Width-10)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyEnter:
			// Handle config/import/settings/stars commands.
			query := strings.TrimSpace(m.input.Value())
			if strings.HasPrefix(query, "/config") {
				return m, func() tea.Msg { return goToConfigMsg{} }
			}
			if strings.HasPrefix(query, "/settings") {
				return m, func() tea.Msg { return goToSettingsMsg{} }
			}
			// Star mode: open the selected starred command.
			if m.starMode {
				if len(m.starList) == 0 {
					return m, nil
				}
				star := m.starList[m.starSelected]
				return m, func() tea.Msg { return selectStarMsg{star: star} }
			}
			if strings.HasPrefix(query, "/import ") {
				target := strings.TrimSpace(strings.TrimPrefix(query, "/import "))
				if target != "" {
					if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
						return m, func() tea.Msg { return importURLMsg{url: target} }
					}
					return m, func() tea.Msg { return importFileMsg{path: target} }
				}
			}
			if len(m.results) == 0 {
				return m, nil
			}
			r := m.results[m.selectedIdx]
			return m, func() tea.Msg {
				return selectOptionMsg{cfg: r.Config, cmd: r.Command, opt: r.Option}
			}

		case tea.KeyUp:
			if m.starMode {
				if m.starSelected > 0 {
					m.starSelected--
				}
				return m, nil
			}
			if m.selectedIdx > 0 {
				m.selectedIdx--
				m.ensureVisible()
			}
			return m, nil

		case tea.KeyDown:
			if m.starMode {
				if m.starSelected < len(m.starList)-1 {
					m.starSelected++
				}
				return m, nil
			}
			if m.selectedIdx < len(m.results)-1 {
				m.selectedIdx++
				m.ensureVisible()
			}
			return m, nil

		case tea.KeyPgUp:
			m.selectedIdx = max(0, m.selectedIdx-m.visibleRows())
			m.ensureVisible()
			return m, nil

		case tea.KeyPgDown:
			m.selectedIdx = min(len(m.results)-1, m.selectedIdx+m.visibleRows())
			if m.selectedIdx < 0 {
				m.selectedIdx = 0
			}
			m.ensureVisible()
			return m, nil

		case tea.KeyTab:
			// In star mode Tab is a no-op (no path completion).
			if m.starMode {
				return m, nil
			}
			query := m.input.Value()
			if !strings.HasPrefix(query, "/import ") {
				return m, nil
			}
			partial := strings.TrimPrefix(query, "/import ")
			// Recompute completions when the path portion changed since last Tab.
			if partial != m.completionBase || len(m.completions) == 0 {
				m.completions = pathCompletions(partial)
				m.completionBase = partial
				m.completionIdx = -1
			}
			if len(m.completions) == 0 {
				return m, nil
			}
			if len(m.completions) == 1 {
				m.input.SetValue("/import " + m.completions[0])
				m.input.CursorEnd()
				m.completionBase = m.completions[0]
				m.completions = nil
				return m, nil
			}
			// Multiple: first Tab fills longest common prefix; subsequent Tabs cycle.
			prefix := longestCommonPrefix(m.completions)
			if m.completionIdx == -1 && prefix != partial {
				m.input.SetValue("/import " + prefix)
				m.input.CursorEnd()
				m.completionBase = prefix
			} else {
				if m.completionIdx == -1 {
					m.completionIdx = 0
				} else {
					m.completionIdx = (m.completionIdx + 1) % len(m.completions)
				}
				m.input.SetValue("/import " + m.completions[m.completionIdx])
				m.input.CursorEnd()
			}
			return m, nil
		}

		// Handle star-mode key actions that need rune inspection.
		if m.starMode {
			if msg.Type == tea.KeyRunes {
				if string(msg.Runes) == "d" || string(msg.Runes) == "D" {
					return m.deleteSelectedStar()
				}
			}
			// Fall through so typing/backspace reach the text input, allowing
			// the user to filter stars ("/s <term>") or exit star mode entirely.
		}

	case importURLMsg:
		cfg, err := m.mgr.ImportConfigFromURL(msg.url)
		if err != nil {
			m.message = StyleError.Render("Import failed: " + err.Error())
		} else {
			m.message = StyleInfo.Render(fmt.Sprintf("Imported config %q", cfg.Name))
		}
		m.input.SetValue("")
		m.results = runSearch("", m.mgr)
		m.selectedIdx = 0
		m.scrollTop = 0
		return m, nil

	case importFileMsg:
		cfg, err := m.mgr.ImportConfigFromFile(msg.path)
		if err != nil {
			m.message = StyleError.Render("Import failed: " + err.Error())
		} else {
			m.message = StyleInfo.Render(fmt.Sprintf("Imported config %q", cfg.Name))
		}
		m.input.SetValue("")
		m.results = runSearch("", m.mgr)
		m.selectedIdx = 0
		m.scrollTop = 0
		return m, nil
	}

	// Delegate to text input.
	prevQuery := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	query := m.input.Value()
	if query != prevQuery {
		isStarQuery := (query == "/s" || strings.HasPrefix(query, "/s ")) &&
			!strings.HasPrefix(query, "/settings")

		if isStarQuery {
			// Enter star mode: load stars and display them inline.
			m.starMode = true
			allStars := config.LoadStars()
			term := strings.TrimSpace(strings.TrimPrefix(query, "/s"))
			if term == "" {
				m.starList = allStars
			} else {
				term = strings.ToLower(term)
				filtered := allStars[:0]
				for _, s := range allStars {
					if strings.Contains(strings.ToLower(s.DisplayName()), term) ||
						strings.Contains(strings.ToLower(s.CommandName), term) ||
						strings.Contains(strings.ToLower(s.OptionName), term) {
						filtered = append(filtered, s)
					}
				}
				m.starList = filtered
			}
			m.starSelected = 0
			m.completions = nil
			m.completionIdx = -1
			m.completionBase = ""
			m.message = ""
		} else {
			m.starMode = false
			m.starList = nil
			// Reset completions when the path portion changes (manual edit).
			if strings.HasPrefix(query, "/import ") {
				partial := strings.TrimPrefix(query, "/import ")
				if partial != m.completionBase {
					m.completions = nil
					m.completionIdx = -1
					m.completionBase = ""
				}
			} else {
				m.completions = nil
				m.completionIdx = -1
				m.completionBase = ""
			}
			// Rerun search only for non-special queries.
			if !strings.HasPrefix(query, "/config") && !strings.HasPrefix(query, "/import") && !strings.HasPrefix(query, "/settings") {
				m.results = runSearch(query, m.mgr)
				m.selectedIdx = 0
				m.scrollTop = 0
				m.message = ""
			}
		}
	}

	return m, cmd
}

// View satisfies tea.Model.
func (m SearchModel) View() string {
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
	title := StyleTitle.Width(w).Render(
		"⚡ " + AppDisplayName,
	)
	b.WriteString(title + "\n")

	// ── Search input ───────────────────────────────────────────────────────
	inputBox := StyleSearchBorderFocused.Width(w - 4).Render(m.input.View())
	b.WriteString(inputBox + "\n")

	// ── Hint line ──────────────────────────────────────────────────────────
	inImport := strings.HasPrefix(m.input.Value(), "/import ")
	var hintText string
	switch {
	case m.starMode:
		hintText = fmt.Sprintf(" %d starred · ↑↓ navigate · Enter open · d delete · Esc or clear to exit", len(m.starList))
	case inImport:
		hintText = " Enter to import · Tab to autocomplete path · Esc clears"
	default:
		hintText = " ↑↓ navigate · Enter select · /default · /all · /<config> · /import <url or path> · /s stars · /settings · Ctrl+C quit"
	}
	hint := StyleResultDesc.Render(hintText)
	b.WriteString(hint + "\n")

	if m.message != "" {
		b.WriteString(m.message + "\n")
	}

	sep := StyleSeparator.Render(strings.Repeat("─", w))
	b.WriteString(sep + "\n")

	// ── Results list ───────────────────────────────────────────────────────
	// Reserve rows: title(1) + input-border(3) + hint(1) + sep(1) + status(1) = 7
	reserved := 7
	if m.message != "" {
		reserved++
	}
	visRows := h - reserved
	if visRows < 1 {
		visRows = 1
	}

	start := m.scrollTop
	end := min(start+visRows, len(m.results))

	// ── Star mode: show starred commands instead of search results ──────────
	if m.starMode {
		if len(m.starList) == 0 {
			b.WriteString(StyleResultDesc.Padding(0, 2).Render(
				"No starred commands yet — fill a form and press * to star it.",
			) + "\n")
			for i := 1; i < visRows; i++ {
				b.WriteString("\n")
			}
		} else {
			shown := min(len(m.starList), visRows)
			for i := 0; i < shown; i++ {
				b.WriteString(m.renderStarResult(m.starList[i], i == m.starSelected, w) + "\n")
			}
			for i := shown; i < visRows; i++ {
				b.WriteString("\n")
			}
		}
		statusLeft := StyleStatus.Render(fmt.Sprintf(" %d starred", len(m.starList)))
		statusRight := StyleStatus.Render(" Ctrl+C quit") + footerVersion()
		b.WriteString(renderFooter(w, statusLeft, statusRight))
		return b.String()
	}

	if len(m.results) == 0 {
		// When in /import mode, show completions (or a hint) instead of "No results".
		if strings.HasPrefix(m.input.Value(), "/import ") {
			if len(m.completions) > 0 {
				const maxShow = 8
				shown := m.completions
				if len(shown) > maxShow {
					shown = shown[:maxShow]
				}
				for i, c := range shown {
					if i == m.completionIdx {
						b.WriteString(StyleResultSelected.Width(w).Render("  "+c) + "\n")
					} else {
						b.WriteString(StyleResultNormal.Render("  "+c) + "\n")
					}
				}
				if len(m.completions) > maxShow {
					b.WriteString(StyleResultDesc.Render(fmt.Sprintf("  … and %d more", len(m.completions)-maxShow)) + "\n")
					for i := len(shown) + 1; i < visRows; i++ {
						b.WriteString("\n")
					}
				} else {
					for i := len(shown); i < visRows; i++ {
						b.WriteString("\n")
					}
				}
			} else {
				b.WriteString(StyleResultDesc.Padding(0, 2).Render("Type a path then press Tab to autocomplete.") + "\n")
				for i := 1; i < visRows; i++ {
					b.WriteString("\n")
				}
			}
		} else {
			// Write the "no results" notice as the first line of the results area,
			// then pad the remainder so the total results area stays exactly visRows.
			b.WriteString(StyleResultDesc.Padding(0, 2).Render("No results. Try a different query.") + "\n")
			for i := 1; i < visRows; i++ {
				b.WriteString("\n")
			}
		}
	} else {
		for i := start; i < end; i++ {
			r := m.results[i]
			line := m.renderResult(r, i == m.selectedIdx, w)
			b.WriteString(line + "\n")
		}
		// Pad remaining space.
		rendered := end - start
		for i := rendered; i < visRows; i++ {
			b.WriteString("\n")
		}
	}

	// ── Status bar ─────────────────────────────────────────────────────────
	statusLeft := StyleStatus.Render(fmt.Sprintf(" %d result(s)", len(m.results)))
	statusRight := StyleStatus.Render(" Ctrl+C quit") + footerVersion()
	b.WriteString(renderFooter(w, statusLeft, statusRight))

	return b.String()
}

func (m SearchModel) renderStarResult(star config.Star, selected bool, width int) string {
	// Build a compact summary of saved non-empty values.
	var parts []string
	for k, v := range star.Values {
		if v != "" {
			parts = append(parts, k+"="+v)
		}
	}
	for k, v := range star.FlagStates {
		if v {
			parts = append(parts, k)
		}
	}
	summary := strings.Join(parts, "  ")
	badge := "[" + star.ConfigName + "]"
	displayName := star.DisplayName()

	if selected {
		line := "★ " + displayName
		if summary != "" {
			line += "  " + summary
		}
		line += "  " + badge
		return StyleResultSelected.Width(width).Render(line)
	}
	label := StyleResultCommand.Render("★ " + displayName)
	line := " " + label
	if summary != "" {
		line += "  " + StyleResultDesc.Render(summary)
	}
	line += "  " + StyleResultConfig.Render(badge)
	return line
}

func (m SearchModel) renderResult(r search.SearchResult, selected bool, width int) string {	// Build the label: "openssl › print-p12"
	cmdPart := StyleResultCommand.Render(r.Command.Name)
	optPart := StyleResultOption.Render(r.Option.Name)
	label := cmdPart + StyleResultDesc.Render(" › ") + optPart

	descPart := "  " + StyleResultDesc.Render(r.Option.Description)
	cfgBadge := " " + StyleResultConfig.Render("["+r.Config.Name+"]")

	line := " " + label + descPart + cfgBadge

	// Truncate if needed.
	if lipgloss.Width(line) > width-2 {
		line = line[:max(0, width-5)] + "…"
	}

	if selected {
		return StyleResultSelected.Width(width).Render(
			r.Command.Name + " › " + r.Option.Name + "  " + r.Option.Description + "  [" + r.Config.Name + "]",
		)
	}
	return line
}

// deleteSelectedStar removes the highlighted star from disk and refreshes the list.
func (m SearchModel) deleteSelectedStar() (tea.Model, tea.Cmd) {
	if len(m.starList) == 0 {
		return m, nil
	}
	star := m.starList[m.starSelected]
	if err := config.DeleteStar(star.ID); err != nil {
		m.message = StyleError.Render("Delete failed: " + err.Error())
		return m, nil
	}
	m.starList = config.LoadStars()
	if m.starSelected >= len(m.starList) && m.starSelected > 0 {
		m.starSelected--
	}
	m.message = StyleInfo.Render(fmt.Sprintf("Removed \"%s\"", star.DisplayName()))
	return m, nil
}

func (m *SearchModel) visibleRows() int {
	if m.height <= 0 {
		return 10
	}
	return max(1, m.height-7)
}

func (m *SearchModel) ensureVisible() {
	vis := m.visibleRows()
	if m.selectedIdx < m.scrollTop {
		m.scrollTop = m.selectedIdx
	}
	if m.selectedIdx >= m.scrollTop+vis {
		m.scrollTop = m.selectedIdx - vis + 1
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
