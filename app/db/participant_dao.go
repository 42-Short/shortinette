package db

type Participant struct {
	IntraLogin  string `db:"intra_login"`
	GitHubLogin string `db:"github_login"`
}

// ParticipantDAO (Data Access Object) provides methods to interact with the participants table.
type ParticipantDAO struct {
}
