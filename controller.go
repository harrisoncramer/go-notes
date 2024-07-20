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

/* Possible prompt types for the user */
var addEntryPrompt = Prompt{Text: "Entry Name", Id: ADD_ENTRY}
var chooseEntryToReadPrompt = Prompt{Text: fmt.Sprintf("Which entry do you want to edit?\n\n"), Id: CHOOSE_EDIT}
var noEntriesFoundPrompt = Prompt{Text: "No entries found!\n\n", Id: NO_ENTRIES}
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
	case ADD_ENTRY:
		return fmt.Sprintf(
			"%s:\n\n%s\n\n%s",
			m.prompt.Text,
			m.textInput.View(),
			"Press <C-c> to quit.\nPress <esc> to go back",
		) + "\n"
	case CHOOSE_EDIT:
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
		return "What would you like to do with this entry?\n\nContinue editing: 'e'\nSave Changes: 'w'\nQuit: 'q'\n"
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
				switch choice.Text {
				case "Add Entry":
					m.prompt = addEntryPrompt
				case "Edit Entries":
					m.readAllData()
					if len(m.choices) > 0 {
						m.prompt = chooseEntryToReadPrompt
					} else {
						m.prompt = noEntriesFoundPrompt
					}
				}
				return m, nil
			}
		}
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
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "e":
				return m.editEntry()
			case "w":
				return m, m.persistEntry()
			case "ctrl+c", "q":
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
	m.cursor.idx = 0
}
