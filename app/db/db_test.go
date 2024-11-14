package db

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close DB: %v", err)
	}
}

//TODO: write tests for ExecWithTimeout

func verifySchemaTableExists(db *DB, targetTable string) error {
	var count int
	query := fmt.Sprintf("SELECT count(*) FROM sqlite_master WHERE type='table' AND name='%s'", targetTable)
	err := db.Connection.QueryRow(query).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query %s table: %v", targetTable, err)
	}
	if count != 1 {
		return fmt.Errorf("expected 1 '%s' table, but found %d", targetTable, count)
	}
	return nil
}

func TestInitialize(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("failed to close DB: %v", err)
		}
	}()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
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
	_, err := NewDB("/path/to/invalid-dsn", 2*time.Second)
	if err == nil {
		t.Fatalf("expected error when opening DB with invalid DSN, but got nil")
	}
}

func TestClose(t *testing.T) {
	testFile := "test.db"

	defer func() {
		err := os.Remove(testFile)
		if err != nil {
			t.Fatalf("failed to remove test database file: %v", err)
		}
	}()

	db, err := NewDB("file:"+testFile+"?cache=shared&_mode=rwc", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close DB: %v", err)
	}
}

func TestInitializeExistingTables(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
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
