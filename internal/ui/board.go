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
	repo          storage.TaskRepository
	columns       []storage.TaskStatus
	focusCol      int
	focusTask     [4]int // Focus index per column
	tasksMap      map[storage.TaskStatus][]storage.Task
	width         int
	height        int
	err           error
	trackedTaskID int
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
		
		// If we are tracking a task, find its new index in the active column
		if m.trackedTaskID != 0 {
			activeCol := m.columns[m.focusCol]
			activeTasks := m.tasksMap[activeCol]
			for idx, t := range activeTasks {
				if t.ID == m.trackedTaskID {
					m.focusTask[m.focusCol] = idx
					break
				}
			}
			m.trackedTaskID = 0 // Clear tracking ID
		} else {
			// Adjust cursor boundaries
			for i, col := range m.columns {
				limit := len(m.tasksMap[col])
				if m.focusTask[i] >= limit {
					m.focusTask[i] = max(0, limit-1)
				}
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
			activeCol := m.columns[m.focusCol]
			return m, func() tea.Msg {
				return BoardActionMsg{Action: "add", Task: storage.Task{Status: activeCol}}
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

		case "D": // Move task directly to Done board
			activeTasks := m.tasksMap[m.columns[m.focusCol]]
			if len(activeTasks) > 0 {
				selected := activeTasks[m.focusTask[m.focusCol]]
				selected.Status = storage.Done
				selected.CompletedAt = time.Now()
				m.trackedTaskID = selected.ID
				m.focusCol = 3 // Done column index is 3
				
				cmd := func() tea.Msg {
					err := m.repo.UpdateTask(selected)
					if err != nil {
						return err
					}
					return m.ReloadTasksCmd()()
				}
				return m, cmd
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

	// Update focus target in the new column and track task ID
	m.focusCol = newColIndex
	m.trackedTaskID = selected.ID

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
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorRed)).Render(fmt.Sprintf("Error: %v\n", m.err)))
	}

	// Dynamically calculate column dimensions
	sidebarWidth := 26
	numCols := len(m.columns)
	colWidth := (m.width - sidebarWidth - 10) / numCols
	if colWidth < 18 {
		colWidth = 18
	}

	// Total height minus padding/margins for header and help bar
	colHeight := m.height - 6
	if colHeight < 8 {
		colHeight = 8
	}

	colStyle := columnStyle.Copy().Width(colWidth).Height(colHeight)
	colFocusedStyle := columnFocusedStyle.Copy().Width(colWidth).Height(colHeight)

	// Available vertical space for card rendering (each card takes approx 3-4 lines)
	visibleLimit := (colHeight - 2) / 4
	if visibleLimit < 1 {
		visibleLimit = 1
	}

	var cols []string
	for i, col := range m.columns {
		var colContent strings.Builder
		header := strings.ToUpper(string(col))
		tasks := m.tasksMap[col]

		var headerStr string
		if i == m.focusCol {
			headerStr = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorMauve)).Render(fmt.Sprintf("❯ %s (%d)", header, len(tasks)))
		} else {
			headerStr = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorSubtext)).Render(fmt.Sprintf("  %s (%d)", header, len(tasks)))
		}
		colContent.WriteString(headerStr + "\n")

		if len(tasks) == 0 {
			colContent.WriteString("\n" + lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Italic(true).Render("  (No tasks)") + "\n")
		} else {
			// Scrolling logic: budget heights by counting the lines of virtually rendered cards
			maxLines := colHeight - 4
			if maxLines < 3 {
				maxLines = 3
			}

			focusedIdx := m.focusTask[i]
			if focusedIdx >= len(tasks) {
				focusedIdx = len(tasks) - 1
			}
			if focusedIdx < 0 {
				focusedIdx = 0
			}

			cardLines := make([]int, len(tasks))
			for j, t := range tasks {
				isSelected := (i == m.focusCol && j == focusedIdx)
				cardLines[j] = strings.Count(m.renderTaskCard(t, isSelected, colWidth), "\n") + 1
			}

			startIdx := focusedIdx
			endIdx := focusedIdx
			totalLines := cardLines[focusedIdx]

			for {
				expanded := false
				if startIdx > 0 && totalLines+cardLines[startIdx-1] <= maxLines {
					startIdx--
					totalLines += cardLines[startIdx]
					expanded = true
				}
				if endIdx < len(tasks)-1 && totalLines+cardLines[endIdx+1] <= maxLines {
					endIdx++
					totalLines += cardLines[endIdx]
					expanded = true
				}
				if !expanded {
					break
				}
			}
			endIdx++ // Exclusive bound for slicing

			// Scroll Up Indicator
			if startIdx > 0 {
				colContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorMauve)).Align(lipgloss.Center).Width(colWidth - 2).Render("▲") + "\n")
			} else {
				colContent.WriteString("\n")
			}

			for j := startIdx; j < endIdx; j++ {
				t := tasks[j]
				isSelected := (i == m.focusCol && j == focusedIdx)
				colContent.WriteString(m.renderTaskCard(t, isSelected, colWidth) + "\n")
			}

			// Scroll Down Indicator
			if endIdx < len(tasks) {
				colContent.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color(colorMauve)).Align(lipgloss.Center).Width(colWidth - 2).Render("▼") + "\n")
			}
		}

		if i == m.focusCol {
			cols = append(cols, colFocusedStyle.Render(colContent.String()))
		} else {
			cols = append(cols, colStyle.Render(colContent.String()))
		}
	}

	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, cols...))
	s.WriteString("\n\n")

	// Helper badges in footer
	badges := []string{
		renderHelpKey("h/l", "Focus Column"),
		renderHelpKey("j/k", "Cursor"),
		renderHelpKey("H/L", "Move"),
		renderHelpKey("a", "Add"),
		renderHelpKey("e", "Edit"),
		renderHelpKey("D", "Done"),
		renderHelpKey("d", "Delete"),
		renderHelpKey("q", "Menu"),
	}
	s.WriteString(lipgloss.JoinHorizontal(lipgloss.Left, badges...))

	return s.String()
}

