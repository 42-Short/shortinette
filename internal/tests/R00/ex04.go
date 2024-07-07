package R00

import (
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

func testNmReleaseMode(exercise Exercise.Exercise) bool {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	_, err := testutils.RunCommandLine(workingDirectory, "cargo build --release")
	if err != nil {
		logger.File.Printf("[%s KO]: compilation error %v", exercise.Name, err)
		return false
	}
	output, err := testutils.RunCommandLine(workingDirectory, "nm target/release/module00-ex04")
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error %v", exercise.Name, err)
		return false
	}
	if output != "" {
		logger.File.Println(testutils.AssertionErrorString(exercise.Name, "", output))
		return false
	}
	return true
}

func testCargoRunBinOtherReleaseMode(exercise Exercise.Exercise) bool {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo run --release --bin other")
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error %v", exercise.Name, err)
		return false
	}
	if output != "Hey! I'm the other bin target!\nI'm in release mode!\n" {
		logger.File.Println(testutils.AssertionErrorString(exercise.Name, "Hey! I'm the other bin target!\nI'm in release mode!\n", output))
		return false
	}
	return true
}

func testCargoRunBinOther(exercise Exercise.Exercise) bool {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo run --bin other")
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error %v", exercise.Name, err)
		return false
	}
	if output != "Hey! I'm the other bin target!\n" {
		logger.File.Println(testutils.AssertionErrorString(exercise.Name, "Hey! I'm the other bin target!\n", output))
		return false
	}
	return true
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
	if !testCargoRunBinOther(*exercise) {
		return false
	}
	if !testCargoRunBinOtherReleaseMode(*exercise) {
		return false
	}
	if !testNmReleaseMode(*exercise) {
		return false
	}

	return testutils.CheckCargoTomlContent(*exercise, expectedTomlContent)
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("EX04", "studentcode", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
