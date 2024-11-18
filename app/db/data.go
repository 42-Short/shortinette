package db

import (
	"reflect"
	"strings"
	"time"
)

type Module struct {
	Id             int       `db:"module_id" primaryKey:"module_id"`
	IntraLogin     string    `db:"intra_login" primaryKey:"intra_login"`
	Attempts       int       `db:"attempts"`
	Score          int       `db:"score"`
	LastGraded     time.Time `db:"last_graded"`
	WaitTime       int       `db:"wait_time"`
	GradingOngoing bool      `db:"grading_ongoing"`
}

type Participant struct {
	IntraLogin  string `db:"intra_login" primaryKey:"intra_login"`
	GitHubLogin string `db:"github_login"`
}

func deriveSchemaNameFromStruct[T any](dummy T) string {
	return strings.ToLower(reflect.TypeOf(dummy).Name())
}
