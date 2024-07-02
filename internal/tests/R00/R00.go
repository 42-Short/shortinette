package R00

import (
	"github.com/42-Short/shortinette/internal/interfaces/module"
	"github.com/42-Short/shortinette/internal/interfaces/exercise"
)

func R00(repoId string, codeDirectory string) {
	r00 := modulebuilder.NewModuleBuilder().
		SetName("R00").
		SetExercises([]exercisebuilder.ExerciseBuilder{ex00(), ex01()})
	r00.SetUp(repoId, codeDirectory)
	r00.Run()
}
