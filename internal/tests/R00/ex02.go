package R00

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/42-Short/shortinette/internal/logger"
	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
	"github.com/42-Short/shortinette/pkg/testutils"
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

func yes() Exercise.Result {
	exercise := Exercise.NewExercise("02", "studentcode", "ex02", []string{"yes.rs"}, "function", "yes()", []string{"println"}, nil, nil, nil)
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(exercise)
	if err := testutils.AppendStringToFile(YesMain, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return Exercise.Result{Passed: false, Output: "internal error"}	
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("compilation error: %s", err)}
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err == nil || !errors.Is(err, testutils.ErrTimeout) {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: %v", err)}
	}
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line != "y" && line != "" {
			assertionError := testutils.AssertionErrorString("y", line)
			return Exercise.Result{Passed: false, Output: assertionError}
		}
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func collatzInfiniteLoopTest(exercise Exercise.Exercise) Exercise.Result {
	main := fmt.Sprintf(CollatzMain, "0")
	if err := testutils.AppendStringToFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return Exercise.Result{Passed: false, Output: "internal error"}
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("compilation error: %v", err)}
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	_, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: %v", err)}
	}
	if err := testutils.DeleteStringFromFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return Exercise.Result{Passed: false, Output: "internal error"}
	}
	return Exercise.Result{Passed: true}
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

func collatzAssertionTest(exercise Exercise.Exercise) Exercise.Result {
	main := fmt.Sprintf(CollatzMain, "42")
	if err := testutils.AppendStringToFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return Exercise.Result{Passed: false, Output: "internal error"}
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("compilation error: %v", err)}
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: %v", err)}
	}
	expectedOutput := doCollatz(42)

	if output != expectedOutput {
		assertionError := testutils.AssertionErrorString(expectedOutput, output)
		logger.File.Printf(assertionError)
		return Exercise.Result{Passed: false, Output: assertionError}
	}
	return Exercise.Result{Passed: false, Output: ""}
}

func collatz() Exercise.Result {
	exercise := Exercise.NewExercise("02", "studentcode", "ex02", []string{"collatz.rs"}, "function", "collatz(42)", []string{"println"}, nil, nil, nil)
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(exercise)
	if result := collatzInfiniteLoopTest(exercise); !result.Passed {
		return result
	}
	if result := collatzAssertionTest(exercise); !result.Passed {
		return result
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func doPrintBytes(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		result.WriteString(fmt.Sprintf("%d\n", s[i]))
	}
	return result.String()
}

func printBytesAssertionTest(exercise Exercise.Exercise) Exercise.Result {
	main := fmt.Sprintf(PrintBytesMain, "Hello, World!")
	if err := testutils.AppendStringToFile(main, exercise.TurnInFiles[0]); err != nil {
		logger.Error.Printf("internal error: %v", err)
		return Exercise.Result{Passed: false, Output: "internal error"}
	}
	if err := testutils.CompileWithRustc(exercise.TurnInFiles[0]); err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("compilation error: %v", err)}
	}
	executablePath := testutils.ExecutablePath(exercise.TurnInFiles[0], ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		return Exercise.Result{Passed: false, Output: fmt.Sprintf("runtime error: %v", err)}
	}
	expectedOutput := doPrintBytes("Hello, World!")

	if output != expectedOutput {
		assertionError := testutils.AssertionErrorString(expectedOutput, output)
		return Exercise.Result{Passed: false, Output: assertionError}
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func printBytes() Exercise.Result {
	exercise := Exercise.NewExercise("02", "studentcode", "ex02", []string{"print_bytes.rs"}, "function", "print_bytes(\"\")", []string{"println", "bytes"}, nil, nil, nil)
	exercise.TurnInFiles = testutils.FullTurnInFilesPath(exercise)
	return printBytesAssertionTest(exercise)
}

func ex02Test(exercise *Exercise.Exercise) Exercise.Result {
	if !testutils.TurnInFilesCheck(*exercise) {
		return Exercise.Result{Passed: false, Output: "invalid files found in turn in directory"}
	}
	if result := yes(); !result.Passed {
		return result
	} 
	if result := collatz(); !result.Passed {
		return result
	} 
	if result := printBytes(); !result.Passed {
		return result
	}
	return Exercise.Result{Passed: true, Output: ""}
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("02", "studentcode", "ex02", []string{"collatz.rs", "print_bytes.rs", "yes.rs"}, "", "", nil, nil, nil, ex02Test)
}
