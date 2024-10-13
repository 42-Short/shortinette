//go:build ignore
package R06

import (
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString01 = `
disallowed-methods = ["std::mem::transmute_copy", "std::ptr::read", "std::ptr::read_unaligned", "std::mem::replace", "std::slice::from_raw_parts_mut", "std::mem::size_of_val", "std::mem::align_of_val", "std::slice::from_raw_parts_mut", "std::ptr::copy_nonoverlapping"]
`

func ex01Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString01, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	tests01, err := os.ReadFile("internal/tests/R06/ex01.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(string(tests01), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"src/lib.rs", "Cargo.toml"}, 10, ex01Test)
}
