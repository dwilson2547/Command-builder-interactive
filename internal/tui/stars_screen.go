package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dwilson2547/command-builder/internal/config"
)

// ---- messages ---------------------------------------------------------------

// goToStarsMsg navigates to the starred commands screen.
type goToStarsMsg struct{}

// selectStarMsg opens the form screen pre-filled with a saved star's values.
type selectStarMsg struct{ star config.Star }

// ---- model ------------------------------------------------------------------

// StarsScreenModel is the /s starred-commands screen.
type StarsScreenModel struct {
	mgr      *config.Manager
	stars    []config.Star
	selected int
	message  string
	width    int
	height   int
}

// NewStarsScreenModel creates a stars screen bound to the given manager.
func NewStarsScreenModel(mgr *config.Manager, w, h int) StarsScreenModel {
	return StarsScreenModel{
		mgr:    mgr,
		stars:  config.LoadStars(),
		width:  w,
		height: h,
	}
}

// Init satisfies tea.Model.
func (m StarsScreenModel) Init() tea.Cmd { return nil }

// Update satisfies tea.Model.
func (m StarsScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
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
			if m.selected < len(m.stars)-1 {
				m.selected++
			}

		case tea.KeyEnter:
			if len(m.stars) == 0 {
				return m, nil
			}
			star := m.stars[m.selected]
			return m, func() tea.Msg { return selectStarMsg{star: star} }

		case tea.KeyRunes:
			if string(msg.Runes) == "d" || string(msg.Runes) == "D" {
				return m.deleteSelected()
			}
		}
	}
	return m, nil
}

// deleteSelected removes the currently highlighted star.
func (m StarsScreenModel) deleteSelected() (tea.Model, tea.Cmd) {
	if len(m.stars) == 0 {
		return m, nil
	}
	star := m.stars[m.selected]
	if err := config.DeleteStar(star.ID); err != nil {
		m.message = StyleError.Render("Delete failed: " + err.Error())
		return m, nil
	}
	m.stars = config.LoadStars()
	if m.selected >= len(m.stars) && m.selected > 0 {
		m.selected--
	}
	m.message = StyleInfo.Render(fmt.Sprintf("Removed \"%s\"", star.DisplayName()))
	return m, nil
}

// View satisfies tea.Model.
func (m StarsScreenModel) View() string {
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
		"⚡ " + AppDisplayName + "  " + StyleResultDesc.Render("Starred Commands"),
	)
	b.WriteString(title + "\n")

	// ── Hint line ──────────────────────────────────────────────────────────
	hint := StyleResultDesc.Render(" ↑↓ navigate · Enter open · d delete · Esc back")
	b.WriteString(hint + "\n")
	b.WriteString(StyleSeparator.Render(strings.Repeat("─", w)) + "\n")

	// ── Star list ──────────────────────────────────────────────────────────
	// Reserve rows: title(1) + hint(1) + sep(1) + status(1) = 4
	reserved := 4
	if m.message != "" {
		reserved++
	}
	visRows := h - reserved
	if visRows < 1 {
		visRows = 1
	}

	if len(m.stars) == 0 {
		b.WriteString(StyleResultDesc.Padding(0, 2).Render(
			"No starred commands yet. Open any command form and press * to star it.",
		) + "\n")
		for i := 1; i < visRows; i++ {
			b.WriteString("\n")
		}
	} else {
		for i, star := range m.stars {
			if i >= visRows {
				break
			}
			b.WriteString(m.renderStar(star, i == m.selected, w) + "\n")
		}
		for i := len(m.stars); i < visRows; i++ {
			b.WriteString("\n")
		}
	}

	// ── Message ────────────────────────────────────────────────────────────
	if m.message != "" {
		b.WriteString("  " + m.message + "\n")
	}

	// ── Status bar ─────────────────────────────────────────────────────────
	keys := StyleStatusKey.Render(" Enter") + StyleStatus.Render(" open") +
		StyleStatusKey.Render("  d") + StyleStatus.Render(" delete") +
		StyleStatusKey.Render("  Esc") + StyleStatus.Render(" back")
	b.WriteString(renderFooter(w, keys, footerVersion()))

	return b.String()
}

// renderStar renders a single star row.
func (m StarsScreenModel) renderStar(star config.Star, selected bool, width int) string {
	// Build the saved-values summary (show non-empty string inputs).
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

	cfgBadge := "[" + star.ConfigName + "]"
	displayName := star.DisplayName()

	if selected {
		line := displayName
		if summary != "" {
			line += "  " + summary
		}
		line += "  " + cfgBadge
		return StyleResultSelected.Width(width).Render(line)
	}

	label := StyleResultCommand.Render(displayName)
	badge := StyleResultConfig.Render(cfgBadge)
	line := " " + label
	if summary != "" {
		line += "  " + StyleResultDesc.Render(summary)
	}
	line += "  " + badge
	return line
}
