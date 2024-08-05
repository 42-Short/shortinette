package shortinette

import (
	"os"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/requirements"
	Short "github.com/42-Short/shortinette/pkg/short"
	webhook "github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
)

func dockerExecMode(args []string, short Short.Short) {
	exercise, ok := short.Modules[args[1]].Exercises[args[2]]
	if !ok {
		os.Exit(1)
	}
	if err := logger.InitializeTraceLogger(args[3]); err != nil {
		os.Exit(1)
	}
	result := exercise.Run()
	logger.File.Printf("[MOD%s][EX%s]: %s", args[1], args[2], result.Output)
	if result.Passed {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
}

func Init(modules map[string]Module.Module) {
	if len(os.Args) == 4 {
		logger.InitializeStandardLoggers(os.Args[2])
	} else {
		logger.InitializeStandardLoggers("")
	}
	if err := requirements.ValidateRequirements(); len(os.Args) != 4 && err != nil {
		logger.Error.Println(err.Error())
		return
	}
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	short := Short.NewShort("Rust Piscine 1.0", modules, webhook.NewWebhookTestMode(modules["00"]))
	if len(os.Args) == 4 {
		dockerExecMode(os.Args, short)
	} else if len(os.Args) != 1 {
		logger.Error.Println("invalid number of arguments")
		return
	}

	Short.StartModule(modules["00"], *config)
	short.TestMode.Run()
	Short.EndModule(modules["00"], *config)
}
