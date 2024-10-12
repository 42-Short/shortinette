package R06

import (
	"os"
	"path/filepath"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/42-Short/shortinette/rust/alloweditems"
)

var clippyTomlAsString00 = `
disallowed-methods = ["std::mem::replace", "std::mem::take", "std::mem::swap", "std::ptr::swap_nonoverlapping", "std::mem::size_of", "std::mem::align_of", "std::ptr::copy", "std::ptr::copy_nonoverlapping", "core::mem::replace", "core::mem::take", "core::mem::swap", "core::ptr::swap_nonoverlapping", "core::mem::size_of", "core::mem::align_of", "core::ptr::copy", "core::ptr::copy_nonoverlapping"]
`

func ex00Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString00, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	tests00, err := os.ReadFile("internal/tests/R06/ex00.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(string(tests00), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "ex00", []string{"src/lib.rs", "Cargo.toml"}, 69, ex00Test)
}