func (m BoardModel) renderTaskCard(t storage.Task, selected bool, colWidth int) string {
	var card strings.Builder

	var priorityName string
	var priorityStr string
	switch t.Priority {
	case 1:
		priorityName = "Low"
		priorityStr = lowPriorityStyle.Render(priorityName)
	case 2:
		priorityName = "Med"
		priorityStr = medPriorityStyle.Render(priorityName)
	case 3:
		priorityName = "High"
		priorityStr = highPriorityStyle.Render(priorityName)
	default:
		priorityName = "None"
		priorityStr = lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Render(priorityName)
	}

	prefix := "  "
	if selected {
		prefix = "❯ "
	}

	cardWidth := colWidth - 6
	if cardWidth < 12 {
		cardWidth = 12
	}

	// We reserve prefix + `[` + priority + `] ` space for title width calculation
	indentWidth := len(prefix) + len(priorityName) + 3
	titleWidth := cardWidth - indentWidth
	if titleWidth < 8 {
		titleWidth = 8
	}

	wrappedTitle := wrapText(t.Title, titleWidth)
	titleLines := strings.Split(wrappedTitle, "\n")

	// Render first line next to prefix and priority
	firstTitleLine := titleLines[0]
	if selected {
		firstTitleLine = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorMauve)).Render(firstTitleLine)
	} else {
		firstTitleLine = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Render(firstTitleLine)
	}
	card.WriteString(fmt.Sprintf("%s[%s] %s", prefix, priorityStr, firstTitleLine))

	// Subsequent lines are indented below matching the title alignment
	indent := strings.Repeat(" ", indentWidth)
	for k := 1; k < len(titleLines); k++ {
		line := titleLines[k]
		if selected {
			line = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(colorMauve)).Render(line)
		} else {
			line = lipgloss.NewStyle().Foreground(lipgloss.Color(colorText)).Render(line)
		}
		card.WriteString("\n" + indent + line)
	}

	if t.Description != "" {
		descWidth := cardWidth - 4
		if descWidth < 8 {
			descWidth = 8
		}
		wrappedDesc := wrapText(t.Description, descWidth)
		lines := strings.Split(wrappedDesc, "\n")
		for _, line := range lines {
			card.WriteString("\n    " + lipgloss.NewStyle().Foreground(lipgloss.Color(colorSubtext)).Render(line))
		}
	}
	if !t.DueDate.IsZero() {
		card.WriteString("\n    " + lipgloss.NewStyle().Foreground(lipgloss.Color(colorLavender)).Render("📅 "+t.DueDate.Format("Jan 02")))
	}

	if selected {
		return lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(colorMauve)).
			PaddingLeft(1).
			Render(card.String())
	}
	return lipgloss.NewStyle().
		PaddingLeft(2).
		Render(card.String())
}

func wrapText(text string, width int) string {
	if width <= 0 {
		return text
	}
	var result []string
	words := strings.Fields(text)
	if len(words) == 0 {
		// No spaces, handle single long word
		return wrapLongWord(text, width)
	}

	var currentLine string
	for _, word := range words {
		// If a single word is longer than width, wrap it by character
		if len(word) > width {
			if len(currentLine) > 0 {
				result = append(result, currentLine)
				currentLine = ""
			}
			parts := wrapLongWord(word, width)
			result = append(result, parts)
			continue
		}

		if len(currentLine)+len(word)+1 > width {
			result = append(result, currentLine)
			currentLine = word
		} else {
			if len(currentLine) == 0 {
				currentLine = word
			} else {
				currentLine += " " + word
			}
		}
	}
	if len(currentLine) > 0 {
		result = append(result, currentLine)
	}
	return strings.Join(result, "\n")
}

func wrapLongWord(word string, width int) string {
	var result []string
	runes := []rune(word)
	for i := 0; i < len(runes); i += width {
		end := i + width
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}
	return strings.Join(result, "\n")
}


