package db

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close DB: %v", err)
	}
}

func TestInitialize(t *testing.T) {
	db, err := NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to Initialize DB: %v", err)
	}

	err = verifySchemaTableExists(db, "module")
	if err != nil {
		t.Fatalf("failed to verify table existence: %v", err)
	}
	err = verifySchemaTableExists(db, "participant")
	if err != nil {
		t.Fatalf("failed to verify table existence: %v", err)
	}
}

func TestInvalidNewDB(t *testing.T) {
	db, err := NewDB(context.Background(), "/path/to/invalid-dsn")
	if err == nil {
		db.Close()
		t.Fatalf("expected error when opening DB with invalid DSN, but got nil")
	}
}

func TestClose(t *testing.T) {
	db, err := NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close DB: %v", err)
	}
}

func TestInitializeExistingTables(t *testing.T) {
	db, err := NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to Initialize DB: %v", err)
	}
	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to initialize DB again: %v", err)
	}

	err = verifySchemaTableExists(db, "module")
	if err != nil {
		t.Fatalf("failed to verify table existence: %v", err)
	}
	err = verifySchemaTableExists(db, "participant")
	if err != nil {
		t.Fatalf("failed to verify table existence: %v", err)
	}
}

func TestBackup(t *testing.T) {
	db, err := NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	err = db.Initialize()
	if err != nil {
		db.Close()
		t.Fatalf("failed to Initialize DB: %v", err)
	}

	const backupDir = "./test_backups"
	db.StartBackupScheduler(backupDir, time.Millisecond*300)
	time.Sleep(1 * time.Second)

	db.Close()
	err = os.RemoveAll(backupDir)
	if err != nil {
		t.Fatalf("failed to remove backup directory: %v", err)
	}
}

func verifySchemaTableExists(db *DB, targetTable string) error {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%s'", targetTable)
	err := db.Conn.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query %s table: %v", targetTable, err)
	}
	if count != 1 {
		return fmt.Errorf("expected 1 '%s' table, but found %d", targetTable, count)
	}
	return nil
}
