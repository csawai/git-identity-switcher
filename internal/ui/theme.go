package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Theme defines the color scheme and styling for gitx
type Theme struct {
	Primary    lipgloss.Color
	Success    lipgloss.Color
	Warning    lipgloss.Color
	Error      lipgloss.Color
	Info       lipgloss.Color
	Border     lipgloss.Color
	Text       lipgloss.Color
	Muted      lipgloss.Color
	Background lipgloss.Color
}

var DefaultTheme = Theme{
	Primary:    lipgloss.Color("39"),    // Bright blue
	Success:    lipgloss.Color("46"),    // Bright green
	Warning:    lipgloss.Color("226"),   // Yellow
	Error:      lipgloss.Color("196"),   // Red
	Info:       lipgloss.Color("51"),    // Cyan
	Border:     lipgloss.Color("240"),   // Gray
	Text:       lipgloss.Color("255"),   // White
	Muted:      lipgloss.Color("244"),   // Light gray
	Background: lipgloss.Color("235"),   // Dark gray
}

// Styles for common UI elements
var (
	// Box styles
	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(DefaultTheme.Border).
		Padding(1, 2)

	SuccessBox = BoxStyle.Copy().
		BorderForeground(DefaultTheme.Success).
		Foreground(DefaultTheme.Success)

	ErrorBox = BoxStyle.Copy().
		BorderForeground(DefaultTheme.Error).
		Foreground(DefaultTheme.Error)

	WarningBox = BoxStyle.Copy().
		BorderForeground(DefaultTheme.Warning).
		Foreground(DefaultTheme.Warning)

	InfoBox = BoxStyle.Copy().
		BorderForeground(DefaultTheme.Info).
		Foreground(DefaultTheme.Info)

	// Text styles
	TitleStyle = lipgloss.NewStyle().
		Foreground(DefaultTheme.Primary).
		Bold(true).
		Margin(1, 0)

	HeaderStyle = lipgloss.NewStyle().
		Foreground(DefaultTheme.Primary).
		Bold(true).
		Margin(0, 0, 1, 0)

	SuccessText = lipgloss.NewStyle().
		Foreground(DefaultTheme.Success)

	ErrorText = lipgloss.NewStyle().
		Foreground(DefaultTheme.Error)

	WarningText = lipgloss.NewStyle().
		Foreground(DefaultTheme.Warning)

	InfoText = lipgloss.NewStyle().
		Foreground(DefaultTheme.Info)

	MutedText = lipgloss.NewStyle().
		Foreground(DefaultTheme.Muted)

	// Table styles
	TableHeaderStyle = lipgloss.NewStyle().
		Foreground(DefaultTheme.Primary).
		Bold(true).
		Padding(0, 1)

	TableRowStyle = lipgloss.NewStyle().
		Padding(0, 1)

	// Status indicators
	StatusBound = lipgloss.NewStyle().
		Foreground(DefaultTheme.Success).
		Render("🟢")

	StatusUnbound = lipgloss.NewStyle().
		Foreground(DefaultTheme.Warning).
		Render("🟡")

	StatusError = lipgloss.NewStyle().
		Foreground(DefaultTheme.Error).
		Render("🔴")
)

// GetStatusIcon returns the appropriate status icon
func GetStatusIcon(bound bool) string {
	if bound {
		return StatusBound
	}
	return StatusUnbound
}

// GetAuthIcon returns icon for auth method
func GetAuthIcon(authMethod string) string {
	switch authMethod {
	case "ssh":
		return "🔑"
	case "pat":
		return "🔐"
	default:
		return "❓"
	}
}

