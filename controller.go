package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

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
		default:
			switch m.view {
			case addEntryView:
				return m.createEntryController(msg)
			default:
				return m.mainController(msg)
			}
		}
	}

	return m, nil
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
				m.changeView(addEntryView)
			case "Edit Entry":
				m.changeView(addEntryView)
			case "Rename Entry":
				m.changeView(addEntryView)
			}
		}
	}

	return m, nil
}

func (m *model) createEntryController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			content := m.textInput.Value()
			if content == "" {
				return m, tea.Quit
			}
			id, err := m.createEntry(Entry{Title: content, Content: ""})
			if err != nil {
				m.err = err
				return m, nil
			}
			m.currentEntryId = id
			m.textInput.SetValue("")
			return m.editEntry() // TODO: Fix
		case tea.KeyEsc:
			m.changeView(mainView)
			m.viewData.choices = initialChoices
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}

		_, cmd := m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

/**** Common Helpers ****/
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

func (m *model) changeView(view string) {
	m.cursor.idx = 0
	m.viewData = ViewData{}
	m.view = view
}
