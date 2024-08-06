package shortinette

import (
	"os"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/requirements"
	Short "github.com/42-Short/shortinette/pkg/short"
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

func Start(short Short.Short, module string) {
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	Short.StartModule(short.Modules[module], *config)
	short.TestMode.Run(module)
	if len(os.Args) == 4 {
		dockerExecMode(os.Args, short)
	} else if len(os.Args) != 1 {
		logger.Error.Println("invalid number of arguments")
		return
	}
}

func Init() {
	if len(os.Args) == 4 {
		logger.InitializeStandardLoggers(os.Args[2])
	} else {
		logger.InitializeStandardLoggers("")
	}
	if err := requirements.ValidateRequirements(); len(os.Args) != 4 && err != nil {
		logger.Error.Println(err.Error())
		return
	}
}
