package R06

import (
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString06 = `
disallowed-types = ["std::collections::Vec", "std::collections::VecDeque", "std::collections::LinkedList", "Box<T>", "Rc<T>", "Arc<T>", "std::cell::RefCell", "std::sync::Mutex", "std::mem::ManuallyDrop", "std::alloc::System"]
disallowed-methods = ["std::ptr::copy", "std::ptr::write", "std::ptr::replace", "std::mem::transmute"]
`

func ex06Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString06, nil, "#![allow(clippy::needless_borrows_for_generic_args)]"); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	tests06, err := os.ReadFile("internal/tests/R06/ex06.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(string(tests06), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex06() Exercise.Exercise {
	return Exercise.NewExercise("06", "ex06", []string{"src/lib.rs", "Cargo.toml", "awesome.c", "build.rs"}, 15, ex06Test)
}
