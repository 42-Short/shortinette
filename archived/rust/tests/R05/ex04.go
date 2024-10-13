package R05

import (
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	Ex04TestMod, err := os.ReadFile("internal/tests/R05/ex04.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := testutils.AppendStringToFile(string(Ex04TestMod), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"nextest", "run", "--release", "shortinette_rust_test_module05_ex04_0001"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "ex04", []string{"src/lib.rs", "Cargo.toml"}, 10, ex04Test)
}
