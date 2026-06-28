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
		var content strings.Builder
		for _, sess := range m.sessions {
			duration := sess.EndTime.Sub(sess.StartTime).Round(time.Second)
			date := sess.StartTime.Format("Jan 02 15:04")
			content.WriteString(fmt.Sprintf("%s │ %-10s │ %s\n", date, lipgloss.NewStyle().Foreground(lipgloss.Color(accentColor)).Render(sess.Label), duration))
		}
		s.WriteString(menuStyle.Render(content.String()))
	}

	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("  q: back to menu"))

	return s.String()
}
