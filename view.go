package main

import "fmt"

const addEntryView = "addEntry"
const editEntryView = "editEntry"
const mainView = "main"

/* The view function is responsible for rendering different screens */
func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	s := ""

	switch m.view {
	case mainView:
		s += m.choiceRenderer(s)
	case addEntryView:
		s += m.addEntryRenderer(s)
	case editEntryView:
		s += m.choiceRenderer(s)
	}

	s += "\nPress <C-c> to quit\n"
	return s
}

func (m model) addEntryRenderer(s string) string {
	s += fmt.Sprintf(
		"%s:\n\n%s\n\n",
		"Add Entry",
		m.textInput.View())
	return s
}

func (m model) choiceRenderer(s string) string {
	for i, choice := range m.viewData.choices {
		prefix := " "
		if m.cursor.idx == i {
			prefix = ">"
		}
		s += fmt.Sprintf("%s %s\n", prefix, choice.Text)
	}

	return s
}
