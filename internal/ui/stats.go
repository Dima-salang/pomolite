package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/stats"
	"github.com/Dima-salang/pomolite/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type StatsModel struct {
	stats  *stats.PomoStats
	err    error
	width  int
	height int
}

func NewStatsModel(repo storage.Repository) StatsModel {
	pomoStats, err := stats.Compute(repo, "all")
	return StatsModel{
		stats: pomoStats,
		err:   err,
	}
}

func (m StatsModel) Init() tea.Cmd {
	return nil
}

func (m StatsModel) Update(msg tea.Msg) (StatsModel, tea.Cmd) {
	return m, nil
}

func (m StatsModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" Pomolite "))
	s.WriteString("\n\n")
	s.WriteString(headerStyle.Render("Statistics Overview"))
	s.WriteString("\n")

	if m.err != nil {
		s.WriteString(fmt.Sprintf("  Error loading stats: %v\n", m.err))
	} else if m.stats != nil {
		cardStyle := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(colorSurface)).
			Padding(1, 2).
			Margin(0, 1).
			Width(24)

		totalSessionsVal := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorMauve)).Render(fmt.Sprintf("%d", m.stats.TotalSessions))
		totalWorkVal := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorGreen)).Render(m.stats.TotalWorkDuration.Round(time.Second).String())
		avgSessionVal := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorLavender)).Render(m.stats.AverageSessionDuration.Round(time.Second).String())
		longestSessionVal := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorPeach)).Render(m.stats.LongestSession.Round(time.Second).String())

		c1 := cardStyle.Render(fmt.Sprintf("Total Sessions\n\n%s", totalSessionsVal))
		c2 := cardStyle.Render(fmt.Sprintf("Total Focus Work\n\n%s", totalWorkVal))
		c3 := cardStyle.Render(fmt.Sprintf("Average Session\n\n%s", avgSessionVal))
		c4 := cardStyle.Render(fmt.Sprintf("Longest Session\n\n%s", longestSessionVal))

		row1 := lipgloss.JoinHorizontal(lipgloss.Top, c1, c2)
		row2 := lipgloss.JoinHorizontal(lipgloss.Top, c3, c4)

		s.WriteString(row1 + "\n\n" + row2)
	}

	s.WriteString("\n\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, renderHelpKey("q", "Back to Menu")))

	return s.String()
}
