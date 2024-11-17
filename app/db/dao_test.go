package db

import (
	"testing"
)

func TestGetAll(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	moduleDAO := newModuleDAO(db)
	participantDAO := newParticipantDAO(db)
	defer db.Close()

	retrievedModules, err := moduleDAO.GetAll()
	retrievedParticipants, err := participantDAO.GetAll()
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	for i, participant := range participants {
		assertParticipant(t, &participant, &retrievedParticipants[i])
		for j, module := range modules {
			assertModule(t, &module, &retrievedModules[j])
		}
	}
}

func assertParticipant(t *testing.T, participant *Participant, retrievedParticipant *Participant) {
	if retrievedParticipant == nil {
		t.Fatalf("participant not found in DB")
	}

	if retrievedParticipant.GitHubLogin != participant.GitHubLogin || retrievedParticipant.IntraLogin != participant.IntraLogin {
		t.Fatalf("retrieved participants does not match the inserted participant")
	}
}

func assertModule(t *testing.T, module *Module, retrievedModule *Module) {
	if retrievedModule == nil {
		t.Fatalf("module not found in DB")
	}

	if retrievedModule.IntraLogin != module.IntraLogin {
		t.Fatalf("retrieved module does not match the inserted module")
	}
}
