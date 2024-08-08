package Exercise

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/pkg/logger"
)

type Result struct {
	Passed bool
	Output string
}

type Exercise struct {
	Name            string
	RepoDirectory   string
	TurnInDirectory string
	TurnInFiles     []string
	AllowedSymbols  []string
	AllowedKeywords map[string]int
	Score           int
	Executer        func(test *Exercise) Result
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
//   - exerciseType (TO BE DEPRECATED): function/program/package, used for exercises which do not use any
//     package managers
//   - prototype (TO BE DEPRECATED): function prototype used for compiling single functions
//   - allowedMacros: list of macros to be allowed in this exercise
//   - allowedFunctions: list of functions to be allowed in this exercise
//   - allowedKeywords: list of keywords to be allowed in this exercise
//   - executer: testing function with this signature: "func(test *Exercise) bool", will be run by the module for grading
func NewExercise(
	name string,
	repoDirectory string,
	turnInDirectory string,
	turnInFiles []string,
	allowedSymbols []string,
	allowedKeywords map[string]int,
	score int,
	executer func(test *Exercise) Result,
) Exercise {

	return Exercise{
		Name:            name,
		RepoDirectory:   repoDirectory,
		TurnInDirectory: turnInDirectory,
		TurnInFiles:     turnInFiles,
		AllowedSymbols:  allowedSymbols,
		AllowedKeywords: allowedKeywords,
		Score:           score,
		Executer:        executer,
	}
}

func searchForKeyword(keywords map[string]int, word string) (keyword string, found bool) {
	for keyword := range keywords {
		if word == keyword {
			return keyword, true
		}
	}
	return keyword, false
}

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

func (e *Exercise) fullTurnInFilesPath() []string {
	var fullFilePaths []string

	for _, path := range e.TurnInFiles {
		fullPath := filepath.Join(e.RepoDirectory, e.TurnInDirectory, path)
		fullFilePaths = append(fullFilePaths, fullPath)
	}
	return fullFilePaths
}

func containsString(hayStack []string, needle string) bool {
	for _, str := range hayStack {
		if str == needle {
			return true
		}
	}
	return false
}

func extractAfterExerciseName(exerciseName string, fullPath string) string {
	index := strings.Index(fullPath, exerciseName)
	if index == -1 {
		return "" // or handle the error as needed
	}
	return "'" + fullPath[index+len(exerciseName)+1:] + "'"
}

func (e *Exercise) turnInFilesCheck() Result {
	var foundTurnInFiles []string
	var errors []string
	fullTurnInFilesPaths := e.fullTurnInFilesPath()
	parentDirectory := filepath.Join(e.RepoDirectory, e.TurnInDirectory)
	err := filepath.Walk(parentDirectory, func(path string, info os.FileInfo, err error) error {
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

func (e *Exercise) forbiddenItemsCheck() (result Result) {
	exercisePath := fmt.Sprintf("compile-environment/%s/%s", e.TurnInDirectory, e.TurnInFiles[0])
	err := lintStudentCode(exercisePath, *e)
	if err != nil {
		return Result{Passed: false, Output: err.Error()}
	}

	logger.Info.Printf("no forbidden items/keywords found in %s", e.TurnInDirectory+"/"+e.TurnInFiles[0])
	return Result{Passed: true, Output: ""}
}

// Runs the Executer function and returns the result
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
