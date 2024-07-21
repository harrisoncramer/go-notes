package main

import (
	"database/sql"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Cursor struct {
	idx int
}

type ViewData struct {
	choices []Choice
}

type Entry struct {
	Id      int64  `db:"id"`
	Title   string `db:"title"`
	Content string `db:"content"`
}

type model struct {
	cursor               Cursor
	conn                 *sql.DB
	err                  error
	entries              []Entry
	viewData             ViewData
	textInput            textinput.Model
	dbName               string
	currentEntryId       int64
	currentEntryFilePath string
	goHome               bool
	view                 interface{}
}

type errMsg struct{ err error }

func (e errMsg) Error() string {
	return e.err.Error()
}

func (m model) Init() tea.Cmd {
	return nil
}

type Choice struct {
	Text string
	Id   int64
}

var addEntryChoice = Choice{Text: "Add Entry"}
var editEntryChoice = Choice{Text: "Edit Entry"}
var renameEntryChoice = Choice{Text: "Rename Entry"}
var initialChoices = []Choice{addEntryChoice, editEntryChoice, renameEntryChoice}

func initialModel() model {

	m := model{
		viewData: ViewData{
			choices: initialChoices,
		},
		err:     nil,
		entries: []Entry{},
		view:    mainView,
	}

	/* Text input component */
	ti := textinput.New()
	ti.Placeholder = "Learning about Go"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20
	m.textInput = ti

	err := m.initDB()
	if err != nil {
		m.err = err
	}

	if m.conn == nil {
		os.Exit(1)
	}

	return m
}
