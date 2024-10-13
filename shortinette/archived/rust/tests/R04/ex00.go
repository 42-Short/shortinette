//go:build ignore
package R04

import (
	"path/filepath"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString00 = `
disallowed-macros = ["std::println", "std::print"]
`

func ex00Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString00, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result := testNoInput(workingDirectory); !result.Passed {
		return result
	}
	_, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "cargo run | true"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "ex00", []string{"Cargo.toml", "src/main.rs"}, 10, ex00Test)
}
