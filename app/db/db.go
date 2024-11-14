package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Connection *sql.DB
}

// Creates a new database connection using the provided DSN.
// It returns a pointer to a DB struct and an error if the connection cannot be established.
func NewDB(dsn string) (*DB, error) {
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB with DSN '%s': %v", dsn, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Cant ping DB: %v", err)
	}

	fmt.Println("Database successfully created.")
	return &DB{db}, nil
}

// Sets up the necessary schema in the database and enabling foreign key.
// It returns an error if any of the schema operations fail.
func (db *DB) Initialize() error {
	_, err := db.Connection.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("Error enabling foreign keys: %v", err)
	}

	_, err = db.Connection.Exec(`
		CREATE TABLE IF NOT EXISTS participant (
			intra_login TEXT PRIMARY KEY NOT NULL,
			github_login TEXT NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("Error creating Participant table: %v", err)
	}

	_, err = db.Connection.Exec(`
		CREATE TABLE IF NOT EXISTS module (
			module_id TEXT NOT NULL,
			intra_login TEXT NOT NULL,
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

	fmt.Println("Schema Tables successfully created.")
	return nil
}

// Closes the database connection.
func (db *DB) Close() error {
	return db.Connection.Close()
}
