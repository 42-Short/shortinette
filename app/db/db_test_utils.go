package db

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func newDummyDB(t *testing.T) (*DB, []Module, []Participant) {
	t.Helper()
	db, err := NewDB("file::memory:?cache=shared", 2*time.Second)
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	err = db.Initialize()
	if err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}

	moduleDao := ModuleDAO{DB: db}
	participantDao := ParticipantDAO{DB: db}
	participants := []Participant{}
	modules := []Module{}

	for i := 0; i < 5; i++ {
		participant := newDummyParticipant()
		module := newDummyModule(participant.IntraLogin)

		err = participantDao.InsertParticipant(participant)
		if err != nil {
			t.Fatalf("failed to insert participant into DB: %v", err)
		}
		participants = append(participants, *participant)
		err = moduleDao.InsertModule(module)
		if err != nil {
			t.Fatalf("failed to insert module into DB: %v", err)
		}
		modules = append(modules, *module)
	}

	return db, modules, participants
}

func newDummyModule(intraLogin string) *Module {
	return &Module{
		ID:             rand.Int(),
		IntraLogin:     intraLogin,
		Attempts:       42,
		Score:          42,
		LastGraded:     time.Now(),
		WaitTime:       42,
		GradingOngoing: false,
	}
}

func newDummyParticipant() *Participant {
	intraLogin := strconv.Itoa(rand.Int())
	return &Participant{IntraLogin: intraLogin, GitHubLogin: "dummy_git_" + intraLogin}
}
