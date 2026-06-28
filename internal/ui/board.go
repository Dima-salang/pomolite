package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/Dima-salang/pomolite/internal/storage"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type BoardActionMsg struct {
	Action string
	Task   storage.Task
}

type BoardModel struct {
	repo       storage.TaskRepository
	columns    []storage.TaskStatus
	focusCol   int
	focusTask  [4]int // Focus index per column
	tasksMap   map[storage.TaskStatus][]storage.Task
	width      int
	height     int
	err        error
}

func NewBoardModel(repo storage.TaskRepository) BoardModel {
	return BoardModel{
		repo:     repo,
		columns:  []storage.TaskStatus{storage.Backlog, storage.Todo, storage.InProgress, storage.Done},
		tasksMap: make(map[storage.TaskStatus][]storage.Task),
		width:    80,
		height:   20,
	}
}


func (m BoardModel) Init() tea.Cmd {
	return m.ReloadTasksCmd()
}

func (m BoardModel) ReloadTasksCmd() tea.Cmd {
	return func() tea.Msg {
		tasks, err := m.repo.ListTasks(nil) // List all tasks
		if err != nil {
			return err
		}
		return tasks
	}
}

func (m BoardModel) Update(msg tea.Msg) (BoardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.err = msg
		return m, nil
	case []storage.Task:
		m.err = nil
		// Clear tasks
		for _, col := range m.columns {
			m.tasksMap[col] = []storage.Task{}
		}
		// Distribute tasks
		for _, t := range msg {
			m.tasksMap[t.Status] = append(m.tasksMap[t.Status], t)
		}
		// Adjust cursor boundaries
		for i, col := range m.columns {
			limit := len(m.tasksMap[col])
			if m.focusTask[i] >= limit {
				m.focusTask[i] = max(0, limit-1)
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, func() tea.Msg { return "home" }

		case "h": // Move focus left
			m.focusCol = (m.focusCol - 1 + len(m.columns)) % len(m.columns)
			return m, nil

		case "l": // Move focus right
			m.focusCol = (m.focusCol + 1) % len(m.columns)
			return m, nil

		case "j": // Move task cursor down
			activeTasks := m.tasksMap[m.columns[m.focusCol]]
			if len(activeTasks) > 0 {
				m.focusTask[m.focusCol] = (m.focusTask[m.focusCol] + 1) % len(activeTasks)
			}
			return m, nil

		case "k": // Move task cursor up
			activeTasks := m.tasksMap[m.columns[m.focusCol]]
			if len(activeTasks) > 0 {
				m.focusTask[m.focusCol] = (m.focusTask[m.focusCol] - 1 + len(activeTasks)) % len(activeTasks)
			}
			return m, nil

		case "H": // Move selected task left
			return m.moveTask(-1)

		case "L": // Move selected task right
			return m.moveTask(1)

		case "a": // Add task
			return m, func() tea.Msg {
				return BoardActionMsg{Action: "add"}
			}

		case "e": // Edit task
			activeTasks := m.tasksMap[m.columns[m.focusCol]]
			if len(activeTasks) > 0 {
				selected := activeTasks[m.focusTask[m.focusCol]]
				return m, func() tea.Msg {
					return BoardActionMsg{Action: "edit", Task: selected}
				}
			}

		case "d": // Delete task
			activeTasks := m.tasksMap[m.columns[m.focusCol]]
			if len(activeTasks) > 0 {
				selected := activeTasks[m.focusTask[m.focusCol]]
				return m, m.deleteTaskCmd(selected.ID)
			}
		}
	}

	return m, nil

}

func (m *BoardModel) moveTask(direction int) (BoardModel, tea.Cmd) {
	colIndex := m.focusCol
	activeCol := m.columns[colIndex]
	activeTasks := m.tasksMap[activeCol]

	if len(activeTasks) == 0 {
		return *m, nil
	}

	selected := activeTasks[m.focusTask[colIndex]]
	newColIndex := colIndex + direction
	if newColIndex < 0 || newColIndex >= len(m.columns) {
		return *m, nil // Out of bounds
	}

	newStatus := m.columns[newColIndex]
	selected.Status = newStatus

	if newStatus == storage.Done {
		selected.CompletedAt = time.Now()
	} else {
		selected.CompletedAt = time.Time{}
	}

	// Create command to update in SQLite and reload
	cmd := func() tea.Msg {
		err := m.repo.UpdateTask(selected)
		if err != nil {
			return err
		}
		return m.ReloadTasksCmd()()
	}

	// Update focus target in the new column
	m.focusCol = newColIndex
	m.focusTask[newColIndex] = len(m.tasksMap[newStatus]) // Append to end

	return *m, cmd
}

func (m BoardModel) deleteTaskCmd(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.repo.DeleteTask(id)
		if err != nil {
			return err
		}
		return m.ReloadTasksCmd()()
	}
}


func (m BoardModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(titleStyle.Render(" Pomolite Task Board "))
	s.WriteString("\n\n")

	if m.err != nil {
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Render(fmt.Sprintf("Error: %v\n", m.err)))
	}

	var cols []string
	for i, col := range m.columns {
		var colContent strings.Builder
		header := strings.ToUpper(string(col))
		tasks := m.tasksMap[col]

		colHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(primaryColor))
		colContent.WriteString(colHeaderStyle.Render(fmt.Sprintf("=== %s (%d) ===", header, len(tasks))) + "\n\n")

		if len(tasks) == 0 {
			colContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(gray)).Italic(true).Render("  (No tasks)") + "\n")
		} else {
			for j, t := range tasks {
				isSelected := (i == m.focusCol && j == m.focusTask[i])
				colContent.WriteString(m.renderTaskCard(t, isSelected) + "\n")
			}
		}

		if i == m.focusCol {
			cols = append(cols, columnFocusedStyle.Render(colContent.String()))
		} else {
			cols = append(cols, columnStyle.Render(colContent.String()))
		}
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cols...))
	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("  h/l: focus col • j/k: focus task • H/L: move task left/right\n  a: add • e: edit • d: delete • q: main menu"))

	return s.String()
}

func (m BoardModel) renderTaskCard(t storage.Task, selected bool) string {
	var card strings.Builder

	var priorityStr string
	switch t.Priority {
	case 1:
		priorityStr = lowPriorityStyle.Render("Low")
	case 2:
		priorityStr = medPriorityStyle.Render("Medium")
	case 3:
		priorityStr = highPriorityStyle.Render("High")
	default:
		priorityStr = lipgloss.NewStyle().Foreground(lipgloss.Color(gray)).Render("Unknown")
	}

	prefix := "  "
	if selected {
		prefix = "▶ "
	}

	card.WriteString(fmt.Sprintf("%s[%s] %s\n", prefix, priorityStr, t.Title))
	if t.Description != "" {
		card.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(gray)).Render("    "+t.Description) + "\n")
	}
	if !t.DueDate.IsZero() {
		card.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(accentColor)).Render("    📅 "+t.DueDate.Format("Jan 02")) + "\n")
	}

	if selected {
		return cardSelectedStyle.Render(card.String())
	}
	return cardStyle.Render(card.String())
}


