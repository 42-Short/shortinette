package R01

import (
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func R01() *Module.Module {
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
	r01 := Module.NewModule("01", 50, exercises, "subjects/module-01.md", "shortinette-testenv")
	return &r01
}
