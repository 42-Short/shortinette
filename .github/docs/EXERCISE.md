Here’s the updated documentation with the additional functions included:

---

# Exercise Package Documentation

The `Exercise` package defines the structure and behavior of an exercise. It includes fields for exercise metadata, allowed resources, and a single function for executing the tests associated with the exercise. The `Run()` method handles checking for forbidden items and incorrect files, while you only need to implement the test logic.

## Structs and Functions

### `Result`
Represents the result of an exercise execution.
- **Fields**:
  - **`Passed`**: A boolean indicating whether the test was successful.
  - **`Output`**: A string containing output information, which will be logged in the trace for the module.

```go
type Result struct {
    Passed bool
    Output string
}
```

### `Exercise`
Represents an exercise with various metadata fields.
- **Fields**:
  - **`Name`**: The exercise's display name.
  - **`RepoDirectory`**: The target directory for cloning repositories, used to construct file paths.
  - **`TurnInDirectory`**: The directory where the exercise's file(s) can be found, relative to the repository's root (e.g., `ex00/`).
  - **`TurnInFiles`**: A list of all files allowed to be submitted.
  - **`AllowedSymbols`**: A list of symbols (functions, macros, etc.) allowed in this exercise. (Note: The enforcement of allowed symbols is the user's responsibility since it is highly language-specific. If you wish to simply have symbols linted out of the submissions, use `AllowedSymbols` - please consider it might not be as robust as you want it to be.)
  - **`AllowedKeywords`**: A map of keywords allowed in this exercise, with an associated integer indicating the maximum number of times each keyword may appear.
  - **`Score`**: The score assigned to the exercise if passed.
  - **`Executer`**: A function used for testing the exercise, which should be implemented by the user.

```go
type Exercise struct {
    Name             string
    RepoDirectory    string
    TurnInDirectory  string
    TurnInFiles      []string
    AllowedSymbols   []string
    AllowedKeywords  map[string]int
    Score            int
    Executer         func(test *Exercise) Result
}
```

### `NewExercise`
Initializes and returns an `Exercise` struct with the necessary data for submission grading.

- **Parameters**:
  - `name`: The exercise's display name.
  - `repoDirectory`: The target directory for cloning repositories, used to construct file paths.
  - `turnInDirectory`: The directory where the exercise's file(s) can be found, relative to the repository's root (e.g., `ex00/`).
  - `turnInFiles`: A list of all files allowed to be submitted.
  - `allowedSymbols`: A list of symbols (functions, macros, etc.) allowed in this exercise.
  - `allowedKeywords`: A map of keywords allowed in this exercise, with an associated integer indicating the maximum number of times each keyword may appear.
  - `score`: The score assigned to the exercise if passed.
  - `executer`: A testing function that will be run by the module for grading.

```go
func NewExercise(
    name string,
    repoDirectory string,
    turnInDirectory string,
    turnInFiles []string,
    allowedSymbols []string,
    allowedKeywords map[string]int,
    score int,
    executer func(test *Exercise) Result,
) Exercise
```

### Example Exercise Setup
Below is an implementation example for a simple exercise that compiles a Rust program and checks its output.

```go
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
    return Exercise.NewExercise(
        "00", 
        "studentcode", 
        "ex00", 
        []string{"hello.rs"}, 
        []string{"println"}, 
        map[string]int{"unsafe": 0}, 
        10, 
        ex00Test,
    )
}
```

### `Run` Method
The `Run` method executes the exercise's tests after checking for forbidden items and ensuring the correct files are submitted.

- **Returns**: A `Result` struct with the outcome of the exercise execution.

```go
func (e *Exercise) Run() (result Result) {
    if result = e.forbiddenItemsCheck(); !result.Passed {
        return result
    }
    if result = e.turnInFilesCheck(); !result.Passed {
        return result
    }
    e.TurnInFiles = e.fullTurnInFilesPath()

    if e.Executer != nil {
        return e.Executer(e)
    }
    return Result{Passed: false, Output: fmt.Sprintf("no executer found for exercise %s", e.Name)}
}
```

### Error Handling and Result Functions

These functions are used to generate specific types of `Result` instances based on different scenarios encountered during the execution of the exercise.

#### `RuntimeError`
Returns a `Result` indicating a runtime error occurred.

- **Parameters**:
  - `errorMessage`: The error message to include in the output.

- **Returns**: A `Result` with `Passed` set to `false` and the error message.

```go
func RuntimeError(errorMessage string) Result {
    return Result{Passed: false, Output: fmt.Sprintf("runtime error: %s", errorMessage)}
}
```

#### `CompilationError`
Returns a `Result` indicating a compilation error occurred.

- **Parameters**:
  - `errorMessage`: The error message to include in the output.

- **Returns**: A `Result` with `Passed` set to `false` and the error message.

```go
func CompilationError(errorMessage string) Result {
    return Result{Passed: false, Output: fmt.Sprintf("compilation error: %s", errorMessage)}
}
```

#### `InvalidFileError`
Returns a `Result` indicating that invalid files were found in the submission.

- **Returns**: A `Result` with `Passed` set to `false` and a message indicating the presence of invalid files.

```go
func InvalidFileError() Result {
    return Result{Passed: false, Output: "invalid file(s) found in turn in directory"}
}
```

#### `AssertionError`
Returns a `Result` indicating that the output of the student's code did not match the expected output.

- **Parameters**:
  - `expected`: The expected output.
  - `got`: The actual output produced by the student's code.

- **Returns**: A `Result` with `Passed` set to `false` and a message detailing the discrepancy.

```go
func AssertionError(expected string, got string) Result {
    expectedReplaced := strings.ReplaceAll(expected, "\n", "\\n")
    gotReplaced := strings.ReplaceAll(got, "\n", "\\n")
    return Result{Passed: false, Output: fmt.Sprintf("invalid output: expected '%s', got '%s'", expectedReplaced, gotReplaced)}
}
```

#### `InternalError`
Returns a `Result` indicating an internal error occurred during the execution of the exercise.

- **Parameters**:
  - `errorMessage`: The error message to include in the output.

- **Returns**: A `Result` with `Passed` set to `false` and the error message.

```go
func InternalError(errorMessage string) Result {
    return Result{Passed: false, Output: fmt.Sprintf("internal error: %v", errorMessage)}
}
```

#### `Passed`
Returns a `Result` indicating that the exercise was successfully completed.

- **Parameters**:
  - `message`: A success message to include in the output.

- **Returns**: A `Result` with `Passed` set to `true` and the success message.

```go
func Passed(message string) Result {
    return Result{Passed: true, Output: message}
}
```

### Helper Functions

#### `searchForKeyword`
Searches for a keyword in the provided map of allowed keywords.

- **Parameters**:
  - `keywords`: The map of allowed keywords.
  - `word`: The word to search for.

- **Returns**: The keyword and a boolean indicating whether it was found.

```go
func searchForKeyword(keywords map[string]int, word string) (keyword string, found bool)
```

#### `checkKeywordAmount`
Checks if any keywords are used more often than allowed.

- **Parameters**:
  - `keywordCounts`: A map of keyword counts found in the student's code.
  - `keywords`: A map of allowed keywords.

- **Returns**: An error if any keyword is used more than allowed.

```go
func checkKeywordAmount(keywordCounts map[string]int, keywords map[string]int) (err error)
```

#### `scanStudentFile`
Scans a student's file and counts the occurrences of each allowed keyword.

- **Parameters**:
  - `scanner`: A `bufio.Scanner` to read the file.
  - `allowedKeywords`: A map of allowed keywords.

- **Returns**: An error if any keyword is used more than allowed.

```go
func scanStudentFile(scanner *bufio.Scanner, allowedKeywords map[string]int) (err error)
```



#### `lintStudentCode`
Lints the student's code to ensure no forbidden items or keywords are present.

- **Parameters**:
  - `exercisePath`: The path to the exercise file.
  - `test`: The `Exercise` struct.

- **Returns**: An error if any forbidden items or keywords are found.

```go
func lintStudentCode(exercisePath string, test Exercise) (err error)
```

#### `fullTurnInFilesPath`
Constructs the full file paths for the files to be turned in.

- **Returns**: A slice of strings containing the full file paths.

```go
func (e *Exercise) fullTurnInFilesPath() []string
```

#### `containsString`
Checks if a string is present in a slice of strings.

- **Parameters**:
  - `hayStack`: The slice of strings.
  - `needle`: The string to search for.

- **Returns**: A boolean indicating whether the string was found.

```go
func containsString(hayStack []string, needle string) bool
```

#### `extractAfterExerciseName`
Extracts a portion of the file path after the exercise name.

- **Parameters**:
  - `exerciseName`: The name of the exercise.
  - `fullPath`: The full file path.

- **Returns**: A string containing the portion of the file path after the exercise name.

```go
func extractAfterExerciseName(exerciseName string, fullPath string) string
```

#### `turnInFilesCheck`
Checks if the correct files have been turned in.

- **Returns**: A `Result` struct indicating whether the check passed or failed.

```go
func (e *Exercise) turnInFilesCheck() Result
```

#### `forbiddenItemsCheck`
Checks for forbidden items in the student's code.

- **Returns**: A `Result` struct indicating whether the check passed or failed.

```go
func (e *Exercise) forbiddenItemsCheck() (result Result)
```
