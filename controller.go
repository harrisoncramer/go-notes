package main

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) getController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.view {
	case addEntryView:
		return m.createEntryController(msg)
	case editEntryView:
		return m.editEntryController(msg)
	case settingsView:
		return m.settingsController(msg)
	case settingsEditorView:
		return m.editSettingsController(msg)
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
		switch data := msg.data.(type) {
		case []Entry:
			for _, entry := range data {
				m.viewData.entries = append(m.viewData.entries, entry)
			}
		case []Setting:
			for _, setting := range data {
				m.viewData.entries = append(m.viewData.entries, Entry{Title: setting.Key, Content: setting.Value})
			}
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m.handleCtrlC()
		default:
			return m.getController(msg)
		}
	}
	return m.getController(msg)
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
			choice := m.viewData.entries[m.cursor.idx]
			switch choice {
			case addEntryChoice:
				m.textInput.Placeholder = "Learning about Go"
				return m, m.changeView(addEntryView)
			case editEntryChoice:
				return m, m.changeView(editEntryView)
			case settingsEntryChoice:
				return m, m.changeView(settingsView)
			}
		}
	}

	return m, nil
}

/* Responsible for creating new entries */
func (m model) createEntryController(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			title := m.textInput.Value()
			if title == "" {
				return m, tea.Quit
			}
			id, err := m.createEntry(title, "")
			if err != nil {
				m.err = err
				return m, nil
			}
			m.currentEntryId = id
			m.textInput.SetValue("")
			return m.editEntry()
		case tea.KeyEsc:
			return m, m.changeView(mainView)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

/* Responsible for all update and delete operations on existing entries */
func (m model) editEntryController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, m.changeView(mainView)
		case tea.KeyEnter:
			choice := m.viewData.entries[m.cursor.idx]
			m.currentEntryId = choice.Id
			return m.editEntry()
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

func (m model) settingsController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, m.changeView(mainView)
		case tea.KeyEnter:
			choice := m.viewData.entries[m.cursor.idx]
			value := choice.Content
			m.textInput.SetValue(value)
			m.textInput.Placeholder = value
			return m, m.changeView(settingsEditorView)
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

func (m model) editSettingsController(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			key := m.viewData.entries[m.cursor.idx].Title
			value := m.textInput.Value()
			if value == "" {
				return m, m.changeView(mainView)
			}
			_, err := m.updateSetting(key, value)
			if err != nil {
				m.err = err
				return m, nil
			}
			return m, m.changeView(mainView)
		case tea.KeyEsc:
			m.textInput.SetValue("")
			return m, m.changeView(mainView)
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

/****************/
/*** Helpers 🤝 */
/****************/

func (m model) handleCtrlC() (model, tea.Cmd) {
	return m, tea.Quit
}

func (m *model) handleUpKey() {
	if m.cursor.idx > 0 {
		m.cursor.idx--
	}
}

func (m *model) handleDownKey() {
	if m.cursor.idx < len(m.viewData.entries)-1 {
		m.cursor.idx++
	}
}

func (m *model) changeView(view string) tea.Cmd {
	m.view = view
	return m.loadData(view)
}

type dataLoader interface{}
type dataLoaded struct{ data dataLoader }

/* Loads the data required for the view and returns it in the dataLoaded message */
func (m *model) loadData(view string) tea.Cmd {
	return func() tea.Msg {
		switch view {
		case addEntryView, settingsEditorView:
			return nil
		case mainView:
			return dataLoaded{data: initialChoices}
		case settingsView:
			settings, err := m.readAllSettings()
			if err != nil {
				m.err = err
			}
			return dataLoaded{data: settings}
		case editEntryView:
			entries, err := m.readAllEntries()
			if err != nil {
				m.err = err
			}
			return dataLoaded{data: entries}
		}

		err := errors.New("Invalid data load")
		return errMsg{err: err}
	}
}
