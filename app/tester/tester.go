package tester

import (
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
)

type Result struct {
	ExerciseID int
	Passed     bool
	Score      int
	ErrorCode  int
	output     string
}

func failed(err error, exerciseID int) Result {
	var customError *GradingError
	if errors.As(err, &customError) {
		return Result{
			Passed:     false,
			ExerciseID: exerciseID,
			ErrorCode:  customError.code,
			output:     customError.err,
		}
	} else {
		return Result{
			Passed:     false,
			ExerciseID: exerciseID,
			ErrorCode:  InternalError,
			output:     err.Error(),
		}
	}
}

func allowedFilesCheck(exercise config.Exercise, exerciseDirectory string) error {
	infos, err := os.Stat(exerciseDirectory)
	if err != nil {
		if os.IsNotExist(err) {
			return TestingError(NothingTurnedIn, "")
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

func GradeExercise(exercise *config.Exercise, exerciseID int, exerciseDirectory string) Result {
	if err := allowedFilesCheck(*exercise, exerciseDirectory); err != nil {
		return failed(err, exerciseID)
	}

	dockerClient, err := docker.NewClient()
	if err != nil {
		return failed(fmt.Errorf("error connecting to docker socket: %s", err), exerciseID)
	}

	// This should probably be moved to the startup phase in the future (and maybe also be able to trigger through the server/cli)
	if err := docker.BuildImage(dockerClient, nil); err != nil {
		return failed(fmt.Errorf("error building Docker image: %s", err), exerciseID)
	}

	container, err := docker.ContainerCreate(dockerClient, []string{filepath.Join("/root", filepath.Base(exercise.ExecutablePath))})
	if err != nil {
		return failed(fmt.Errorf("error creating Docker container: %s", err), exerciseID)
	}

	if err := container.CopyFilesToContainer(*exercise, exerciseDirectory); err != nil {
		container.Kill()
		return failed(fmt.Errorf("error copying files into container: %s", err), exerciseID)
	}

	if err := container.Exec(time.Second); err != nil {
		container.Kill()
		return failed(err, exerciseID)
	}
	if container.ExitCode == 0 {
		return Result{
			Passed:     true,
			Score:      exercise.Score,
			ExerciseID: exerciseID,
			output:     container.Logs,
		}
	} else {
		return Result{
			Passed:     false,
			Score:      exercise.Score,
			ExerciseID: exerciseID,
			ErrorCode:  RuntimeError,
			output:     container.Logs,
		}
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

	firstFailed := -1
	output := ""
	for i, result := range results {
		if result.Passed {
			output += fmt.Sprintf("Exercise %02d: OK\n", i)
		} else {
			if firstFailed == -1 {
				firstFailed = i
			}
			output += fmt.Sprintf("Exercise %02d: KO\n", i)
		}
	}

	if firstFailed != -1 {
		output += fmt.Sprintf("\n\n=====Trace for Exercise %02d=====\n", firstFailed)
		output += results[firstFailed].output
	}

	if _, err := file.WriteString(output); err != nil {
		return err
	}
	return nil
}

// Clones the specified repo and grades it according to the passed module.
// Returns an error if the start time hasn't been reached, repo cloning or
// upload failed, or if there was an issue writing the logs.
// Returns true if the module was passed (enough points), false if not.
func GradeModule(module config.Module, repo string) (bool, error) {
	if time.Now().Before(module.StartTime) {
		return false, TestingError(EarlyGrading, fmt.Sprintf("start time for repo '%s' not reached yet", repo))
	}

	// This should probably be moved to the startup phase in the future,
	// as there is not really a point in validating this for each grading.
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
			result := GradeExercise(e, i, path.Join(repo, exercise.TurnInDirectory))
			resultsChan <- result
		}(&exercise, i)
	}

	wg.Wait()
	close(resultsChan)

	results := sortResults(module, resultsChan)
	timestamp := time.Now().Local().Format("20060102_150405")
	traceName := fmt.Sprintf("%s-%s.log", repo, timestamp)
	if err := writeTrace(results, traceName); err != nil {
		return false, err
	}
	// TODO: Trace should be uploaded to a different branch (maybe add a parameter to the UploadFiles function)
	if err := git.UploadFiles(repo, "Trace", traceName); err != nil {
		return false, err
	}
	totalPoints, maxPoints := calculateTotalPoints(results)
	// TODO: Release Notes should contain some information (at least a link to the trace file)
	if err := git.NewRelease(repo, fmt.Sprintf("grade-%s", timestamp), fmt.Sprintf("%d/%d", totalPoints, maxPoints), ""); err != nil {
		return false, err
	}

	return totalPoints >= module.MinimumScore, nil
}
