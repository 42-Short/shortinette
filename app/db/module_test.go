package db

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func constructDummyModule() Module {
	return Module{
		ModuleID:       strconv.Itoa(rand.Intn(1000000)),
		IntraLogin:     strconv.Itoa(rand.Intn(1000000)),
		Attempts:       42,
		Score:          42,
		LastGraded:     time.Now(),
		WaitTime:       42,
		GradingOngoing: false,
	}
}

func TestInsertModule(t *testing.T) {
	t.Skip("TestInsertModule not implemented yet")
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
