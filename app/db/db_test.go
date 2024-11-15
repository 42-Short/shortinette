package db

import (
	"fmt"
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

func TestInitialize(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
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
	db, err := NewDB("/path/to/invalid-dsn", 2*time.Second)
	if err == nil {
		db.Close()
		t.Fatalf("expected error when opening DB with invalid DSN, but got nil")
	}
}

func TestClose(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
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
