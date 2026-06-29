package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)


type HomeModel struct {
	choices  []string
	cursor   int
	selected map[int]struct{}
	width    int
	height   int
}

func NewHomeModel() HomeModel {
	return HomeModel{
		choices:  []string{"Start Timer", "Task Board", "Statistics", "History", "Quit"},
		selected: make(map[int]struct{}),
	}
}

func (m HomeModel) Init() tea.Cmd {
	return nil
}

func (m HomeModel) Update(msg tea.Msg) (HomeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			return m, m.selectChoice()
		}
	}
	return m, nil
}

func (m HomeModel) selectChoice() tea.Cmd {
	return func() tea.Msg {
		return m.choices[m.cursor]
	}
}

var choiceIcons = map[string]string{
	"Start Timer": "⏱  Start Timer",
	"Task Board":  "📋  Task Board",
	"Statistics":  "📊  Statistics",
	"History":     "📜  History",
	"Quit":        "✖  Quit",
}

func (m HomeModel) View() string {
	var s strings.Builder

	s.WriteString(headerStyle.Render("Main Menu"))
	s.WriteString("\n\n")

	for i, choice := range m.choices {
		iconChoice, ok := choiceIcons[choice]
		if !ok {
			iconChoice = choice
		}

		if m.cursor == i {
			s.WriteString(selectedItemStyle.Render(fmt.Sprintf("❯ %s", iconChoice)))
		} else {
			s.WriteString(itemStyle.Render(iconChoice))
		}
		s.WriteString("\n")
	}

	return s.String()
}

