package tester

import (
	"os"
	"testing"
	"time"

	"github.com/42-Short/shortinette/config"
)

func TestGradeModuleBeforeStarttime(t *testing.T) {
	startTime := time.Now().Add(time.Hour)
	module := config.Module{
		StartTime: startTime,
	}
	_, err := GradeModule(module, "repo")
	if err == nil || !matchesCustomError(err, EarlyGrading) {
		t.Fatalf("Grading before starttime shouldn't be possible")
	}
}

func TestGradeModuleAfterStarttime(t *testing.T) {
	startTime := time.Now().Add(-1 * time.Hour)
	module := config.Module{
		StartTime: startTime,
	}
	_, err := GradeModule(module, "repo")
	if err != nil && matchesCustomError(err, EarlyGrading) {
		t.Fatalf("Grading after starttime should be possible")
	}
}

func TestGradeExerciseOk(t *testing.T) {
	if err := os.Mkdir("test", 0755); err != nil {
		t.Fatalf("Unable to create test folder: %s", err)
	}
	if _, err := os.Create("test/test.rs"); err != nil {
		if err2 := os.Remove("test"); err2 != nil {
			t.Fatalf("unable to create test/test.rs and unable to remove test/ folder")
		}
		t.Fatalf("unable to create test/test.rs file")
	}
	exercise := config.Exercise{
		ExecutablePath:  "executables/testexecutable.sh",
		AllowedFiles:    []string{"test.rs"},
		TurnInDirectory: "test",
	}
	result := GradeExercise(&exercise, 0, "test")
	if err := os.RemoveAll("test"); err != nil {
		t.Fatalf("unable to remove test/ directory: %s", err)
	}

	if !result.Passed {
		t.Fatalf("Not passed: %v", result)
	}
	t.Log(result.output)
}

func TestGradeExerciseFail(t *testing.T) {
	if err := os.Mkdir("test", 0755); err != nil {
		t.Fatalf("Unable to create test folder: %s", err)
	}
	if _, err := os.Create("test/test.rs"); err != nil {
		if err2 := os.Remove("test"); err2 != nil {
			t.Fatalf("unable to create test/test.rs and unable to remove test/ folder")
		}
		t.Fatalf("unable to create test/test.rs file")
	}
	exercise := config.Exercise{
		ExecutablePath:  "executables/testexecutable_fail.sh",
		AllowedFiles:    []string{"test/test.rs"},
		TurnInDirectory: "test",
	}
	result := GradeExercise(&exercise, 0, "test")
	if err := os.RemoveAll("test"); err != nil {
		t.Fatalf("unable to remove test/ directory: %s", err)
	}

	if result.Passed {
		t.Fatalf("Exercise passed but shouldn't: %v", result)
	}
}

func TestGradeExerciseNoPermission(t *testing.T) {
	if err := os.Mkdir("test", 0755); err != nil {
		t.Fatalf("Unable to create test folder: %s", err)
	}
	if _, err := os.Create("test/test.rs"); err != nil {
		if err2 := os.Remove("test"); err2 != nil {
			t.Fatalf("unable to create test/test.rs and unable to remove test/ folder")
		}
		t.Fatalf("unable to create test/test.rs file")
	}
	exercise := config.Exercise{
		ExecutablePath:  "executables/testexecutable_noperm.sh",
		AllowedFiles:    []string{"test/test.rs"},
		TurnInDirectory: "test",
	}
	result := GradeExercise(&exercise, 0, "test")
	if err := os.RemoveAll("test"); err != nil {
		t.Fatalf("unable to remove test/ directory: %s", err)
	}

	if result.Passed {
		t.Fatalf("Exercise passed but shouldn't: %v", result)
	}
}

func TestGradeModule(t *testing.T) {
	// TODO: Create repo with some files
	exercises := make([]config.Exercise, 3)
	exercises[0] = config.Exercise{
		ExecutablePath:  "executables/slow_executable.sh",
		Score:           10,
		AllowedFiles:    []string{"test.rs"},
		TurnInDirectory: "ex00",
	}
	exercises[1] = config.Exercise{
		ExecutablePath:  "executables/testexecutable_fail.sh",
		Score:           10,
		AllowedFiles:    []string{"test.rs"},
		TurnInDirectory: "ex00",
	}
	exercises[2] = config.Exercise{
		ExecutablePath:  "executables/testexecutable.sh",
		Score:           10,
		AllowedFiles:    []string{"test.rs"},
		TurnInDirectory: "ex00",
	}

	module := config.Module{
		Exercises:    exercises,
		MinimumScore: 20,
		StartTime:    time.Now(),
	}
	passed, err := GradeModule(module, "test2")
	if err2 := os.RemoveAll("test2"); err2 != nil {
		t.Fatalf("unable to remove test2/ directory: %s", err2)
	}

	// TODO: Delete repo again

	if err != nil {
		t.Fatal(err)
	}

	if passed {
		t.Fatalf("Module passed but shouldn't: %v", module)
	}
}

/*
Missing Tests rn:
- Allowed files check
	- one correct test
	- one test with a missing file
	- two tests with an additional file
		- one of them a regular file
		- one of them a hidden file like .gitignore which shouldn't be failed
	- missing exercise folder to simulate "Nothing turned in"
	- existing but empty exercise folder
	- tests with existing files which don't have permission (folder needs r+w ig, files only read)
- Test executables missing or no permission
- Docker errors (Docker not implemented yet)
*/