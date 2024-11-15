package db

import (
	"fmt"
	"testing"
)

func TestInsertModule(t *testing.T) {
	db, module, _ := newDummyDB(t)
	dao := ModuleDAO{DB: db}

	var retrievedModule Module
	err := dao.DB.getWithTimeout(&retrievedModule, "SELECT * FROM module WHERE module_id = ?", module.ID)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	err = validateModule(module, &retrievedModule)
	if err != nil {
		t.Fatalf("module validation failed: %v", err)
	}
}

func TestGetModulesByID(t *testing.T) {
	t.Skip("TestGetModulesByID not implemented yet")
}

func TestGetModulesByLogin(t *testing.T) {
	t.Skip("TestGetModulesByLogin not implemented yet")
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
