package timer

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	primaryColor = "#7D56F4"
	accentColor  = "#EE6FF8"
	successColor = "#04B575"
	white        = "#FAFAFA"
	gray         = "#626262"

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
			Foreground(lipgloss.Color(white)).
			Background(lipgloss.Color(primaryColor)).
			Padding(0, 1).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(successColor)).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(accentColor)).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(gray)).
			MarginTop(1)

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(accentColor)).
			Bold(true).
			MarginTop(1).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(primaryColor)).
			Bold(true).
			MarginLeft(2).
			MarginBottom(1)

	menuStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginTop(1).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryColor))

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(white)).
			PaddingLeft(2)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(primaryColor)).
				Bold(true).
				PaddingLeft(0)

	timerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(primaryColor)).
			Bold(true).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(primaryColor)).
			Padding(1, 4)
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
