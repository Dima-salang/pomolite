package ui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	// Catppuccin Mocha palette semantic mapping
	colorBase      = "#1e1e2e" // Background
	colorSurface   = "#313244" // Inactive panels/borders
	colorOverlay   = "#45475a" // Subtle grid lines
	colorSubtext   = "#a6adc8" // Non-focused text/description
	colorText      = "#cdd6f4" // Main text
	colorMauve     = "#cba6f7" // Primary accent (Focus)
	colorLavender  = "#b4befe" // Secondary accent
	colorGreen     = "#a6e3a1" // Success / Done / Work
	colorPeach     = "#fab387" // Warning / In Progress / Break
	colorRed       = "#f38ba8" // Danger / High Priority
	colorBlue      = "#89b4fa" // Info / Low Priority
	colorYellow    = "#f9e2af" // Medium Priority

	// Compatibility aliases for older components
	primaryColor = colorMauve
	accentColor  = colorLavender
	successColor = colorGreen
	white        = colorText
	gray         = colorSubtext

	ASCIIArt = `
██████╗  ██████╗ ███╗   ███╗ ██████╗ ██╗     ██╗████████╗███████╗
██╔══██╗██╔═══██╗████╗ ████║██╔═══██╗██║     ██║╚══██╔══╝██╔════╝
██████╔╝██║   ██║██╔████╔██║██║   ██║██║     ██║   ██║   █████╗  
██╔═══╝ ██║   ██║██║╚██╔╝██║██║   ██║██║     ██║   ██║   ██╔══╝  
██║     ╚██████╔╝██║ ╚═╝ ██║╚██████╔╝███████╗██║   ██║   ███████╗
╚═╝      ╚═════╝ ╚═╝     ╚═╝ ╚═════╝ ╚══════╝╚═╝   ╚═╝   ╚══════╝
`
)

var (
	centerStyle = lipgloss.NewStyle().
			Align(lipgloss.Center, lipgloss.Center)

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colorBase)).
			Background(lipgloss.Color(colorMauve)).
			Padding(0, 1).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorGreen)).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorLavender)).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorSubtext)).
			MarginTop(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorLavender)).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMauve)).
			Bold(true).
			MarginLeft(2).
			MarginBottom(1)

	menuStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorMauve))

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorText)).
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMauve)).
			Bold(true).
			PaddingLeft(0)

	timerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorMauve)).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorMauve)).
			Padding(1, 4)

	lowPriorityStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(colorBlue))
	medPriorityStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color(colorYellow))
	highPriorityStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed))

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorSurface)).
			Padding(1, 2)

	panelFocusedStyle = panelStyle.Copy().
				BorderForeground(lipgloss.Color(colorMauve))

	sidebarStyle = lipgloss.NewStyle()

	workspaceStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	columnStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorSurface)).
			Padding(0, 1).
			Margin(0, 1).
			Width(30)

	columnFocusedStyle = columnStyle.Copy().
				BorderForeground(lipgloss.Color(colorMauve))

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorSurface)).
			Padding(0, 1).
			MarginBottom(1)

	cardSelectedStyle = cardStyle.Copy().
				Bold(true).
				BorderForeground(lipgloss.Color(colorMauve)).
				Foreground(lipgloss.Color(colorText))
)

var bigDigits = map[rune][]string{
	'0': {" ███ ", "█   █", "█   █", "█   █", " ███ "},
	'1': {"  █  ", " ██  ", "  █  ", "  █  ", " ███ "},
	'2': {" ███ ", "    █", "  ██ ", " █   ", " ████"},
	'3': {" ███ ", "    █", "  ██ ", "    █", " ███ "},
	'4': {"█   █", "█   █", " ████", "    █", "    █"},
	'5': {"█████", "█    ", " ███ ", "    █", " ███ "},
	'6': {" ███ ", "█    ", "████ ", "█   █", " ███ "},
	'7': {"█████", "    █", "   █ ", "  █  ", " █   "},
	'8': {" ███ ", "█   █", " ███ ", "█   █", " ███ "},
	'9': {" ███ ", "█   █", " ████", "    █", " ███ "},
	':': {"     ", "  █  ", "     ", "  █  ", "     "},
}

func renderBigText(text string) string {
	lines := make([]string, 5)
	for _, r := range text {
		digit, ok := bigDigits[r]
		if !ok {
			continue
		}
		for i := 0; i < 5; i++ {
			lines[i] += digit[i] + "  "
		}
	}
	return strings.Join(lines, "\n")
}
