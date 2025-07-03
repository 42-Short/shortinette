package tester

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/logger"
	"github.com/42-Short/shortinette/tester/docker"
	"github.com/bmatcuk/doublestar/v4"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Result struct {
	ExerciseID int
	Passed     bool
	Score      int
	ErrorCode  int
	output     string
}

type GradingResult struct {
	Passed   bool
	Score    int
	MaxScore int
	Trace    string
}

func failed(err error, exerciseID int, exercise *config.Exercise) Result {
	var customError *GradingError
	var errorcode int
	var output string

	if errors.As(err, &customError) {
		errorcode = customError.code
		output = customError.err
	} else {
		errorcode = InternalError
		output = err.Error()
	}

	return Result{
		Passed:     false,
		Score:      exercise.Score,
		ExerciseID: exerciseID,
		ErrorCode:  errorcode,
		output:     output,
	}
}

func allowedFilesCheck(exercise config.Exercise, exerciseDirectory string) error {
	infos, err := os.Stat(exerciseDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return TestingError(NothingTurnedIn, "Nothing turned in")
		} else {
			return TestingError(InternalError, err.Error())
		}
	}

	if !infos.IsDir() {
		return TestingError(InvalidFiles, fmt.Sprintf("'%s' is not a directory", exerciseDirectory))
	}

	submittedFiles := make(map[string]struct{})
	err = filepath.Walk(exerciseDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return TestingError(InternalError, fmt.Sprintf("error while iterating through directory: %s", err))
		}

		if !info.IsDir() && filepath.Base(path)[0] != '.' {
			relativePath, err := filepath.Rel(exerciseDirectory, path)
			if err != nil {
				return TestingError(InternalError, fmt.Sprintf("error removing exerciseDirectory from file path: %s", err))
			}
			submittedFiles[relativePath] = struct{}{}
		}

		return nil
	})

	if err != nil {
		return err
	}

	missingFiles := []string{}
	for _, allowedFilePattern := range exercise.AllowedFiles {
		exists := false
		for submittedFile := range submittedFiles {
			// Pattern was validated by config checker
			// Cannot break loop after this, since it could still match other submitted files
			if doublestar.PathMatchUnvalidated(allowedFilePattern, submittedFile) {
				exists = true
				delete(submittedFiles, submittedFile)
			}
		}

		if !exists {
			missingFiles = append(missingFiles, allowedFilePattern)
		}
	}

	output := ""
	if len(missingFiles) != 0 {
		output += "Missing files: " + strings.Join(missingFiles, ", ")
	}

	var additionalFiles []string
	for key := range submittedFiles {
		additionalFiles = append(additionalFiles, key)
	}

	if len(additionalFiles) != 0 {
		if output != "" {
			output += "; "
		}
		output += "Additional files: " + strings.Join(additionalFiles, ", ")
	}

	if output != "" {
		return TestingError(InvalidFiles, output)
	}

	return nil
}

func GradeExercise(exercise *config.Exercise, module *config.Module, exerciseDirectory string) Result {
	if err := allowedFilesCheck(*exercise, exerciseDirectory); err != nil {
		return failed(err, exercise.ID, exercise)
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		return failed(fmt.Errorf("error connecting to docker socket: %s", err), exercise.ID, exercise)
	}

	env := []string{fmt.Sprintf("MODULE=0%d", module.ID), fmt.Sprintf("EXERCISE=0%d", exercise.ID)}
	containerName := fmt.Sprintf("shortinette-grade-%d-%s", exercise.ID, exerciseDirectory)
	if removeErr := dockerClient.ContainerRemove(context.Background(), containerName, container.RemoveOptions{
		Force: true,
	}); removeErr != nil {
		logger.Warning.Printf("error removing container %s: %s", containerName, removeErr)
	}

	cont, err := docker.ContainerCreate(dockerClient, exercise.DockerImage, containerName, env)
	if err != nil {
		return failed(fmt.Errorf("error creating Docker container: %s", err), exercise.ID, exercise)
	}
	defer func() {
		if killErr := cont.Kill(); killErr != nil {
			logger.Warning.Printf("error killing container %s: %s", cont.ID, killErr)
		}

		ctx := context.Background()

		removeOptions := container.RemoveOptions{
			Force: true,
		}
		if removeErr := dockerClient.ContainerRemove(ctx, cont.ID, removeOptions); removeErr != nil {
			logger.Warning.Printf("error removing container %s: %s", cont.ID, removeErr)
		}
	}()

	if err := cont.CopyFilesToContainer(*exercise, exerciseDirectory); err != nil {
		return failed(fmt.Errorf("error copying files into container: %s", err), exercise.ID, exercise)
	}

	// Hard cap at 5 minutes just in case the test executable doesn't handle endless loops correctly
	if err := cont.Exec(5 * time.Minute); err != nil {
		return failed(err, exercise.ID, exercise)
	}

	passed := false
	var errorcode int

	switch cont.ExitCode {
	case 0:
		errorcode = Passed
		passed = true
	case 1:
		errorcode = Failed
	case 2:
		errorcode = Timeout
	case 3:
		errorcode = CompilationError
	case 4:
		errorcode = ForbiddenFunction
	case 137:
		errorcode = Cancelled

	default:
		errorcode = RuntimeError
	}

	if cont.Timeout {
		errorcode = Timeout
	}

	return Result{
		Passed:     passed,
		Score:      exercise.Score,
		ExerciseID: exercise.ID,
		ErrorCode:  errorcode,
		output:     cont.Logs,
	}
}

