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

var clippyTomlAsString03 = `
disallowed-types = ["std::cell::Cell", "std::cell::RefCell", "std::sync::Mutex", "std::mem::ManuallyDrop", "std::rc::Rc", "std::sync::Arc", "std::sync::RwLock"]
disallowed-methods = ["std::ptr::null", "std::ptr::null_mut", "std::mem::transmute"]
`

func ex03Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString03, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	tests03, err := os.ReadFile("internal/tests/R06/ex03.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(string(tests03), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "ex03", []string{"src/lib.rs", "Cargo.toml"}, 10, ex03Test)
}
