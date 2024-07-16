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

func setUpEnvironment(repoId string, testDirectory string) error {
	repoLink := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId)
	if err := git.Clone(repoLink, testDirectory); err != nil {
		errorMessage := fmt.Sprintf("failed to clone repository: %v", err)
		return errors.NewInternalError(errors.ErrInternal, errorMessage)
	}
	if err := git.Clone(fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoId), "compile-environment/src/"); err != nil {
		return err
	}
	return nil
}

// Executes the exercises, returns the results and the path to the traces
func (m *Module) Run(repoId string, testDirectory string) (map[string]bool, string) {
	defer func() {
		if err := os.RemoveAll("compile-environment"); err != nil {
			logger.Error.Printf("could not tear down testing environment: %v", err)
		}
		if err := os.RemoveAll("studentcode"); err != nil {
			logger.Error.Printf("could not tear down testing environment: %v", err)
		}
	}()
	err := setUpEnvironment(repoId, testDirectory)
	if err != nil {
		logger.Error.Println(err)
		return nil, ""
	}
	results := make(map[string]bool)
	tracesPath := logger.GetNewTraceFile(repoId)
	if m.Exercises != nil {
		for _, exercise := range m.Exercises {
			command := "docker"
			args := []string{
				"run",
				"-i",
				"--rm",
				"-v",
				"/root/shortinette:/app",
				"testenv",
				"sh",
				"-c",
				fmt.Sprintf("go run . \"%s\" \"%s\" \"%s\"", m.Name, exercise.Name, tracesPath),
			}
			_, err := testutils.RunCommandLine(".", command, args)
			if err != nil {
				results[exercise.Name] = false
			} else {
				results[exercise.Name] = true
			}
		}
	}
	return results, tracesPath
}
