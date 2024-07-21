package db

import (
	"database/sql"
	"fmt"
	"os"
)

func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s'", tableName)
	var name string
	err := db.QueryRow(query).Scan(&name)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func CreateTable(tableName string) (bool, error) {
	created := false

	if err := os.MkdirAll("./sqlite3", 0755); err != nil {
		return created, err
	}

	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return created, err
	}
	defer db.Close()

	exists, err := tableExists(db, tableName)
	if err != nil {
		return created, err
	}

	if !exists {
		query := fmt.Sprintf("CREATE TABLE `%s` (`id` TEXT PRIMARY KEY, `first_attempt` BOOLEAN DEFAULT 1, `last_graded` TEXT DEFAULT (datetime('now', 'localtime')), `wait_time` INTEGER DEFAULT 0, `score` INTEGER DEFAULT 0)", tableName)
		_, err = db.Exec(query)
		if err != nil {
			return created, err
		}
		created = true
	}
	return created, nil
}

func InitModuleTable(participants [][]string, moduleName string) (err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories-%s", moduleName)

	for _, participant := range participants {
		repoID := fmt.Sprintf("%s-%s", moduleName, participant[1])
		query := fmt.Sprintf("INSERT INTO `%s` (id) VALUES(%s)", tableName, repoID)
		if _, err = db.Exec(query); err != nil {
			return err
		}
	}
	return nil
}
