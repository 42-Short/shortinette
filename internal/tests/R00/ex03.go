package R00

import (
	"fmt"
	"strings"
	"time"

	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
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

func ex03Test(exercise *Exercise.Exercise) bool {
	directory := testutils.FullTurnInDirectory("studentcode", *exercise)
	if err := testutils.CompileWithRustc(directory, exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s KO]: %v", exercise.Name, err)
		return false
	}
	fullTurnInFilePath := testutils.FullTurnInFilePath("studentcode", *exercise, exercise.TurnInFiles[0])
	executablePath := testutils.ExecutablePath(fullTurnInFilePath, ".rs")
	output, err := testutils.RunCode(executablePath, testutils.WithTimeout(500*time.Millisecond))
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

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("EX03", "studentcode", "ex03", []string{"fizzbuzz.rs"}, "program", "", []string{"println"}, nil, map[string]int{"match": 1, "for": 1}, ex03Test)
}
