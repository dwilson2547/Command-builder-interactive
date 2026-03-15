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

const appVersion = "v1.0.0"

// SearchModel is the main search screen.
type SearchModel struct {
	mgr         *config.Manager
	input       textinput.Model
	results     []search.SearchResult
	selectedIdx int
	scrollTop   int
	width       int
	height      int
	message     string // transient status message
}

// NewSearchModel creates a new search screen bound to the given manager.
func NewSearchModel(mgr *config.Manager) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search commands… (e.g. 'openssl print p12')"
	ti.Width = 60
	ti.Focus()

	m := SearchModel{
		mgr:   mgr,
		input: ti,
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
			// Handle config/import commands.
			query := strings.TrimSpace(m.input.Value())
			if strings.HasPrefix(query, "/config") {
				return m, func() tea.Msg { return goToConfigMsg{} }
			}
			if strings.HasPrefix(query, "/import ") {
				url := strings.TrimSpace(strings.TrimPrefix(query, "/import "))
				if url != "" {
					return m, func() tea.Msg { return importURLMsg{url: url} }
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
			if m.selectedIdx > 0 {
				m.selectedIdx--
				m.ensureVisible()
			}
			return m, nil

		case tea.KeyDown:
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
	}

	// Delegate to text input.
	prevQuery := m.input.Value()
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)

	// Rerun search only when the query actually changes (not on blink ticks etc).
	// Exceptions: /config and /import are handled on Enter only.
	query := m.input.Value()
	if query != prevQuery && !strings.HasPrefix(query, "/config") && !strings.HasPrefix(query, "/import") {
		m.results = runSearch(query, m.mgr)
		m.selectedIdx = 0
		m.scrollTop = 0
		m.message = ""
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
	title := StyleTitle.Copy().Width(w).Render(
		"⚡ Command Builder " + StyleTitleVersion.Render(appVersion),
	)
	b.WriteString(title + "\n")

	// ── Search input ───────────────────────────────────────────────────────
	inputBox := StyleSearchBorderFocused.Copy().Width(w - 4).Render(m.input.View())
	b.WriteString(inputBox + "\n")

	// ── Hint line ──────────────────────────────────────────────────────────
	hint := StyleResultDesc.Render(
		" ↑↓ navigate · Enter select · /default · /all · /<config> · /import <url> · Ctrl+C quit",
	)
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

	if len(m.results) == 0 {
		// Write the "no results" notice as the first line of the results area,
		// then pad the remainder so the total results area stays exactly visRows.
		b.WriteString(StyleResultDesc.Padding(0, 2).Render("No results. Try a different query.") + "\n")
		for i := 1; i < visRows; i++ {
			b.WriteString("\n")
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
	statusLeft := fmt.Sprintf(" %d result(s)", len(m.results))
	statusRight := " Ctrl+C quit "
	gap := w - len(statusLeft) - len(statusRight)
	if gap < 0 {
		gap = 0
	}
	status := StyleStatus.Copy().Width(w).Render(
		statusLeft + strings.Repeat(" ", gap) + statusRight,
	)
	b.WriteString(status)

	return b.String()
}

func (m SearchModel) renderResult(r search.SearchResult, selected bool, width int) string {
	// Build the label: "openssl › print-p12"
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
		return StyleResultSelected.Copy().Width(width).Render(
			r.Command.Name + " › " + r.Option.Name + "  " + r.Option.Description + "  [" + r.Config.Name + "]",
		)
	}
	return line
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
