package R00

import (
	"fmt"
	"path/filepath"

	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

var expectedTomlContent = map[string]string{
	"package.name":        "module00-ex04",
	"package.edition":     "2021",
	"package.description": "my answer to the fifth exercise of the first module of 42's Rust Piscine",
}

func testCargoRun(exercise Exercise.Exercise) bool {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo run")
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error %v", exercise.Name, err)
		return false
	}
	if output != "Hello, cargo!\n" {
		logger.File.Println(testutils.AssertionErrorString(exercise.Name, "Hello, cargo!", output))
		return false
	}
	return true
}

func ex04Test(exercise *Exercise.Exercise) bool {
	if !testutils.TurnInFilesCheck(*exercise) {
		return false
	}
	if err := testutils.ForbiddenItemsCheck(*exercise, "shortinette-test-R00"); err != nil {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)

	if !testCargoRun(*exercise) {
		return false
	}

	return testutils.CheckCargoTomlContent(*exercise, expectedTomlContent)
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("EX04", "studentcode", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
