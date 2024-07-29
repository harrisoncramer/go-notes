package main

import (
	"fmt"
	"os"

	"github.com/harrisoncramer/go-notes/internal/db"
	_ "github.com/mattn/go-sqlite3"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	database, err := db.InitSqliteDb()
	if err != nil {
		fmt.Printf("Could not start database: %v", err)
		os.Exit(1)
	}

	p := tea.NewProgram(initialModel(database))
	if _, err := p.Run(); err != nil {
		fmt.Printf("There's been an error: %v", err)
		os.Exit(1)
	}
}
