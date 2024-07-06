package R00

import (
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

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
	return testutils.CheckCargoTomlContent(*exercise, expectedContent)
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("EX04", "studentcode", "ex04", []string{"src/main.rs", "src/overflow.rs", "src/other.rs", "Cargo.toml"}, "", "", []string{"println"}, nil, nil, ex04Test)
}
