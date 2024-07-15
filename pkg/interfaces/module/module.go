package Module

import (
	"fmt"
	"os"

	"github.com/42-Short/shortinette/internal/errors"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/pkg/git"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

type Module struct {
	Name      string
	Exercises map[string]Exercise.Exercise
}

// NewModule initializes and returns a Module struct
//
//   - name: module display name
//   - exercises: list of all Exercise.Exercise objects belonging into the module
func NewModule(name string, exercises map[string]Exercise.Exercise) (Module, error) {
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

// Executes the exercises, returns the results and the path to the traces
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
			commandLine := fmt.Sprintf("docker run -i --rm -v /root/shortinette:/app test-env sh -c 'go run . %s %s test'", m.Name, exercise.Name)
			output, err := testutils.RunCommandLine(".", commandLine)
			if err != nil {
				logger.Error.Printf("error running containerized test: %v: %s", err, output)
			}
		}
	}
	return results, tracesPath
}
