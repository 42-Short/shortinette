# Exercise Package
## `exercise.go`
The `Exercise` package defines the structure and behavior of an exercise. 
It includes fields for exercise metadata, allowed resources, and 1 (!) function for executing
the tests for the exercise.
The `Run()` method handles checking for forbidden items and wrong turn in files, you only need to implement the tests.

### Structs and Functions
* **Result**: Represents the result of an exercise execution, with a **`Passed`** field indicating whether the test was successful and an **`Output`** field, which will be written in the trace for the module.
```go
type Result struct {
	Passed bool
	Output string
}
```
* **Exercise**: Represents an exercise with various metadata fields such as:
  * **Name**: The exercise's display name.
  * **RepoDirectory**: Target directory for cloning repositories, used to construct filepaths.
  * **TurnInDirectory**: Directory in which the exercise's file can be found, relative to the repository's root (e.g., ex00/).
  * **TurnInFiles**: List of all files allowed to be turned in.
  * **ExerciseType**: Function/program/package, used for exercises which do not use any package managers.
  * **Prototype**: Function prototype used for compiling single functions.
  * **AllowedMacros**: List of macros to be allowed in this exercise.
  * **AllowedFunctions**: List of functions to be allowed in this exercise.
  * **AllowedKeywords**: List of keywords to be allowed in this exercise.
  * **Score**: Score given to students for passing this exercise.
  * **Executer**: Testing function which is to be run by the module for grading.
```go
type Exercise struct {
	Name             string
	RepoDirectory    string
	TurnInDirectory  string
	TurnInFiles      []string
	ExerciseType     string
	Prototype        string
	AllowedMacros    []string
	AllowedFunctions []string
	AllowedKeywords  map[string]int
	Score            int
	Executer         func(test *Exercise) Result
}
```
## Exercise Setup Example
Below is an implementation example for a simple exercise compiling a Rust exercise and checking its output. 
```go
// The Exercise (import "github.com/42-Short/shortinette/pkg/interfaces/exercise")
// contains some functions returning premade Result structs, like the ones used below
func ex00Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := compile(exercise); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	output, err := runExecutable(executablePath)
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}

	if output != "Hello, World!\n" {
		return Exercise.AssertionError("Hello, World!\n", output)
	}

	return Exercise.Passed("OK")
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "studentcode", "ex00", []string{"hello.rs"}, "program", "", []string{"println"}, nil, map[string]int{"unsafe": 0}, 10, ex00Test)
}

```
