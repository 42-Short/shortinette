package R00

import (
	"fmt"

	"github.com/42-Short/shortinette/internal/functioncheck"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/testbuilder"
)

func ex01Test(test *testbuilder.Test) bool {
	fmt.Println(test.ExerciseType)
	fmt.Println(test.TurnInDirectory)
	if err := functioncheck.Execute(*test, "shortinette-test-R00"); err != nil {
		logger.File.Printf("[%s KO]: %v", test.Name, err)
		return false
	}
	logger.File.Printf("[%s OK]", test.Name)
	return true
}

func ex01() testbuilder.TestBuilder {
	return testbuilder.NewTestBuilder().
		SetName("EX01").
		SetTurnInDirectory("ex01").
		SetTurnInFile("min.rs").
		SetExerciseType("function").
		SetPrototype("min(0, 0)").
		SetAllowedMacros([]string{"println"}).
		SetAllowedFunctions(nil).
		SetAllowedKeywords(map[string]int{"unsafe": 0}).
		SetExecuter(ex01Test)
}
