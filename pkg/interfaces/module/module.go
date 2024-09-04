// Module provides structures and functions for managing and executing modules,
// which consist of multiple exercises. The package includes features for setting up
// environments, running exercises in isolated containers, and grading the results.
package Module

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/42-Short/shortinette/pkg/git"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
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
func NewModule(name string, minimumGrade int, exercises map[string]Exercise.Exercise, subjectPath string) (module Module) {
	return Module{
		Name:         name,
		MinimumGrade: minimumGrade,
		Exercises:    exercises,
		SubjectPath:  subjectPath,
	}
}

// setUpEnvironment sets up the environment by cloning the student's repository and
// preparing the compile environment.
//
//   - repoID: the ID of the repository to be cloned
//   - testDirectory: the directory where the repository will be cloned
//
// Returns an error if the environment setup fails.
func setUpEnvironment(repoID string) (err error) {
	repoLink := fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("GITHUB_ORGANISATION"), repoID)

	if err = git.Clone(repoLink, repoID); err != nil {
		return fmt.Errorf("failed to clone repository: %v", err)
	}
	if _, err = testutils.RunCommandLine(".", "sh", []string{"-c", fmt.Sprintf("chmod -R 777 %s", repoID)}); err != nil {
		return err
	}
	return nil
}

// tearDownEnvironment removes the environment set up for grading, including the compile
// environment and the student's code directory.
//
// Returns an error if the environment teardown fails.
func tearDownEnvironment(repoID string) error {
	if err := os.RemoveAll(repoID); err != nil {
		return fmt.Errorf("remove clone directory: %v", err)
	}
	return nil
}

type GradingConfig struct {
	ModuleName     string
	ExerciseName   string
	TracesPath     string
	CloneDirectory string
}

// runContainerized runs an exercise within a Docker container to prevent running malicious
// code on the host machine.
//
//   - config: GradingConfig object filled with the metadata needed for grading execution
//
// Returns a boolean indicating whether the exercise passed or failed.
func runContainerized(config GradingConfig) bool {
	configJSON, err := json.Marshal(config)
	if err != nil {
		logger.Error.Printf("marshal config: %v", err)
		return false
	}

	ctx := context.Background()
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.Error.Printf("Docker client creation: %v", err)
		return false
	}

	dir, _ := os.Getwd()
	containerConfig := &container.Config{
		Image:      "shortinette-testenv",
		Cmd:        []string{"sh", "-c", fmt.Sprintf("go run . '%s'", string(configJSON))},
		WorkingDir: "/app",
	}
	hostConfig := &container.HostConfig{
		Binds: []string{fmt.Sprintf("%s/traces:/app/traces", dir)},
	}

	response, err := client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, "")
	if err != nil {
		logger.Error.Printf("container creation: %v", err)
		return false
	}

	if err := copyToContainer(ctx, client, response.ID, config.CloneDirectory, "/app"); err != nil {
		logger.Error.Printf("copying files to container: %v", err)
		return false
	}
	if err := copyToContainer(ctx, client, response.ID, "./go.mod", "/app/go.mod"); err != nil {
		logger.Error.Printf("copying files to container: %v", err)
		return false
	}
	if err := copyToContainer(ctx, client, response.ID, "./internal", "/app/internal"); err != nil {
		logger.Error.Printf("copying files to container: %v", err)
		return false
	}
	if err := copyToContainer(ctx, client, response.ID, "./go.sum", "/app/go.sum"); err != nil {
		logger.Error.Printf("copying files to container: %v", err)
		return false
	}
	if err := copyToContainer(ctx, client, response.ID, "./main.go", "/app/main.go"); err != nil {
		logger.Error.Printf("copying files to container: %v", err)
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
	if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, output); err != nil {
		logger.Error.Printf("copying logs: %v", err)
		return false
	}

	inspect, err := client.ContainerInspect(ctx, response.ID)
	if err != nil {
		logger.Error.Printf("inspecting container: %v", err)
		return false
	}

	if inspect.State.ExitCode != 0 {
		logger.Error.Printf("container exited with non-zero status: %d", inspect.State.ExitCode)
		return false
	}

	return true
}

func copyToContainer(ctx context.Context, cli *client.Client, containerID, srcPath, destPath string) error {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)

	err := filepath.Walk(srcPath, func(file string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.Mode().IsRegular() {
			data, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			header := &tar.Header{
				Name:    filepath.ToSlash(file),
				Mode:    int64(fi.Mode().Perm()),
				Size:    fi.Size(),
				ModTime: fi.ModTime(),
			}
			if err := tw.WriteHeader(header); err != nil {
				return err
			}
			if _, err := tw.Write(data); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if err := tw.Close(); err != nil {
		return err
	}

	return cli.CopyToContainer(ctx, containerID, destPath, buf, container.CopyToContainerOptions{})
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
func gradingRoutine(module Module, tracesPath string, repoID string) (results map[string]bool) {
	resultsChannel := make(chan exerciseResult, len(module.Exercises))
	var waitGroup sync.WaitGroup
	results = make(map[string]bool)

	for _, exercise := range module.Exercises {
		waitGroup.Add(1)
		conf := GradingConfig{module.Name, exercise.Name, tracesPath, repoID}
		go func(ex Exercise.Exercise) {
			defer waitGroup.Done()
			result := runContainerized(conf)
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
func (m *Module) Run(repoID string) (results map[string]bool, tracesPath string, err error) {
	defer func() {
		if err := tearDownEnvironment(repoID); err != nil {
			logger.Error.Printf("error tearing down grading environment: %s", err.Error())
		}
	}()
	err = setUpEnvironment(repoID)
	if err != nil {
		return nil, "", fmt.Errorf("grading environment setup: %v", err)
	}
	tracesPath = logger.GetNewTraceFile(repoID)
	if err := logger.InitializeTraceLogger(tracesPath); err != nil {
		return nil, "", err
	}
	if m.Exercises != nil {
		results = gradingRoutine(*m, tracesPath, repoID)
	}
	return results, tracesPath, nil
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
