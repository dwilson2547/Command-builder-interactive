package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// renderFooter builds a full-width status bar with left-aligned key hints and
// a right-aligned version badge. left and right may already contain ANSI
// styling; lipgloss.Width is used (ANSI-aware) for the gap calculation.
//
// The outer wrapper uses Padding(0,0) so it does NOT add extra side-padding
// on top of any padding already present inside the styled left/right content.
// Width(w) then fills/clips to exactly w columns.
func renderFooter(w int, left, right string) string {
	lw := lipgloss.Width(left)
	rw := lipgloss.Width(right)
	gap := w - lw - rw
	if gap < 0 {
		gap = 0
	}
	bar := StyleStatus.Padding(0, 0).Width(w)
	return bar.Render(left + strings.Repeat(" ", gap) + right)
}

// footerVersion returns a consistently styled version badge for use on the
// right side of every screen's status bar.
func footerVersion() string {
	return StyleTitleVersion.Render(" " + AppVersion + " ")
}
