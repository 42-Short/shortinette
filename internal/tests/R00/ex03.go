package R00

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
)

func doFizzBuzz() string {
	var result strings.Builder
	for number := 1; number <= 100; number++ {
		switch {
		case number%3 == 0 && number%5 == 0:
			result.WriteString("fizzbuzz\n")
		case number%3 == 0:
			result.WriteString("fizz\n")
		case number%5 == 0:
			result.WriteString("buzz\n")
		case number%11 == 3:
			result.WriteString("FIZZ\n")
		case number%11 == 5:
			result.WriteString("BUZZ\n")
		default:
			result.WriteString(fmt.Sprintf("%d\n", number))
		}
	}
	return result.String()
}

func fizzBuzzOutputTest(exercise Exercise.Exercise) bool {
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], filepath.Ext(exercise.TurnInFiles[0]))
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		logger.File.Printf("[%s KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	expectedOutput := doFizzBuzz()

	if output != expectedOutput {
		assertionError := testutils.AssertionErrorString(exercise.Name, expectedOutput, output)
		logger.File.Printf(assertionError)
		return false
	}
	return true
}

func ex03Test(exercise *Exercise.Exercise) bool {
	if !testutils.TurnInFilesCheck(*exercise) {
		return false
	}
	if err := testutils.ForbiddenItemsCheck(*exercise, "shortinette-test-R00"); err != nil {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(*exercise)

	return fizzBuzzOutputTest(*exercise)
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("EX03", "studentcode", "ex03", []string{"fizzbuzz.rs"}, "program", "", []string{"println"}, nil, map[string]int{"match": 1, "for": 1}, ex03Test)
}