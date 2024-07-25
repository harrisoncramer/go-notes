package main

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/harrisoncramer/go-notes/internal/db"
)

/* The update function is responsible for getting a message from the Init function, and updating state in the model */
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case editorFinishedMsg:
		return m, m.persistEntry()
	case errMsg:
		m.err = msg
		return m, nil
	case dataLoaded:
		m.cursor.idx = 0
		m.state = State{}
		switch data := msg.data.(type) {
		case []db.Entry:
			for _, entry := range data {
				m.state.entries = append(m.state.entries, entry)
			}
		case []db.Setting:
			for _, setting := range data {
				m.state.entries = append(m.state.entries, db.Entry{Title: setting.Key, Content: setting.Value})
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
func (m Model) mainController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.handleUpKey()
		case "down", "j":
			m.handleDownKey()
		case "enter":
			choice := m.state.entries[m.cursor.idx]
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
func (m Model) createEntryController(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			title := m.textInput.Value()
			if title == "" {
				return m, Quitter
			}
			entry, err := m.db.CreateEntry("title", "content")
			if err != nil {
				m.err = err
				return m, nil
			}
			m.currentEntryId = entry.Id
			m.textInput.SetValue("")
			return m.editEntry()
		case tea.KeyEsc:
			return m, m.changeView(mainView)
		case tea.KeyCtrlC:
			return m, Quitter
		}
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

/* Responsible for all update and delete operations on existing entries */
func (m Model) editEntryController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, m.changeView(mainView)
		case tea.KeyEnter:
			choice := m.state.entries[m.cursor.idx]
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

func (m Model) settingsController(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, m.changeView(mainView)
		case tea.KeyEnter:
			choice := m.state.entries[m.cursor.idx]
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

func (m Model) editSettingsController(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			key := m.state.entries[m.cursor.idx].Title
			value := m.textInput.Value()
			if value == "" {
				return m, m.changeView(mainView)
			}
			_, err := m.db.UpdateSetting(key, value)
			if err != nil {
				m.err = err
				return m, nil
			}
			return m, m.changeView(mainView)
		case tea.KeyEsc:
			m.textInput.SetValue("")
			return m, m.changeView(mainView)
		case tea.KeyCtrlC:
			return m, Quitter
		}
		m.textInput, cmd = m.textInput.Update(msg)
	}
	return m, cmd
}

/****************/
/*** Helpers 🤝 */
/****************/

func (m Model) handleCtrlC() (Model, tea.Cmd) {
	return m, Quitter
}

func (m *Model) handleUpKey() {
	if m.cursor.idx > 0 {
		m.cursor.idx--
	}
}

func (m *Model) handleDownKey() {
	if m.cursor.idx < len(m.state.entries)-1 {
		m.cursor.idx++
	}
}

func (m *Model) changeView(view View) tea.Cmd {
	m.view = view
	return m.loadData(view)
}

type dataLoader interface{}
type dataLoaded struct{ data dataLoader }

/* Loads the data required for the view and returns it in the dataLoaded message */
func (m *Model) loadData(view View) tea.Cmd {
	return func() tea.Msg {
		switch view {
		case addEntryView, settingsEditorView:
			return nil
		case mainView:
			return dataLoaded{data: initialChoices}
		case settingsView:
			settings, err := m.db.ReadAllSettings()
			if err != nil {
				m.err = err
			}
			return dataLoaded{data: settings}
		case editEntryView:
			entries, err := m.db.ReadAllEntries()
			if err != nil {
				m.err = err
			}
			return dataLoaded{data: entries}
		}

		err := errors.New("Invalid data load")
		return errMsg{err: err}
	}
}

func Quitter() tea.Msg {
	return tea.Quit()
}

func (m Model) getController(msg tea.Msg) (tea.Model, tea.Cmd) {
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
