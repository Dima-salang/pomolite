package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/stats"
	"github.com/Dima-salang/pomolite/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
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
	s.WriteString(headerStyle.Render("Statistics"))
	s.WriteString("\n")

	if m.err != nil {
		s.WriteString(fmt.Sprintf("  Error loading stats: %v\n", m.err))
	} else if m.stats != nil {
		content := fmt.Sprintf(
			"  Total Sessions:     %d\n"+
				"  Total Work:         %s\n"+
				"  Avg Session:        %s\n"+
				"  Longest Session:    %s\n",
			m.stats.TotalSessions,
			m.stats.TotalWorkDuration.Round(time.Second),
			m.stats.AverageSessionDuration.Round(time.Second),
			m.stats.LongestSession.Round(time.Second),
		)
		s.WriteString(menuStyle.Render(content))
	}

	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("  q: back to menu"))

	return s.String()
}
