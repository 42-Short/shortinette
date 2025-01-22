package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/42-Short/shortinette/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsert(t *testing.T) {
	db, _, participants := newDummyDB(t)
	participantDAO := NewDAO[Participant](db)
	defer db.Close()

	participant := NewDummyParticipant(100)
	err := participantDAO.Insert(context.Background(), *participant)
	require.NoError(t, err)

	retrievedParticipants, err := participantDAO.GetAll(context.Background())
	require.NoError(t, err)
	expectedSize := len(participants) + 1
	actualSize := len(retrievedParticipants)
	assert.Equal(t, actualSize, expectedSize, fmt.Sprintf("expected %d participants after insertion but got %d", expectedSize, actualSize))
}

func TestUpdate(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	defer db.Close()

	modules[0].Score += 100
	modules[0].Attempts += 1
	err := moduleDAO.Update(context.Background(), modules[0])
	require.NoError(t, err)

	retrievedModule, err := moduleDAO.Get(context.Background(), modules[0].Id, modules[0].IntraLogin)
	require.NoError(t, err)

	assert.Equal(t, modules[0].IntraLogin, retrievedModule.IntraLogin)
	assert.Equal(t, modules[0].Id, retrievedModule.Id)
	assert.Equal(t, modules[0].Score, retrievedModule.Score, "failed to update module score in DB")
	assert.Equal(t, modules[0].Attempts, retrievedModule.Attempts, "failed to update module score in DB")

}

func TestGet(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	participantDAO := NewDAO[Participant](db)
	defer db.Close()

	retrievedParticipant, err := participantDAO.Get(context.Background(), participants[0].IntraLogin)
	require.NoError(t, err)
	retrievedModule, err := moduleDAO.Get(context.Background(), modules[0].Id, modules[0].IntraLogin)
	require.NoError(t, err)
	assert.Equal(t, retrievedModule.IntraLogin, modules[0].IntraLogin)
	assert.Equal(t, retrievedModule.Id, modules[0].Id)
	assert.Equal(t, retrievedParticipant.IntraLogin, participants[0].IntraLogin)
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
	require.NoError(t, err)
	for _, retrievedModule := range retrievedModules {
		assert.Equal(t, retrievedModule.WaitTime, modules[0].WaitTime, "incorrect WaitTime in filtered fetch")
		assert.Equal(t, retrievedModule.Attempts, modules[0].Attempts, "incorrect WaitTime in filtered fetch")
	}
}

func TestGetAll(t *testing.T) {
	db, _, participants := newDummyDB(t)
	participantDAO := NewDAO[Participant](db)
	defer db.Close()

	retrievedParticipants, err := participantDAO.GetAll(context.Background())
	require.NoError(t, err)
	for i, participant := range participants {
		assert.Equal(t, &participant, &retrievedParticipants[i])
	}
}

func TestDelete(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	moduleDAO := NewDAO[Module](db)
	defer db.Close()

	err := moduleDAO.Delete(context.Background(), modules[0].Id, modules[0].IntraLogin)
	require.NoError(t, err)
	retrievedModules, err := moduleDAO.GetAll(context.Background())
	require.NoError(t, err)
	assert.Equal(t, len(retrievedModules), len(modules)-1, "failed to delete module from DB")
}

func newDummyDB(t *testing.T) (*db.DB, []Module, []Participant) {
	t.Helper()

	db, err := db.NewDB(context.Background(), "file::memory:?cache=shared")
	require.NoError(t, err)
	err = db.Initialize("../db/schema.sql")
	require.NoError(t, err)

	data, err := SeedDB(db)
	require.NoError(t, err, "failed to seed db")
	return db, data.modules, data.participants
}
