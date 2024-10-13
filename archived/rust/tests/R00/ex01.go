package R00

import (
	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

const CargoTest = `
#[cfg(test)]
mod shortinette_tests_rust_0001 {
	use super::*;

	#[test]
	fn test_0() {
		assert_eq!(min(1i32, 2i32), 1i32);
	}

	#[test]
	fn test_1() {
		assert_eq!(min(2i32, 1i32), 1i32);
	}

	#[test]
	fn test_2() {
		assert_eq!(min(1i32, 1i32), 1i32);
	}

	#[test]
	fn test_3() {
		assert_eq!(min(-1i32, 0i32), -1i32);
	}
}
`

var clippyTomlAsString01 = `
disallowed-methods = ["std::cmp::min"]
`

func clippyCheck01(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"init", "--lib"}); err != nil {
		return Exercise.InternalError("cargo init failed")
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "cp", []string{"min.rs", "src/lib.rs"}); err != nil {
		return Exercise.InternalError("unable to copy file to src/ folder")
	}
	tmp := Exercise.Exercise{
		CloneDirectory:  exercise.CloneDirectory,
		TurnInDirectory: exercise.TurnInDirectory,
		TurnInFiles:     []string{filepath.Join(workingDirectory, "src/lib.rs")},
	}
	if err := alloweditems.Check(tmp, clippyTomlAsString01, map[string]int{"unsafe": 0, "return": 0}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	return Exercise.Passed("")
}

func ex01Test(exercise *Exercise.Exercise) Exercise.Result {
	if err := testutils.AppendStringToFile(CargoTest, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	if result := clippyCheck01(exercise); !result.Passed {
		return result
	}
	if err := CompileWithRustcTest(exercise.TurnInFiles[0]); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	if output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond)); err != nil {
		return Exercise.RuntimeError(output)
	}
	if err := CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		return Exercise.CompilationError("main function missing")
	}
	if output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond)); err != nil {
		return Exercise.RuntimeError(output)
	}
	return Exercise.Passed("OK")
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"min.rs"}, 10, ex01Test)
}
