# Exercise Package
## `exercise.go`
The `Exercise` package defines the structure and behavior of an exercise. 
It includes fields for exercise metadata, allowed resources, and 1 (!) function for executing
the tests for the exercise.

### Structs and Functions
* **Result**: Represents the result of an exercise execution, with a **`Passed`** field indicating if the test was successful.
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
  * **Executer**: Testing function which is to be run by the module for grading.

## Exercise Setup Example
Below is an implementation example for a simple exercise compiling a Rust exercise with Cargo and checking its output. 

```go
func doTest(*exercise) bool {
    workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
    output, err := testutils.RunCommandLine(workingDirectory, "cargo run")
    if err != nil {
        logger.File.Printf("[%s KO]: runtime error %v", exercise.Name, err)
        return false
    }
    if output != "Hello, cargo!\n" {
        logger.File.Println(testutils.AssertionErrorString(exercise.Name, "Hello, cargo!", output))
        return false
    }
    return true
}

// Implementation of the test function. Runs tests and returns a boolean
// indicating whether the test was successful.
func ex04Test(exercise *Exercise.Exercise) bool {
    // Ensures only the allowed files are present in the turn in directory.
    if !testutils.TurnInFilesCheck(*exercise) {
        return false
    }

    // Ensures no forbidden items have been used.
    if err := testutils.ForbiddenItemsCheck(*exercise, "shortinette-test-R00"); err != nil {
        return false
    }

    // Converts all turn in file paths to be relative to the project's root directory.
    // This is an important step to avoid later undefined behavior.
    exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)

    // Returns a boolean value indicating whether the exercise is passed.
    return doTest(*exercise)
}

// Creates the exercise object to be passed to the module.
func ex04() Exercise.Exercise {
    return Exercise.NewExercise("EX04", "studentcode", "ex04", []string{"src/main.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
```
