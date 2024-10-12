package R06

import (
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func R06() *Module.Module {
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
	r00 := Module.NewModule("06", 50, exercises, "subjects/module-06.md", "shortinette-testenv")
	return &r00
}
