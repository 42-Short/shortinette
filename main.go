package main

import (
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	"github.com/42-Short/shortinette/internal/tests/testutils"
	"github.com/42-Short/shortinette/internal/utils"
)

func main() {
	logger.InitializeStandardLoggers()
	if err := utils.RequireEnv(); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	if _, err := testutils.RunExecutable("./scripts/check_dependencies.sh"); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	logger.Info.Println("all dependencies are already installed")
	R00.R00("shortinette-test-R00", "studentcode")
}
