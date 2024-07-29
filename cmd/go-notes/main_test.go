package main

import (
	"io"
	"testing"

	"github.com/charmbracelet/x/exp/teatest"
	"github.com/harrisoncramer/go-notes/internal/db"
)

func TestMain(t *testing.T) {
	testDb := TestDb{}
	m := initialModel(testDb)
	tm := teatest.NewTestModel(
		t,
		m,
		teatest.WithInitialTermSize(300, 100),
	)

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
