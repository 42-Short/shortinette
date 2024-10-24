//go:build ignore
// Package db provides functions and structures for interacting with the SQLite database
// used to store information about student repositories, including grading attempts and scores.
package db

// Templates contains SQL query templates used for interacting with the SQLite database.
type Templates struct {
	RepoById        string // RepoById is the SQL template for querying repository data by ID.
	TableByName     string // TableByName is the SQL template for checking if a table exists by its name.
	CreateRepoTable string // CreateRepoTable is the SQL template for creating a new repository table.
	UpdateRepo      string // UpdateRepo is the SQL template for updating repository data after a grading attempt.
	NewRepoRow      string // NewRepoRow is the SQL template for inserting a new repository record into the table.
}

// sqliteTemplates holds the SQL query templates used in the db package for SQLite operations.
var sqliteTemplates = Templates{
	RepoById:        "SELECT id, first_attempt, last_graded, wait_time, score, attempts FROM %s WHERE id = ?",
	TableByName:     "SELECT name FROM sqlite_master WHERE type='table' AND name='%s'",
	CreateRepoTable: "CREATE TABLE %s (id TEXT PRIMARY KEY, first_attempt BOOLEAN DEFAULT 1, last_graded TEXT DEFAULT (datetime('now')), wait_time INTEGER DEFAULT 0, score INTEGER DEFAULT 0, attempts INTEGER DEFAULT 1)",
	UpdateRepo:      "UPDATE %s SET first_attempt = ?, last_graded = ?, wait_time = ?, score = ?, attempts = ? WHERE id = ?",
	NewRepoRow:      "INSERT OR IGNORE INTO %s (id) VALUES(?)",
}
