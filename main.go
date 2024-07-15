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

// Problem for tomorrow's Arthur:
// The repo needs to be cloned only once before spawning the containers, otherwise it
// adds a whole new layer of complexity (generating ssh-keys from the containers) and 
// huge additional overhead.
// This will likely require some restructuring - have fun with that. Idea would be:
//
// 	1. shortinette receives a grading trigger (webhook payload)
// 	2. shortinette clones the repo
//	3. For each exercise in the module, shortinette spawns a container, mounting the current directory into it
//	4. The entrypoint command of the container will be something like:
//	
//	`docker run -it --rm -v /root/shortinette:/app <IMAGE NAME> sh -c "go run . <MODULE> <EXERCISE>"`
//	5. We need some way to keep track of which repo we are currently grading across all different containers, 
//	and al results should be appended to the same trace file, which will then be handled by the main program on
//	the host machine.
func dockerExecMode(args []string, short Short.Short) error {
	exercise, ok := short.Modules[args[1]].Exercises[args[2]]
	if !ok {
		return fmt.Errorf("could not find exercise")
	}
	logger.InitializeTraceLogger(args[3])
	exercise.Run()
	return nil
}

func main() {
	logger.InitializeStandardLoggers()
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
	if err := requirements.ValidateRequirements(); err != nil {
		logger.Error.Println(err.Error())
		return
	}
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	Short.StartModule(*R00.R00(), *config)
	short.TestMode.Run()
	Short.EndModule(*R00.R00(), *config)
}
