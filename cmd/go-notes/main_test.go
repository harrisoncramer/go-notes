package main

import (
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/harrisoncramer/go-notes/internal/db"
)

type TestDb struct{}

func (t TestDb) ReadAllEntries() ([]db.Entry, error) {
	return []db.Entry{
		{Title: "Fake title #1", Content: "Fake content #1", Id: 1},
		{Title: "Fake title #2", Content: "Fake content #2", Id: 2},
		{Title: "Fake title #3", Content: "Fake content #3", Id: 3},
	}, nil
}
func (t TestDb) CreateEntry(title string, content string) (db.Entry, error) {
	return db.Entry{}, nil
}
func (t TestDb) ReadEntry(id int64) (db.Entry, error) {
	return db.Entry{}, nil
}
func (t TestDb) UpdateEntryText(id int64, content string) (db.Entry, error) {
	return db.Entry{}, nil
}
func (t TestDb) RenameEntry(id int64, title string) (db.Entry, error) {
	return db.Entry{}, nil
}
func (t TestDb) ReadAllSettings() ([]db.Setting, error) {
	return []db.Setting{}, nil
}
func (t TestDb) UpdateSetting(key string, value string) (db.Setting, error) {
	return db.Setting{}, nil
}
func (t TestDb) GetName() string {
	return "TestDb"
}

/* Tests that the inital view renders correctly */
func TestMain(t *testing.T) {
	tm := commonSetup(t)
	tm.Quit()
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

/* Tests that you can navigate downwards */
func TestDown(t *testing.T) {
	tm := commonSetup(t)
	tm.Send(k(tea.KeyDown))
	tm.Send(k(tea.KeyDown))
	tm.Quit()
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

/* Tests that you can navigate downwards */
func TestUp(t *testing.T) {
	tm := commonSetup(t)
	tm.Send(k(tea.KeyDown))
	tm.Send(k(tea.KeyDown))
	tm.Send(k(tea.KeyUp))
	tm.Quit()
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

/* Tests that you can type an entry */
func TestTypeNewEntry(t *testing.T) {
	tm := commonSetup(t)
	tm.Send(k(tea.KeyEnter))
	tm.Type("Some new entry...")
	tm.Quit()
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

/* Tests that you can open and navigate around to edit existing entries */
func TestEditExistingEntries(t *testing.T) {
	tm := commonSetup(t)
	tm.Send(k(tea.KeyDown))
	tm.Send(k(tea.KeyEnter))
	tm.Send(k(tea.KeyEnter)) // TODO: Why does this require two keypresses?
	tm.Quit()
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

/* Helpers and variables 🤝 */
func commonSetup(t *testing.T) *teatest.TestModel {
	testDb := TestDb{}
	m := initialModel(testDb)
	tm := teatest.NewTestModel(
		t,
		m,
		teatest.WithInitialTermSize(300, 100),
	)

	return tm
}

func k(key tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg(tea.Key{
		Type: key,
	})
}
