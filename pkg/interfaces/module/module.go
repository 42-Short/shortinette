package Module

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/pkg/git"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

type Module struct {
	Name      string
	Exercises []Exercise.Exercise
}

// NewModule initializes and returns a Module struct
func NewModule(name string, exercises []Exercise.Exercise) (Module, error) {
	return Module{
		Name:      name,
		Exercises: exercises,
	}, nil
}

func setUpEnvironment(repoId string, testDirectory string) (tracesPath string, err error) {
	repoLink := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId)
	if err := git.Clone(repoLink, testDirectory); err != nil {
		errorMessage := fmt.Sprintf("failed to clone repository: %v", err)
		return "", errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if tracesPath, err = logger.InitializeTraceLogger(repoId); err != nil {
		errorMessage := fmt.Sprintf("failed to initalize logging system (%v), does the ./traces directory exist?", err)
		return "", errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if err := git.Clone(fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId), "compile-environment/src/"); err != nil {
		return "", err
	}
	return tracesPath, nil
}

// Run executes the exercises and returns the results and the path to the traces
func (m *Module) Run(repoId string, testDirectory string) (results []Exercise.Result, tracesPath string) {
	defer func() {
		if err := os.RemoveAll("compile-environment"); err != nil {
			logger.Error.Printf("could not tear down testing environment: %v", err)
		}
		if err := os.RemoveAll("studentcode"); err != nil {
			logger.Error.Printf("could not tear down testing environment: %v", err)
		}
	}()
	tracesPath, err := setUpEnvironment(repoId, testDirectory)
	if err != nil {
		return nil, tracesPath
	}
	if m.Exercises != nil {
		for _, exercise := range m.Exercises {
			res := exercise.Run()
			results = append(results, res)
			if res.Passed {
				logger.File.Printf("[%s OK]", exercise.Name)
			}
		}
	}
	return results, tracesPath
}
