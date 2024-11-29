package api

import (
	"context"
	"fmt"
	"time"

	"github.com/42-Short/shortinette/data"
	"github.com/42-Short/shortinette/logger"
)

func updateModuleGradingState(dao *data.DAO[data.Module], module *data.Module) error {
	module.LastGraded = time.Now()
	module.WaitTime = time.Duration(1<<module.Attempts) * time.Minute // 1, 2, 4, 8, 16, 32,... //TODO: maybe too much?
	module.Attempts++

	err := dao.Update(context.TODO(), *module)
	if err != nil {
		return fmt.Errorf("failed to update module in DB")
	}

	return nil
}

func processGrading(dao *data.DAO[data.Module], intra_login string, module_id int) {
	module, err := dao.Get(context.TODO(), module_id, intra_login)
	if err != nil {
		logger.Error.Printf("failed to get target module for %s%d: %v", intra_login, module_id, err)
		return
	}
	if time.Since(module.LastGraded) < module.WaitTime {
		logger.Warning.Printf("noticed early grading attempt for %s%d, aborting grading", intra_login, module_id)
		return
	}
	err = updateModuleGradingState(dao, module)
	if err != nil {
		logger.Error.Printf("failed to set grading state for %s%d, aborting grading: %v", intra_login, module_id, err)
		return
	}
	logger.Info.Printf("[STATE] id: %d name: %s waittime: %d attempts: %d", module.Id, module.IntraLogin, module.WaitTime, module.Attempts)
	time.Sleep(time.Second * 3)
}
