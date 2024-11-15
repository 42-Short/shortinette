package db

import (
	"fmt"
	"time"

	_ "github.com/jmoiron/sqlx"
)

type Module struct {
	ID             int       `db:"module_id"`
	IntraLogin     string    `db:"intra_login"`
	Attempts       int       `db:"attempts"`
	Score          int       `db:"score"`
	LastGraded     time.Time `db:"last_graded"`
	WaitTime       int       `db:"wait_time"`
	GradingOngoing bool      `db:"grading_ongoing"`
}

// ModuleDAO (Data Access Object) provides methods to interact with the module table.
type ModuleDAO struct {
	DB *DB
}

// InsertModule adds a new module to the modules table.
func (dao *ModuleDAO) InsertModule(module *Module) error {
	_, err := dao.DB.namedExecWithTimeout(`
		INSERT INTO module (module_id, intra_login, attempts, score, last_graded, wait_time, grading_ongoing)
		VALUES (:module_id, :intra_login, :attempts, :score, :last_graded, :wait_time, :grading_ongoing);
	`, module)

	return err
}

// GetModule retrieves a unique module by ID and intraLogin
func (dao *ModuleDAO) GetModule(moduleID int, intraLogin string) (*Module, error) {
	var module Module
	err := dao.DB.getWithTimeout(&module, "SELECT * FROM module WHERE module_id = ? and intra_login = ?;", moduleID, intraLogin)
	if err != nil {
		return nil, fmt.Errorf("failed to get module with timeout: %v", err)
	}

	return &module, nil
}

// GetModulesByID retrieves all modules for a specific moduleID.
func (dao *ModuleDAO) GetModulesByID(moduleID int) ([]Module, error) {
	var modules []Module
	err := dao.DB.selectWithTimeout(&modules, "SELECT * FROM module WHERE module_id = ?;", moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules by ID with timeout: %v", err)
	}

	return modules, nil
}

// GetModulesByLogin retrieves all modules for a specific intraLogin.
func (dao *ModuleDAO) GetModulesByLogin(intraLogin string) ([]Module, error) {
	var modules []Module
	err := dao.DB.selectWithTimeout(&modules, "SELECT * FROM module WHERE intra_login = ?;", intraLogin)
	if err != nil {
		return nil, fmt.Errorf("failed to get modules by login with timeout: %v", err)
	}

	return modules, nil
}

// GetAllModules retrieves all modules.
func (dao *ModuleDAO) GetAllModules() ([]Module, error) {
	var modules []Module
	err := dao.DB.selectWithTimeout(&modules, "SELECT * FROM module;")
	if err != nil {
		return nil, fmt.Errorf("failed to get modules by login with timeout: %v", err)
	}

	return modules, nil
}

// UpdateModule updates an existing module.
func (dao *ModuleDAO) UpdateModule(module *Module) error {
	panic("UpdateModule not implemented yet")
}

// DeleteModule removes a module by its ID.
func (dao *ModuleDAO) DeleteModule(moduleID string, intraLogin string) error {
	panic("DeleteModule not implemented yet")
}
