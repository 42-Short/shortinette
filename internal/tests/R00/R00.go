package R00

import (
	"github.com/42-Short/shortinette/internal/modulebuilder"
	"github.com/42-Short/shortinette/internal/testbuilder"
)

func R00(repoId string, codeDirectory string) {
	r00 := modulebuilder.NewModuleBuilder().
		SetName("R00").
		SetExercises([]testbuilder.TestBuilder{ex00(), ex01()})
	r00.SetUp(repoId, codeDirectory)
	r00.Run()
}
