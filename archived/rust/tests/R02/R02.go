package R02

import (
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"

	"path/filepath"
	"time"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/42-Short/shortinette/rust/alloweditems"
)

func runDefaultTest(exercise *Exercise.Exercise, cargoTestModAsString string, clippyTomlAsString string, allowedKeywords map[string]int) Exercise.Result {
	if err := alloweditems.Check(*exercise, clippyTomlAsString, allowedKeywords); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if err := testutils.AppendStringToFile(cargoTestModAsString, exercise.TurnInFiles[0]); err != nil {
		logger.Exercise.Printf("could not write to %s: %v", exercise.TurnInFiles[0], err)
		return Exercise.InternalError(err.Error())
	}
	_, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"valgrind", "test"}, testutils.WithTimeout(100*time.Second)) //TODO: maybe adjust the time
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	return Exercise.Passed("OK")
}

func R02() *Module.Module {
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
	r02 := Module.NewModule("02", 50, exercises, "subjects/module-02.md", "shortinette-testenv")
	return &r02
}
