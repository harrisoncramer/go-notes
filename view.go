package main

import (
	"fmt"
)

type View string

const (
	addEntryView       View = "Entries - Add"
	editEntryView           = "Entries - Edit"
	mainView                = "Main"
	settingsView            = "Settings"
	settingsEditorView      = "Settings -> Edit"
)

/* The view function is responsible for rendering different screens */
func (m Model) View() string {
	if m.err != nil {
		return errStyle.Render(m.err.Error())
	}

	s := titleStyle.Render(fmt.Sprintf("%s 📓", db.name))

	if m.view != mainView {
		s += "\n"
		s += navStyle.Render(fmt.Sprintf("%s", m.view))
	}

	s += "\n\n"

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

func (m Model) textInputRenderer() string {
	s := ""
	s += fmt.Sprintf(
		"%s:\n\n",
		m.textInput.View())
	return s
}

func (m Model) choiceRenderer() string {
	s := ""
	if len(m.state.entries) == 0 {
		s += "No entries found!\n"
		return s
	}

	for i, choice := range m.state.entries {
		prefix := " "
		if m.cursor.idx == i {
			prefix = ">"
		}
		s += fmt.Sprintf("%s %s\n", prefix, choice.Title)
	}

	return s
}
