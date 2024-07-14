package main

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/R00"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/requirements"
	Short "github.com/42-Short/shortinette/pkg/short"
	webhook "github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
)

var ModuleOne = map[string]bool{
	"00": true,
	"01": true,
	"02": true,
	"03": true,
	"04": true,
	"05": true,
}

var ModulesLookupTable = map[string]interface{}{
	"00": ModuleOne,
}

func dockerExecMode(args []string) error {
	module, ok := ModulesLookupTable[args[1]]
	if !ok {
		return fmt.Errorf("module not found")
	}

	moduleMap, ok := module.(map[string]bool)
	if !ok {
		return fmt.Errorf("invalid module type")
	}

	if _, ok := moduleMap[args[2]]; !ok {
		return fmt.Errorf("exercise not found in module")
	}
	return nil
}

func main() {
	logger.InitializeStandardLoggers()
	if len(os.Args) == 4 {
		if err := dockerExecMode(os.Args); err != nil {
			logger.Error.Println(err)
			return
		}
		return

	} else if len(os.Args) != 1 {
		logger.Error.Println("invalid number of arguments")
		return
	}
	if err := requirements.ValidateRequirements(); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	short := Short.NewShort("Rust Piscine 1.0", map[string]Module.Module{"00": *R00.R00()}, webhook.NewWebhookTestMode())
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	Short.StartModule(*R00.R00(), *config)
	short.TestMode.Run()
	Short.EndModule(*R00.R00(), *config)
}
