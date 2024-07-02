package R00

import exercisebuilder "github.com/42-Short/shortinette/internal/interfaces/exercise"

func collatzTest(test *exercisebuilder.Test) bool {

}

func collatz() exercisebuilder.ExerciseBuilder {
	return exercisebuilder.NewExerciseBuilder().
		SetName("EX02.0").
		SetTurnInDirectory("ex02").
		SetTurnInFile("collatz.rs").
		SetExerciseType("function").
		SetPrototype("collatz(1)").
		SetAllowedMacros([]string{"println"}).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(nil).
		SetExecuter(collatzTest)
	}
