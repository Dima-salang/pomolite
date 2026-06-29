package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
)

type tickMsg time.Time

type PomoModel struct {
	duration       time.Duration
	remaining      time.Duration
	label          string
	workLabel      string
	paused         bool
	isBreak        bool
	workDuration   time.Duration
	breakDuration  time.Duration
	startTime      time.Time
	storage        storage.Repository
	progress       progress.Model
	quitting       bool
	waitingForNext bool
}

func NewPomoModel(workDuration, breakDuration time.Duration, label string, store storage.Repository) PomoModel {
	return PomoModel{
		duration:      workDuration,
		remaining:     workDuration,
		label:         label,
		workLabel:     label,
		workDuration:  workDuration,
		breakDuration: breakDuration,
		startTime:     time.Now(),
		storage:       store,
		progress:      progress.New(progress.WithScaledGradient(colorMauve, colorLavender), progress.WithWidth(maxWidth)),
	}
}

func (m PomoModel) Init() tea.Cmd {
	return tick()
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m PomoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg)
	case tea.WindowSizeMsg:
		m.progress.Width = min(msg.Width-padding*2-4, maxWidth)
		return m, nil
	case tickMsg:
		return m.handleTick()
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		m.progress = newModel.(progress.Model)
		return m, cmd
	}
	return m, nil
}

func (m PomoModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.waitingForNext {
		switch msg.String() {
		case "y", "enter":
			m.waitingForNext = false
			m.paused = false
			return m, nil
		case "n", "q", "esc":
			m.quitting = true
			return m, nil
		}
	}

	switch msg.String() {
	case "q", "esc":
		if !m.isBreak {
			err := m.storage.Save(m.label, m.startTime, time.Now())
			if err != nil {
				fmt.Println("Error saving timer data:", err)
			}
		}
		m.quitting = true
		return m, nil
	case "ctrl+c":
		m.quitting = true
		return m, tea.Quit
	case "p", " ":
		m.paused = !m.paused
	case "r":
		m.paused = false
	}
	return m, nil
}

func (m PomoModel) handleTick() (tea.Model, tea.Cmd) {
	if m.paused {
		return m, tick()
	}

	m.remaining -= time.Second
	if m.remaining <= 0 {
		if !m.isBreak {
			err := m.storage.Save(m.label, m.startTime, time.Now())
			if err != nil {
				fmt.Println("Error saving timer data:", err)
			}
			m.isBreak = true
			m.remaining = m.breakDuration
			m.duration = m.breakDuration
			m.label = "Break"
			notify("Work completed!", "Good job! Take a break.")
		} else {
			// resume timer
			m.isBreak = false
			m.remaining = m.workDuration
			m.duration = m.workDuration
			m.label = m.workLabel
			notify("Break completed!", "Focus session finished.")
		}
		m.paused = true
		m.waitingForNext = true
	}
	return m, tick()
}

func (m PomoModel) View() string {
	percent := 1.0 - float64(m.remaining)/float64(m.duration)
	
	var statusBadge string
	if m.waitingForNext {
		statusBadge = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(colorLavender)).Foreground(lipgloss.Color(colorBase)).Padding(0, 2).Render(" ◌ WAITING ")
	} else if m.paused {
		statusBadge = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(colorPeach)).Foreground(lipgloss.Color(colorBase)).Padding(0, 2).Render(" ⏸ PAUSED ")
	} else {
		statusBadge = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(colorGreen)).Foreground(lipgloss.Color(colorBase)).Padding(0, 2).Render(" ● RUNNING ")
	}

	mins := int(m.remaining.Minutes())
	secs := int(m.remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)
	bigTime := renderBigText(timeStr)

	var s strings.Builder
	contentStyle := lipgloss.NewStyle().Align(lipgloss.Center)

	var headerStyle lipgloss.Style
	if m.isBreak {
		headerStyle = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(colorPeach)).Foreground(lipgloss.Color(colorBase)).Padding(0, 3)
	} else {
		headerStyle = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color(colorGreen)).Foreground(lipgloss.Color(colorBase)).Padding(0, 3)
	}
	header := headerStyle.Render(fmt.Sprintf(" %s ", strings.ToUpper(m.label)))
	
	s.WriteString(contentStyle.Width(m.progress.Width+4).Render(header) + "\n\n")
	s.WriteString(contentStyle.Width(m.progress.Width+4).Render(timerStyle.Render(bigTime)) + "\n\n")
	s.WriteString("  " + m.progress.ViewAs(percent) + "\n\n")

	if m.waitingForNext {
		nextType := "work"
		if m.isBreak {
			nextType = "break"
		}
		prompt := fmt.Sprintf("Start %s session? (y/n)", nextType)
		s.WriteString(contentStyle.Width(m.progress.Width+4).Render(promptStyle.Render(prompt)) + "\n\n")
	} else {
		statusInfo := fmt.Sprintf("Status: %s", statusBadge)
		s.WriteString(contentStyle.Width(m.progress.Width+4).Render(statusInfo) + "\n\n")
	}

	helpText := "p: pause • r: resume • q: quit"
	if m.waitingForNext {
		helpText = "y: start • n: main menu"
	}
	s.WriteString(contentStyle.Width(m.progress.Width + 4).Render(helpStyle.Render(helpText)))

	return s.String()
}

func notify(title, message string) {
	_ = beeep.Notify(title, message, "")
	fmt.Print("\a")
	_ = beeep.Beep(500, 200)
}
