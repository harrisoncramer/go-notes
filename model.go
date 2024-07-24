package main

import (
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

type State struct {
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

type Model struct {
	cursor               Cursor
	err                  error
	state                State
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

func (m Model) Init() tea.Cmd {
	return nil
}

var db = Orm{}

var addEntryChoice = Entry{Title: "Add Entry"}
var editEntryChoice = Entry{Title: "Edit Entry"}
var settingsEntryChoice = Entry{Title: "Settings"}
var initialChoices = []Entry{addEntryChoice, editEntryChoice, settingsEntryChoice}

func initialModel() Model {
	m := Model{
		state: State{
			entries: initialChoices,
		},
		view: mainView,
	}

	err := db.Init(m)
	if err != nil {
		m.err = err
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
