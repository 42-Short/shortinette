package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sqlx.DB
	dsn  string
}

// Creates a new database connection using the provided DSN.
func NewDB(ctx context.Context, dsn string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(dsn), os.FileMode(0755)); err != nil {
		return nil, fmt.Errorf("could not create directory '%s' for sqlite database: %v", filepath.Dir(dsn), err)
	}

	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB with DSN '%s': %v", dsn, err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("can not ping DB: %v", err)
	}

	logger.Info.Println("Database successfully created.")
	return &DB{db, dsn}, nil
}

// Sets up the necessary schema in the database
func (db *DB) Initialize(schemaPath string) error {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("error reading schema '%s': %v", schemaPath, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = db.Conn.ExecContext(ctx, string(data))
	if err != nil {
		return fmt.Errorf("error creating Module schema: %v", err)
	}

	logger.Info.Println("Schema Tables successfully created.")
	return nil
}

// creates a backup of the DB in the specified directory
func (db *DB) Backup(backupDir string) error {
	err := os.MkdirAll(backupDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create backup directory %s: %v", backupDir, err)
	}

	backupFile := fmt.Sprintf("%s/%s-backup.db", backupDir, time.Now().Format(time.RFC3339))
	cmd := exec.Command("sqlite3", db.dsn, fmt.Sprintf(".backup %s", backupFile))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute sqlite3 backup: %v, output: %s", err, string(output))
	}

	logger.Info.Printf("successfully created DB backup at %s\n", backupFile)
	return nil
}

// schedules backups at specified intervals
func (db *DB) StartBackupScheduler(backupDir string, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			err := db.Backup(backupDir)
			if err != nil {
				logger.Error.Printf("failed to backup DB: %v\n", err)
			}
		}
	}()
}

// Closes the database connection.
func (db *DB) Close() error {
	return db.Conn.Close()
}
