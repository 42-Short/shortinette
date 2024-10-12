package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Main struct, containing all necessary information for running the Short
type Short struct {
	Modules        []Module
	ModuleDuration time.Duration
	StartTime      time.Time
}

// Group of exercises
type Module struct {
	Exercises    []Exercise
	MinimumScore int
	StartTime    time.Time // Set for each Module in NewShort based on the short's start time and the module duration
}

// Single exercise
type Exercise struct {
	ExecutablePath  string
	Score           int
	AllowedFiles    []string
	TurnInDirectory string
}

type Participant struct {
	GithubUserName string
}

// Reads the participants list (json) from participantsListPath.
//
// Returns a slice of Participant structs containing the GitHub usernames of the participants.
func NewParticipants(participantsListPath string) (participants []Participant, err error) {
	data, err := os.ReadFile(participantsListPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read config file %s: %v", participantsListPath, err)
	}

	var rawConfig struct {
		Participants []struct {
			GithubUserName string `json:"github_username"`
		} `json:"participants"`
	}

	if err := json.Unmarshal(data, &rawConfig); err != nil {
		return nil, fmt.Errorf("unable to parse config file: %v", err)
	}

	for _, p := range rawConfig.Participants {
		participants = append(participants, Participant{
			GithubUserName: p.GithubUserName,
		})
	}

	return participants, nil
}

// Initializes a new Short.
//
// Arguments:
//
//   - modules: slice of Module structs
//   - moduleDuration: how long one module should be active
//   - startTime: when the Short starts
func NewShort(modules []Module, moduleDuration time.Duration, startTime time.Time) (short *Short, err error) {
	if modules == nil || len(modules) < 1 {
		return nil, fmt.Errorf("you need at least one module to initialize a short")
	}
	if moduleDuration < 0 {
		return nil, fmt.Errorf("moduleDuration cannot be negative")
	}

	currentStartTime := startTime
	updatedModules := []Module{}
	for _, mod := range modules {
		updatedModules = append(updatedModules, Module{
			Exercises:    mod.Exercises,
			MinimumScore: mod.MinimumScore,
			StartTime:    currentStartTime,
		})
		currentStartTime = currentStartTime.Add(moduleDuration)
	}

	return &Short{
		Modules:        updatedModules,
		ModuleDuration: moduleDuration,
		StartTime:      startTime,
	}, nil
}

// Initializes a new Module (group of Exercise structs).
//
// Arguments:
//
//   - exercises: list of single exercises
//   - minimumScore: minimum score needed to pass the module
func NewModule(exercises []Exercise, minimumScore int) (mod *Module, err error) {
	if exercises == nil || len(exercises) < 1 {
		return nil, fmt.Errorf("you need at least one exercise to initialize a module")
	}
	totalScore := 0
	for _, ex := range exercises {
		totalScore += ex.Score
	}
	if totalScore < minimumScore {
		return nil, fmt.Errorf("the total score of all exercises (%d) adds up to less than expected minimum score (%d)", totalScore, minimumScore)
	}
	if minimumScore < 0 {
		return nil, fmt.Errorf("minimumScore cannot be negative")
	}

	return &Module{
		Exercises:    exercises,
		MinimumScore: minimumScore,
	}, nil
}

// Initializes a new Exercise (data structure for single exercises).
//
// Arguments:
//   - executablePath: path to the executable for running the tests
//   - score: score given when passing this exercise
//   - allowedFiles: files allowed to be found in this exercise's directory
//   - the repository's in which the exercise files are expected to be found
func NewExercise(executablePath string, score int, allowedFiles []string, turnInDirectory string) (ex *Exercise, err error) {
	if allowedFiles == nil || len(allowedFiles) < 1 {
		return nil, fmt.Errorf("at least one allowed file required")
	}
	if score < 0 {
		return nil, fmt.Errorf("score cannot be negative")
	}
	if len(executablePath) < 1 {
		return nil, fmt.Errorf("executablePath cannot be empty")
	}
	if len(turnInDirectory) < 1 {
		return nil, fmt.Errorf("turnInDirectory cannot be empty")
	}

	return &Exercise{
		ExecutablePath:  executablePath,
		Score:           score,
		AllowedFiles:    allowedFiles,
		TurnInDirectory: turnInDirectory,
	}, nil
}
