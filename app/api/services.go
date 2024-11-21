package api

import (
	"context"
	"os"
	"time"

	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/logger"
)

func processGrading(dao *data.DAO[data.Module], intra_login string, moduleId int) {
	traceFile := logger.GetNewTraceFile(moduleId)
	err := logger.InitializeTraceLogger(traceFile)
	if err != nil {
		logger.Error.Printf("cant Initialize trace logger: %v", err)
		return
	}
	defer os.Remove(traceFile)

	module, err := dao.Get(context.TODO(), moduleId, intra_login)
	if err != nil {
		logger.Error.Printf("cant get module%d for user %s from DB: %v", moduleId, intra_login, err)
		return
	}

	module.Attempts += 1
	module.LastGraded = time.Now()
	module.WaitTime = (module.Attempts - 1) * 5

	logger.Info.Printf("grading %s...\n", intra_login)
	time.Sleep(time.Second * 3)

	err = dao.Update(context.TODO(), *module)
	if err != nil {
		logger.Error.Printf("cant get module%d for user %s from DB: %v", moduleId, intra_login, err)
		return
	}
}
