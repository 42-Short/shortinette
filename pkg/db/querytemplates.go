package db

type Templates struct {
	RepoById        string
	TableByName     string
	CreateRepoTable string
	UpdateRepo      string
	NewRepoRow      string
}

var sqliteTemplates = Templates{
	RepoById:        "SELECT id, first_attempt, last_graded, wait_time, score, attempts FROM %s WHERE id = ?",
	TableByName:     "SELECT name FROM sqlite_master WHERE type='table' AND name='%s'",
	CreateRepoTable: "CREATE TABLE %s (id TEXT PRIMARY KEY, first_attempt BOOLEAN DEFAULT 1, last_graded TEXT DEFAULT (datetime('now')), wait_time INTEGER DEFAULT 0, score INTEGER DEFAULT 0, attempts INTEGER DEFAULT 0)",
	UpdateRepo:      "UPDATE %s SET first_attempt = ?, last_graded = ?, wait_time = ?, score = ?, attempts = ? WHERE id = ?",
	NewRepoRow:      "INSERT INTO %s (id) VALUES(?)",
}
