package R04

import (
	"fmt"
	"strings"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/testutils"
)

type outputChannel struct {
	out []byte
	err error
}

// Runs `cargo run` with no arguments and checks for panic.
func testNoInput(workingDirectory string) Exercise.Result {
	commandLine := "cargo run"
	if _, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine}); err != nil {
		if strings.Contains(err.Error(), "thread 'main' panicked") {
			return Exercise.RuntimeError(fmt.Sprintf("i said don't panic :(\n%s", err.Error()), commandLine)
		}
	}
	return Exercise.Passed("OK")
}

func doTest(workingDirectory string, expectedOutput string, commandLine string) Exercise.Result {
	output, err := testutils.RunCommandLine(workingDirectory, "sh", []string{"-c", commandLine})
	if err != nil {
		return Exercise.RuntimeError(err.Error(), commandLine)
	}
	if output != expectedOutput {
		return Exercise.AssertionError(expectedOutput, output, commandLine)
	}
	return Exercise.Passed("OK")
}

func R04() *Module.Module {
	exercises := map[string]Exercise.Exercise{
		"00": ex00(),
		"01": ex01(),
		"02": ex02(),
		"03": ex03(),
		"04": ex04(),
		"05": ex05(),
		"06": ex06(),
		"07": ex07(),
	}
	r00 := Module.NewModule("04", 50, exercises, "subjects/module-04.md", "shortinette-testenv")
	return &r00
}
