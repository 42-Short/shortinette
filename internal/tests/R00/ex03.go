package R00

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
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

func fizzBuzzOutputTest(exercise Exercise.Exercise) Exercise.Result {
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], filepath.Ext(exercise.TurnInFiles[0]))
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	expectedOutput := doFizzBuzz()

	if output != expectedOutput {
		return Exercise.AssertionError(expectedOutput, output)
	}
	return Exercise.Passed("OK")
}

func ex03Test(exercise *Exercise.Exercise) Exercise.Result {
	return fizzBuzzOutputTest(*exercise)
}

func ex03() Exercise.Exercise {
	return Exercise.NewExercise("03", "studentcode", "ex03", []string{"fizzbuzz.rs"}, "program", "", []string{"println"}, nil, map[string]int{"match": 1, "for": 1}, 10, ex03Test)
}
