package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func createDummyDB(t *testing.T) *DB {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to Initialize DB: %v", err)
	}

	return db
}

func constructDummyModule() *Module {
	return &Module{
		ID:             strconv.Itoa(rand.Int()),
		IntraLogin:     strconv.Itoa(rand.Int()),
		Attempts:       42,
		Score:          42,
		LastGraded:     time.Now(),
		WaitTime:       42,
		GradingOngoing: false,
	}
}

func insertDummyParticipant(db *DB, intraLogin string) error {
	githubLogin := "dummy_github_" + intraLogin

	query := fmt.Sprintf(`
		INSERT INTO participant (intra_login, github_login)
		VALUES ('%s', '%s')
	`, intraLogin, githubLogin)

	_, err := db.execWithTimeout(query)
	if err != nil {
		return fmt.Errorf("failed to insert dummy participant: %v", err)
	}

	return nil
}
