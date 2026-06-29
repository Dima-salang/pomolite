package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HistoryModel struct {
	sessions []storage.Session
	err      error
	width    int
	height   int
}

func NewHistoryModel(repo storage.Repository) HistoryModel {
	sessions, err := repo.List(10) // Show last 10
	return HistoryModel{
		sessions: sessions,
		err:      err,
	}
}

func (m HistoryModel) Init() tea.Cmd {
	return nil
}

func (m HistoryModel) Update(msg tea.Msg) (HistoryModel, tea.Cmd) {
	return m, nil
}

func (m HistoryModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" Pomolite "))
	s.WriteString("\n\n")
	s.WriteString(headerStyle.Render("Recent Sessions"))
	s.WriteString("\n")

	if m.err != nil {
		s.WriteString(fmt.Sprintf("  Error loading history: %v\n", m.err))
	} else if len(m.sessions) == 0 {
		s.WriteString(itemStyle.Render("  No sessions found yet. Get to work!"))
	} else {
		var tableBuilder strings.Builder
		
		// Header row
		headerDate := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorMauve)).Width(18).Render("Date & Time")
		headerLabel := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorLavender)).Width(15).Render("Session Label")
		headerDuration := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorGreen)).Width(12).Render("Duration")
		
		tableBuilder.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, headerDate, headerLabel, headerDuration) + "\n")
		tableBuilder.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorSurface)).Render(strings.Repeat("─", 45)) + "\n")

		for _, sess := range m.sessions {
			duration := sess.EndTime.Sub(sess.StartTime).Round(time.Second)
			dateStr := sess.StartTime.Format("Jan 02 15:04")
			
			colDate := lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Width(18).Render(dateStr)
			colLabel := lipgloss.NewStyle().Foreground(lipgloss.Color(colorLavender)).Width(15).Render(sess.Label)
			colDuration := lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Width(12).Render(duration.String())
			
			tableBuilder.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, colDate, colLabel, colDuration) + "\n")
		}
		s.WriteString(menuStyle.Render(tableBuilder.String()))
	}

	s.WriteString("\n\n")
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, renderHelpKey("q", "Back to Menu")))

	return s.String()
}
