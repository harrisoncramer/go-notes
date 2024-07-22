package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/user"
)

/* Initializes the database and all tables */
func (m *model) initDB() error {
	user, err := user.Current()
	if err != nil {
		return err
	}

	if len(os.Args) < 2 {
		return errors.New("Must provide a database name!")
	}

	m.dbName = os.Args[1]

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", user.HomeDir, m.dbName))
	if err != nil {
		return err
	}

	_, err = db.Exec(`
      CREATE TABLE IF NOT EXISTS entries (
        id integer primary key autoincrement,
        title TEXT,
        content TEXT
      );
    `)

	_, err = db.Exec(`
      CREATE TABLE IF NOT EXISTS settings (
          key VARCHAR(255) UNIQUE,
          value VARCHAR(255)
      );
    `)

	_, err = db.Exec(`
    INSERT OR IGNORE INTO settings (key, value) VALUES 
        ('backup_url', 'http://example.com/backup'),
        ('another_setting', 'some_value');
  `)

	m.conn = db
	return err
}

/* Reads all entries from the SQL database */
func (m *model) readAllEntries() ([]Entry, error) {
	rows, err := m.conn.Query("SELECT id, title FROM entries")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []Entry
	for rows.Next() {
		var id int64
		var title string
		if err := rows.Scan(&id, &title); err != nil {
			return nil, err
		}
		results = append(results, Entry{Title: title, Id: id})
	}

	return results, nil
}

/* Adds a record to the SQL database */
func (m *model) createEntry(title string, content string) (int64, error) {
	result, err := m.conn.Exec("INSERT INTO entries(title, content) VALUES(?, ?)", title, content)
	if err != nil {
		return 0, err
	}

	id, _ := result.LastInsertId()
	return id, nil
}

/* Updates an entry's (by ID) text */
func (m *model) updateEntryText(id int64, text string) error {
	_, err := m.conn.Exec("UPDATE entries SET content = ? WHERE id = ?", text, id)
	return err
}

/* Reads an entry from the database by it's ID */
func (m model) readEntryById(id int64) (*Entry, error) {
	var entry Entry
	err := m.conn.QueryRow("SELECT id, title, content FROM entries WHERE id = ?", id).Scan(&entry.Id, &entry.Title, &entry.Content)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

/* Renames an entry's (by ID) title */
func (m *model) renameEntry(id int64, title string) (*Entry, error) {
	_, err := m.conn.Exec("UPDATE entries SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return nil, err
	}

	var entry Entry
	err = m.conn.QueryRow("SELECT id, title, content FROM entries WHERE id = ?", id).Scan(&entry.Id, &entry.Title, &entry.Content)
	if err != nil {
		return nil, err
	}

	return &entry, nil
}

func (m *model) readAllSettings() ([]Setting, error) {
	rows, err := m.conn.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []Setting
	for rows.Next() {
		var key string
		var value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		results = append(results, Setting{Key: key, Value: value})
	}

	return results, nil
}

func (m *model) updateSetting(key string, value string) (*Setting, error) {
	_, err := m.conn.Exec("UPDATE settings SET value = ? WHERE key = ?", value, key)
	if err != nil {
		return nil, err
	}

	var setting Setting
	err = m.conn.QueryRow("SELECT key, value FROM settings WHERE key = ?", key).Scan(&setting.Key, &setting.Value)
	if err != nil {
		return nil, err
	}

	return &setting, nil
}