func sortResults(module config.Module, resultsChan chan Result) []Result {
	results := make([]Result, len(module.Exercises))

	for result := range resultsChan {
		results[result.ExerciseID] = result
	}
	return results
}

func calculateTotalPoints(results []Result) (int, int) {
	totalPoints := 0
	maxPoints := 0
	exerciseFailed := false

	for _, result := range results {
		maxPoints += result.Score
		if result.Passed && !exerciseFailed {
			totalPoints += result.Score
		} else if !exerciseFailed {
			exerciseFailed = true
		}
	}
	return totalPoints, maxPoints
}

func getTraceContent(results []Result) string {
	var traceIDs []int
	output := ""
	for i, result := range results {
		matchError := func(errorcode int) string {
			switch errorcode {
			case Passed:
				return "OK"
			case Cancelled:
				return "Cancelled"
			case CompilationError:
				return "Compilation Error"
			case EarlyGrading:
				return "Grading time for module hasn't started yet"
			case Failed:
				return "KO"
			case ForbiddenFunction:
				return "Forbidden Function"
			case InternalError:
				return "Internal Error"
			case InvalidFiles:
				return "Invalid Files"
			case NothingTurnedIn:
				return "Nothing turned in"
			case RuntimeError:
				return "KO"
			case Timeout:
				return "Timeout"

			default:
				return "Unknown error"
			}
		}
		output += fmt.Sprintf("Exercise %02d: %s\n", i, matchError(result.ErrorCode))
		if !result.Passed && result.output != "" {
			traceIDs = append(traceIDs, i)
		}
	}

	for _, id := range traceIDs {
		output += fmt.Sprintf("\n\n=====Trace for Exercise %02d=====\n", id)
		output += results[id].output
	}

	return output
}

func checkGradingCancelled(results []Result) bool {
	for _, result := range results {
		if result.ErrorCode == Cancelled {
			return true
		}
	}
	return false
}

// Grades the specified folder according to the passed module.
// Returns an error if the start time hasn't been reached, or if there
// is an issue with the test executables. Returns the reached points,
// as well as the output that should be written in the logs.
func GradeModule(module config.Module, folder string, dockerImage string) (*GradingResult, error) {
	if time.Now().Before(module.StartTime) {
		return nil, TestingError(EarlyGrading, fmt.Sprintf("start time for repo '%s' not reached yet", folder))
	}

	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(module.Exercises))
	for i, exercise := range module.Exercises {
		wg.Add(1)
		go func(e *config.Exercise, exerciseID int) {
			defer wg.Done()
			result := GradeExercise(e, &module, path.Join(folder, exercise.TurnInDirectory))
			resultsChan <- result
		}(&exercise, i)
	}

	wg.Wait()
	close(resultsChan)

	results := sortResults(module, resultsChan)

	if checkGradingCancelled(results) {
		return nil, fmt.Errorf("grading for repo %s was cancelled", folder)
	}

	traceContent := getTraceContent(results)
	totalPoints, maxPoints := calculateTotalPoints(results)

	gradingResult := GradingResult{
		Passed:   totalPoints >= module.MinimumScore,
		Score:    totalPoints,
		MaxScore: maxPoints,
		Trace:    traceContent,
	}

	return &gradingResult, nil
}

// Stops all containers which are used for grading atm
// Returns an error if creating the docker client  or
// listing the containers fails. Sleeps 3 seconds to
// make sure everything afterwards is handled correctly.
func StopAllGradings() error {
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	containers, err := dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return fmt.Errorf("could not get docker containers: %s", err)
	}

	var killedContainers []string
	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, "shortinette-grade-") {
				dockerClient.ContainerKill(ctx, container.ID, "SIGKILL") //nolint:errcheck
				killedContainers = append(killedContainers, container.ID)
				break
			}
		}
	}

	for _, containerID := range killedContainers {
		for {
			_, err := dockerClient.ContainerInspect(ctx, containerID)
			if client.IsErrNotFound(err) {
				break
			}
			time.Sleep(time.Second)
		}
	}
	return nil
}

func HandleSignals(done chan bool, exit bool) {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		StopAllGradings() //nolint:errcheck
		done <- true
		if exit {
			os.Exit(0)
		}
	}()
}
