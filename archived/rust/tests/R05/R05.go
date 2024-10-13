package R05

import (
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func R05() *Module.Module {
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
	r05 := Module.NewModule("05", 50, exercises, "subjects/module-05.md", "shortinette-testenv")
	return &r05
}
