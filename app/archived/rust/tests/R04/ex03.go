//go:build ignore
package R04

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

var clippyTomlAsString03 = `
disallowed-methods = ["std::process::Command::exec"]
`

func testInputFileBadPermissions(workingDirectory string) Exercise.Result {
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "chmod -R 777 ../target && touch donotpanic && chmod 000 donotpanic"}); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", "su -c 'echo donotpanic | cargo run -- cat'"}); err != nil {
		if strings.Contains(err.Error(), "panicked") {
			return Exercise.RuntimeError(err.Error(), "touch donotpanic", "chmod 000 donotpanic", "echo donotpanic | cargo run -- cat")
		}
	}
	return Exercise.Passed("OK")
}

func testHelloWorld(workingDirectory string) Exercise.Result {
	commandLine := "echo 'Hello, World!' | cargo run -- echo"
	expectedOutput := "Hello, World!\n"

	return doTest(workingDirectory, expectedOutput, commandLine)
}

func testCatMultipleFiles(workingDirectory string) Exercise.Result {
	if err := os.WriteFile("/tmp/a", []byte("a"), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := os.WriteFile("/tmp/b", []byte("b"), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := os.WriteFile("/tmp/c", []byte("c"), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := os.WriteFile("/tmp/d", []byte("d"), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}
	if err := os.WriteFile("/tmp/e", []byte("e"), fs.FileMode(os.O_WRONLY)); err != nil {
		return Exercise.InternalError(err.Error())
	}

	commandLine := "<< EOF cargo run -- cat\n/tmp/a\n/tmp/b\n/tmp/c\n/tmp/d\n/tmp/e\n"
	expectedOutput := "abcde"

	return doTest(workingDirectory, expectedOutput, commandLine)
}

func testMultiLineInput(workingDirectory string) Exercise.Result {
	commandLine := "<< EOF cargo run -- echo -n\nHello\n,\nWorld\n!\nEOF"
	expectedOutput := "Hello , World !"

	return doTest(workingDirectory, expectedOutput, commandLine)
}

func ex03Test(exercise *Exercise.Exercise) (result Exercise.Result) {
	if err := alloweditems.Check(*exercise, clippyTomlAsString03, nil); err != nil {
		return Exercise.CompilationError(err.Error())
	}

	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)

	if result = testNoInput(workingDirectory); !result.Passed {
		return result
	}
	if result = testInputFileBadPermissions(workingDirectory); !result.Passed {
		return result
	}
	if result = testHelloWorld(workingDirectory); !result.Passed {
		return result
	}
	if result = testMultiLineInput(workingDirectory); !result.Passed {
		return result
	}
	if result = testCatMultipleFiles(workingDirectory); !result.Passed {
		return result
	}

	return Exercise.Passed("OK")
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "ex03", []string{"Cargo.toml", "src/main.rs"}, 10, ex03Test)
}
