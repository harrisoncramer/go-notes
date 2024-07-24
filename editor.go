package main

import (
	"io"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }
type fileSavedMsg bool

/* Opens an editor, which upon closure, will return the "editorFinishedMsg" message */
func (m *Model) editEntry() (tea.Model, tea.Cmd) {
	currentEntryId := m.currentEntryId
	entry, err := db.readEntry(currentEntryId)
	if err != nil {
		m.err = err
		return m, nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	tmpfile, err := os.CreateTemp("", "entry.txt")
	if err != nil {
		m.err = err
		return m, nil
	}

	name := tmpfile.Name()
	m.currentEntryFilePath = name

	if _, err := tmpfile.Write([]byte(entry.Content)); err != nil {
		m.err = err
		return m, nil
	}

	if err := tmpfile.Close(); err != nil {
		return m, nil
	}

	c := exec.Command(editor, name)
	return m, tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

/* Gets the contents of the file at the current temporary file location and saves it to the database */
func (m *Model) persistEntry() tea.Cmd {
	file, err := os.Open(m.currentEntryFilePath)
	if err != nil {
		m.err = err
		return nil
	}

	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		m.err = err
		return nil
	}

	err = db.updateEntryText(m.currentEntryId, string(content))
	if err != nil {
		m.err = err
		return nil
	}

	return func() tea.Msg {
		return fileSavedMsg(true)
	}
}
