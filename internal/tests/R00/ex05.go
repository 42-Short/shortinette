package R00

import (
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

func ex05Test(exercise *Exercise.Exercise) bool {
	if !testutils.TurnInFilesCheck(*exercise) {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)
	return true
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("EX05", "studentcode", "ex05", []string{"src/main.rs", "Cargo.toml"}, "", "", []string{"assert", "assert_eq", "assert_ne", "panic", "print", "println"}, nil, nil, ex05Test)
}
