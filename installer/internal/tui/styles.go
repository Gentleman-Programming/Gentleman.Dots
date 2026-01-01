package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors - Gentleman Theme (from opencode theme)
	// Base colors
	Background        = lipgloss.Color("#06080f")
	BackgroundPanel   = lipgloss.Color("#06080f")
	BackgroundElement = lipgloss.Color("#06080f")

	// Text colors
	Text      = lipgloss.Color("#F3F6F9")
	TextMuted = lipgloss.Color("#5C6170")

	// Accent colors
	Primary   = lipgloss.Color("#7FB4CA") // Blue-ish
	Secondary = lipgloss.Color("#A3B5D6") // Light blue
	Accent    = lipgloss.Color("#E0C15A") // Gold/Yellow

	// Status colors
	Error   = lipgloss.Color("#CB7C94") // Pink-red
	Warning = lipgloss.Color("#DEBA87") // Orange-tan
	Success = lipgloss.Color("#B7CC85") // Green
	Info    = lipgloss.Color("#7FB4CA") // Blue

	// Border colors
	Border       = lipgloss.Color("#313342")
	BorderActive = lipgloss.Color("#7FB4CA")
	BorderSubtle = lipgloss.Color("#232A40")

	// Syntax colors (for code display)
	SyntaxComment     = lipgloss.Color("#8394A3")
	SyntaxKeyword     = lipgloss.Color("#C99AD6") // Purple
	SyntaxFunction    = lipgloss.Color("#B99BF2") // Light purple
	SyntaxVariable    = lipgloss.Color("#F3F6F9")
	SyntaxString      = lipgloss.Color("#DFBD76") // Gold
	SyntaxNumber      = lipgloss.Color("#A4DAA7") // Light green
	SyntaxType        = lipgloss.Color("#8FB8DD") // Light blue
	SyntaxOperator    = lipgloss.Color("#DEBA87")
	SyntaxPunctuation = lipgloss.Color("#96A2B0")

	// Markdown colors
	MarkdownHeading = lipgloss.Color("#B5B2D0") // Lavender

	// Text styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	WarningStyle = lipgloss.NewStyle().
			Foreground(Warning)

	InfoStyle = lipgloss.NewStyle().
			Foreground(Info)

	MutedStyle = lipgloss.NewStyle().
			Foreground(TextMuted)

	// Selection styles
	SelectedStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true).
			PaddingLeft(2)

	UnselectedStyle = lipgloss.NewStyle().
			Foreground(Text).
			PaddingLeft(4)

	// Box styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderActive).
			Padding(1, 2)

	// Progress bar styles
	ProgressBarFilled = lipgloss.NewStyle().
				Foreground(Success).
				Background(Success)

	ProgressBarEmpty = lipgloss.NewStyle().
				Foreground(Border).
				Background(Border)

	// Logo style
	LogoStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true)

	// Step indicator
	StepActiveStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	StepDoneStyle = lipgloss.NewStyle().
			Foreground(Success)

	StepPendingStyle = lipgloss.NewStyle().
				Foreground(TextMuted)

	// Help style
	HelpStyle = lipgloss.NewStyle().
			Foreground(TextMuted).
			Italic(true).
			MarginTop(1)

	// Keymaps style
	KeyStyle = lipgloss.NewStyle().
			Foreground(SyntaxKeyword).
			Bold(true)

	// Code style
	CodeStyle = lipgloss.NewStyle().
			Foreground(SyntaxString)

	// Additional styles for backup screens
	BackupItemStyle = lipgloss.NewStyle().
			Foreground(Secondary)

	DangerStyle = lipgloss.NewStyle().
			Foreground(Error).
			Bold(true)

	HighlightStyle = lipgloss.NewStyle().
			Foreground(Accent).
			Bold(true)

	// Vim Trainer cursor styles
	StartCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#06080f")).
				Background(Warning).
				Bold(true)

	CurrentCursorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#06080f")).
				Background(Success).
				Bold(true)

	// Visual selection style (like Vim's visual mode)
	SelectionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06080f")).
			Background(lipgloss.Color("#7aa2f7")).
			Bold(false)

	// Dimmed code style for parts already passed
	DimmedCodeStyle = lipgloss.NewStyle().
			Foreground(TextMuted)
)

// CenterHorizontally centers text horizontally within a given width
func CenterHorizontally(text string, width int) string {
	return lipgloss.NewStyle().Width(width).Align(lipgloss.Center).Render(text)
}

// CenterVertically centers text vertically within a given height
func CenterVertically(text string, height int) string {
	lines := lipgloss.Height(text)
	if lines >= height {
		return text
	}

	topPadding := (height - lines) / 2
	return lipgloss.NewStyle().PaddingTop(topPadding).Render(text)
}

// CenterBoth centers text both horizontally and vertically using lipgloss.Place
func CenterBoth(text string, width, height int) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, text)
}
