package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors.
	colorPrimary    = lipgloss.Color("39")  // bright blue
	colorAccent     = lipgloss.Color("213") // pink/magenta
	colorSuccess    = lipgloss.Color("76")  // green
	colorWarning    = lipgloss.Color("220") // yellow
	colorError      = lipgloss.Color("196") // red
	colorMuted      = lipgloss.Color("241") // grey
	colorText       = lipgloss.Color("252") // near-white
	colorBackground = lipgloss.Color("235") // dark bg
	colorSelected   = lipgloss.Color("24")  // dark blue bg for selection

	// Title bar.
	StyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Background(lipgloss.Color("236")).
			Padding(0, 2).
			Width(0) // width set at runtime

	StyleTitleVersion = lipgloss.NewStyle().
				Foreground(colorMuted).
				Background(lipgloss.Color("236"))

	// Search input border.
	StyleSearchBorder = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(0, 1)

	StyleSearchBorderFocused = lipgloss.NewStyle().
					Border(lipgloss.RoundedBorder()).
					BorderForeground(colorAccent).
					Padding(0, 1)

	// Result list items.
	StyleResultNormal = lipgloss.NewStyle().
				Foreground(colorText).
				Padding(0, 2)

	StyleResultSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(colorSelected).
				Bold(true).
				Padding(0, 2)

	StyleResultCommand = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true)

	StyleResultOption = lipgloss.NewStyle().
				Foreground(colorAccent)

	StyleResultDesc = lipgloss.NewStyle().
			Foreground(colorMuted)

	StyleResultConfig = lipgloss.NewStyle().
				Foreground(lipgloss.Color("243")).
				Italic(true)

	// Form styles.
	StyleFormHeader = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			Border(lipgloss.NormalBorder(), false, false, true, false).
			BorderForeground(colorPrimary).
			Padding(0, 1).
			MarginBottom(1)

	StyleInputLabel = lipgloss.NewStyle().
			Foreground(colorText).
			Bold(true)

	StyleInputLabelRequired = lipgloss.NewStyle().
				Foreground(colorAccent).
				Bold(true)

	StyleInputDesc = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)

	StyleInputFocused = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)

	StyleInputBlurred = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorMuted).
				Padding(0, 1)

	// Command preview.
	StylePreviewBox = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorSuccess).
			Padding(0, 1).
			Foreground(colorSuccess)

	StylePreviewLabel = lipgloss.NewStyle().
				Foreground(colorSuccess).
				Bold(true)

	// Completion overlay.
	StyleCompletionBox = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorWarning).
				Padding(0, 1)

	StyleCompletionItem = lipgloss.NewStyle().
				Foreground(colorText).
				Padding(0, 1)

	StyleCompletionSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(lipgloss.Color("94")).
				Padding(0, 1)

	// Status bar.
	StyleStatus = lipgloss.NewStyle().
			Foreground(colorMuted).
			Background(lipgloss.Color("236")).
			Padding(0, 1)

	StyleStatusKey = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Background(lipgloss.Color("236")).
			Bold(true)

	// Config screen.
	StyleConfigItem = lipgloss.NewStyle().
			Foreground(colorText).
			Padding(0, 2)

	StyleConfigItemSelected = lipgloss.NewStyle().
				Foreground(lipgloss.Color("255")).
				Background(colorSelected).
				Bold(true).
				Padding(0, 2)

	StyleConfigHeader = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary).
				Padding(0, 1).
				MarginBottom(1)

	// Section separator.
	StyleSeparator = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Error / info banners.
	StyleError = lipgloss.NewStyle().
			Foreground(colorError).
			Bold(true).
			Padding(0, 1)

	StyleInfo = lipgloss.NewStyle().
			Foreground(colorSuccess).
			Padding(0, 1)
)
