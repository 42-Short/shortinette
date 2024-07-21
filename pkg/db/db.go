package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	ID              string
	FirstAttempt    bool
	LastGradingTime time.Time
	WaitingTime     time.Duration
	Score           int
}

func GetRepositoryData(moduleName string, repoId string) (repo Repository, err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return Repository{}, err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories_%s", moduleName)

	query := fmt.Sprintf("SELECT id, first_attempt, last_graded, wait_time, score FROM `%s` WHERE id = ?", tableName)
	row := db.QueryRow(query, repoId)

	var lastGraded string
	var waitTime int64

	err = row.Scan(&repo.ID, &repo.FirstAttempt, &lastGraded, &waitTime, &repo.Score)
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
		query := fmt.Sprintf("CREATE TABLE `%s` (`id` TEXT PRIMARY KEY, `first_attempt` BOOLEAN DEFAULT 1, `last_graded` TEXT DEFAULT (datetime('now')), `wait_time` INTEGER DEFAULT 0, `score` INTEGER DEFAULT 0)", tableName)
		_, err = db.Exec(query)
		if err != nil {
			return created, err
		}
		created = true
	}
	logger.Info.Printf("table %s created", tableName)
	return created, nil
}

func InitModuleTable(participants [][]string, moduleName string) (err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories_%s", moduleName)

	for _, participant := range participants {
		repoID := fmt.Sprintf("%s-%s", participant[1], moduleName)
		query := fmt.Sprintf("INSERT INTO `%s` (id) VALUES(?)", tableName)
		if _, err = db.Exec(query, repoID); err != nil {
			return err
		}
	}
	logger.Info.Printf("initialized table %s with participant data", tableName)
	return nil
}

func UpdateRepository(moduleName string, repo Repository) (err error) {
	db, err := sql.Open("sqlite3", "./sqlite3/repositories.db")
	if err != nil {
		return err
	}
	defer db.Close()

	tableName := fmt.Sprintf("repositories_%s", moduleName)

	updateQuery := `
			UPDATE %s
			SET first_attempt = ?, last_graded = ?, wait_time = ?, score = ?
			WHERE id = ?
		`
	_, err = db.Exec(fmt.Sprintf(updateQuery, tableName),
		repo.FirstAttempt,
		repo.LastGradingTime.Format("2006-01-02 15:04:05"),
		int(repo.WaitingTime.Seconds()),
		repo.Score,
		repo.ID)
	if err != nil {
		logger.Error.Printf("could not update repo: %v", err)
		return err
	}
	logger.Info.Printf("updated row with id %s in table %s", repo.ID, tableName)
	return nil
}
