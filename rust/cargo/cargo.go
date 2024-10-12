package cargo

import (
	"path/filepath"
	"time"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func CargoTest(exercise *Exercise.Exercise, timeout time.Duration, flags []string) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", append([]string{"test", "--no-run"}, flags...)); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	_, err := testutils.RunCommandLine(workingDirectory, "cargo", append([]string{"test"}, flags...), testutils.WithTimeout(timeout))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	_, err = testutils.RunCommandLine(workingDirectory, "cargo", append([]string{"valgrind", "test"}, flags...))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}
