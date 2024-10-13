package R05

import (
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	"github.com/42-Short/shortinette/pkg/logger"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func ex06Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if err := alloweditems.Check(*exercise, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	Ex06TestMod, err := os.ReadFile("internal/tests/R05/ex06.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := testutils.AppendStringToFile(string(Ex06TestMod), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}

	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"test", "--release", "shortinette_rust_test_module05_ex06_0001"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex06() Exercise.Exercise {
	return Exercise.NewExercise("06", "ex06", []string{"src/main.rs", "Cargo.toml"}, 15, ex06Test)
}
