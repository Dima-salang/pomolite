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

func (m HomeModel) View() string {
	var s strings.Builder

	s.WriteString(headerStyle.Render("Main Menu"))
	s.WriteString("\n")

	var menu strings.Builder
	for i, choice := range m.choices {
		if m.cursor == i {
			menu.WriteString(selectedItemStyle.Render(fmt.Sprintf("▶ %s", choice)))
		} else {
			menu.WriteString(itemStyle.Render(choice))
		}
		menu.WriteString("\n")
	}

	s.WriteString(menu.String())
	return s.String()
}

