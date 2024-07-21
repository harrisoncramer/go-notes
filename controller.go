package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

const otherView = "other"

/* The view function is responsible for rendering different screens */
func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	s := ""

	switch m.view {
	case otherView:
		s += "Other View\n"
	default:
		s += m.mainRenderer(s)
	}

	s += "\nPress <C-c> to quit\n"
	return s
}

/* The update function is responsible for getting a message from the Init function, and updating state in the model */
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case errMsg:
		m.err = msg
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m.handleCtrlC()
		}
	}

	return m.mainController(msg)
}

/* The main controller is the root of the application */
func (m *model) mainController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.handleUpKey()
		case "down", "j":
			m.handleDownKey()
		case "enter":
			choice := m.viewData.choices[m.cursor.idx]
			switch choice.Text {
			case "Add Entry":
				return m.changeView(otherView)
			case "Edit Entry":
				return m.changeView(otherView)
			case "Rename Entry":
				return m.changeView(otherView)
			}
		}
	}

	return m, nil /* Return the new model state */
}

func (m *model) mainRenderer(s string) string {
	for i, choice := range m.viewData.choices {
		prefix := " "
		if m.cursor.idx == i {
			prefix = ">"
		}
		s += fmt.Sprintf("%s %s\n", prefix, choice.Text)
	}

	return s
}

func (m *model) handleCtrlC() (model, tea.Cmd) {
	return *m, tea.Quit
}

func (m *model) handleUpKey() {
	if m.cursor.idx > 0 {
		m.cursor.idx--
	}
}

func (m *model) handleDownKey() {
	if m.cursor.idx < len(m.viewData.choices)-1 {
		m.cursor.idx++
	}
}

func (m *model) changeView(view string) (tea.Model, tea.Cmd) {
	m.viewData = ViewData{}
	m.view = view
	return m, nil
}
