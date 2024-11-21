package api

import (
	"time"

	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/logger"
)

func processGrading(_ *data.DAO[data.Module], intra_login string, _ int) {
	logger.Info.Printf("grading %s...\n", intra_login)
	time.Sleep(time.Second * 3)
}
