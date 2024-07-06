package R00

import (
	"path/filepath"

	toml "github.com/42-Short/shortinette/internal/datastructures"
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

func checkCargoTomlContent(exercise Exercise.Exercise, expectedContent map[string]string) bool {
	tomlPath := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory, "Cargo.toml")
	fieldMap, err := toml.ReadToml(tomlPath)
	if err != nil {
		logger.Error.Printf("internal error: %s", err)
		return false
	}
	for key, expectedValue := range expectedContent {
		value, ok := fieldMap[key]
		if !ok {
			logger.File.Printf("[%s KO]: '%s' not found in Cargo.toml", exercise.Name, key)
		} else if value != expectedValue {
			logger.File.Printf("[%s KO]: Cargo.toml content mismatch, expected '%s', got '%s'", exercise.Name, expectedValue, value)
		}
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
	expectedContent := map[string]string{
		"package.name":        "module00-ex04",
		"package.edition":     "2021",
		"package.description": "my answer to the fifth exercise of the first module of 42's Rust Piscine",
	}
	return checkCargoTomlContent(*exercise, expectedContent)
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("EX04", "studentcode", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
