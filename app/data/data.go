package data

import (
	"time"
)

type Module struct {
	Id             int       `db:"module_id" json:"module_id" primaryKey:"module_id"`
	IntraLogin     string    `db:"intra_login" json:"intra_login" primaryKey:"intra_login"`
	Attempts       int       `db:"attempts" json:"attempts"`
	Score          int       `db:"score" json:"score"`
	LastGraded     time.Time `db:"last_graded" json:"last_graded"`
	WaitTime       int       `db:"wait_time" json:"wait_time"`
	GradingOngoing bool      `db:"grading_ongoing" json:"grading_ongoing"`
}

type Participant struct {
	IntraLogin  string `db:"intra_login" json:"intra_login" primaryKey:"intra_login"`
	GitHubLogin string `db:"github_login" json:"github_login"`
}
