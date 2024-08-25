// Package Module provides structures and functions for managing and executing modules,
// which consist of multiple exercises. The package includes features for setting up
// environments, running exercises in isolated containers, and grading the results.
package Module

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/42-Short/shortinette/pkg/git"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

// Module represents a module containing multiple exercises. It includes the module's
// name, the minimum grade required to pass, a map of exercises, and the path to the
// subject file.
type Module struct {
	Name         string                       // Name is the module's display name.
	MinimumGrade int                          // MinimumGrade is the minimum score required to pass the module.
	Exercises    map[string]Exercise.Exercise // Exercises is a map of all exercises belonging to the module.
	SubjectPath  string                       // SubjectPath is the path to the module's subject file.
}

// NewModule initializes and returns a Module struct.
//
//   - name: module display name
//   - minimumGrade: the minimum score required to pass the module
//   - exercises: map of all Exercise.Exercise objects belonging to the module
//   - subjectPath: path to the module's subject file
func NewModule(name string, minimumGrade int, exercises map[string]Exercise.Exercise, subjectPath string) (Module, error) {
	return Module{
		Name:         name,
		MinimumGrade: minimumGrade,
		Exercises:    exercises,
		SubjectPath:  subjectPath,
	}, nil
}

// setUpEnvironment sets up the environment by cloning the student's repository and
// preparing the compile environment.
//
//   - repoID: the ID of the repository to be cloned
//   - testDirectory: the directory where the repository will be cloned
//
// Returns an error if the environment setup fails.
func setUpEnvironment(repoID string, testDirectory string) error {
	repoLink := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoID)
	if err := git.Clone(repoLink, testDirectory); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	return nil
}

// tearDownEnvironment removes the environment set up for grading, including the compile
// environment and the student's code directory.
//
// Returns an error if the environment teardown fails.
func tearDownEnvironment() error {
	if err := os.RemoveAll("studentcode"); err != nil {
		return fmt.Errorf("failed to tear down code directory: %v", err)
	}
	return nil
}

// runContainerized runs an exercise within a Docker container to prevent running malicious
// code on the host machine.
//
//   - module: the Module containing the exercise
//   - exercise: the Exercise to be run
//   - tracesPath: the path to store the trace logs
//
// Returns a boolean indicating whether the exercise passed or failed.
func runContainerized(module Module, exercise Exercise.Exercise, tracesPath string) bool {
	ctx := context.Background()
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		logger.Error.Printf("Docker client creation: %v", err)
		return false
	}

	dir, _ := os.Getwd()
	config := &container.Config{
		Image: "shortinette-testenv",
		Cmd:   []string{"sh", "-c", fmt.Sprintf("go run . \"%s\" \"%s\" \"%s\"", module.Name, exercise.Name, tracesPath)},
	}
	hostConfig := &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s:/app", dir)},
	}

	response, err := client.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		logger.Error.Printf("container creation: %v", err)
		return false
	}

	if err := client.ContainerStart(ctx, response.ID, container.StartOptions{}); err != nil {
		logger.Error.Printf("container startup: %v", err)
		return false
	}

	statusChannel, errorChannel := client.ContainerWait(ctx, response.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errorChannel:
		if err != nil {
			logger.Error.Printf("waiting for container: %v", err)
			return false
		}
	case <-statusChannel:
	}
	output, err := client.ContainerLogs(ctx, response.ID, container.LogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		logger.Error.Printf("fetching container logs: %v", err)
		return false
	}
	stdcopy.StdCopy(os.Stdout, os.Stderr, output)
	return true
}

// exerciseResult represents the result of an individual exercise run, including the
// exercise's name and whether it passed or failed.
type exerciseResult struct {
	name   string // name is the name of the exercise.
	result bool   // result indicates whether the exercise passed or failed.
}

// gradingRoutine runs all exercises in the module concurrently within Docker containers
// and collects the results.
//
//   - module: the Module containing the exercises
//   - tracesPath: the path to store the trace logs
//
// Returns a map of exercise names to their pass/fail results.
func gradingRoutine(module Module, tracesPath string) (results map[string]bool) {
	resultsChannel := make(chan exerciseResult, len(module.Exercises))
	var waitGroup sync.WaitGroup
	results = make(map[string]bool)

	for _, exercise := range module.Exercises {
		waitGroup.Add(1)
		go func(ex Exercise.Exercise) {
			defer waitGroup.Done()
			result := runContainerized(module, ex, tracesPath)
			resultsChannel <- exerciseResult{name: ex.Name, result: result}
		}(exercise)
	}
	go func() {
		waitGroup.Wait()
		close(resultsChannel)
	}()
	for result := range resultsChannel {
		results[result.name] = result.result
	}
	return results
}

// Run executes the exercises, spawning a Docker container for each of them to prevent running
// malicious code on your machine. It returns the results and the path to the trace logs.
//
//   - repoID: the ID of the repository to be cloned
//   - testDirectory: the directory where the repository will be cloned
//
// Returns a map of exercise names to their pass/fail results and the path to the trace logs.
func (m *Module) Run(repoID string, testDirectory string) (results map[string]bool, tracesPath string) {
	defer func() {
		if err := tearDownEnvironment(); err != nil {
			logger.Error.Printf("error tearing down grading environment: %s", err.Error())
		}
	}()
	err := setUpEnvironment(repoID, testDirectory)
	if err != nil {
		logger.Error.Println(err)
		return nil, ""
	}
	tracesPath = logger.GetNewTraceFile(repoID)
	if m.Exercises != nil {
		results = gradingRoutine(*m, tracesPath)
	}
	return results, tracesPath
}

// GetScore calculates the total score based on the results of the exercises and determines
// if the module is passed.
//
//   - results: a map of exercise names to their pass/fail results
//
// Returns the total score and a boolean indicating whether the module is passed.
func (m *Module) GetScore(results map[string]bool) (score int, passed bool) {
	var keys []string
	for key := range results {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	score = 0

	for _, key := range keys {
		if !results[key] {
			break
		} else {
			score += m.Exercises[key].Score
		}
	}
	passed = score >= m.MinimumGrade

	return score, passed
}
