package main

import (
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/42-Short/shortinette/internal/utils"
	Short "github.com/42-Short/shortinette/pkg/short"
	webhook "github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
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
	short := Short.NewShort("Rust Piscine 1.0", webhook.NewWebhookTestMode())
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	Short.StartModule(*R00.R00(), *config)
	short.TestMode.Run()
	Short.EndModule(*R00.R00(), *config)
}
