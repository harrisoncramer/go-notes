package main

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/harrisoncramer/go-notes/internal/db"
)

type Cursor struct {
	idx int
}

type KeyVal struct {
	Key   string
	Value string
}

type State struct {
	entries []db.Entry
	keyVals []KeyVal
}

type Model struct {
	cursor               Cursor
	err                  error
	state                State
	textInput            textinput.Model
	currentEntryId       int64
	currentEntryFilePath string
	goHome               bool
	view                 View
	db                   db.Database
}

type errMsg struct{ err error }

func (e errMsg) Error() string {
	return e.err.Error()
}

func (m Model) Init() tea.Cmd {
	return nil
}

var addEntryChoice = db.Entry{Title: "Add Entry"}
var editEntryChoice = db.Entry{Title: "Edit Entry"}
var settingsEntryChoice = db.Entry{Title: "Settings"}
var initialChoices = []db.Entry{addEntryChoice, editEntryChoice, settingsEntryChoice}

func initialModel() Model {

	database, err := db.InitSqliteDb()
	if err != nil {
		return Model{err: err}
	}

	m := Model{
		state: State{
			entries: initialChoices,
		},
		view: mainView,
		db:   database,
	}

	/* Text input component */
	ti := textinput.New()
	ti.Placeholder = "Learning about Go"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 60
	m.textInput = ti

	return m
}
