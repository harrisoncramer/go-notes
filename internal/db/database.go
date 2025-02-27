package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/user"
)

type Entry struct {
	Id      int64  `db:"id"`
	Title   string `db:"title"`
	Content string `db:"string"`
}

type Setting struct {
	Key   string `db:"key"`
	Value string `db:"value"`
}

type Database interface {
	ReadAllEntries() ([]Entry, error)
	CreateEntry(title string, content string) (Entry, error)
	ReadEntry(id int64) (Entry, error)
	UpdateEntryText(id int64, content string) (Entry, error)
	RenameEntry(id int64, title string) (Entry, error)
	ReadAllSettings() ([]Setting, error)
	UpdateSetting(key string, value string) (Setting, error)
	GetName() string
}

type SqlLite struct {
	conn *sql.DB
	name string
}

/* Initializes the database and all tables */
func InitSqliteDb() (Database, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	if len(os.Args) < 2 {
		return nil, errors.New("Must provide a database name!")
	}

	db := SqlLite{}
	db.name = os.Args[1]
	conn, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s.db", user.HomeDir, db.name))

	if err != nil {
		return nil, err
	}

	db.conn = conn

	_, err = db.conn.Exec(`
      CREATE TABLE IF NOT EXISTS entries (
        id integer primary key autoincrement,
        title TEXT,
        content TEXT
      );
    `)

	if err != nil {
		return nil, err
	}

	_, err = db.conn.Exec(`
      CREATE TABLE IF NOT EXISTS settings (
          key VARCHAR(255) UNIQUE,
          value VARCHAR(255)
      );
    `)
	if err != nil {
		return nil, err
	}

	_, err = db.conn.Exec(`
    INSERT OR IGNORE INTO settings (key, value) VALUES 
        ('backup_url', 'http://example.com/backup'),
        ('another_setting', 'some_value');
  `)

	if err != nil {
		return nil, err
	}

	return db, nil
}

/* Returns the name of the database */
func (db SqlLite) GetName() string {
	return db.name
}

/* Reads all entries from the SQL database */
func (db SqlLite) ReadAllEntries() ([]Entry, error) {
	rows, err := db.conn.Query("SELECT id, title FROM entries")
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
func (db SqlLite) CreateEntry(title string, content string) (Entry, error) {
	result, err := db.conn.Exec("INSERT INTO entries(title, content) VALUES(?, ?)", title, content)
	if err != nil {
		return Entry{}, err
	}

	id, _ := result.LastInsertId()
	return db.ReadEntry(id)
}

/* Updates an entry's (by ID) text */
func (db SqlLite) UpdateEntryText(id int64, text string) (Entry, error) {
	_, err := db.conn.Exec("UPDATE entries SET content = ? WHERE id = ?", text, id)
	if err != nil {
		return Entry{}, err
	}

	return db.ReadEntry(id)
}

/* Reads an entry from the database by it's ID */
func (db SqlLite) ReadEntry(id int64) (Entry, error) {
	var entry Entry
	err := db.conn.QueryRow("SELECT id, title, content FROM entries WHERE id = ?", id).Scan(&entry.Id, &entry.Title, &entry.Content)
	if err != nil {
		return Entry{}, err
	}

	return entry, nil
}

/* Renames an entry's (by ID) title */
func (db SqlLite) RenameEntry(id int64, title string) (Entry, error) {
	_, err := db.conn.Exec("UPDATE entries SET title = ? WHERE id = ?", title, id)
	if err != nil {
		return Entry{}, err
	}

	return db.ReadEntry(id)
}

func (db SqlLite) ReadAllSettings() ([]Setting, error) {
	rows, err := db.conn.Query("SELECT key, value FROM settings")
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

func (db SqlLite) UpdateSetting(key string, value string) (Setting, error) {
	_, err := db.conn.Exec("UPDATE settings SET value = ? WHERE key = ?", value, key)
	if err != nil {
		return Setting{}, err
	}

	var setting Setting
	err = db.conn.QueryRow("SELECT key, value FROM settings WHERE key = ?", key).Scan(&setting.Key, &setting.Value)
	if err != nil {
		return Setting{}, err
	}

	return setting, nil
}
