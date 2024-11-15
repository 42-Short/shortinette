package db

import (
	"fmt"
	"testing"
)

func TestInsertModule(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	dao := ModuleDAO{DB: db}
	defer db.Close()

	var count int
	err := dao.DB.getWithTimeout(&count, "SELECT COUNT(*) FROM module;")
	if err != nil {
		t.Fatalf("failed to count modules in DB: %v", err)
	}
	if count != len(modules) {
		t.Fatalf("expected %d modules, but found %d in DB", len(modules), count)
	}

	retrievedModules := []Module{}
	err = dao.DB.selectWithTimeout(&retrievedModules, "SELECT * FROM module;")
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	for i, module := range modules {
		err = validateModule(&module, &retrievedModules[i])
		if err != nil {
			t.Fatalf("failed to fetch module from DB: %v", err)
		}
	}
}

func TestGetModule(t *testing.T) {
	t.Skip("TestGetModulesByLogin not implemented yet")
}

func TestGetModulesByID(t *testing.T) {
	t.Skip("TestGetModulesByID not implemented yet")
}

func TestModuleByIDAndLogin(t *testing.T) {
	t.Skip("TestModuleByIDAndLogin not implemented yet")
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

func validateModule(module *Module, retrievedModule *Module) error {
	if retrievedModule == nil {
		return fmt.Errorf("module not found in DB")
	}

	if retrievedModule.ID != module.ID || retrievedModule.IntraLogin != module.IntraLogin {
		return fmt.Errorf("retrieved module does not match the inserted module")
	}
	return nil
}
