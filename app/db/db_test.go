package db

import (
	"fmt"
	"os"
	"testing"
)

func TestNewDB(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close DB: %v", err)
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

func TestInitializeDB(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			t.Errorf("failed to close DB: %v", err)
		}
	}()

	err = db.InitializeDB()
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
	_, err := NewDB("/path/to/invalid-dsn")
	if err == nil {
		t.Fatalf("expected error when opening DB with invalid DSN, but got nil")
	}
}

func TestCloseDB(t *testing.T) {
	testFile := "test.db"

	defer func() {
		err := os.Remove(testFile)
		if err != nil {
			t.Fatalf("failed to remove test database file: %v", err)
		}
	}()

	db, err := NewDB("file:" + testFile + "?cache=shared&_mode=rwc")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Close()
	if err != nil {
		t.Fatalf("failed to close DB: %v", err)
	}
}

func TestInitializeExistingTablesDB(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	defer db.Close()

	err = db.InitializeDB()
	if err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}

	err = db.InitializeDB()
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
