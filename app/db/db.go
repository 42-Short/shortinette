package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sql.DB
}

// Creates a new db connection
func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("Cant open DB from dsn: %s: %v", dsn, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Cant ping DB: %v", err)
	}
	return &DB{db}, nil
}

// Initializes the database
func (db *DB) InitializeDB() error {
	_, err := db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS module (
			module_id TEXT,
			intra_login TEXT,
			attempts INTEGER,
			score INTEGER,
			last_graded DATETIME,
			wait_time INTEGER,
			grading_ongoing BOOLEAN,
			PRIMARY KEY (module_id, intra_login),
			FOREIGN KEY (intra_login) REFERENCES participant(intra_login) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return fmt.Errorf("Error creating Module table: %v", err)
	}

	_, err = db.Conn.Exec(`
		CREATE TABLE IF NOT EXISTS participant (
			intra_login TEXT PRIMARY KEY,
			github_login TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("Error creating Participant table: %v", err)
	}
	return nil
}

// Closes the database connection.
func (db *DB) Close() error {
	return db.Conn.Close()
}
