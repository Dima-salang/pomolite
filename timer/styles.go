package timer

import "github.com/charmbracelet/lipgloss"

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
)
