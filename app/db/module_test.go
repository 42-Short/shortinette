package db

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func constructDummyModule() *Module {
	return &Module{
		ID:             strconv.Itoa(rand.Intn(1000000)),
		IntraLogin:     strconv.Itoa(rand.Intn(1000000)),
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

	_, err := db.ExecWithTimeout(query)
	if err != nil {
		return fmt.Errorf("failed to insert dummy participant: %v", err)
	}

	return nil
}

func TestInsertModule(t *testing.T) {
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}

	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to Initialize DB: %v", err)
	}

	dao := ModuleDAO{DB: db}
	module := constructDummyModule()

	err = insertDummyParticipant(dao.DB, module.IntraLogin)
	if err != nil {
		t.Fatalf("failed to insert participant into DB: %v", err)
	}
	err = dao.InsertModule(module)
	if err != nil {
		t.Fatalf("failed to insert module into DB: %v", err)
	}

	retrievedModule, err := dao.GetModuleByID(module.ID) //TODO: maybe get rid of this
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}

	if retrievedModule == nil {
		t.Fatalf("module not found in DB")
	}

	if retrievedModule.ID != module.ID || retrievedModule.IntraLogin != module.IntraLogin {
		t.Errorf("retrieved module does not match the inserted module")
	}
}

func TestGetModuleByID(t *testing.T) {
	t.Skip("TestGetModuleByID not implemented yet")
}

func TestGetModulesByLogin(t *testing.T) {
	t.Skip("TestGetModulesByLogin not implemented yet")
}

func TestGetAllModules(t *testing.T) {
	t.Skip("TestGetAllModules not implemented yet")
}

func TestUpdateModule(t *testing.T) {
	t.Skip("TestUpdateModule not implemented yet")
}

func TestDeleteModule(t *testing.T) {
	t.Skip("TestDeleteModule not implemented yet")
}
