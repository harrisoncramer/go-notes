package main

import "fmt"

const addEntryView = "Add Entry"
const editEntryView = "Edit Entry"
const mainView = "Main"
const settingsView = "Settings"
const settingsEditorView = "Settings/Edit"

/* The view function is responsible for rendering different screens */
func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	s := fmt.Sprintf("%s 📓\n", m.dbName)
	s += fmt.Sprintf("%s\n\n", m.view)

	switch m.view {
	case mainView:
		s += m.choiceRenderer()
	case addEntryView:
		s += m.textInputRenderer()
	case editEntryView:
		s += m.choiceRenderer()
	case settingsView:
		s += m.choiceRenderer()
	case settingsEditorView:
		s += m.textInputRenderer()
	}

	s += "\nPress <C-c> to quit\n"
	return s
}

func (m model) textInputRenderer() string {
	s := ""
	s += fmt.Sprintf(
		"%s:\n\n",
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
