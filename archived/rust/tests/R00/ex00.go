//go:build ignore
package R00

import (
	"path/filepath"
	"github.com/42-Short/shortinette/rust/alloweditems"
	"time"

	"github.com/42-Short/shortinette/pkg/testutils"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func clippyCheck00(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"init"}); err != nil {
		return Exercise.InternalError("cargo init failed")
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "cp", []string{"hello.rs", "src/main.rs"}); err != nil {
		return Exercise.InternalError("unable to copy file to src/ folder")
	}
	tmp := Exercise.Exercise{
		CloneDirectory:  exercise.CloneDirectory,
		TurnInDirectory: exercise.TurnInDirectory,
		TurnInFiles:     []string{filepath.Join(workingDirectory, "src/main.rs")},
	}
	if err := alloweditems.Check(tmp, "", map[string]int{"unsafe": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	return Exercise.Passed("")
}

func ex00Test(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	fileName := filepath.Base(exercise.TurnInFiles[0])
	if _, err := testutils.RunCommandLine(workingDirectory, "rustc", []string{fileName}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if result := clippyCheck00(exercise); !result.Passed {
		return result
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		return Exercise.RuntimeError(output)
	}
	if output != "Hello, World!\n" {
		return Exercise.AssertionError("Hello, World!\n", output)
	}
	return Exercise.Passed("OK")
}

func ex00() Exercise.Exercise {
	return Exercise.NewExercise("00", "ex00", []string{"hello.rs"}, 10, ex00Test)
}
