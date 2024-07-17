package Exercise

import "fmt"

type Result struct {
	Passed bool
	Output string
}

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
	exerciseType string,
	prototype string,
	allowedMacros []string,
	allowedFunctions []string,
	allowedKeywords map[string]int,
	score int,
	executer func(test *Exercise) Result,
) Exercise {

	return Exercise{
		Name:             name,
		RepoDirectory:    repoDirectory,
		TurnInDirectory:  turnInDirectory,
		TurnInFiles:      turnInFiles,
		ExerciseType:     exerciseType,
		Prototype:        prototype,
		AllowedMacros:    allowedMacros,
		AllowedFunctions: allowedFunctions,
		AllowedKeywords:  allowedKeywords,
		Score:            score,
		Executer:         executer,
	}
}

// Run runs the Executer function and returns the result
func (e *Exercise) Run() Result {
	if e.Executer != nil {
		return e.Executer(e)
	}
	return Result{Passed: false, Output: fmt.Sprintf("no executer found for exercise %s", e.Name)}
}
