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

	s := fmt.Sprintf("%s 📓 \n\n", m.dbName)

	switch m.view {
	case mainView:
		s += m.choiceRenderer()
	case addEntryView:
		s += m.addEntryRenderer()
	case editEntryView:
		s += m.choiceRenderer()
	}

	s += "\nPress <C-c> to quit\n"
	return s
}

func (m model) addEntryRenderer() string {
	s := ""
	s += fmt.Sprintf(
		"%s:\n\n%s\n\n",
		"Add Entry",
		m.textInput.View())
	return s
}

func (m model) choiceRenderer() string {
	s := ""
	if len(m.viewData.entries) == 0 {
		s += "No entries found!\n"
		return s
	}

	for i, choice := range m.viewData.entries {
		prefix := " "
		if m.cursor.idx == i {
			prefix = ">"
		}
		s += fmt.Sprintf("%s %s\n", prefix, choice.Title)
	}

	return s
}
