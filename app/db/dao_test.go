package db

import (
	"testing"
)

func TestInsert(t *testing.T) {
	db, _, participants := newDummyDB(t)
	participantDAO := newParticipantDAO(db)
	defer db.Close()

	participant := newDummyParticipant()
	participantDAO.Insert(participant)

	retrievedParticipants, err := participantDAO.GetAll()
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

func TestGet(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	moduleDAO := newModuleDAO(db)
	participantDAO := newParticipantDAO(db)
	defer db.Close()

	retrievedParticipant, err := participantDAO.Get(participants[0].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch participant from DB: %v", err)
	}
	retrievedModule, err := moduleDAO.Get(modules[0].Id, modules[0].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	assertModule(t, &modules[0], retrievedModule)
	assertParticipant(t, &participants[0], retrievedParticipant)
}

func TestGetAll(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	moduleDAO := newModuleDAO(db)
	participantDAO := newParticipantDAO(db)
	defer db.Close()

	retrievedModules, err := moduleDAO.GetAll()
	if err != nil {
		t.Fatalf("failed to fetch modules from DB: %v", err)
	}
	retrievedParticipants, err := participantDAO.GetAll()
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
