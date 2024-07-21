package db

type Templates struct {
	RepoById string
}

var templates = Templates{
	RepoById: `
		SELECT id, first_attempt, last_graded, wait_time, score FROM %s WHERE id = ?`,
}