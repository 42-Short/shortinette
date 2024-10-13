package R06

import (
	"os"
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString04 = `
disallowed-types = ["std::fs::File", "std::fs::OpenOptions", "std::io::Read", "std::io::Write", "std::mem::ManuallyDrop", "Box<T>", "Rc<T>", "Arc<T>"]
disallowed-methods = ["std::process::exit"]
`

func ex04Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString04, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	tests04, err := os.ReadFile("internal/tests/R06/ex04.rs")
	if err != nil {
		return Exercise.InternalError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(string(tests04), exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex04() Exercise.Exercise {
	return Exercise.NewExercise("04", "ex04", []string{"src/lib.rs", "Cargo.toml"}, 20, ex04Test)
}
