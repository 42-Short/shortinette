package db

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func newDummyDB(t *testing.T) (*DB, []Module, []Participant) {
	t.Helper()

	db, err := NewDB(context.Background(), "file::memory:?cache=shared")
	if err != nil {
		t.Fatalf("failed to open DB: %v", err)
	}
	if err := db.Initialize(); err != nil {
		t.Fatalf("failed to initialize DB: %v", err)
	}

	modules, participants := createDummyData(t, db, 7, 40)

	return db, modules, participants
}

func createDummyData(t *testing.T, db *DB, moduleAmount int, participantAmount int) ([]Module, []Participant) {
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
