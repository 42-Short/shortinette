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

var clippyTomlAsString05 = `
disallowed-types = ["std::collections::Vec", "std::collections::VecDeque", "std::collections::LinkedList", "Box<T>", "Rc<T>", "Arc<T>", "std::cell::RefCell", "std::sync::Mutex", "std::mem::ManuallyDrop"]
disallowed-methods = ["std::slice::from_raw_parts", "std::slice::from_raw_parts_mut"]
`

func ex05Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString05, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	tests05, err := os.ReadFile("internal/tests/R06/ex05.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(string(tests05), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex05() Exercise.Exercise {
	return Exercise.NewExercise("05", "ex05", []string{"src/lib.rs", "Cargo.toml"}, 15, ex05Test)
}
