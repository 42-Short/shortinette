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
	"github.com/42-Short/shortinette/pkg/testutils"
)

func dockerExecMode(args []string, short Short.Short) error {
	exercise, ok := short.Modules[args[1]].Exercises[args[2]]
	if !ok {
		return fmt.Errorf("could not find exercise")
	}
	if err := logger.InitializeTraceLogger(args[3]); err != nil {
		return err
	}
	result := exercise.Run()
	logger.File.Printf("[MOD%s][EX%s]: %s", args[1], args[2], result.Output)
	if result.Passed {
		os.Exit(0)
	} else {
		os.Exit(1)
	}
	return nil
}

func verifyImage() error {
	cmd := "docker image ls | grep testenv"

	if _, err := testutils.RunCommandLine(".", "sh", []string{"-c", cmd}); err != nil {
		return err
	}
	return nil
}

func main() {
	logger.InitializeStandardLoggers()
	if err := requirements.ValidateRequirements(); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	short := Short.NewShort("Rust Piscine 1.0", map[string]Module.Module{"00": *R00.R00()}, webhook.NewWebhookTestMode())
	if len(os.Args) == 4 {
		if err := dockerExecMode(os.Args, short); err != nil {
			logger.Error.Println(err)
			return
		}
		return
	} else if len(os.Args) != 1 {
		logger.Error.Println("invalid number of arguments")
		return
	}
	if err := verifyImage(); err != nil {
		fmt.Println(err.Error())
		logger.Info.Printf("in order to compile and test submissions in a safe environment, you will need a pre-built Docker image containing all language-specific dependencies needed to compile the code which is to be tested - see http://github.com/42-Short/shortinette/tree/main/.github/docs")
		return
	}

	Short.StartModule(*R00.R00(), *config)
	short.TestMode.Run()
	Short.EndModule(*R00.R00(), *config)
}
