package R00

import (
	"path/filepath"

	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
	"github.com/42-Short/shortinette/pkg/testutils"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func CompileWithRustc(turnInFile string) error {
	workingDirectory := filepath.Dir(turnInFile)
	if _, err := testutils.RunCommandLine(workingDirectory, "rustc", []string{filepath.Base(turnInFile)}); err != nil {
		return err
	}
	return nil
}

func CompileWithRustcTest(turnInFile string) error {
	workingDirectory := filepath.Dir(turnInFile)
	if _, err := testutils.RunCommandLine(workingDirectory, "rustc", []string{"--test", filepath.Base(turnInFile)}); err != nil {
		return err
	}
	return nil
}

func R00() *Module.Module {
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
	r00 := Module.NewModule("00", 50, exercises, "subjects/module-00.md", "shortinette-testenv")
	return &r00
}
