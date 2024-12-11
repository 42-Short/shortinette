package data

import (
	"time"
)

type Module struct {
	Id         int           `db:"id" json:"id" primaryKey:"id"`
	IntraLogin string        `db:"intra_login" json:"intra_login" primaryKey:"intra_login"`
	Attempts   int           `db:"attempts" json:"attempts"`
	Score      int           `db:"score" json:"score"`
	LastGraded time.Time     `db:"last_graded" json:"last_graded"`
	WaitTime   time.Duration `db:"wait_time" json:"wait_time"`
}

type Participant struct {
	IntraLogin  string `db:"intra_login" json:"intra_login" primaryKey:"intra_login"`
	GitHubLogin string `db:"github_login" json:"github_login"`
}
