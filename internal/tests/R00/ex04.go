package R00

import (
	"fmt"
	"path/filepath"

	toml "github.com/42-Short/shortinette/internal/datastructures"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var expectedTomlContent = map[string]string{
	"package.name":        "module00-ex04",
	"package.edition":     "2021",
	"package.description": "my answer to the fifth exercise of the first module of 42's Rust Piscine",
}

func testNmReleaseMode(exercise Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"build", "--release"})
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("compilation error: %v", err)}
	}
	output, err := testutils.RunCommandLine(workingDirectory, "nm", []string{"target/release/module00-ex04"})
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: nm did not execute as expected: %v", err)}
	}
	if output != "" {
		return Exercise.AssertionError("", output)
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func testCargoRunBinOtherReleaseMode(exercise Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"run", "--release", "--bin", "other"})
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: cargo run: %s", err)}
	}
	if output != "Hey! I'm the other bin target!\nI'm in release mode!\n" {
		return Exercise.AssertionError("Hey! I'm the other bin target!\nI'm in release mode!\n", output)
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func testCargoRunBinOther(exercise Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"run", "--bin", "other"})
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: %v", err)}
	}
	if output != "Hey! I'm the other bin target!\n" {
		return Exercise.AssertionError("Hey! I'm the other bin target!\n", output)
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func testCargoRun(exercise Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"run"})
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: %v", err)}
	}
	if output != "Hello, Cargo!\n" {
		return Exercise.AssertionError("Hello, Cargo!\n", output)
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
	if result := testCargoRun(*exercise); !result.Passed {
		return result
	}
	if result := testCargoRunBinOther(*exercise); !result.Passed {
		return result
	}
	if result := testCargoRunBinOtherReleaseMode(*exercise); !result.Passed {
		return result
	}
	if result := testNmReleaseMode(*exercise); !result.Passed {
		return result
	}

	return toml.CheckCargoTomlContent(*exercise, expectedTomlContent)
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "studentcode", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, 25, ex04Test)
}
