package tester

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/git"
	"github.com/42-Short/shortinette/tester/docker"
	"github.com/docker/docker/api/types/container"
)

type Result struct {
	ExerciseID int
	Passed     bool
	Score      int
	ErrorCode  int
	output     string
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
	for _, allowedFile := range exercise.AllowedFiles {
		if _, exists := submittedFiles[allowedFile]; !exists {
			missingFiles = append(missingFiles, allowedFile)
		} else {
			delete(submittedFiles, allowedFile)
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

func GradeExercise(exercise *config.Exercise, exerciseID int, exerciseDirectory string, dockerImage string) Result {
	if err := allowedFilesCheck(*exercise, exerciseDirectory); err != nil {
		return failed(err, exerciseID, exercise)
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		return failed(fmt.Errorf("error connecting to docker socket: %s", err), exerciseID, exercise)
	}

	container, err := docker.ContainerCreate(dockerClient, []string{filepath.Join("/root", filepath.Base(exercise.ExecutablePath))}, dockerImage, fmt.Sprintf("shortinette-grade-%d-%s", exerciseID, exerciseDirectory))
	if err != nil {
		return failed(fmt.Errorf("error creating Docker container: %s", err), exerciseID, exercise)
	}
	defer container.Kill()

	if err := container.CopyFilesToContainer(*exercise, exerciseDirectory); err != nil {
		return failed(fmt.Errorf("error copying files into container: %s", err), exerciseID, exercise)
	}

	// Hard cap at 5 minutes just in case the test executable doesn't handle endless loops correctly
	if err := container.Exec(5 * time.Minute); err != nil {
		return failed(err, exerciseID, exercise)
	}

	var passed bool
	var errorcode int

	if container.ExitCode == 0 {
		errorcode = Passed
		passed = true
	} else if container.Timeout {
		errorcode = Timeout
		passed = false
	} else if container.ExitCode == 137 {
		errorcode = Cancelled
		passed = false
	} else {
		errorcode = RuntimeError
		passed = false
	}

	return Result{
		Passed:     passed,
		Score:      exercise.Score,
		ExerciseID: exerciseID,
		ErrorCode:  errorcode,
		output:     container.Logs,
	}
}

func checkTestExecutable(executable string) error {
	infos, err := os.Stat(executable)
	if err != nil {
		return fmt.Errorf("error getting infos of test executable: '%s'", err)
	}

	mode := infos.Mode()
	if mode.IsDir() || mode.Perm()&0111 == 0 {
		return fmt.Errorf("test executable '%s' isn't an executable file", executable)
	}

	return nil
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

func writeTrace(results []Result, filename string) error {
	if _, err := os.Stat(filename); err == nil {
		return TestingError(InternalError, fmt.Sprintf("trace file %s already exists", filename))
	}
	file, err := os.Create(filename)
	if err != nil {
		return TestingError(InternalError, fmt.Sprintf("error creating trace file: %s", err.Error()))
	}
	defer file.Close()

	var traceIDs []int
	output := ""
	for i, result := range results {
		matchError := func(errorcode int) string {
			switch errorcode {
			case Passed:
				return "OK"
			case RuntimeError:
				return "KO"
			case InternalError:
				return "Internal Error"
			case EarlyGrading:
				return "Grading time for module hasn't started yet"
			case InvalidFiles:
				return "Invalid Files"
			case NothingTurnedIn:
				return "Nothing turned in"
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

	if _, err := file.WriteString(output); err != nil {
		return err
	}
	return nil
}

//nolint:errcheck
func deleteRepo(repo string) error {
	return os.RemoveAll(repo)
}

func checkGradingCancelled(results []Result) bool {
	for _, result := range results {
		if result.ErrorCode == Cancelled {
			return true
		}
	}
	return false
}

// Clones the specified repo and grades it according to the passed module.
// Returns an error if the start time hasn't been reached, repo cloning or
// upload failed, or if there was an issue writing the logs.
// Returns true if the module was passed (enough points), false if not.
func GradeModule(module config.Module, repo string, dockerImage string) (bool, error) {
	if time.Now().Before(module.StartTime) {
		return false, TestingError(EarlyGrading, fmt.Sprintf("start time for repo '%s' not reached yet", repo))
	}

	for _, exercise := range module.Exercises {
		if err := checkTestExecutable(exercise.ExecutablePath); err != nil {
			return false, err
		}
	}

	if err := git.Clone(repo); err != nil {
		return false, err
	}

	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(module.Exercises))
	for i, exercise := range module.Exercises {
		wg.Add(1)
		go func(e *config.Exercise, exerciseID int) {
			defer wg.Done()
			result := GradeExercise(e, i, path.Join(repo, exercise.TurnInDirectory), dockerImage)
			resultsChan <- result
		}(&exercise, i)
	}

	wg.Wait()
	close(resultsChan)

	results := sortResults(module, resultsChan)
	defer deleteRepo(repo)

	if checkGradingCancelled(results) {
		return false, fmt.Errorf("grading for repo %s was cancelled", repo)
	}
	timestamp := time.Now().Local().Format("20060102_150405")
	traceName := fmt.Sprintf("%s-%s.log", repo, timestamp)
	if err := writeTrace(results, traceName); err != nil {
		return false, err
	}

	if err := git.UploadFiles(repo, "Trace", "traces", false, traceName); err != nil {
		return false, err
	}

	if err := os.Remove(traceName); err != nil {
		return false, err
	}
	totalPoints, maxPoints := calculateTotalPoints(results)
	// TODO: Release Notes should contain some information (at least a link to the trace file)
	if err := git.NewRelease(repo, fmt.Sprintf("grade-%s", timestamp), fmt.Sprintf("%d/%d", totalPoints, maxPoints), ""); err != nil {
		return false, err
	}

	return totalPoints >= module.MinimumScore, nil
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

	for _, container := range containers {
		for _, name := range container.Names {
			if strings.Contains(name, "shortinette-grade-") {
				dockerClient.ContainerKill(ctx, container.ID, "SIGKILL") //nolint:errcheck
				break
			}
		}
	}
	time.Sleep(3 * time.Second)
	return nil
}
