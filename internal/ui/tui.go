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
				m.taskForm = NewTaskFormModel(taskRepo, msg.Task, false)
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

func renderHelpKey(key, desc string) string {
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorLavender)).
		Padding(0, 1).
		Bold(true)
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorText)).
		Padding(0, 1)

	return lipgloss.JoinHorizontal(lipgloss.Left, keyStyle.Render(key), descStyle.Render(desc)) + "  "
}

func (m MainModel) View() string {
	if m.quitting {
		return "\n  See you later!\n\n"
	}

	// Calculate layout dimensions precisely to prevent terminal scrolling/chopping
	headerHeight := 1
	helpBarHeight := 1
	// Total padding/spacers/borders overhead = 4 lines (1 top margin, 1 header spacer, 1 help spacer, 1 bottom padding)
	contentHeight := m.height - headerHeight - helpBarHeight - 4
	if contentHeight < 10 {
		contentHeight = 10
	}

	sidebarWidth := 26
	workspaceWidth := m.width - sidebarWidth - 4
	if workspaceWidth < 40 {
		workspaceWidth = 40
	}

	// 1. Top Header Banner (Tight 1-line budget)
	appTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(colorBase)).
		Background(lipgloss.Color(colorMauve)).
		Padding(0, 1).
		Render("POMOLITE")
	
	appStatus := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorSubtext)).
		PaddingLeft(1).
		Render("⏱ Focus Companion")

	headerBanner := lipgloss.JoinHorizontal(lipgloss.Center, appTitle, appStatus)

	// 2. Sidebar Layout
	var sidebarBuilder strings.Builder
	
	// Menu is exactly 5 choices + "Main Menu" title. Inner height is 6. With borders/padding, it's 8 lines.
	menuInnerHeight := 5
	var menuPanel string
	if m.activePanel == "sidebar" {
		menuPanel = panelFocusedStyle.Width(sidebarWidth - 4).Height(menuInnerHeight).Render(m.home.View())
	} else {
		menuPanel = panelStyle.Width(sidebarWidth - 4).Height(menuInnerHeight).Render(m.home.View())
	}
	sidebarBuilder.WriteString(menuPanel + "\n")

	// Timer widget gets the remaining sidebar height
	widgetInnerHeight := contentHeight - menuInnerHeight - 4
	if widgetInnerHeight < 2 {
		widgetInnerHeight = 2
	}

	var widgetContent string
	if m.state == viewTimer {
		percent := 1.0 - float64(m.timer.remaining)/float64(m.timer.duration)
		mins := int(m.timer.remaining.Minutes())
		secs := int(m.timer.remaining.Seconds()) % 60
		
		status := "RUN"
		statusColor := colorGreen
		if m.timer.paused {
			status = "PAUSE"
			statusColor = colorPeach
		}
		if m.timer.waitingForNext {
			status = "WAIT"
			statusColor = colorLavender
		}
		
		statusBadge := lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color(statusColor)).
			Foreground(lipgloss.Color(colorBase)).
			Padding(0, 1).
			Render(status)

		widgetContent = fmt.Sprintf(
			"Label: %s\nTime:  %02d:%02d\nState: %s\nProgress: %d%%",
			lipgloss.NewStyle().Foreground(lipgloss.Color(colorLavender)).Bold(true).Render(m.timer.label), 
			mins, secs, statusBadge, int(percent*100),
		)
	} else {
		widgetContent = "Timer: Idle\n\nNo active session."
	}
	
	timerWidget := panelStyle.Width(sidebarWidth - 4).Height(widgetInnerHeight).Render(widgetContent)
	sidebarBuilder.WriteString(timerWidget)

	// 3. Workspace Layout
	var workspaceContent string
	switch m.state {
	case "welcome":
		welcomeTitle := lipgloss.NewStyle().Foreground(lipgloss.Color(colorMauve)).Render(ASCIIArt)
		welcomeText := "\nWelcome to PomoLite!\n\nA premium terminal companion.\n\nUse [Tab] to switch focus."
		workspaceContent = welcomeTitle + "\n" + welcomeText
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
		workspacePanel = panelFocusedStyle.Width(workspaceWidth).Height(contentHeight).Render(workspaceContent)
	} else {
		workspacePanel = panelStyle.Width(workspaceWidth).Height(contentHeight).Render(workspaceContent)
	}

	leftCol := sidebarBuilder.String()
	rightCol := workspacePanel
	mainLayout := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	// 4. Bottom Status Help Bar (Tight 1-line budget)
	var helpElements []string
	if m.activePanel == "sidebar" {
		helpElements = []string{
			renderHelpKey("Tab", "Focus Workspace"),
			renderHelpKey("↑/↓", "Navigate"),
			renderHelpKey("Enter", "Select"),
			renderHelpKey("Ctrl+C", "Quit"),
		}
	} else {
		helpElements = []string{
			renderHelpKey("Tab", "Focus Menu"),
			renderHelpKey("Esc", "Return"),
			renderHelpKey("Ctrl+C", "Quit"),
		}
	}
	bottomHelp := lipgloss.JoinHorizontal(lipgloss.Left, helpElements...)
	helpBar := helpStyle.Render(bottomHelp)

	return "\n" + headerBanner + "\n" + mainLayout + "\n" + helpBar
}

