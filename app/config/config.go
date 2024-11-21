package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/42-Short/shortinette/logger"
	"github.com/joho/godotenv"
)

//TODO: should be part of the short package
//	ModuleDuration time.Duration
//	StartTime      time.Time

// Main struct, containing all necessary metadata
type Config struct {
	Participants []Participant
	Modules      []Module

	TEMPLATE_REPO string
	GITHUB_TOKEN  string
	GITHUB_ORGA   string
	SERVER_ADDR   string
	API_TOKEN     string
	CONFIG_PATH   string
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
	IntraLogin     string `json:"intra_login"`
	GithubUserName string `json:"github_username"`
}

var C Config

func init() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Warning.Printf(".env file not found, this is expected in the GitHub Actions environment, this is a problem if you are running this locally\n")
	}

	requiredEnvVars := map[string]*string{
		"TEMPLATE_REPO": &C.TEMPLATE_REPO,
		"GITHUB_TOKEN":  &C.GITHUB_TOKEN,
		"GITHUB_ORGA":   &C.GITHUB_ORGA,
		"API_TOKEN":     &C.API_TOKEN,
		"SERVER_ADDR":   &C.SERVER_ADDR,
		"CONFIG_PATH":   &C.CONFIG_PATH,
	}

	missingEnvVars := make([]string, 0, len(requiredEnvVars))
	for key, value := range requiredEnvVars {
		*value = os.Getenv(key)
		if *value == "" {
			missingEnvVars = append(missingEnvVars, key)
		}
	}

	if len(missingEnvVars) > 0 {
		logger.Error.Fatalf("missing environment variables: %s", strings.Join(missingEnvVars, ", "))
	}
}

// Reads the participants list (json) from participantsListPath.
//
// Returns a slice of Participant structs containing the GitHub usernames of the participants.
func (config *Config) LoadParticipants(participantsConfigPath string) error {
	data, err := os.ReadFile(participantsConfigPath)
	if err != nil {
		return fmt.Errorf("unable to read config file %s: %v", participantsConfigPath, err)
	}

	if err := json.Unmarshal(data, &config.Participants); err != nil {
		return fmt.Errorf("unable to parse participants list: %v", err)
	}

	if len(config.Participants) < 1 {
		return fmt.Errorf("you need at least one participant")
	}
	for _, participant := range config.Participants {
		if participant.IntraLogin == "" || participant.GithubUserName == "" {
			return fmt.Errorf("participant information incomplete")
		}
	}

	return nil
}

// Initializes a new Module (group of Exercise structs).
//
// Arguments:
//
//   - exercises: list of single exercises
//   - minimumScore: minimum score needed to pass the module
func (config *Config) NewModule(exercises []Exercise, minimumScore int) (mod *Module, err error) {
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
func (config *Config) NewExercise(executablePath string, score int, allowedFiles []string, turnInDirectory string) (ex *Exercise, err error) {
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
