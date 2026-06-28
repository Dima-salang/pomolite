package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80

	viewHome     = "home"
	viewForm     = "form"
	viewTimer    = "timer"
	viewStats    = "stats"
	viewHistory  = "history"
	viewBoard    = "board"
	viewTaskForm = "taskForm"
)

type MainModel struct {
	state       string
	activePanel string // "sidebar" or "workspace"
	home        HomeModel
	form        FormModel
	timer       PomoModel
	stats       StatsModel
	history     HistoryModel
	board       BoardModel
	taskForm    TaskFormModel
	storage     storage.Repository
	width       int
	height      int
	quitting    bool
	workDur     time.Duration
	breakDur    time.Duration
	workLabel   string
}

func NewMainModel(workDur, breakDur time.Duration, label string, store storage.Repository) MainModel {
	taskRepo, _ := store.(storage.TaskRepository)
	return MainModel{
		state:       "welcome",
		activePanel: "sidebar",
		home:        NewHomeModel(),
		storage:     store,
		workDur:     workDur,
		breakDur:    breakDur,
		workLabel:   label,
		board:       NewBoardModel(taskRepo),
		width:       100,
		height:      24,
	}
}


func (m MainModel) Init() tea.Cmd {
	return nil
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		if msg.String() == "tab" {
			if m.activePanel == "sidebar" {
				if m.state == viewBoard || m.state == viewForm || m.state == viewTimer || m.state == viewTaskForm {
					m.activePanel = "workspace"
				}
			} else {
				m.activePanel = "sidebar"
			}
			return m, nil
		}

		if msg.String() == "esc" && m.activePanel == "workspace" {
			m.activePanel = "sidebar"
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.home.width, m.home.height = msg.Width, msg.Height
		sidebarWidth := 26
		workspaceWidth := m.width - sidebarWidth - 4
		workspaceHeight := m.height - 4
		m.board.width = workspaceWidth
		m.board.height = workspaceHeight

		var cmd tea.Cmd
		m.home, cmd = m.home.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		m.board, cmd = m.board.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
		var newTimer tea.Model
		newTimer, cmd = m.timer.Update(msg)
		m.timer = newModelToPomo(newTimer)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}



	if m.activePanel == "sidebar" {
		var cmd tea.Cmd
		m.home, cmd = m.home.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		selectedChoice := m.home.choices[m.home.cursor]
		switch selectedChoice {
		case "Start Timer":
			if m.state != viewTimer && m.state != viewForm {
				m.form = NewFormModel(int(m.workDur.Minutes()), int(m.breakDur.Minutes()), m.workLabel)
				m.state = viewForm
				cmds = append(cmds, m.form.Init())
			}
		case "Task Board":
			if m.state != viewBoard {
				m.state = viewBoard
				cmds = append(cmds, m.board.ReloadTasksCmd())
			}
		case "Statistics":
			if m.state != viewStats {
				m.stats = NewStatsModel(m.storage)
				m.state = viewStats
				cmds = append(cmds, m.stats.Init())
			}
		case "History":
			if m.state != viewHistory {
				m.history = NewHistoryModel(m.storage)
				m.state = viewHistory
				cmds = append(cmds, m.history.Init())
			}
		case "Quit":
			m.state = "quit_screen"
		}

		if k, ok := msg.(tea.KeyMsg); ok && k.String() == "enter" {
			if selectedChoice == "Quit" {
				m.quitting = true
				return m, tea.Quit
			}
			if m.state == viewBoard || m.state == viewForm || m.state == viewTimer {
				m.activePanel = "workspace"
			}
		}

	} else {
		var cmd tea.Cmd
		switch m.state {
		case viewForm:
			m.form, cmd = m.form.Update(msg)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		case viewTimer:
			var newTimer tea.Model
			newTimer, cmd = m.timer.Update(msg)
			m.timer = newModelToPomo(newTimer)
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			if m.timer.quitting {
				m.timer.quitting = false
				m.state = viewForm
				m.activePanel = "sidebar"
			}
		case viewBoard:
			var newBoard BoardModel
			newBoard, cmd = m.board.Update(msg)
			m.board = newBoard
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			if msgStr, ok := msg.(string); ok && msgStr == "home" {
				m.activePanel = "sidebar"
			}
		case viewTaskForm:
			var newForm TaskFormModel
			newForm, cmd = m.taskForm.Update(msg)
			m.taskForm = newForm
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
		}
	}

	switch msg := msg.(type) {
	case FormResultMsg:
		m.workDur = time.Duration(msg.WorkMin) * time.Minute
		m.breakDur = time.Duration(msg.BreakMin) * time.Minute
		m.workLabel = msg.Label
		m.timer = NewPomoModel(m.workDur, m.breakDur, m.workLabel, m.storage)
		m.state = viewTimer
		m.activePanel = "workspace"
		cmds = append(cmds, m.timer.Init())

	case BoardActionMsg:
		if taskRepo, ok := m.storage.(storage.TaskRepository); ok {
			if msg.Action == "add" {
				m.taskForm = NewTaskFormModel(taskRepo, storage.Task{}, false)
				m.state = viewTaskForm
				m.activePanel = "workspace"
				cmds = append(cmds, m.taskForm.Init())
			} else if msg.Action == "edit" {
				m.taskForm = NewTaskFormModel(taskRepo, msg.Task, true)
				m.state = viewTaskForm
				m.activePanel = "workspace"
				cmds = append(cmds, m.taskForm.Init())
			}
		}
	case string:
		if msg == "board" {
			m.state = viewBoard
			m.activePanel = "workspace"
			cmds = append(cmds, m.board.ReloadTasksCmd())
		}
	case []storage.Task:
		var cmd tea.Cmd
		m.board, cmd = m.board.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func newModelToPomo(m tea.Model) PomoModel {
	if p, ok := m.(PomoModel); ok {
		return p
	}
	return PomoModel{}
}

func (m MainModel) View() string {
	if m.quitting {
		return "\n  See you later!\n\n"
	}

	sidebarWidth := 26
	
	helpBarHeight := 2
	contentHeight := m.height - helpBarHeight - 2
	if contentHeight < 12 {
		contentHeight = 12
	}

	menuHeight := 7
	widgetHeight := contentHeight - menuHeight - 2
	if widgetHeight < 4 {
		widgetHeight = 4
	}

	workspaceWidth := m.width - sidebarWidth - 4
	if workspaceWidth < 40 {
		workspaceWidth = 40
	}

	var sidebarBuilder strings.Builder

	var menuPanel string
	if m.activePanel == "sidebar" {
		menuPanel = panelFocusedStyle.Width(sidebarWidth - 4).Render(m.home.View())
	} else {
		menuPanel = panelStyle.Width(sidebarWidth - 4).Render(m.home.View())
	}
	sidebarBuilder.WriteString(menuPanel + "\n\n")

	var widgetContent string
	if m.state == viewTimer {
		percent := 1.0 - float64(m.timer.remaining)/float64(m.timer.duration)
		mins := int(m.timer.remaining.Minutes())
		secs := int(m.timer.remaining.Seconds()) % 60
		status := "RUNNING"
		if m.timer.paused {
			status = "PAUSED"
		}
		if m.timer.waitingForNext {
			status = "WAITING"
		}
		widgetContent = fmt.Sprintf(
			"Session: %s\nTime:    %02d:%02d\nStatus:  %s\nProgress: %d%%",
			m.timer.label, mins, secs, status, int(percent*100),
		)
	} else {
		widgetContent = "Timer Status: Idle\n\nNo active session.\nSelect 'Start Timer'\nto begin."
	}
	timerWidget := panelStyle.Width(sidebarWidth - 4).Render(widgetContent)
	sidebarBuilder.WriteString(timerWidget)

	var workspaceContent string
	switch m.state {
	case "welcome":
		workspaceContent = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor)).Render(ASCIIArt) + "\n\nWelcome to PomoLite TUI!\nUse Tab to focus the workspace panel."
	case viewForm:
		workspaceContent = m.form.View()
	case viewTimer:
		workspaceContent = m.timer.View()
	case viewStats:
		workspaceContent = m.stats.View()
	case viewHistory:
		workspaceContent = m.history.View()
	case viewBoard:
		workspaceContent = m.board.View()
	case viewTaskForm:
		workspaceContent = m.taskForm.View()
	case "quit_screen":
		workspaceContent = "\nReady to quit? Press Enter on sidebar choice to exit."
	default:
		workspaceContent = "Unknown workspace view"
	}

	var workspacePanel string
	if m.activePanel == "workspace" {
		workspacePanel = panelFocusedStyle.Width(workspaceWidth).Render(workspaceContent)
	} else {
		workspacePanel = panelStyle.Width(workspaceWidth).Render(workspaceContent)
	}

	leftCol := sidebarBuilder.String()
	rightCol := workspacePanel


	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)


	var bottomHelp string
	if m.activePanel == "sidebar" {
		bottomHelp = "  tab: focus workspace • ↑/↓: navigate menu • enter: select/activate • ctrl+c: quit"
	} else {
		bottomHelp = "  tab: focus sidebar • esc: return to sidebar • ctrl+c: quit"
	}
	helpBar := helpStyle.Render(bottomHelp)

	return "\n" + mainLayout + "\n" + helpBar
}

