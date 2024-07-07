package R00

import (
	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	Module "github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/logger"
)

func R00(repoId string, codeDirectory string) {
	r00, err := Module.NewModule("R00", []Exercise.Exercise{ex00(), ex01(), ex02(), ex03(), ex04(), ex05()}, "shortinette-test-R00", "studentcode")
	if err != nil {
		logger.Error.Printf("internal error: %v", err)
		return
	}
	r00.Run()
}
