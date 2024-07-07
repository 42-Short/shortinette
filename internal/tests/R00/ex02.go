package R00

import (
	"errors"
	"fmt"
	"strings"
	"time"

	Exercise "github.com/42-Short/shortinette/internal/interfaces/exercise"
	"github.com/42-Short/shortinette/internal/logger"
	"github.com/42-Short/shortinette/internal/tests/testutils"
)

const YesMain = `
fn main() {
	yes();
}
`

const CollatzMain = `
fn main() {
	collatz(%s);
}
`

const PrintBytesMain = `
fn main() {
	print_bytes("%s")
}
`

func yes() bool {
	exercise := Exercise.NewExercise("EX02", "studentcode", "ex02", []string{"yes.rs"}, "function", "yes()", []string{"println"}, nil, nil, nil)
	if err := testutils.ForbiddenItemsCheck(exercise, "shortinette-test-R00"); err != nil {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(exercise)
	if err := testutils.AppendStringToFile(YesMain, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s.0 KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err == nil || !errors.Is(err, testutils.ErrTimeout) {
		logger.File.Printf("[%s.0 KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "y" && line != "" {
			assertionError := testutils.AssertionErrorString(exercise.Name, "y", line)
			logger.File.Printf(assertionError)
			return false
		}
	}
	return true
}

func collatzInfiniteLoopTest(exercise Exercise.Exercise) bool {
	main := fmt.Sprintf(CollatzMain, "0")
	if err := testutils.AppendStringToFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s.1 KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	_, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		logger.File.Printf("[%s.1 KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	if err := testutils.DeleteStringFromFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	return true
}

func doCollatz(n int) string {
	if n <= 0 {
		return ""
	}
	var results []string
	for n != 1 {
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
		results = append(results, fmt.Sprintf("%d", n))
	}
	return strings.Join(results, "\n") + "\n"
}

func collatzAssertionTest(exercise Exercise.Exercise) bool {
	main := fmt.Sprintf(CollatzMain, "42")
	if err := testutils.AppendStringToFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s.1 KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		logger.File.Printf("[%s.1 KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	expectedOutput := doCollatz(42)

	if output != expectedOutput {
		assertionError := testutils.AssertionErrorString(exercise.Name, expectedOutput, output)
		logger.File.Printf(assertionError)
		return false
	}
	return true
}

func collatz() bool {
	exercise := Exercise.NewExercise("EX02", "studentcode", "ex02", []string{"collatz.rs"}, "function", "collatz(42)", []string{"println"}, nil, nil, nil)
	if err := testutils.ForbiddenItemsCheck(exercise, "shortinette-test-R00"); err != nil {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(exercise)
	if !collatzInfiniteLoopTest(exercise) {
		return false
	}
	if !collatzAssertionTest(exercise) {
		return false
	}
	return true
}

func doPrintBytes(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		result.WriteString(fmt.Sprintf("%d\n", s[i]))
	}
	return result.String()
}

func printBytesAssertionTest(exercise Exercise.Exercise) bool {
	main := fmt.Sprintf(PrintBytesMain, "Hello, World!")
	if err := testutils.AppendStringToFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return false
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		logger.File.Printf("[%s.2 KO]: %v", exercise.Name, err)
		return false
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		logger.File.Printf("[%s.2 KO]: runtime error: %v", exercise.Name, err)
		return false
	}
	expectedOutput := doPrintBytes("Hello, World!")

	if output != expectedOutput {
		assertionError := testutils.AssertionErrorString(exercise.Name, expectedOutput, output)
		logger.File.Printf(assertionError)
		return false
	}
	return true
}

func printBytes() bool {
	exercise := Exercise.NewExercise("EX02", "studentcode", "ex02", []string{"print_bytes.rs"}, "function", "print_bytes(\"\")", []string{"println", "bytes"}, nil, nil, nil)
	if err := testutils.ForbiddenItemsCheck(exercise, "shortinette-test-R00"); err != nil {
		return false
	}
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(exercise)
	return printBytesAssertionTest(exercise)
}

func ex02Test(exercise *Exercise.Exercise) bool {
	if !testutils.TurnInFilesCheck(*exercise) {
		return false
	}
	if yes() && collatz() && printBytes() {
		return true
	}
	return false
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("EX02", "studentcode", "ex02", []string{"collatz.rs", "print_bytes.rs", "yes.rs"}, "", "", nil, nil, nil, ex02Test)
}
