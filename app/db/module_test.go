package db

import (
	"fmt"
	"math/rand"
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
			t.Fatalf("failed to validate module: %v", err)
		}
	}
}

func TestGetModule(t *testing.T) {
	db, modules, participants := newDummyDB(t)
	dao := ModuleDAO{DB: db}
	defer db.Close()

	idx := rand.Intn(len(modules))
	retrievedModule, err := dao.GetModule(modules[idx].ID, participants[idx].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	err = validateModule(&modules[idx], retrievedModule)
	if err != nil {
		t.Fatalf("failed to validate module: %v", err)
	}
}

func TestGetModulesByID(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	dao := ModuleDAO{DB: db}
	defer db.Close()

	idx := rand.Intn(len(modules))
	retrievedModules, err := dao.GetModulesByID(modules[idx].ID)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	for _, retrievedModule := range retrievedModules {
		err = validateModule(&modules[idx], &retrievedModule)
		if err != nil {
			t.Fatalf("failed to validate module: %v", err)
		}
	}
}

func TestGetModulesByLogin(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	dao := ModuleDAO{DB: db}
	defer db.Close()

	idx := rand.Intn(len(modules))
	retrievedModules, err := dao.GetModulesByLogin(modules[idx].IntraLogin)
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	for _, retrievedModule := range retrievedModules {
		err = validateModule(&modules[idx], &retrievedModule)
		if err != nil {
			t.Fatalf("failed to validate module: %v", err)
		}
	}
}

func TestGetAllModules(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	dao := ModuleDAO{DB: db}
	defer db.Close()

	retrievedModules, err := dao.GetAllModules()
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	for i, module := range modules {
		err = validateModule(&module, &retrievedModules[i])
		if err != nil {
			t.Fatalf("failed to validate module: %v", err)
		}
	}
}

func TestUpdateModule(t *testing.T) {
	t.Skip("TestUpdateModule not implemented yet")
}

func TestDeleteModule(t *testing.T) {
	db, modules, _ := newDummyDB(t)
	dao := ModuleDAO{DB: db}
	defer db.Close()

	retrievedModules, err := dao.GetAllModules()
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	if len(retrievedModules) != len(modules) || len(modules) == 0 {
		t.Fatalf("failed to fetch all modules: %v", err)
	}

	for _, module := range modules {
		err = dao.DeleteModule(module.ID, module.IntraLogin)
		if err != nil {
			t.Fatalf("failed to delete module: %v", err)
		}
	}

	retrievedModules, err = dao.GetAllModules()
	if err != nil {
		t.Fatalf("failed to fetch module from DB: %v", err)
	}
	if len(retrievedModules) != 0 {
		t.Fatalf("expected modules to be delted but got %d modules: %v", len(retrievedModules), err)
	}

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
