package R00

import (
	"path/filepath"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func ex05Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.RepoDirectory, exercise.TurnInDirectory)
	output, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"test"})
	if err != nil {
		logger.Exercise.Printf("%v", err)
		return Exercise.AssertionError("", output)
	}
	return Exercise.Passed("OK")
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("05", "studentcode", "ex05", []string{"src/main.rs", "Cargo.toml"}, "", "", []string{"assert", "assert_eq", "assert_ne", "panic", "print", "println"}, nil, nil, 25, ex05Test)
}
