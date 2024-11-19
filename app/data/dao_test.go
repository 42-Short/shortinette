package data

import (
	"context"
	"strconv"
	"testing"
	"time"

	"math/rand"

	"github.com/42-Short/shortinette/db"
)

func TestInsert(t *testing.T) {
	db, _, participants := newDummyDB(t)
	participantDAO := NewDAO[Participant](db)
	defer db.Close()

	participant := newDummyParticipant()
	err := participantDAO.Insert(context.Background(), participant)
	if err != nil {
		t.Fatalf("failed to insert participant into DB: %v", err)
	}

	retrievedParticipants, err := participantDAO.GetAll(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch participant from DB: %v", err)
	}
	expectedSize := len(participants) + 1
	actualSize := len(retrievedParticipants)
	if actualSize != expectedSize {
		t.Fatalf("expected %d participants after insertion but got %d", expectedSize, actualSize)
	}
	assertParticipant(t, participant, &retrievedParticipants[actualSize-1])
}

func TestUpdate(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	defer db.Close()

	modules[0].Score += 100
	modules[0].Attempts += 1
	err := moduleDAO.Update(context.Background(), &modules[0])
	if err != nil {
		t.Fatalf("failed to update module in DB %v", err)
	}

	retrievedModule, err := moduleDAO.Get(context.Background(), modules[0].Id, modules[0].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}

	assertModule(t, &modules[0], retrievedModule)
	if modules[0].Score != retrievedModule.Score {
		t.Fatalf("failed to update module score in DB: %v", err)
	}
	if modules[0].Attempts != retrievedModule.Attempts {
		t.Fatalf("failed to update module Attempts in DB: %v", err)
	}

}

func TestGet(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	participantDAO := NewDAO[Participant](db)
	defer db.Close()

	retrievedParticipant, err := participantDAO.Get(context.Background(), participants[0].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch participant from DB: %v", err)
	}
	retrievedModule, err := moduleDAO.Get(context.Background(), modules[0].Id, modules[0].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	assertModule(t, &modules[0], retrievedModule)
	assertParticipant(t, &participants[0], retrievedParticipant)
}

func TestGetFiltered(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	defer db.Close()

	filters := map[string]any{
		"wait_time": modules[0].WaitTime,
		"attempts":  modules[0].Attempts,
	}
	retrievedModules, err := moduleDAO.GetFiltered(context.Background(), filters)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	for _, retrievedModule := range retrievedModules {
		if retrievedModule.WaitTime != modules[0].WaitTime {
			t.Fatalf("incorrect WaitTime in filtered fetch")
		}
		if retrievedModule.Attempts != modules[0].Attempts {
			t.Fatalf("incorrect Score in filtered fetch")
		}
	}
}

func TestGetAll(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	participantDAO := NewDAO[Participant](db)
	defer db.Close()

	retrievedModules, err := moduleDAO.GetAll(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch modules from DB: %v", err)
	}
	retrievedParticipants, err := participantDAO.GetAll(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch participants from DB: %v", err)
	}
	for i, participant := range participants {
		assertParticipant(t, &participant, &retrievedParticipants[i])
		for j, module := range modules {
			assertModule(t, &module, &retrievedModules[j])
		}
	}
}

func TestDelete(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	defer db.Close()

	err := moduleDAO.Delete(context.Background(), modules[0].Id, modules[0].IntraLogin)
	if err != nil {
		t.Fatalf("failed to delete modules from DB: %v", err)
	}
	retrievedModules, err := moduleDAO.GetAll(context.Background())
	if err != nil {
		t.Fatalf("failed to fetch modules from DB: %v", err)
	}
	if len(retrievedModules) != len(modules)-1 {
		t.Fatalf("failed to delete module from DB")
	}
}

func assertParticipant(t *testing.T, participant *Participant, retrievedParticipant *Participant) {
	t.Helper()
	if retrievedParticipant == nil {
		t.Fatalf("participant not found in DB")
	}

	if retrievedParticipant.IntraLogin != participant.IntraLogin {
		t.Fatalf("retrieved participants does not match the inserted participant")
	}
}

func assertModule(t *testing.T, module *Module, retrievedModule *Module) {
	t.Helper()
	if retrievedModule == nil {
		t.Fatalf("module not found in DB")
	}

	if retrievedModule.IntraLogin != module.IntraLogin {
		t.Fatalf("retrieved module does not match the inserted module")
	}
}

func newDummyDB(t *testing.T) (*db.DB, []Module, []Participant) {
	t.Helper()

	db, err := db.NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	if err := db.Initialize("../db/schema.sql"); err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}

	modules, participants := createDummyData(t, db, 7, 40)
	return db, modules, participants
}

func createDummyData(t *testing.T, db *db.DB, moduleAmount int, participantAmount int) ([]Module, []Participant) {
	t.Helper()
	moduleDao := NewDAO[Module](db)
	participantDao := NewDAO[Participant](db)

	participants := make([]Participant, 0, participantAmount)
	modules := make([]Module, 0, moduleAmount*participantAmount)

	for i := 0; i < participantAmount; i++ {
		participant := newDummyParticipant()

		if err := participantDao.Insert(context.Background(), participant); err != nil {
			t.Fatalf("failed to insert participant into DB: %v", err)
		}
		participants = append(participants, *participant)

		for j := 0; j < moduleAmount; j++ {
			module := newDummyModule(j, participant.IntraLogin)
			if err := moduleDao.Insert(context.Background(), module); err != nil {
				t.Fatalf("failed to insert module into DB: %v", err)
			}
			modules = append(modules, *module)
		}
	}
	return modules, participants
}

func newDummyModule(moduleID int, intraLogin string) *Module {
	return &Module{
		Id:             moduleID,
		IntraLogin:     intraLogin,
		Attempts:       rand.Int(),
		Score:          rand.Int(),
		LastGraded:     time.Now(),
		WaitTime:       rand.Int(),
		GradingOngoing: rand.Intn(2) == 0,
	}
}

func newDummyParticipant() *Participant {
	intraLogin := strconv.Itoa(rand.Int())
	return &Participant{
		IntraLogin:  intraLogin,
		GitHubLogin: "dummy_git_" + intraLogin,
	}
}
