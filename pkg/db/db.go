// Package `db` provides functions for interacting with the SQLite database, which stores
// information about student repositories, including grading attempts and scores.
package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/42-Short/shortinette/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
)

// Repository represents the data related to a student's repository, including grading
// history and current score.
type Repository struct {
	ID              string        // ID is the unique identifier for the repository.
	FirstAttempt    bool          // FirstAttempt indicates if this is the student's first grading attempt.
	LastGradingTime time.Time     // LastGradingTime stores the timestamp of the last grading attempt.
	WaitingTime     time.Duration // WaitingTime is the required duration to wait before the next grading attempt.
	Score           int           // Score is the student's score in the current module.
	Attempts        int           // Attempts tracks the number of grading attempts made by the student.
}

// GetRepositoryData retrieves the repository data for a specific module and repository ID.
//
//   - moduleName: the name of the module
//   - repoID: the unique identifier for the repository
//
// Returns the repository data and an error if the operation fails.
func GetRepositoryData(moduleName string, repoID string) (repo Repository, err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return Repository{}, err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories_%s", moduleName)

	query := fmt.Sprintf(sqliteTemplates.RepoById, tableName)
	row := db.QueryRow(query, repoID)

	var lastGraded string
	var waitTime int64

	err = row.Scan(&repo.ID, &repo.FirstAttempt, &lastGraded, &waitTime, &repo.Score, &repo.Attempts)
	if err != nil {
		return Repository{}, err
	}

	repo.LastGradingTime, err = time.Parse("2006-01-02 15:04:05", lastGraded)
	if err != nil {
		return Repository{}, err
	}

	repo.WaitingTime = time.Duration(waitTime) * time.Second

	return repo, nil
}

// tableExists checks if a specific table exists in the database.
//
//   - db: the database connection
//   - tableName: the name of the table to check for existence
//
// Returns a boolean indicating if the table exists and an error if the operation fails.
func tableExists(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf(sqliteTemplates.TableByName, tableName)
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

// CreateTable creates a new table in the database for storing repository data if it does
// not already exist.
//
//   - tableName: the name of the table to create
//
// Returns a boolean indicating if the table was created and an error if the operation fails.
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
		query := fmt.Sprintf(sqliteTemplates.CreateRepoTable, tableName)
		_, err = db.Exec(query)
		if err != nil {
			return created, err
		}
		created = true
	}
	logger.Info.Printf("table %s created", tableName)
	return created, nil
}

// InitModuleTable initializes a table for a module with participant data.
//
//   - participants: a slice of slices, where each inner slice contains participant details
//   - moduleName: the name of the module
//
// Returns an error if the initialization fails.
func InitModuleTable(participants [][]string, moduleName string) (err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories_%s", moduleName)

	for _, participant := range participants {
		repoID := fmt.Sprintf("%s-%s", participant[1], moduleName)
		query := fmt.Sprintf(sqliteTemplates.NewRepoRow, tableName)
		if _, err = db.Exec(query, repoID); err != nil {
			return err
		}
	}
	logger.Info.Printf("initialized table %s with participant data", tableName)
	return nil
}

// UpdateRepository updates the repository data in the database after a grading attempt.
//
//   - moduleName: the name of the module
//   - repo: the repository data to be updated
//
// Returns an error if the update fails.
func UpdateRepository(moduleName string, repo Repository) (err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories_%s", moduleName)

	_, err = db.Exec(fmt.Sprintf(sqliteTemplates.UpdateRepo, tableName),
		repo.FirstAttempt,
		(time.Now()).Add(-2*time.Hour).Format("2006-01-02 15:04:05"),
		int(repo.WaitingTime.Seconds()),
		repo.Score,
		repo.Attempts+1,
		repo.ID)
	if err != nil {
		logger.Error.Printf("could not update repo: %v", err)
		return err
	}
	logger.Info.Printf("updated row with id %s in table %s", repo.ID, tableName)
	return nil
}
