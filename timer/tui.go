package timer

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/beeep"
)

const (
	padding  = 2
	maxWidth = 80

	viewHome    = "home"
	viewForm    = "form"
	viewTimer   = "timer"
	viewStats   = "stats"
	viewHistory = "history"
)

type tickMsg time.Time

type MainModel struct {
	state     string
	home      HomeModel
	form      FormModel
	timer     PomoModel
	stats     StatsModel
	history   HistoryModel
	storage   Storage
	width     int
	height    int
	quitting  bool
	workDur   time.Duration
	breakDur  time.Duration
	workLabel string
}

func NewMainModel(workDur, breakDur time.Duration, label string, storage Storage) MainModel {
	return MainModel{
		state:     viewHome,
		home:      NewHomeModel(),
		storage:   storage,
		workDur:   workDur,
		breakDur:  breakDur,
		workLabel: label,
	}
}

func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		if msg.String() == "q" && m.state == viewHome {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.home.width, m.home.height = msg.Width, msg.Height
	case string:
		return m.handleMenuSelection(msg)
	case FormResultMsg:
		return m.handleFormResult(msg)
	}

	return m.updateViewState(msg)
}

func (m MainModel) handleMenuSelection(choice string) (tea.Model, tea.Cmd) {
	switch choice {
	case "Start Timer":
		m.form = NewFormModel(int(m.workDur.Minutes()), int(m.breakDur.Minutes()), m.workLabel)
		m.state = viewForm
		return m, m.form.Init()
	case "Statistics":
		m.stats = NewStatsModel(m.storage)
		m.state = viewStats
		return m, m.stats.Init()
	case "History":
		m.history = NewHistoryModel(m.storage)
		m.state = viewHistory
		return m, m.history.Init()
	case "Quit":
		m.quitting = true
		return m, tea.Quit
	case "home":
		m.state = viewHome
		return m, nil
	}
	return m, nil
}

func (m MainModel) handleFormResult(msg FormResultMsg) (tea.Model, tea.Cmd) {
	m.workDur = time.Duration(msg.WorkMin) * time.Minute
	m.breakDur = time.Duration(msg.BreakMin) * time.Minute
	m.workLabel = msg.Label
	m.timer = NewPomoModel(m.workDur, m.breakDur, m.workLabel, m.storage)
	m.state = viewTimer
	return m, m.timer.Init()
}

func (m MainModel) updateViewState(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch m.state {
	case viewHome:
		m.home, cmd = m.home.Update(msg)
	case viewForm:
		m.form, cmd = m.form.Update(msg)
	case viewTimer:
		var newModel tea.Model
		newModel, cmd = m.timer.Update(msg)
		m.timer = newModel.(PomoModel)
		if m.timer.quitting {
			m.timer.quitting = false
			m.state = viewHome
			return m, nil
		}
	case viewStats, viewHistory:
		if k, ok := msg.(tea.KeyMsg); ok && (k.String() == "q" || k.String() == "esc") {
			m.state = viewHome
			return m, nil
		}
		if m.state == viewStats {
			m.stats, cmd = m.stats.Update(msg)
		} else {
			m.history, cmd = m.history.Update(msg)
		}
	}
	return m, cmd
}

func (m MainModel) View() string {
	if m.quitting {
		return "\n  See you later!\n\n"
	}
	var content string
	switch m.state {
	case viewHome:
		content = m.home.View()
	case viewForm:
		content = m.form.View()
	case viewTimer:
		content = m.timer.View()
	case viewStats:
		content = m.stats.View()
	case viewHistory:
		content = m.history.View()
	default:
		content = "Unknown state"
	}
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}

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
		progress:      progress.New(progress.WithDefaultGradient(), progress.WithWidth(maxWidth)),
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
	switch msg.String() {
	case "q", "esc":
		if !m.isBreak {
			m.storage.SaveTimerData(m.label, m.startTime, time.Now())
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
			m.storage.SaveTimerData(m.label, m.startTime, time.Now())
			m.isBreak = true
			m.remaining = m.breakDuration
			m.duration = m.breakDuration
			m.label = "Break"
			notify("Work completed!", "Good job! Take a break.")
		} else {
			// Instead of auto-restart, let's go back to menu or ask
			notify("Break completed!", "Focus session finished.")
			m.quitting = true
			return m, nil
		}
	}
	return m, tick()
}

func (m PomoModel) View() string {
	percent := 1.0 - float64(m.remaining)/float64(m.duration)
	status := "Running"
	if m.paused {
		status = "Paused"
	}

	// Format time: MM:SS
	mins := int(m.remaining.Minutes())
	secs := int(m.remaining.Seconds()) % 60
	timeStr := fmt.Sprintf("%02d:%02d", mins, secs)
	bigTime := renderBigText(timeStr)

	var s strings.Builder

	// Center the entire block
	contentStyle := lipgloss.NewStyle().Align(lipgloss.Center)

	header := titleStyle.Render(fmt.Sprintf(" %s ", strings.ToUpper(m.label)))
	s.WriteString(contentStyle.Width(m.progress.Width+4).Render(header) + "\n\n")

	s.WriteString(contentStyle.Width(m.progress.Width+4).Render(timerStyle.Render(bigTime)) + "\n\n")

	s.WriteString("  " + m.progress.ViewAs(percent) + "\n\n")

	statusInfo := fmt.Sprintf("Status: %s", statusStyle.Render(status))
	s.WriteString(contentStyle.Width(m.progress.Width+4).Render(statusInfo) + "\n\n")

	s.WriteString(contentStyle.Width(m.progress.Width + 4).Render(helpStyle.Render("p: pause • r: resume • q: quit")))

	return s.String()
}

func notify(title, message string) {
	_ = beeep.Notify(title, message, "")
	fmt.Print("\a")
	_ = beeep.Beep(500, 200)
}
