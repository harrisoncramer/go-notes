package main

import (
	"io"
	"os"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"
)

type editorFinishedMsg struct{ err error }

func (m *model) openEditor(entry *Entry) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	tmpfile, err := os.CreateTemp("", "entry.txt")
	if err != nil {
		m.err = err
		return nil
	}

	name := tmpfile.Name()
	m.currentEntryFilePath = name

	if _, err := tmpfile.Write([]byte(entry.Content)); err != nil {
		m.err = err
		return nil
	}

	if err := tmpfile.Close(); err != nil {
		return nil
	}

	c := exec.Command(editor, name)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err}
	})
}

func (m *model) editEntry() (tea.Model, tea.Cmd) {
	currentEntryId := m.currentEntryId
	entry, err := m.readEntryById(currentEntryId)
	if err != nil {
		m.err = err
		return m, nil
	}

	return m, m.openEditor(entry)
}

func (m *model) persistEntry() tea.Cmd {
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

	err = m.updateEntryText(m.currentEntryId, string(content))
	if err != nil {
		m.err = err
		return nil
	}

	return tea.Quit
}
