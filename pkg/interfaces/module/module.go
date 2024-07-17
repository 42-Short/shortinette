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

func tearDownEnvironment() error {
	if err := os.RemoveAll("compile-environment"); err != nil {
		return fmt.Errorf("failed to tear down compiling environment: %v", err)
	}
	if err := os.RemoveAll("studentcode"); err != nil {
		return fmt.Errorf("failed to tear down code directory: %v", err)
	}
	return nil
}

func runContainerized(module Module, exercise Exercise.Exercise, tracesPath string) bool {
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
		fmt.Sprintf("go run . \"%s\" \"%s\" \"%s\"", module.Name, exercise.Name, tracesPath),
	}
	_, err := testutils.RunCommandLine(".", command, args)
	return err == nil
}

// Executes the exercises, spawning a Docker container for each of them to prevent running
// malicious code on your machine, returns the results and the path to the traces
func (m *Module) Run(repoId string, testDirectory string) (map[string]bool, string) {
	defer func() {
		if err := tearDownEnvironment(); err != nil {
			logger.Error.Printf(err.Error())
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
			results[exercise.Name] = runContainerized(*m, exercise, tracesPath)
		}
	}
	return results, tracesPath
}
