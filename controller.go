package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type Entry struct {
	Id      int64  `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
}

type Prompt struct {
	Text string
	Id   string
}

type Choice struct {
	Text string
	Id   int64
}

const MAIN = "main"
const ADD_ENTRY = "add_entry"
const CHOOSE_EDIT = "choose_edit"
const NO_ENTRIES = "no_entries"
const EDITOR = "editor"
const CHOOSE_RENAME = "choose_rename"
const RENAME_ENTRY = "rename_entry"

/* Possible prompt types for the user */
var addEntryPrompt = Prompt{Text: "Entry Name", Id: ADD_ENTRY}
var chooseEntryToReadPrompt = Prompt{Text: fmt.Sprintf("Which entry do you want to edit?\n\n"), Id: CHOOSE_EDIT}
var noEntriesFoundPrompt = Prompt{Text: "No entries found!\n\n", Id: NO_ENTRIES}
var chooseRenamePrompt = Prompt{Text: fmt.Sprintf("Which entry do you want to rename?\n\n"), Id: CHOOSE_RENAME}
var renameEntryPrompt = Prompt{Text: fmt.Sprintf("New Entry Name"), Id: RENAME_ENTRY}
var editEntryPrompt = Prompt{Text: "", Id: EDITOR}

/* The view function is responsible for rendering different screens depending on the Prompt ID */
func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}

	switch m.prompt.Id {
	case MAIN:
		for i, choice := range m.choices {
			prefix := " "
			if m.cursor.idx == i {
				prefix = ">"
			}
			m.prompt.Text += fmt.Sprintf("%s %s\n", prefix, choice.Text)
		}
		m.prompt.Text += "\nPress <C-c> to quit\n"
		return m.prompt.Text
	case ADD_ENTRY, RENAME_ENTRY:
		return fmt.Sprintf(
			"%s:\n\n%s\n\n%s",
			m.prompt.Text,
			m.textInput.View(),
			"Press <C-c> to quit.\nPress <esc> to go back",
		) + "\n"
	case CHOOSE_EDIT, CHOOSE_RENAME:
		for i, choice := range m.choices {
			prefix := " "
			if m.cursor.idx == i {
				prefix = ">"
			}
			m.prompt.Text += fmt.Sprintf("%s %s\n", prefix, choice.Text)
		}
		m.prompt.Text += "\nPress <C-c> to quit\n"
		m.prompt.Text += "Press <esc> to go back\n"
		return m.prompt.Text
	case EDITOR:
		m.prompt.Text = "Press 'w' to save this entry, or 'e' to continue editing\n\n"
		m.prompt.Text += "Press <C-c> to quit (no save)\n"
		m.prompt.Text += "Press <esc> to go back (no save)\n"
		return m.prompt.Text
	case NO_ENTRIES:
		return m.prompt.Text
	}

	return "Invalid prompt ID chosen!"
}

/* The update function is responsible for updating state in the model and choosing a prompt */
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.prompt.Id {
	case MAIN:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor.idx > 0 {
					m.cursor.idx--
				}
			case "down", "j":
				if m.cursor.idx < len(m.choices)-1 {
					m.cursor.idx++
				}
			case "enter", " ":
				choice := m.choices[m.cursor.idx]
				switch choice.Id {
				case addEntryChoice.Id:
					m.prompt = addEntryPrompt
				case editEntryChoice.Id:
					m.cursor.idx = 0
					m.readAllData()
					if len(m.choices) > 0 {
						m.prompt = chooseEntryToReadPrompt
					} else {
						m.prompt = noEntriesFoundPrompt
					}
				case renameEntryChoice.Id:
					m.cursor.idx = 0
					m.readAllData()
					if len(m.choices) > 0 {
						m.prompt = chooseRenamePrompt
					} else {
						m.prompt = noEntriesFoundPrompt
					}
				}
				return m, nil
			}
		}
	case RENAME_ENTRY:
		var cmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				content := m.textInput.Value()
				if content == "" {
					return m, tea.Quit
				}
				m.renameEntry(m.currentEntryId, content)
				m.textInput.SetValue("")
				m.returnHome()
			case tea.KeyEsc:
				m.returnHome()
			case tea.KeyCtrlC:
				return m, tea.Quit
			}

			m.textInput, cmd = m.textInput.Update(msg)
		}
		return m, cmd
	case ADD_ENTRY:
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

				m.cursor.idx = 0
				m.currentEntryId = id
				m.prompt = editEntryPrompt
				m.textInput.SetValue("")
				return m.editEntry()
			case tea.KeyEsc:
				m.returnHome()
				return m, nil
			case tea.KeyCtrlC:
				return m, tea.Quit
			}

			m.textInput, cmd = m.textInput.Update(msg)
			return m, cmd
		}
	case CHOOSE_RENAME:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyEsc {
				m.returnHome()
				return m, nil
			}
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor.idx > 0 {
					m.cursor.idx--
				}
			case "down", "j":
				if m.cursor.idx < len(m.choices)-1 {
					m.cursor.idx++
				}
			case "enter", " ":
				choice := m.choices[m.cursor.idx]
				m.currentEntryId = choice.Id
				m.textInput.Placeholder = choice.Text
				m.prompt = renameEntryPrompt
				return m, nil
			}
		}
	case CHOOSE_EDIT:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyEsc {
				m.returnHome()
				return m, nil
			}
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.cursor.idx > 0 {
					m.cursor.idx--
				}
			case "down", "j":
				if m.cursor.idx < len(m.choices)-1 {
					m.cursor.idx++
				}
			case "enter", " ":
				choice := m.choices[m.cursor.idx]
				if m.prompt.Id == CHOOSE_EDIT {
					m.currentEntryId = choice.Id
					m.prompt = editEntryPrompt
					return m.editEntry()
				}

				return m, nil
			}
		}
	case NO_ENTRIES:
		return m, tea.Quit
	case EDITOR:
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.Type == tea.KeyEsc {
				m.returnHome()
				return m, nil
			}
		}
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "e":
				return m.editEntry()
			case "w":
				m.persistEntry()
				m.returnHome()
				return m, nil
			case "ctrl+c":
				return m, tea.Quit
			}
		case editorFinishedMsg:
			if msg.err != nil {
				m.err = msg.err
				return m, tea.Quit
			}
		}
		return m, nil
	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m *model) returnHome() {
	m.prompt = Prompt{Text: m.dbName + "\n\n", Id: MAIN}
	m.choices = initialChoices
	m.currentEntryId = -1
	m.cursor.idx = 0
}
