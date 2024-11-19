package db

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	Conn *sqlx.DB
	dsn  string
}

// Creates a new database connection using the provided DSN.
func NewDB(ctx context.Context, dsn string) (*DB, error) {
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB with DSN '%s': %v", dsn, err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("Cant ping DB: %v", err)
	}

	fmt.Println("Database successfully created.")
	return &DB{db, dsn}, nil
}

// Sets up the necessary schema in the database
func (db *DB) Initialize(schemaPath string) error {
	data, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("Error reading sql file: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	_, err = db.Conn.ExecContext(ctx, string(data))
	if err != nil {
		return fmt.Errorf("Error creating Module schema: %v", err)
	}

	fmt.Println("Schema Tables successfully created.")
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

	fmt.Printf("successfully created DB backup at %s\n", backupFile)
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
				fmt.Printf("failed to backup DB: %v\n", err)
			}
		}
	}()
}

// Closes the database connection.
func (db *DB) Close() error {
	return db.Conn.Close()
}
