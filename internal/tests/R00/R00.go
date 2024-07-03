package R00

import (
	exercisebuilder "github.com/42-Short/shortinette/internal/interfaces/exercise"
	modulebuilder "github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/logger"
)

func R00(repoId string, codeDirectory string) {
	r00 := modulebuilder.NewModuleBuilder().
		SetName("R00").
		SetExercises([]exercisebuilder.ExerciseBuilder{ex00(), ex01(), ex02()})
	if err := r00.SetUp(repoId, codeDirectory); err != nil {
		logger.Error.Printf("internal error: %v", err)
	}
	r00.Run()
}
