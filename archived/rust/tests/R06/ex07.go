package R06

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/42-Short/shortinette/rust/attributes"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
)

func ex07Test(exercise *Exercise.Exercise) Exercise.Result {
	requiredAttributes := map[string]bool{
		"#![no_std]":  true,
		"#![no_main]": true,
	}
	if err := attributes.Check(filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory, "ft_putchar.rs"), requiredAttributes, map[string]bool{}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	if output, err := testutils.RunCommandLine(filepath.Dir(exercise.TurnInFiles[0]), "rustc", []string{"-C", "panic=abort", "-C", "link-args=-nostartfiles", "-o", "ft_putchar", "ft_putchar.rs"}); err != nil {
		return Exercise.CompilationError(fmt.Sprintf("%s: %s", err.Error(), output))
	}
	executablePath := strings.TrimSuffix(exercise.TurnInFiles[0], ".rs")

	output, err := testutils.RunExecutable(executablePath)
	if err != nil {
		if err.Error() != "exit status 42" {
			return Exercise.RuntimeError("invalid exit code:" + err.Error())
		}
	}

	if !strings.Contains(output, "42") {
		return Exercise.AssertionError("42", output)
	}

	return Exercise.Passed("OK")
}

func ex07() Exercise.Exercise {
	return Exercise.NewExercise("07", "ex07", []string{"ft_putchar.rs"}, 20, ex07Test)
}
