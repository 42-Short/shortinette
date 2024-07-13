# Module Package
## module.go
The `Module` package groups multiple exercises into a single module and provides functions
to execute and manage these exercises.
### **Structs and Functions**
* **Module**: Represents a collection of exercises, with fields `Name` and `Exercises`.
* **NewModule**: Initializes and returns a **`Module`** struct.
* **setUpEnvironment**: Prepares the environment for running the exercises by cloning the neessary repositories.
* **Run**: Executes all exercises in the module, returning the results and the path to the trace logs.
