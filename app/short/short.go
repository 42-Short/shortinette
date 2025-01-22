package short

import (
	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/dao"
	"github.com/42-Short/shortinette/logger"
)

type Short struct {
	Participants []dao.Participant
	Config       config.Config
}

func NewShort(participants []dao.Participant, config config.Config) (sh Short) {
	return Short{
		Participants: participants,
		Config:       config,
	}
}

func (sh *Short) Launch() (err error) {
	logger.Info.Println("launching Short now!")
	return nil
}
