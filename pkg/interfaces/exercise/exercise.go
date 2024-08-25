// Package Exercise provides structures and functions for defining and running exercises.
package Exercise

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/pkg/logger"
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
	AllowedKeywords map[string]int              // AllowedKeywords is a map of keywords allowed in this exercise, with an associated integer indicating the maximum number of times each keyword may appear. This will be linted by shortinette and the exercise will not pass if the submission does not respect those constraints.
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
//   - allowedKeywords: list of keywords to be allowed in this exercise
//   - score: score assigned to the exercise if passed
//   - executer: testing function with this signature: "func(test *Exercise) Result", will be run by the module for grading
func NewExercise(
	name string,
	turnInDirectory string,
	turnInFiles []string,
	allowedKeywords map[string]int,
	score int,
	executer func(test *Exercise) Result,
) Exercise {
	return Exercise{
		Name:            name,
		TurnInDirectory: turnInDirectory,
		TurnInFiles:     turnInFiles,
		AllowedKeywords: allowedKeywords,
		Score:           score,
		Executer:        executer,
	}
}

// searchForKeyword searches for a keyword in the provided map of allowed keywords.
//
//   - keywords: The map of allowed keywords.
//   - word: The word to search for.
//
// Returns the keyword and a boolean indicating whether it was found.
func searchForKeyword(keywords map[string]int, word string) (keyword string, found bool) {
	for keyword := range keywords {
		if word == keyword {
			return keyword, true
		}
	}
	return keyword, false
}

// checkKeywordAmount checks if any keywords are used more often than allowed.
//
//   - keywordCounts: A map of keyword counts found in the student's code.
//   - keywords: A map of allowed keywords.
//
// Returns an error if any keyword is used more than allowed.
func checkKeywordAmount(keywordCounts map[string]int, keywords map[string]int) (err error) {
	foundKeywords := make([]string, 0, len(keywords))
	for keyword, allowedAmount := range keywords {
		if count, inMap := keywordCounts[keyword]; inMap {
			if count > allowedAmount {
				foundKeywords = append(foundKeywords, keyword)
			}
		}
	}
	if len(foundKeywords) > 0 {
		return fmt.Errorf("keywords %s are used more often than allowed", strings.Join(foundKeywords, ", "))
	}
	return nil
}

// scanStudentFile scans a student's file and counts the occurrences of each allowed keyword.
//
//   - scanner: A bufio.Scanner to read the file.
//   - allowedKeywords: A map of allowed keywords.
//
// Returns an error if any keyword is used more than allowed.
func scanStudentFile(scanner *bufio.Scanner, allowedKeywords map[string]int) (err error) {
	keywordCounts := make(map[string]int)
	for scanner.Scan() {
		word := scanner.Text()
		keyword, found := searchForKeyword(allowedKeywords, word)
		if found {
			keywordCounts[keyword]++
		}
	}
	err = checkKeywordAmount(keywordCounts, allowedKeywords)
	if err != nil {
		return err
	}
	return nil
}

// lintStudentCode lints the student's code to ensure no forbidden items or keywords are present.
//
//   - exercisePath: The path to the exercise file.
//   - test: The Exercise struct.
//
// Returns an error if any forbidden items or keywords are found.
func lintStudentCode(exercisePath string, test Exercise) (err error) {
	file, err := os.Open(exercisePath)
	if err != nil {
		return fmt.Errorf("could not open %s: %w", exercisePath, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanWords)
	return scanStudentFile(scanner, test.AllowedKeywords)
}

// fullTurnInFilesPath constructs the full file paths for the files to be turned in.
//
// Returns a slice of strings containing the full file paths.
func (e *Exercise) fullTurnInFilesPath() []string {
	var fullFilePaths []string

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
func containsString(hayStack []string, needle string) bool {
	for _, str := range hayStack {
		if str == needle {
			return true
		}
	}
	return false
}

// extractAfterExerciseName extracts a portion of the file path after the exercise name.
//
//   - exerciseName: The name of the exercise.
//   - fullPath: The full file path.
//
// Returns a string containing the portion of the file path after the exercise name.
func extractAfterExerciseName(exerciseName string, fullPath string) string {
	index := strings.Index(fullPath, exerciseName)
	if index == -1 {
		return "" // or handle the error as needed
	}
	return "'" + fullPath[index+len(exerciseName)+1:] + "'"
}

// turnInFilesCheck checks if the correct files have been turned in.
//
// Returns a Result struct indicating whether the check passed or failed.
func (e *Exercise) turnInFilesCheck() Result {
	var foundTurnInFiles []string
	var errors []string
	fullTurnInFilesPaths := e.fullTurnInFilesPath()
	parentDirectory := filepath.Join(e.CloneDirectory, e.TurnInDirectory)
	_, err := os.Stat(parentDirectory)
	if os.IsNotExist(err) {
		return Result{Passed: false, Output: err.Error()}
	}
	err = filepath.Walk(parentDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			errors = append(errors, err.Error())
			return nil
		}
		if filepath.Base(path)[0] == '.' || path == parentDirectory || info.IsDir() {
			return nil
		} else if !containsString(fullTurnInFilesPaths, path) {
			errors = append(errors, extractAfterExerciseName(e.Name, path))
		} else {
			foundTurnInFiles = append(foundTurnInFiles, extractAfterExerciseName(e.Name, path))
		}
		return nil
	})
	if err != nil {
		errors = append(errors, err.Error())
	}
	if len(errors) > 0 {
		return Result{Passed: false, Output: fmt.Sprintf("invalid files found in %s/:\n%s\nnot in allowed turn in files", e.TurnInDirectory, strings.Join(errors, "\n"))}
	} else if len(foundTurnInFiles) != len(fullTurnInFilesPaths) {
		return Result{Passed: false, Output: fmt.Sprintf("missing files in %s/; found: %v", e.TurnInDirectory, foundTurnInFiles)}
	}
	return Result{Passed: true, Output: ""}
}

// forbiddenItemsCheck checks for forbidden items in the student's code.
//
// Returns a Result struct indicating whether the check passed or failed.
func (e *Exercise) forbiddenItemsCheck() (result Result) {
	pathsToCheck := []string{}

	for _, path := range e.TurnInFiles {
		exercisePath := filepath.Join(e.CloneDirectory, e.TurnInDirectory, path)
		if strings.HasSuffix(exercisePath, ".rs") {
			pathsToCheck = append(pathsToCheck, exercisePath)
		}
	}

	output := ""
	for _, path := range pathsToCheck {
		if err := lintStudentCode(path, *e); err != nil {
			output = fmt.Sprintf("%s\n%s", output, err.Error())
		}
	}
	if output != "" {
		return Result{Passed: false, Output: output}
	}

	logger.Info.Printf("no forbidden items/keywords found in %s", e.TurnInDirectory+"/"+e.TurnInFiles[0])
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
	if result = e.forbiddenItemsCheck(); !result.Passed {
		return result
	}
	e.TurnInFiles = e.fullTurnInFilesPath()

	if e.Executer != nil {
		return e.Executer(e)
	}
	return Result{Passed: false, Output: fmt.Sprintf("no executer found for exercise %s", e.Name)}
}
