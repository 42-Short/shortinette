# shortinette
shortinette is the core framework for managing and automating the process of
grading coding bootcamps (Shorts). It provides a comprehensive set of tools for
running and testing student submissions across various programming languages.
Grading is currently triggered by a [webhook](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/testmodes/webhook) - when a participant 
pushes to the main branch of their repository with `grademe` as a commit message, their submission
will be graded and the results uploaded to their repo.

The shortinette package is composed of several sub-packages, each responsible for a specific
aspect of the grading pipeline:

- [logger](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/logger): Handles logging for the framework, including general informational messages,
  error reporting, and trace logging for feedback on individual submissions. This package ensures
  that all important events and errors are captured for debugging and auditing purposes.

- [requirements](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/requirements): Validates the necessary environment variables and dependencies required
  by the framework. This includes checking for essential configuration values in a `.env` file
  and ensuring that all necessary tools (e.g., Docker images) are available before grading begins.

- [testutils](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/testutils): Provides utility functions for compiling and running code submissions.
  This includes functions for compiling Rust code, running executables with various
  options (such as timeouts and real-time output), and manipulating files. The utility
  functions are designed to handle the intricacies of running untrusted student code
  safely and efficiently.

- [db](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/db): Provides utility functions
  and query templates for interacting with the SQLite database.

- [git](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/git): Manages interactions with GitHub, including cloning repositories, managing
  collaborators, and uploading files. This package abstracts the GitHub API to simplify
  common tasks such as adding collaborators to repositories, creating branches, and
  pushing code or data to specific locations in a repository.

- [exercise](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/exercise): Defines the structure and behavior of individual coding exercises.
  This includes specifying the files that students are allowed to submit, the expected
  output, and the functions to be tested. The `exercise` package provides the framework
  for setting up exercises, running tests, and reporting results.

- [module](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/module): Organizes exercises into modules, allowing for the grouping of related exercises
  into a coherent curriculum. The `module` package handles the execution of all exercises
  within a module, aggregates results, and manages the overall grading process.

- [webhook](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/testmodes/webhook): Enables automatic grading triggered by GitHub webhooks. This allows for a
  fully automated workflow where student submissions are graded as soon as they are
  pushed to a specific branch in a GitHub repository.

- [short](https://pkg.go.dev/github.com/42-Short/shortinette/pkg/short): The central orchestrator of the grading process, integrating all sub-packages
  into a cohesive system. The `short` package handles the setup and teardown of grading
  environments, manages the execution of modules and exercises, and ensures that all
  results are properly recorded and reported.

Overall, shortinette is designed to streamline the grading of programming assignments
in a secure, automated, and scalable manner. It leverages Docker for sandboxed execution
of code, GitHub for version control and collaboration, and a flexible logging system
for detailed tracking of all grading activities. By using shortinette, educators can
focus on teaching and mentoring, while the framework handles the repetitive and error-prone
tasks of compiling, running, and grading student code.
## Implementation Example

### Example Overview
- **Module Name**: `example-module`
- **Exercise Name**: `example-exercise`
- **Programming Language**: Rust
- **Objective**: The exercise expects the student to write a Rust program that prints "Hello, World!" to the console.

### Step 1: Define the Exercise

First, you need to define the exercise. The exercise will specify the files to be submitted, the allowed keywords, and how to test the submission.

```go
package main

import (
	"fmt"

	"github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

// Define the test function for the exercise
func helloWorldTest(ex *exercise.Exercise) exercise.Result {
	// Compile the Rust file
	if err := testutils.CompileWithRustc(turnInFile); err != nil {
		return exercise.CompilationError(err.Error())
	}

	// Run the compiled executable
	output, err := testutils.RunExecutable(testutils.ExecutablePath(turnInFile, ".rs"))
	if err != nil {
		return exercise.RuntimeError(err.Error())
	}

	// Check if the output is correct
	expectedOutput := "Hello, World!\n"
	if output != expectedOutput {
		return exercise.AssertionError(expectedOutput, output)
	}

	return exercise.Passed("OK")
}

// Create the exercise
func createExampleExercise() exercise.Exercise {
	return exercise.NewExercise(
		"example",            // Name of the exercise
		"studentcode",        // Repo directory
		"ex00",               // TurnIn directory
		[]string{"main.rs"},  // Allowed turn-in files
		nil,                  // Allowed symbols (not used)
		nil,                  // Allowed keywords (no restrictions)
		10,                   // Score for the exercise
		helloWorldTest,       // Executer function to test the exercise
	)
}
```

### Step 2: Define the Module

Next, create the module and add the exercise to it.

```go
package main

import (
	"github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/logger"
)

func createExampleModule() module.Module {
	exercises := map[string]exercise.Exercise{
		"hello-world": createHelloWorldExercise(),
	}

	module, err := module.NewModule(
		"example-module",  // Module name
		50,                // Minimum grade to pass the module
		exercises,         // Exercises map
		"subjects/ex00",   // Path to the module subject file
	)
	if err != nil {
		logger.Error.Fatalf("Failed to create module: %v", err)
	}
	return module
}
```

### Step 3: Initialize and Run the Shortinette

Finally, integrate the module into the `shortinette` system and run it.

```go
package main

import (
	"github.com/42-Short/shortinette/pkg/short"
	"github.com/42-Short/shortinette/pkg/webhook"
	"github.com/42-Short/shortinette/pkg/logger"
)

func main() {
	// Create a map of modules for the shortinette
	modules := map[string]module.Module{
		"example-module": createExampleModule(),
	}

	// Set up webhook test mode
	testMode := webhook.NewWebhookTestMode(modules)

	// Initialize and run shortinette
	s := short.NewShort("Example Shortinette", modules, testMode)
	s.Start("example-module")
}
```

### Step 4: Set Up Environment
#### .env File
Create a `.env` file at the root of your repository and fill it up like below:
```.env
# These are used for identifying you when making requests on the GitHub API.
GITHUB_ADMIN="Your GitHub username"
GITHUB_EMAIL="Your GitHub email"
GITHUB_TOKEN="Your GitHub personal access token"

# We use Webhooks to record events on repositories.
WEBHOOK_URL="<HOST>:8080/webhook"

# This is the organization under which the repositories will be created.
GITHUB_ORGANISATION="Your GitHub organization's name"

CONFIG_PATH="Path to your Short config"

```
#### Configuration File
Now configure the .json file whose path you set in your environment:
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
_note: The intra_login variable is used to build the names of the repos which will be created. You can of course set it to something else if you want the repos to be named differently. The repo naming format is: <intra_login>-<module_name>_

_note 2: The start\_ & end_date from the config are not being taken into account by shortinette just yet, watch for future releases!_
### Step 5: Run shortinette

Run the `shortinette` using the `go run .` command (assuming your Go files are in the current directory).

### Explanation

- **Exercise Definition**: The `createExampleExercise` function defines a Rust exercise that expects a file named `main.rs` to be turned in. The `helloWorldTest` function handles compiling and running the file and checks if the output matches "Hello, World!\n".
- **Module Definition**: The `createExampleModule` function creates a module named `example-module` that includes the `hello-world` exercise. The module requires a minimum score of 50 to pass.
- **Running Shortinette**: The `main` function initializes the logging system, creates the module, and starts the `shortinette` with the webhook test mode enabled.
