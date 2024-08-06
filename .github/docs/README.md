# Project Documentation
This documentation provides an overview and setup instructions for the project.

## Overview
The project is build with various packages to manage exercises, modules, and their configurations.

## Key Packages
* **[Exercise](EXERCISE.md)**: Manages individual exercises.
* **[Module](MODULE.md)**: Manages groups of exercises.
* **[Short](SHORT.md)**: Manages the high-level configuration and execution of exercises in a modular fashion.
* **[Git](GIT.md)**: Handles Git operations such as cloning repositories, adding collaborators, and uploading files.
* **[Test Utils](TESTUTILS.md)**: Testing utilities suite.

## Example
Below is an example from our [Rust Piscine](https://github.com/42-Short/rust) repo.

Configuration:
1. See the **[.env](DOTENV.md)** doc to configure your environment
2. You're back! Now configure the `.json` file whose path you set in your environment:
```json
{
    "start_date": "08.07.2024",
    "end_date": "12.07.2024",
    "participants": [
        {
            "github_username": "shortinette-test",
            "intra_login": "1"
        }
    ]
}
```
_note: The intra\_login variable is used to build the names of the repos which will be created. You can of course set it to something else if you want the repos to be named differently. The repo naming format is: <intra_login>-<module_name>_

_note 2: The start\_ & end\_date from the config are not being taken into account by shortinette just yet, watch for future releases!_

3. You're all set! Running the code below will start a Short with one module containing one exercise. A repo will be created for each participant. Your tests will be run on a repo when the participant pushes on the main branch of their repo with `grademe` as a commit message.
```go
package main

import (
	"github.com/42-Short/shortinette/pkg/logger"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette"
	"github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	Short "github.com/42-Short/shortinette/pkg/short"
)

func ex00Test() Exercise.Result {
    if passed {
        return Exercise.Result{Passed: true, Output: "OK"}
    } else {
        return Exercise.Result{Passed: false, Output: "KO"}
    }
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "studentcode", "ex00", []string{"hello.rs"}, "program", "", []string{"println"}, nil, map[string]int{"unsafe": 0}, 10, ex00Test)
}

func R00() *Module.Module {
	exercises := map[string]Exercise.Exercise{
		"00": ex00(),
	}
	r00, err := Module.NewModule("00", 70, exercises)
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return nil
	}
	return &r00
}

func main() {
	shortinette.Init()
	modules := map[string]Module.Module{
		"00": *R00.R00(),
	}
	short := Short.NewShort("Rust Piscine 1.0", modules, webhook.NewWebhookTestMode(modules))
	short.Start("00")
}
```
