package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/user"
)

func (m *model) initDB() {
	user, err := user.Current()
	if err != nil {
		m.err = err
		return
	}

	if len(os.Args) < 2 {
		m.err = errors.New("Must provide a database name!")
		return
	}

	m.dbName = os.Args[1]

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", user.HomeDir, m.dbName))
	if err != nil {
		m.err = err
		return
	}

	m.conn = db
	_, err = db.Exec(`
      CREATE TABLE IF NOT EXISTS entries (
        id integer primary key autoincrement,
        title TEXT,
        content TEXT
      );
    `)

	m.err = err

}

/* Reads all entries from the SQL database into the model */
func (m *model) readAllData() {
	if m.conn == nil {
		m.err = errors.New("DB Connection not established!")
		return
	}

	rows, err := m.conn.Query("SELECT id, title FROM entries")
	if err != nil {
		m.err = err
		return
	}

	defer rows.Close()

	var results []Entry
	for rows.Next() {
		var id int64
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			m.err = err
			return
		}
		results = append(results, Entry{Title: title, Id: id})
	}

	m.entries = results
	m.choices = []Choice{}
	for _, entry := range m.entries {
		m.choices = append(m.choices, Choice{Text: entry.Title, Id: entry.Id})
	}
}

/* Adds a record to the SQL database */
func (m *model) createEntry(data Entry) (int64, error) {
	if m.conn == nil {
		return 0, errors.New("DB Connection not established!")
	}

	result, err := m.conn.Exec("INSERT INTO entries(title, content) VALUES(?, ?)", data.Title, data.Content)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return id, nil
}

func (m *model) updateEntryText(id int64, text string) error {
	if m.conn == nil {
		return errors.New("DB Connection not established!")
	}

	_, err := m.conn.Exec("UPDATE entries SET content = ? WHERE id = ?", text, id)
	return err
}

func (m *model) readEntryById(id int64) (*Entry, error) {
	if m.conn == nil {
		return nil, errors.New("DB Connection not established!")
	}

	var entry Entry
	err := m.conn.QueryRow("SELECT id, title, content FROM entries WHERE id = ?", id).Scan(&entry.Id, &entry.Title, &entry.Content)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (m *model) renameEntry(id int64, title string) (*Entry, error) {
	if m.conn == nil {
		return nil, errors.New("DB Connection not established!")
	}

	_, err := m.conn.Exec("UPDATE entries SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return nil, err
	}

	var entry Entry
	err = m.conn.QueryRow("SELECT id, title, content FROM entries WHERE id = ?", id).Scan(&entry.Id, &entry.Title, &entry.Content)
	if err != nil {
		return nil, err
	}

	m.readAllData()

	return &entry, nil
}
