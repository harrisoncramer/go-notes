package main

import (
	"database/sql"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Cursor struct {
	idx int
}

type model struct {
	choices              []Choice
	cursor               Cursor
	conn                 *sql.DB
	err                  error
	entries              []Entry
	prompt               Prompt
	textInput            textinput.Model
	dbName               string
	currentEntryId       int64
	currentEntryFilePath string
}

func (m model) Init() tea.Cmd {
	return nil
}

var addEntryChoice = Choice{"Add Entry", 0}
var editEntryChoice = Choice{"Edit Entries", 1}
var renameEntryChoice = Choice{"Rename Entry", 2}
var initialChoices = []Choice{addEntryChoice, editEntryChoice, renameEntryChoice}

func initialModel() model {

	m := model{
		conn:    nil,
		choices: initialChoices,
		err:     nil,
		entries: []Entry{},
	}

	m.initDB()
	m.prompt = Prompt{Text: m.dbName + "\n\n", Id: MAIN}

	/* Text input component */
	ti := textinput.New()
	ti.Placeholder = "Learning about Go"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m.textInput = ti

	return m
}
