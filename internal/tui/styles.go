package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/dwilson2547/command-builder/internal/config"
)

// ---- mutable colour variables -----------------------------------------------
// These are rebuilt by ApplyTheme whenever the user saves new colour settings.

var (
	colorPrimary    lipgloss.Color
	colorAccent     lipgloss.Color
	colorSuccess    lipgloss.Color
	colorWarning    lipgloss.Color
	colorError      lipgloss.Color
	colorMuted      lipgloss.Color
	colorText       lipgloss.Color
	colorSelected   lipgloss.Color
)

// ---- mutable style variables ------------------------------------------------
// Set to zero values here; populated by the init() / ApplyTheme() call below.

var (
	// Title bar.
	StyleTitle        lipgloss.Style
	StyleTitleVersion lipgloss.Style

	// Search input border.
	StyleSearchBorder        lipgloss.Style
	StyleSearchBorderFocused lipgloss.Style

	// Result list items.
	StyleResultNormal   lipgloss.Style
	StyleResultSelected lipgloss.Style
	StyleResultCommand  lipgloss.Style
	StyleResultOption   lipgloss.Style
	StyleResultDesc     lipgloss.Style
	StyleResultConfig   lipgloss.Style

	// Form styles.
	StyleFormHeader         lipgloss.Style
	StyleInputLabel         lipgloss.Style
	StyleInputLabelRequired lipgloss.Style
	StyleInputDesc          lipgloss.Style
	StyleInputFocused       lipgloss.Style
	StyleInputBlurred       lipgloss.Style

	// Command preview.
	StylePreviewBox   lipgloss.Style
	StylePreviewLabel lipgloss.Style

	// Completion overlay.
	StyleCompletionBox      lipgloss.Style
	StyleCompletionItem     lipgloss.Style
	StyleCompletionSelected lipgloss.Style

	// Status bar.
	StyleStatus    lipgloss.Style
	StyleStatusKey lipgloss.Style

	// Config / settings screen.
	StyleConfigItem         lipgloss.Style
	StyleConfigItemSelected lipgloss.Style
	StyleConfigHeader       lipgloss.Style

	// Section separator.
	StyleSeparator lipgloss.Style

	// Error / info banners.
	StyleError lipgloss.Style
	StyleInfo  lipgloss.Style

	// Flag toggle inputs.
	StyleFlagOn      lipgloss.Style // flag enabled (unfocused)
	StyleFlagOff     lipgloss.Style // flag disabled (unfocused)
	StyleFlagFocused lipgloss.Style // flag row when focused
)

func init() {
	ApplyTheme(config.DefaultSettings())
}

// ApplyTheme rebuilds every colour variable and style from the supplied
// settings. Call this on startup (after loading saved settings) and again
// whenever the user changes a colour in the /settings screen.
func ApplyTheme(s config.AppSettings) {
	// Update display name.
	if s.AppName != "" {
		AppDisplayName = s.AppName
	} else {
		AppDisplayName = "Command Builder"
	}

	colorPrimary = lipgloss.Color(s.ColorPrimary)
	colorAccent = lipgloss.Color(s.ColorAccent)
	colorSuccess = lipgloss.Color(s.ColorSuccess)
	colorWarning = lipgloss.Color(s.ColorWarning)
	colorError = lipgloss.Color(s.ColorError)
	colorMuted = lipgloss.Color(s.ColorMuted)
	colorText = lipgloss.Color(s.ColorText)
	colorSelected = lipgloss.Color(s.ColorSelected)

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

	// Config / settings screen.
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

	// Flag toggle inputs.
	StyleFlagOn = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorSuccess).
		Foreground(colorSuccess).
		Bold(true).
		Padding(0, 1)

	StyleFlagOff = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorMuted).
		Foreground(colorMuted).
		Padding(0, 1)

	StyleFlagFocused = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Foreground(colorAccent).
		Bold(true).
		Padding(0, 1)
}
