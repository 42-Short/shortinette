package R05

import (
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func ex01Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	Ex01TestMod, err := os.ReadFile("internal/tests/R05/ex01.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := testutils.AppendStringToFile(string(Ex01TestMod), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module05_ex01_0001"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"src/main.rs", "Cargo.toml"}, 10, ex01Test)
}
