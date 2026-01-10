package timer

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FormModel struct {
	inputs   []textinput.Model
	focused  int
	errorMsg string
	width    int
	height   int
}

func NewFormModel(workMin, breakMin int, label string) FormModel {
	m := FormModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor))
		t.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(primaryColor))

		switch i {
		case 0:
			t.Placeholder = "30"
			t.Prompt = "Work Duration (min):  "
			t.Focus()
			t.SetValue(strconv.Itoa(workMin))
		case 1:
			t.Placeholder = "5"
			t.Prompt = "Break Duration (min): "
			t.SetValue(strconv.Itoa(breakMin))
		case 2:
			t.Placeholder = "Work"
			t.Prompt = "Session Label:        "
			t.CharLimit = 20
			t.SetValue(label)
		}

		m.inputs[i] = t
	}

	return m
}

func (m FormModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m FormModel) Update(msg tea.Msg) (FormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.errorMsg = "" // Clear error on any keypress
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			cmd := m.handleNavigation(msg.String())
			if cmd != nil {
				return m, cmd
			}
		case "q", "esc":
			if m.focused != len(m.inputs)-1 || msg.String() == "esc" {
				return m, func() tea.Msg { return "home" }
			}
		}
	}

	cmd := m.updateInputs(msg)
	return m, cmd
}

func (m *FormModel) handleNavigation(key string) tea.Cmd {
	if key == "enter" && m.focused == len(m.inputs)-1 {
		return m.submit()
	}

	if key == "up" || key == "shift+tab" {
		m.focused--
	} else {
		m.focused++
	}

	if m.focused > len(m.inputs)-1 {
		m.focused = 0
	} else if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}

	cmds := make([]tea.Cmd, len(m.inputs))
	for i := 0; i <= len(m.inputs)-1; i++ {
		if i == m.focused {
			cmds[i] = m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
	return tea.Batch(cmds...)
}

func (m *FormModel) submit() tea.Cmd {
	workStr := m.inputs[0].Value()
	breakStr := m.inputs[1].Value()

	work, errW := strconv.Atoi(workStr)
	breakMin, errB := strconv.Atoi(breakStr)

	if errW != nil || errB != nil || work <= 0 || breakMin <= 0 {
		m.errorMsg = "Please enter valid positive numbers for durations."
		return nil
	}

	if work > 1440 || breakMin > 1440 {
		m.errorMsg = "Duration cannot exceed 24 hours (1440 min)."
		return nil
	}

	return func() tea.Msg {
		return FormResultMsg{
			WorkMin:  work,
			BreakMin: breakMin,
			Label:    m.inputs[2].Value(),
		}
	}
}

func (m *FormModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only allow numeric input for the first two fields
	if k, ok := msg.(tea.KeyMsg); ok && m.focused < 2 {
		s := k.String()
		if len(s) == 1 && (s < "0" || s > "9") {
			return nil
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return tea.Batch(cmds...)
}

func (m FormModel) View() string {
	var s strings.Builder

	s.WriteString("\n")
	s.WriteString(headerStyle.Render("Timer Configuration"))
	s.WriteString("\n")

	var form strings.Builder
	for i := range m.inputs {
		form.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			form.WriteString("\n")
		}
	}

	s.WriteString(menuStyle.Render(form.String()))

	if m.errorMsg != "" {
		s.WriteString("\n")
		s.WriteString(lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true).PaddingLeft(2).Render("❌ " + m.errorMsg))
	}

	s.WriteString("\n\n")
	s.WriteString(helpStyle.Render("  tab: navigate • enter: start • q/esc: back"))

	return s.String()
}

type FormResultMsg struct {
	WorkMin  int
	BreakMin int
	Label    string
}
