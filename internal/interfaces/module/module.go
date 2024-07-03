package Module

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/git"
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
)

type Module struct {
	Name      string
	Exercises []Exercise.Exercise
	// duration
}

// NewModule initializes and returns a Module struct
func NewModule(name string, exercises []Exercise.Exercise, repoId string, testDirectory string) (Module, error) {
	if err := setUpEnvironment(repoId, testDirectory); err != nil {
		return Module{}, err
	}

	return Module{
		Name:      name,
		Exercises: exercises,
	}, nil
}

func setUpEnvironment(repoId string, testDirectory string) error {
	repoLink := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId)
	if err := git.Get(repoLink, testDirectory); err != nil {
		errorMessage := fmt.Sprintf("failed to clone repository: %v", err)
		return errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if err := logger.InitializeTraceLogger(repoId); err != nil {
		errorMessage := fmt.Sprintf("failed to initalize logging system (%v), does the ./traces directory exist?", err)
		return errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if err := git.Get(fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId), "compile-environment/src/"); err != nil {
		return err
	}
	return nil
}

// Run executes the exercises and returns the results
func (m *Module) Run() []Exercise.Result {
	defer func() {
		if err := os.RemoveAll("compile-environment"); err != nil {
			logger.Error.Printf("could not tear down testing environment: %v", err)
		}
		if err := os.RemoveAll("studentcode"); err != nil {
			logger.Error.Printf("could not tear down testing environment: %v", err)
		}
	}()
	var results []Exercise.Result
	if m.Exercises != nil {
		for _, exercise := range m.Exercises {
			res := exercise.Run()
			results = append(results, res)
		}
	}
	return results
}
