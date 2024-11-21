package api

import (
	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/logger"
)

func processGrading(dao *data.DAO[data.Module], intra_login string, moduleId int) {
	err := logger.InitializeTraceLogger(logger.GetNewTraceFile(moduleId))
	if err != nil {
		logger.Error.Printf("cant Initialize trace logger: %v", err)
		return
	}
	logger.Info.Printf("grading %s module %d...\n", intra_login, moduleId)
}
