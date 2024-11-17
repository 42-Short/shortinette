package db

import "time"

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
	*BaseDAO[Module]
}

func newModuleDAO(db *DB) *ModuleDAO {
	return &ModuleDAO{
		BaseDAO: NewBaseDao[Module](db, "module"),
	}
}
