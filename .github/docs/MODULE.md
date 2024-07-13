# Module Package
## module.go
The `Module` package groups multiple exercises into a single module and provides functions
to execute and manage these exercises.
### **Structs and Functions**
* **Module**: Represents a collection of exercises, with fields `Name` and `Exercises`.
* **NewModule**: Initializes and returns a **`Module`** struct.
* **setUpEnvironment**: Prepares the environment for running the exercises by cloning the neessary repositories.
* **Run**: Executes all exercises in the module, returning the results and the path to the trace logs.

### Module Setup Example
Below is an example of how to set up and run a module with exercises, using the example from the [exercise module documentation](EXERCISE.md).

```go
package module00

import (
	"github.com/42-Short/shortinette/internal/logger"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
)

// Initializes and returns a module object.
func module00() *Module.Module {
	r00, err := Module.NewModule("module00", []Exercise.Exercise{ex00(), ..., exXX()})
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return nil
	}
	return &r00
}
```
