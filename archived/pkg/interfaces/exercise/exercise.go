//go:build ignore
// Exercise provides structures and functions for defining and running exercises.
package Exercise

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Result represents the result of an exercise execution, including whether it passed and
// any relevant output or error messages.
type Result struct {
	Passed bool   // Passed indicates whether the test was successful.
	Output string // Output contains the output message or error details.
}

// Exercise represents an exercise with various metadata fields.
type Exercise struct {
	Name            string // Name is the exercise's display name.
	CloneDirectory  string
	TurnInDirectory string                      // TurnInDirectory is the directory where the exercise's file(s) can be found, relative to the repository's root.
	TurnInFiles     []string                    // TurnInFiles is a list of all files allowed to be submitted.
	Score           int                         // Score is the score assigned to the exercise if passed.
	Executer        func(test *Exercise) Result // Executer is a function used for testing the exercise, which should be implemented by the user.
}

// NewExercise initializes and returns an Exercise struct with all the necessary data
// for submission grading.
//
//   - name: exercise's display name
//   - repoDirectory: the target directory for cloning repositories, used to construct
//     filepaths
//   - turnInDirectory: the directory in which the exercise's file can be found, relative
//     to the repository's root (e.g., ex00/)
//   - turnInFiles: list of all files allowed to be turned in
//   - score: score assigned to the exercise if passed
//   - executer: testing function with this signature: "func(test *Exercise) Result", will be run by the module for grading
func NewExercise(
	name string,
	turnInDirectory string,
	turnInFiles []string,
	score int,
	executer func(test *Exercise) Result,
) (exercise Exercise) {
	return Exercise{
		Name:            name,
		TurnInDirectory: turnInDirectory,
		TurnInFiles:     turnInFiles,
		Score:           score,
		Executer:        executer,
	}
}

// fullTurnInFilesPath constructs the full file paths for the files to be turned in, relative to
// shortinette's current working directory.
//
// Returns a slice of strings containing the full file paths.
func (e *Exercise) fullTurnInFilesPath() (fullFilePaths []string) {
	for _, path := range e.TurnInFiles {
		fullPath := filepath.Join(e.CloneDirectory, e.TurnInDirectory, path)
		fullFilePaths = append(fullFilePaths, fullPath)
	}
	return fullFilePaths
}

// containsString checks if a string is present in a slice of strings.
//
//   - hayStack: The slice of strings.
//   - needle: The string to search for.
//
// Returns a boolean indicating whether the string was found.
func containsString(hayStack []string, needle string) (found bool) {
	for _, str := range hayStack {
		if str == needle {
			return true
		}
	}
	return false
}

// Returns a file path relative to exercise.TurnInDirectory
//
//   - exerciseName: The name of the exercise.
//   - fullPath: The full file path.
//
// Returns a string containing the portion of the file path after the exercise name.
func extractAfterExerciseName(exerciseName string, fullPath string) (trimmed string) {
	index := strings.Index(fullPath, exerciseName)
	if index == -1 {
		return ""
	}
	return "'" + fullPath[index+len(exerciseName)+1:] + "'"
}

func walkTurnInDirectory(parentDirectory string, fullTurnInFilesPaths []string, exercise Exercise) (errors []string, foundTurnInFiles []string) {
	err := filepath.Walk(parentDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, err.Error())
			return nil
		}
		if filepath.Base(path)[0] == '.' || path == parentDirectory || info.IsDir() {
			return nil
		} else if !containsString(fullTurnInFilesPaths, path) {
			errors = append(errors, extractAfterExerciseName(exercise.Name, path))
		} else {
			foundTurnInFiles = append(foundTurnInFiles, extractAfterExerciseName(exercise.Name, path))
		}
		return nil
	})
	if err != nil {
		errors = append(errors, err.Error())
	}
	return errors, foundTurnInFiles
}

// turnInFilesCheck checks if the correct files have been turned in.
//
// Returns a Result struct indicating whether the check passed or failed.
func (e *Exercise) turnInFilesCheck() (res Result) {
	fullTurnInFilesPaths := e.fullTurnInFilesPath()
	parentDirectory := filepath.Join(e.CloneDirectory, e.TurnInDirectory)
	_, err := os.Stat(parentDirectory)
	if os.IsNotExist(err) {
		return Result{Passed: false, Output: err.Error()}
	}

	errors, foundTurnInFiles := walkTurnInDirectory(parentDirectory, fullTurnInFilesPaths, *e)
	if len(errors) > 0 {
		return Result{Passed: false, Output: fmt.Sprintf("invalid files found in %s/:\n%s\nnot in allowed turn in files", e.TurnInDirectory, strings.Join(errors, "\n"))}
	} else if len(foundTurnInFiles) != len(fullTurnInFilesPaths) {
		return Result{Passed: false, Output: fmt.Sprintf("missing files in %s/; found: %v", e.TurnInDirectory, foundTurnInFiles)}
	}
	return Result{Passed: true, Output: ""}
}

// Run executes the exercise's tests after checking for forbidden items and ensuring
// the correct files are submitted.
//
// Returns a Result struct with the outcome of the exercise execution.
func (e *Exercise) Run() (result Result) {
	if result = e.turnInFilesCheck(); !result.Passed {
		return result
	}
	e.TurnInFiles = e.fullTurnInFilesPath()

	if e.Executer != nil {
		return e.Executer(e)
	}
	return Result{Passed: false, Output: fmt.Sprintf("no executer found for exercise %s", e.Name)}
}
