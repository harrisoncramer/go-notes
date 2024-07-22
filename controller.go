package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

type dataLoaded struct{ data []Entry }

func (m model) getView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.view {
	case addEntryView:
		return m.createEntryController(msg)
	case editEntryView:
		return m.editEntryViewController(msg)
	default:
		return m.mainController(msg)
	}
}

/* The update function is responsible for getting a message from the Init function, and updating state in the model */
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editorFinishedMsg:
		return m, m.persistEntry()
	case errMsg:
		m.err = msg
		return m, nil
	case dataLoaded:
		m.cursor.idx = 0
		m.viewData = ViewData{}
		for _, entry := range msg.data {
			m.viewData.choices = append(m.viewData.choices, Choice{Text: entry.Title, Id: entry.Id})
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m.handleCtrlC()
		default:
			return m.getView(msg)
		}
	}
	return m.getView(msg)
}

/* The main controller is the root of the application */
func (m model) mainController(msg tea.Msg) (tea.Model, tea.Cmd) {
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
				return m, m.changeView(addEntryView)
			case "Edit Entry":
				return m, m.changeView(editEntryView)
			case "Rename Entry":
				return m, m.changeView(addEntryView)
			}
		}
	}

	return m, nil
}

func (m model) createEntryController(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
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
			return m.editEntry()
		case tea.KeyEsc:
			m.changeView(mainView)
			m.viewData.choices = initialChoices
			return m, nil
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

func (m model) editEntryViewController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.changeView(mainView)
			m.viewData.choices = initialChoices
		}
		switch msg.String() {
		case "up", "k":
			m.handleUpKey()
		case "down", "j":
			m.handleDownKey()
		}
	}
	return m, nil
}

/**** Common Helpers ****/
func (m model) handleCtrlC() (model, tea.Cmd) {
	return m, tea.Quit
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

func (m *model) changeView(view string) tea.Cmd {
	m.view = view
	return m.loadData(view)
}

func (m *model) loadData(view string) tea.Cmd {
	return func() tea.Msg {
		entries, err := m.readAllEntries()
		if err != nil {
			m.err = err
		}
		return dataLoaded{data: entries}
	}
}
