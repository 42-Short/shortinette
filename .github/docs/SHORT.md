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
	"rust-piscine/internal/tests/R00"

	"github.com/42-Short/shortinette"
	"github.com/42-Short/shortinette/pkg/short/testmodes/webhooktestmode"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	Short "github.com/42-Short/shortinette/pkg/short"
)

func main() {
	// Initializes shortinette, checks for requirements
	shortinette.Init()

	// Create a map with all of your Module objects in for shortinette
	// to look them up easily
	modules := map[string]Module.Module{
		"00": *R00.R00(),
	}

	// Initialize the Short object and the WebhookTestMode
	short := Short.NewShort("Rust Piscine 1.0", modules, webhook.NewWebhookTestMode(modules))

	// Start the module with the specified module name
	shortinette.Start(short, "00")
}
```
