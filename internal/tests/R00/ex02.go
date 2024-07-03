package R00

import (
	IExercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

const YesMain = `
fn main() {
	yes();
}
`

func yes(exercise *IExercise.Exercise) bool {
	fullTurnInPath := testutils.FullTurnInFilePath("studentcode", *exercise)
	if err := testutils.AppendStringToFile(YesMain, fullTurnInPath); err != nil {
		logger.File.Printf("internal error: %v", err)
	}
	return true
}

func collatz(exercise *IExercise.Exercise) bool {
	return true
}

func print_bytes(exercise *IExercise.Exercise) bool {
	return true
}

func ex02Test(exercise *IExercise.Exercise) bool {
	if yes(exercise) && collatz(exercise) && print_bytes(exercise) {
		return true
	}
	return false
}

func ex02() IExercise.ExerciseBuilder {
	return IExercise.NewExerciseBuilder().
		SetName("EX02").
		SetTurnInDirectory("ex02").
		SetTurnInFile("").
		SetExerciseType("").
		SetPrototype("").
		SetAllowedMacros(nil).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(nil).
		SetExecuter(ex02Test)
}
