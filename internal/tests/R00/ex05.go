package R00

import (
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func ex05Test(exercise *Exercise.Exercise) Exercise.Result {
	if !testutils.TurnInFilesCheck(*exercise) {
		return Exercise.Result{Passed: false, Output: "invalid files found in turn in directory"}
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)
	return Exercise.Result{Passed: true, Output: ""}
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("05", "studentcode", "ex05", []string{"src/main.rs", "Cargo.toml"}, "", "", []string{"assert", "assert_eq", "assert_ne", "panic", "print", "println"}, nil, nil, ex05Test)
}
