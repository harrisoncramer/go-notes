package main

import (
	"io"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/x/exp/teatest"
	"github.com/harrisoncramer/go-notes/internal/db"
)

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
	tm.Type("j")
	tm.Type("j")
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
	tm.Type("j")
	tm.Type("j")
	tm.Type("k")
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
	enterKey := tea.KeyMsg(tea.Key{
		Type: tea.KeyEnter,
	})

	tm.Send(enterKey)
	tm.Type("Some new entry...")
	tm.Quit()
	out, err := io.ReadAll(tm.FinalOutput(t))
	if err != nil {
		t.Error(err)
	}
	teatest.RequireEqualOutput(t, out)
}

type TestDb struct{}

func (t TestDb) ReadAllEntries() ([]db.Entry, error) {
	return []db.Entry{}, nil
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
