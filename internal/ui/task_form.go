package ui

import (
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TaskFormModel struct {
	repo        storage.TaskRepository
	task        storage.Task
	isEdit      bool
	inputs      []textinput.Model
	focused     int
	priorityIdx int // 0 = Low, 1 = Medium, 2 = High
	errorMsg    string
}

func NewTaskFormModel(repo storage.TaskRepository, task storage.Task, isEdit bool) TaskFormModel {
	m := TaskFormModel{
		repo:   repo,
		task:   task,
		isEdit: isEdit,
		inputs: make([]textinput.Model, 2),
	}

	for i := range m.inputs {
		t := textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor))
		t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor))

		switch i {
		case 0:
			t.Placeholder = "Task Title"
			t.Prompt = "Title:       "
			t.Focus()
			if isEdit {
				t.SetValue(task.Title)
			}
		case 1:
			t.Placeholder = "Task Description (optional)"
			t.Prompt = "Description: "
			if isEdit {
				t.SetValue(task.Description)
			}
		}
		m.inputs[i] = t
	}

	if isEdit {
		m.priorityIdx = task.Priority - 1
		if m.priorityIdx < 0 || m.priorityIdx > 2 {
			m.priorityIdx = 0
		}
	}

	return m
}

func (m TaskFormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TaskFormModel) Update(msg tea.Msg) (TaskFormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = ""
		switch msg.String() {
		case "tab", "shift+tab", "up", "down", "enter":
			cmd := m.handleNavigation(msg.String())
			if cmd != nil {
				return m, cmd
			}
		case "h", "left":
			if m.focused == 2 {
				m.priorityIdx = (m.priorityIdx - 1 + 3) % 3
			}
		case "l", "right":
			if m.focused == 2 {
				m.priorityIdx = (m.priorityIdx + 1) % 3
			}
		case "esc":
			return m, func() tea.Msg { return "board" }
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *TaskFormModel) handleNavigation(key string) tea.Cmd {
	if key == "enter" && m.focused == 2 {
		return m.submit()
	}

	if key == "up" || key == "shift+tab" {
		m.focused--
	} else {
		m.focused++
	}

	if m.focused > 2 {
		m.focused = 0
	} else if m.focused < 0 {
		m.focused = 2
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		if i == m.focused {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *TaskFormModel) submit() tea.Cmd {
	title := m.inputs[0].Value()
	if strings.TrimSpace(title) == "" {
		m.errorMsg = "Title cannot be empty."
		return nil
	}

	m.task.Title = title
	m.task.Description = m.inputs[1].Value()
	m.task.Priority = m.priorityIdx + 1

	return func() tea.Msg {
		var err error
		if m.isEdit {
			err = m.repo.UpdateTask(m.task)
		} else {
			if m.task.Status == "" {
				m.task.Status = storage.Backlog
			}
			m.task.CreatedAt = time.Now()
			err = m.repo.AddTask(m.task)
		}
		if err != nil {
			return err
		}
		return "board"
	}
}

func (m *TaskFormModel) updateInputs(msg tea.Msg) tea.Cmd {
	if m.focused < 2 {
		var cmd tea.Cmd
		m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
		return cmd
	}
	return nil
}

func (m TaskFormModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	if m.isEdit {
		s.WriteString(titleStyle.Render(" Edit Task "))
	} else {
		s.WriteString(titleStyle.Render(" Create Task "))
	}
	s.WriteString("\n\n")

	var form strings.Builder
	for i := range m.inputs {
		if m.focused == i {
			pointer := lipgloss.NewStyle().Foreground(lipgloss.Color(colorMauve)).Render("❯ ")
			form.WriteString(pointer + m.inputs[i].View())
		} else {
			form.WriteString("  " + m.inputs[i].View())
		}
		form.WriteString("\n\n")
	}

	priorityLabel := "  Priority:    "
	if m.focused == 2 {
		priorityLabel = lipgloss.NewStyle().Foreground(lipgloss.Color(colorMauve)).Render("❯ Priority:    ")
	}
	form.WriteString(priorityLabel)

	priorities := []string{"Low", "Medium", "High"}
	for idx, name := range priorities {
		isSel := (idx == m.priorityIdx)
		isFocused := (m.focused == 2)

		var rendered string
		if isSel {
			if isFocused {
				rendered = lipgloss.NewStyle().Background(lipgloss.Color(colorMauve)).Foreground(lipgloss.Color(colorBase)).Bold(true).Render(" " + name + " ")
			} else {
				rendered = lipgloss.NewStyle().Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color(colorMauve)).Bold(true).Render(" " + name + " ")
			}
		} else {
			rendered = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Render(" " + name + " ")
		}
		form.WriteString(rendered + "  ")
	}

	s.WriteString(menuStyle.Render(form.String()))

	if m.errorMsg != "" {
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed)).Bold(true).PaddingLeft(2).Render("❌ " + m.errorMsg))
	}

	s.WriteString("\n\n")
	if m.focused == 2 {
		s.WriteString(helpStyle.Render("  tab: navigate • h/l: change priority • enter: save • esc: cancel"))
	} else {
		s.WriteString(helpStyle.Render("  tab: navigate • enter: next • esc: cancel"))
	}

	return s.String()
}
