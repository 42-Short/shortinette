package data

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"math/rand"

	"github.com/42-Short/shortinette/db"
)

type data struct {
	participants []Participant
	modules      []Module
}

func SeedDB(db *db.DB) (*data, error) {
	moduleDao := NewDAO[Module](db)
	participantDao := NewDAO[Participant](db)

	const (
		participantAmount = 20
		moduleAmount      = 7
	)

	participants := make([]Participant, 0, participantAmount)
	modules := make([]Module, 0, moduleAmount*participantAmount)

	for i := 0; i < participantAmount; i++ {
		participant := newDummyParticipant()

		if err := participantDao.Insert(context.Background(), *participant); err != nil {
			return nil, fmt.Errorf("failed to insert participant into DB: %v", err)
		}
		participants = append(participants, *participant)

		for j := 0; j < moduleAmount; j++ {
			module := newDummyModule(j, participant.IntraLogin)
			if err := moduleDao.Insert(context.Background(), *module); err != nil {
				return nil, fmt.Errorf("failed to insert module into DB: %v", err)
			}
			modules = append(modules, *module)
		}
	}
	return &data{participants: participants, modules: modules}, nil
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