package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/42-Short/shortinette/logger"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/joho/godotenv"
)

// Main struct, containing all necessary metadata
type Config struct {
	Modules        []Module
	ModuleDuration time.Duration
	StartTime      time.Time
	DockerImage    string

	TemplateRepo string
	TokenGithub  string
	OrgaGithub   string
	ServerAddr   string
	ApiToken     string
	BasePath     string
}

// Group of exercises
type Module struct {
	ID           int
	Exercises    []Exercise
	MinimumScore int
	StartTime    time.Time // Set for each Module in NewShort based on the short's start time and the module duration
}

// Single exercise
type Exercise struct {
	ID              int
	DockerImage     string
	Score           int
	AllowedFiles    []string
	TurnInDirectory string
}

// Initializes a new Config (group of all configurations)
//
// Arguments:
//
//   - participants: list of single participant
//   - modules: list of single module
//   - moduleDuration: duration of each module
//   - startTime: time on which to start the short
func NewConfig(modules []Module, moduleDuration time.Duration, startTime time.Time, dockerImage string) (conf *Config) {

	for modIdx := range modules {
		modules[modIdx].StartTime = startTime.Add(time.Duration(modIdx) * moduleDuration)
		modules[modIdx].ID = modIdx

		for exIdx := range modules[modIdx].Exercises {
			modules[modIdx].Exercises[exIdx].DockerImage = dockerImage
			modules[modIdx].Exercises[exIdx].ID = exIdx
		}
	}

	return &Config{
		Modules:        modules,
		ModuleDuration: moduleDuration,
		StartTime:      startTime,
	}
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
func NewExercise(score int, allowedFiles []string, turnInDirectory string) (ex *Exercise, err error) {
	if allowedFiles == nil || len(allowedFiles) < 1 {
		return nil, fmt.Errorf("at least one allowed file required")
	}

	for _, globPattern := range allowedFiles {
		if !doublestar.ValidatePattern(globPattern) {
			return nil, fmt.Errorf("allowedFiles contains an invalid glob pattern: %s", globPattern)
		}
	}

	if score < 0 {
		return nil, fmt.Errorf("score cannot be negative")
	}
	if len(turnInDirectory) < 1 {
		return nil, fmt.Errorf("turnInDirectory cannot be empty")
	}

	return &Exercise{
		Score:           score,
		AllowedFiles:    allowedFiles,
		TurnInDirectory: turnInDirectory,
	}, nil
}

func (config *Config) FetchEnvVariables() error {
	err1 := godotenv.Load("../.env")
	err2 := godotenv.Load("../.env")
	if err1 != nil && err2 != nil {
		logger.Warning.Println(".env file not found, this is expected in the GitHub Actions environment, this is a problem if you are running this locally")
	}
	requiredEnvVars := map[string]*string{
		"TEMPLATE_REPO": &config.TemplateRepo,
		"TOKEN_GITHUB":  &config.TokenGithub,
		"ORGA_GITHUB":   &config.OrgaGithub,
		"API_TOKEN":     &config.ApiToken,
		"SERVER_ADDR":   &config.ServerAddr,
		"BASE_PATH":     &config.BasePath,
	}

	missingEnvVars := make([]string, 0, len(requiredEnvVars))
	for key, value := range requiredEnvVars {
		*value = os.Getenv(key)
		if *value == "" {
			missingEnvVars = append(missingEnvVars, key)
		}
	}

	if len(missingEnvVars) > 0 {
		return fmt.Errorf("missing environment variables: %s", strings.Join(missingEnvVars, ", "))
	}

	return nil
}
