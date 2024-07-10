package R00

import (
	"github.com/42-Short/shortinette/internal/logger"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	Module "github.com/42-Short/shortinette/pkg/interfaces/module"
)

func R00() *Module.Module {
	r00, err := Module.NewModule("R00", []Exercise.Exercise{ex00(), ex01(), ex02(), ex03(), ex04(), ex05()})
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return nil
	}
	return &r00
}
