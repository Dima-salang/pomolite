package timer

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
)

const (
	padding  = 2
	maxWidth = 80
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EE6FF8")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

type tickMsg time.Time

type PomoModel struct {
	duration      time.Duration
	remaining     time.Duration
	label         string
	workLabel     string
	paused        bool
	isBreak       bool
	workDuration  time.Duration
	breakDuration time.Duration
	startTime     time.Time
	storage       Storage
	progress      progress.Model
	quitting      bool
}

func NewPomoModel(workDuration, breakDuration time.Duration, label string, storage Storage) PomoModel {
	return PomoModel{
		duration:      workDuration,
		remaining:     workDuration,
		label:         label,
		workLabel:     label,
		workDuration:  workDuration,
		breakDuration: breakDuration,
		startTime:     time.Now(),
		storage:       storage,
		progress:      progress.New(progress.WithDefaultGradient()),
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
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
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
	switch msg.String() {
	case "q", "ctrl+c":
		m.quitting = true
		if !m.isBreak {
			m.storage.SaveTimerData(m.label, m.startTime, time.Now())
		}
		return m, tea.Quit
	case "p", " ":
		m.paused = !m.paused
		return m, nil
	case "r":
		m.paused = false
		return m, nil
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
			// Work finished, start break
			m.storage.SaveTimerData(m.label, m.startTime, time.Now())
			m.isBreak = true
			m.remaining = m.breakDuration
			m.duration = m.breakDuration
			m.label = "Break"
			notify("Work completed!", "Good job! Take a break.")
		} else {
			// Break finished, back to work
			m.isBreak = false
			m.remaining = m.workDuration
			m.duration = m.workDuration
			m.label = m.workLabel
			m.startTime = time.Now()
			notify("Break completed!", "Back to work.")
		}
	}
	return m, tick()
}

func (m PomoModel) View() string {
	if m.quitting {
		return "\n  See you later!\n\n"
	}

	percent := 1.0 - float64(m.remaining)/float64(m.duration)
	if percent < 0 {
		percent = 0
	}

	status := "Running"
	if m.paused {
		status = "Paused"
	}

	s := "\n"
	s += titleStyle.Render(" Pomolite ") + "\n\n"
	s += fmt.Sprintf("  Status: %s\n", statusStyle.Render(status))
	s += fmt.Sprintf("  Phase:  %s\n", labelStyle.Render(m.label))
	s += fmt.Sprintf("  Time:   %s\n\n", m.remaining.Round(time.Second).String())
	s += "  " + m.progress.ViewAs(percent) + "\n\n"
	s += helpStyle.Render("  p: pause/resume • q: quit") + "\n"

	return s
}

func notify(title, message string) {
	_ = beeep.Notify(title, message, "")
	fmt.Print("\a")
	_ = beeep.Beep(500, 200)
}
