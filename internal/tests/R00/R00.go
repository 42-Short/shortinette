package R00

import (
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	Module "github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/logger"
)

func R00() *Module.Module {
	r00, err := Module.NewModule("R00", []Exercise.Exercise{ex00(), ex01(), ex02(), ex03(), ex04(), ex05()})
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return nil
	}
	return &r00
}
