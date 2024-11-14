package db

import "time"

type Module struct {
	ModuleID       string    `db:"module_id"`
	IntraLogin     string    `db:"intra_login"`
	Attempts       int       `db:"attempts"`
	Score          int       `db:"score"`
	LastGraded     time.Time `db:"last_graded"`
	WaitTime       int       `db:"wait_time"`
	GradingOngoing bool      `db:"grading_ongoing"`
}

// ModuleDAO provides methods to interact with the module table
type ModuleDAO struct {
	DB *DB
}

// InsertModule adds a new module to the modules table.
func (dao *ModuleDAO) InsertModule(module *Module) error {
	panic("InsertModule not implemented yet")
}

// GetModuleByID retrieves a module by its ID.
func (dao *ModuleDAO) GetModuleByID(moduleID string) (*Module, error) {
	panic("GetModuleByID not implemented yet")
}

// GetModulesByLogin retrieves all modules for a specific intraLogin.
func (dao *ModuleDAO) GetModulesByLogin(intraLogin string) ([]Module, error) {
	panic("GetModulesByLogin not implemented yet")
}

// GetAllModules retrieves all modules.
func (dao *ModuleDAO) GetAllModules() ([]Module, error) {
	panic("GetAllModules not implemented yet")
}

// UpdateModule updates an existing module.
func (dao *ModuleDAO) UpdateModule(module *Module) error {
	panic("UpdateModule not implemented yet")
}

// DeleteModule removes a module by its ID.
func (dao *ModuleDAO) DeleteModule(moduleID string) error {
	panic("DeleteModule not implemented yet")
}
