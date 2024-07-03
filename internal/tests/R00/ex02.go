package R00

import (
	"time"

	IExercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

const YesMain = `
fn main() {
	yes();
}
`

func yesBuilder() IExercise.ExerciseBuilder {
	return IExercise.NewExerciseBuilder().
		SetName("EX02").
		SetTurnInDirectory("ex02").
		SetTurnInFile("yes.rs").
		SetExerciseType("function").
		SetPrototype("yes()").
		SetAllowedMacros([]string{"println"}).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(nil).
		SetExecuter(nil)
}

func yes() bool {
	exercise := yesBuilder().Build()
	fullTurnInFilePath := testutils.FullTurnInFilePath("studentcode", exercise)
	if err := testutils.AppendStringToFile(YesMain, fullTurnInFilePath); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	executablePath := testutils.ExecutablePath(fullTurnInFilePath, ".rs")
	_, err := testutils.RunCode(executablePath, testutils.WithTimeout(1*time.Second))
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	return true
}

func collatz() bool {
	return true
}

func print_bytes() bool {
	return true
}

func ex02Test(exercise *IExercise.Exercise) bool {
	if yes() && collatz() && print_bytes() {
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
