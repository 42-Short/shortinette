package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn         *sqlx.DB
	QueryTimeout time.Duration
}

const (
	modulesTableName      string = "modules"
	participantsTableName string = "participants"
)

// Creates a new database connection using the provided DSN.
// It returns a pointer to a DB struct and an error if the connection cannot be established.
func NewDB(dsn string, queryTimeout time.Duration) (*DB, error) {
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB with DSN '%s': %v", dsn, err)
	}

	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Cant ping DB: %v", err)
	}

	fmt.Println("Database successfully created.")
	return &DB{db, queryTimeout}, nil
}

// Sets up the necessary schema in the database and enabling foreign key.
// It returns an error if any of the schema operations fail.
func (db *DB) Initialize() error {
	_, err := db.execWithTimeout("PRAGMA foreign_keys = ON;")
	if err != nil {
		return fmt.Errorf("Error enabling foreign keys: %v", err)
	}

	_, err = db.execWithTimeout(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			intra_login TEXT PRIMARY KEY NOT NULL UNIQUE,
			github_login TEXT NOT NULL UNIQUE,
			FOREIGN KEY (intra_login) REFERENCES participants(intra_login) ON DELETE CASCADE
		);
	`, participantsTableName))
	if err != nil {
		return fmt.Errorf("Error creating Participant schema: %v", err)
	}

	_, err = db.execWithTimeout(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			module_id INTEGER NOT NULL,
			intra_login TEXT NOT NULL,
			attempts INTEGER DEFAULT 0,
			score INTEGER DEFAULT 0,
			last_graded DATETIME,
			wait_time INTEGER DEFAULT 0,
			grading_ongoing BOOLEAN DEFAULT 0,
			PRIMARY KEY (module_id, intra_login),
			FOREIGN KEY (intra_login) REFERENCES participants(intra_login) ON DELETE CASCADE
		);
	`, modulesTableName))
	if err != nil {
		return fmt.Errorf("Error creating Module schema: %v", err)
	}

	fmt.Println("Schema Tables successfully created.")
	return nil
}

// Closes the database connection.
func (db *DB) Close() error {
	return db.Conn.Close()
}

// Executes a query with a timeout specified in the DB struct.
// It returns sql.Result and any error encountered during execution.
func (db *DB) execWithTimeout(query string, args ...any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	result, err := db.Conn.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Executes a Named query with a timeout specified in the DB struct.
// It returns sql.Result and any error encountered during execution.
func (db *DB) namedExecWithTimeout(query string, arg any) (sql.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	result, err := db.Conn.NamedExecContext(ctx, query, arg)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Executes a query with a timeout specified in the DB struct.
// It retrieves a single row from the database and maps it to the provided struct.
func (db *DB) getWithTimeout(dest interface{}, query string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	err := db.Conn.GetContext(ctx, dest, query, args...)
	if err != nil {
		return fmt.Errorf("failed to get data with timeout: %v", err)
	}

	return nil
}

// Executes a query with a timeout specified in the DB struct.
// It retrieves multiple rows from the database and maps it to the provided struct.
func (db *DB) selectWithTimeout(dest interface{}, query string, args ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), db.QueryTimeout)
	defer cancel()

	err := db.Conn.SelectContext(ctx, dest, query, args...)
	if err != nil {
		return fmt.Errorf("failed to get data with timeout: %v", err)
	}

	return nil
}
