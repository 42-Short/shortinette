package db

type Participant struct {
	IntraLogin  string `db:"intra_login"`
	GitHubLogin string `db:"github_login"`
}

// ParticipantDAO (Data Access Object) provides methods to interact with the participants table.
type ParticipantDAO struct {
	DB *DB
}

// InsertParticipant adds a new participant to the participants table.
func (dao *ParticipantDAO) InsertParticipant(participant *Participant) error {
	_, err := dao.DB.namedExecWithTimeout(`
		INSERT INTO participant (intra_login, github_login)
		VALUES (:intra_login, :github_login)
	`, participant)
	return err
}

// GetParticipantByLogin retrieves a participant by their intraLogin.
func (dao *ParticipantDAO) GetParticipantByLogin(intraLogin string) (*Participant, error) {
	panic("GetParticipantByLogin not implemented yet")
}

// GetAllParticipants retrieves all participants.
func (dao *ParticipantDAO) GetAllParticipants() ([]Participant, error) {
	panic("GetAllParticipants not implemented yet")
}

// UpdateParticipant updates an existing participant's information.
func (dao *ParticipantDAO) UpdateParticipant(participant *Participant) error {
	panic("UpdateParticipant not implemented yet")
}

// DeleteParticipant removes a participant by their intraLogin.
func (dao *ParticipantDAO) DeleteParticipant(intraLogin string) error {
	panic("DeleteParticipant not implemented yet")
}
