package main

import (
	"database/sql"
	"errors"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Cursor struct {
	idx int
}

type KeyVal struct {
	Key   string
	Value string
}

type ViewData struct {
	entries []Entry
	keyVals []KeyVal
}

type Entry struct {
	Id      int64  `db:"id"`
	Title   string `db:"title"`
	Content string `db:"string"`
}

type Setting struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

type model struct {
	cursor               Cursor
	conn                 *sql.DB
	err                  error
	viewData             ViewData
	textInput            textinput.Model
	dbName               string
	currentEntryId       int64
	currentEntryFilePath string
	goHome               bool
	view                 View
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

var addEntryChoice = Entry{Title: "Add Entry"}
var editEntryChoice = Entry{Title: "Edit Entry"}
var settingsEntryChoice = Entry{Title: "Settings"}
var initialChoices = []Entry{addEntryChoice, editEntryChoice, settingsEntryChoice}

func initialModel() model {

	m := model{
		viewData: ViewData{
			entries: initialChoices,
		},
		err:  nil,
		view: mainView,
	}

	/* Text input component */
	ti := textinput.New()
	ti.Placeholder = "Learning about Go"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 60
	m.textInput = ti

	err := m.initDB()
	if err != nil {
		m.err = err
	}

	if m.conn == nil && m.err == nil {
		m.err = errors.New("DB Connection not established!")
	}

	return m
}
