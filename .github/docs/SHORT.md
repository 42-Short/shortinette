# Short Package
## `short.go`
The **`Short`** package wraps the configuration and execution logic for managing modules and exercises, 
including functions for grading and starting modules.
### Structs and Functions
* **Short**: Represents the configuration for managing modules, with fields **`Name`** and **`TestMode`**.
* **NewShort**: Initializes and returns a **`Short`** struct.
* **GradeModule**: Grades a single participant's module and uploads the traces logs.
* **GradeAll**: Grades all participant's modules and uploads the trace logs.
* **StartModule**: Initializes the module by creating repositories, adding collaborators, and uploading the module's subject.
* **EndModule**: Finalizes the module by grading all repositories and setting read permissions.

### Short Setup Example
Below is an example of how to set up and manage modules using the Short package and the Webhook testmode, which triggers module grading when a participant pushes on their main branch. First, set up your configuration file.
#### Example Configuration
```json
{
    "start_date": "01.01.2024",
    "end_date": "31.01.2024",
    "participants": [
        {
            "github_username": "participant1",
            "intra_login": "p1login"
        },
        {
            "github_username": "participant2",
            "intra_login": "p2login"
        }
    ]
}
```
Afterwards, the configuration file's path can be set in the [.env file](DOTENV.md), and it can be parsed and used for initialization of the Short.
```go
package main

import (
	"fmt"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/git"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	ITestMode "github.com/42-Short/shortinette/pkg/short/testmodes"
)

func main() {
    // Initializes a new Short object with the Webhook test mode.
	short := Short.NewShort("Rust Piscine 1.0", webhook.NewWebhookTestMode())
    // Gets the short configuration to fetch all participants.
	config, err := Short.GetConfig()
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}

    // Starts the module (repo creations, subject uploads, etc).
	Short.StartModule(*R00.R00(), *config)

    // Starts the grader, which will listen for webhook payloads on port 8080.
	short.TestMode.Run()
}
```
