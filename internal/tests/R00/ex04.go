package R00

import (
	"path/filepath"

	toml "github.com/42-Short/shortinette/internal/datastructures"
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

func checkCargoTomlContent(exercise Exercise.Exercise) bool {
	tomlPath := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory, "Cargo.toml")
	tomlConfig, err := toml.ReadToml(tomlPath)
	if err != nil {
		logger.Error.Printf("internal error: %s", err)
	}
	if tomlConfig.Package.Name != "module00-ex04" {
		logger.File.Printf("[%s KO]: Cargo.toml content mismatch, expected '%s', got '%s'", exercise.Name, "module00-ex04", tomlConfig.Package.Name)
		return false
	}
	if tomlConfig.Package.Edition != "2021" {
		return false
	}
	if tomlConfig.Package.Description != "my answer to the fifth exercise of the first module of 42's Rust Piscine" {
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
	if !checkCargoTomlContent(*exercise) {
		return false
	} else {
		return true
	}
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("EX04", "studentcode", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
