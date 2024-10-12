package R04

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString01 = `
disallowed-methods = ["std::io::copy", "std::fs::write", "core::option::Option::unwrap_or_else", "core::result::Result::unwrap_or_else"]
`

func testRedirectionBadPermissionTargetFile(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "chmod -R 777 ../target && touch foo.txt && chmod 000 foo.txt"}); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "su -c 'echo donotpanic | cargo run -- foo.txt' student"}); err != nil {
		if strings.Contains(err.Error(), "panicked") {
			return Exercise.RuntimeError(fmt.Sprintf("i said don't panic :/\n%s", err.Error()), "touch foo.txt", "chmod 000 foo.txt", "echo donotpanic | cargo run -- foo.txt")
		}
	}
	return Exercise.Passed("OK")
}

func testRedirectionBadPermissionTargetDir(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "chmod -R 777 ../target && mkdir foo && chmod 000 foo"}); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "su -c 'echo donotpanic | cargo run -- foo/foo.txt' student"}); err != nil {
		if strings.Contains(err.Error(), "panicked") {
			return Exercise.RuntimeError(fmt.Sprintf("i said don't panic :/\n%s", err.Error()), "mkdir foo", "chmod 000 foo", "echo donotpanic | cargo run -- foo/foo.txt")
		}
	}
	return Exercise.Passed("OK")
}

func testRedirectionMultipleFile(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "echo 'Hello, World!' | cargo run -- a b c"}); err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	if output, _ := testutils.RunCommandLine(workingDirectory, "cat", []string{"a", "b", "c"}); output != "Hello, World!\nHello, World!\nHello, World!\n" {
		return Exercise.AssertionError("Hello, World!\nHello, World!\nHello, World!\n", output, "echo 'Hello, World!' | cargo run -- a b c", "cat a b c")
	}
	return Exercise.Passed("OK")
}

func testRedirectionOneFile(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "echo 'Hello, World!' | cargo run -- a"}); err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	if output, _ := testutils.RunCommandLine(workingDirectory, "cat", []string{"a"}); output != "Hello, World!\n" {
		return Exercise.AssertionError("Hello, World!\n", output, "echo 'Hello, World!' | cargo run -- a", "cat a")
	}
	return Exercise.Passed("OK")
}

func testStdout(workingDirectory string) Exercise.Result {
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "echo 'Hello, World!' | cargo run"})
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	if output != "Hello, World!\n" {
		return Exercise.AssertionError("Hello, World!", output)
	}
	return Exercise.Passed("OK")
}

func ex01Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString01, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testStdout(workingDirectory); !result.Passed {
		return result
	}
	if result = testRedirectionOneFile(workingDirectory); !result.Passed {
		return result
	}
	if result = testRedirectionMultipleFile(workingDirectory); !result.Passed {
		return result
	}
	if result = testRedirectionBadPermissionTargetDir(workingDirectory); !result.Passed {
		return result
	}
	if result = testRedirectionBadPermissionTargetFile(workingDirectory); !result.Passed {
		return result
	}
	return Exercise.Passed("OK")
}

func ex01() Exercise.Exercise {
	return Exercise.NewExercise("01", "ex01", []string{"Cargo.toml", "src/main.rs"}, 10, ex01Test)
}
