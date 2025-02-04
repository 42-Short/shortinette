package tester

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/42-Short/shortinette/config"
	"github.com/42-Short/shortinette/tester/docker"
)

func pullDebianImage() error {
	dockerClient, err := docker.NewClient()
	if err != nil {
		return err
	}

	if err := docker.PullImage(dockerClient, "42short/rust"); err != nil {
		return err
	}
	return nil
}

func wrapSignalHandlerFunction(testFunc func()) {
	done := make(chan bool, 1)

	HandleSignals(done, false)
	testFunc()
	done <- true

	<-done
}

func TestGradeModuleBeforeStarttime(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		startTime := time.Now().Add(time.Hour)
		module := config.Module{
			StartTime: startTime,
		}
		time.Sleep(5 * time.Second)
		_, err := GradeModule(module, "repo", "shortinette-testenv")
		if err == nil || !matchesCustomError(err, EarlyGrading) {
			t.Fatalf("Grading before starttime shouldn't be possible")
		}
	})

}

func TestGradeModuleAfterStarttime(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		if err := pullDebianImage(); err != nil {
			t.Fatal(err)
		}

		startTime := time.Now().Add(-1 * time.Hour)
		module := config.Module{
			StartTime: startTime,
		}
		_, err := GradeModule(module, "repo", "42short/rust")
		if err != nil && matchesCustomError(err, EarlyGrading) {
			t.Fatalf("Grading after starttime should be possible")
		}
	})
}

func TestGradeExerciseFail(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		if err := pullDebianImage(); err != nil {
			t.Fatal(err)
		}

		if err := os.Mkdir("test", 0755); err != nil {
			t.Fatalf("Unable to create test folder: %s", err)
		}
		defer os.RemoveAll("test")

		if _, err := os.Create("test/test.rs"); err != nil {
			t.Fatalf("unable to create test/test.rs file")
		}

		exercise := config.Exercise{
			DockerImage:     "42short/rust",
			AllowedFiles:    []string{"test/test.rs"},
			TurnInDirectory: "test",
		}

		result := GradeExercise(&exercise, &config.Module{
			ID: 0,
		}, "test")

		if result.Passed {
			t.Fatalf("Exercise passed but shouldn't: %v", result)
		}
	})
}

func TestGradeExerciseNoPermission(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		if err := pullDebianImage(); err != nil {
			t.Fatal(err)
		}

		if err := os.Mkdir("test", 0755); err != nil {
			t.Fatalf("Unable to create test folder: %s", err)
		}
		defer os.RemoveAll("test")

		if _, err := os.Create("test/test.rs"); err != nil {
			t.Fatalf("unable to create test/test.rs file")
		}
		exercise := config.Exercise{
			DockerImage:     "42short/rust",
			AllowedFiles:    []string{"test/test.rs"},
			TurnInDirectory: "test",
		}

		result := GradeExercise(&exercise, &config.Module{
			ID: 0,
		}, "test")

		if result.Passed {
			t.Fatalf("Exercise passed but shouldn't: %v", result)
		}
	})
}


func TestGradeModuleMissingFile(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		if err := pullDebianImage(); err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll("testrepo/ex00", 0755); err != nil {
			t.Fatalf("Unable to create testrepo folder: %s", err)
		}
		defer os.RemoveAll("testrepo")

		if _, err := os.Create("testrepo/ex00/test.rs"); err != nil {
			t.Fatalf("unable to create testrepo/ex00/test.rs file")
		}

		if _, err := os.Create("testrepo/ex00/.gitignore"); err != nil {
			t.Fatalf("unable to create testrepo/ex00/test.rs file")
		}

		exercises := make([]config.Exercise, 1)
		exercises[0] = config.Exercise{
			// ExecutablePath:  "executables/testexecutable.sh",
			Score:           10,
			AllowedFiles:    []string{"test.rs", "test2.rs"},
			TurnInDirectory: "ex00",
		}

		module := config.Module{
			Exercises:    exercises,
			MinimumScore: 10,
			StartTime:    time.Now(),
		}
		result, err := GradeModule(module, "testrepo", "42short/rust")

		if err != nil {
			t.Fatal(err)
		}

		if result.Passed {
			t.Fatalf("Module passed but shouldn't: %v", module)
		}

		if !strings.Contains(result.Trace, "Missing") {
			t.Fatalf("Expected missing file error in trace: %s", result.Trace)
		}
	})
}

func TestGradeModuleAdditionalFiles(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		if err := pullDebianImage(); err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll("testrepo/ex00", 0755); err != nil {
			t.Fatalf("Unable to create testrepo folder: %s", err)
		}
		defer os.RemoveAll("testrepo")

		if _, err := os.Create("testrepo/ex00/test.rs"); err != nil {
			t.Fatalf("unable to create testrepo/ex00/test.rs file")
		}

		if _, err := os.Create("testrepo/ex00/.gitignore"); err != nil {
			t.Fatalf("unable to create testrepo/ex00/test.rs file")
		}

		exercises := make([]config.Exercise, 1)
		exercises[0] = config.Exercise{
			// ExecutablePath:  "executables/testexecutable.sh",
			Score:           10,
			AllowedFiles:    []string{},
			TurnInDirectory: "ex00",
		}

		module := config.Module{
			Exercises:    exercises,
			MinimumScore: 10,
			StartTime:    time.Now(),
		}
		result, err := GradeModule(module, "testrepo", "42short/rust")

		if err != nil {
			t.Fatal(err)
		}

		if result.Passed {
			t.Fatalf("Module passed but shouldn't: %v", module)
		}

		if !strings.Contains(result.Trace, "Additional") {
			t.Fatalf("Expected additional file error in trace: %s", result.Trace)
		}
	})
}

func TestGradeModuleNothingTurnedIn(t *testing.T) {
	wrapSignalHandlerFunction(func() {
		if err := pullDebianImage(); err != nil {
			t.Fatal(err)
		}

		if err := os.MkdirAll("testrepo", 0755); err != nil {
			t.Fatalf("Unable to create testrepo folder: %s", err)
		}
		defer os.RemoveAll("testrepo")

		exercises := make([]config.Exercise, 1)
		exercises[0] = config.Exercise{
			// ExecutablePath:  "executables/testexecutable.sh",
			Score:           10,
			AllowedFiles:    []string{},
			TurnInDirectory: "ex00",
		}

		module := config.Module{
			Exercises:    exercises,
			MinimumScore: 10,
			StartTime:    time.Now(),
		}
		result, err := GradeModule(module, "testrepo", "42short/rust")

		if err != nil {
			t.Fatal(err)
		}

		if result.Passed {
			t.Fatalf("Module passed but shouldn't: %v", module)
		}

		if !strings.Contains(result.Trace, "Nothing") {
			t.Fatalf("Expected Nothing turned in error in trace: %s", result.Trace)
		}
	})
}
