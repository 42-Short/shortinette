package R00

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/42-Short/shortinette/pkg/logger"
	"github.com/42-Short/shortinette/pkg/testutils"
	"github.com/42-Short/shortinette/rust/alloweditems"

	Exercise "github.com/42-Short/shortinette/pkg/interfaces/exercise"
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

func yes(filename string) Exercise.Result {
	if err := testutils.AppendStringToFile(YesMain, filename); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if err := CompileWithRustc(filename); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	executablePath := testutils.ExecutablePath(filename, ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil && !errors.Is(err, testutils.ErrTimeout) {
		return Exercise.RuntimeError(err.Error())
	}
	count := 0
	lines := strings.Split(output, "\n")
	nl_found := false // Just to avoid empty lines being graded as correct, except at EOF
	for _, line := range lines {
		if nl_found {
			return Exercise.AssertionError("y", "")
		}
		if line != "y" && line != "" {
			return Exercise.AssertionError("y", line)
		}
		if line == "" {
			nl_found = true
		}
		count++
	}
	if count < 1000 {
		return Exercise.RuntimeError("Expected 'y' to be printed more often")
	}
	return Exercise.Passed("OK")
}

func collatzInfiniteLoopTest(filename string) Exercise.Result {
	main := fmt.Sprintf(CollatzMain, "0")
	if err := testutils.AppendStringToFile(main, filename); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if err := CompileWithRustc(filename); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	executablePath := testutils.ExecutablePath(filename, ".rs")

	if _, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond)); err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	if err := testutils.DeleteStringFromFile(main, filename); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	return Exercise.Passed("OK")
}

func doCollatz(n int) string {
	if n <= 0 {
		return ""
	}
	var results []string
	for n != 1 {
		results = append(results, fmt.Sprintf("%d", n))
		if n%2 == 0 {
			n /= 2
		} else {
			n = 3*n + 1
		}
	}
	results = append(results, "1")
	return strings.Join(results, "\n") + "\n"
}

func collatzAssertionTest(filename string, number int) Exercise.Result {
	main := fmt.Sprintf(CollatzMain, fmt.Sprintf("%du32", number))
	if err := testutils.AppendStringToFile(main, filename); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	if err := CompileWithRustc(filename); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	executablePath := testutils.ExecutablePath(filename, ".rs")
	output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
	if err != nil {
		return Exercise.RuntimeError(err.Error())
	}
	expectedOutput := doCollatz(number)

	if output != expectedOutput {
		return Exercise.AssertionError(expectedOutput, output)
	}
	if err := testutils.DeleteStringFromFile(main, filename); err != nil {
		logger.Exercise.Printf("internal error: %v", err)
		return Exercise.InternalError(err.Error())
	}
	return Exercise.Passed("OK")
}

func collatz(filename string) Exercise.Result {
	if result := collatzInfiniteLoopTest(filename); !result.Passed {
		return result
	}
	testNumbers := []int{42, 1, 524287}
	for _, number := range testNumbers {
		if result := collatzAssertionTest(filename, number); !result.Passed {
			return result
		}
	}
	return Exercise.Passed("OK")
}

func doPrintBytes(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		result.WriteString(fmt.Sprintf("%d\n", s[i]))
	}
	return result.String()
}

func printBytesAssertionTest(filename string) Exercise.Result {
	testStrings := []string{"Hello, World", "", "Rust is awesome! 🦀", string([]byte{0})}
	for _, testString := range testStrings {
		main := fmt.Sprintf(PrintBytesMain, testString)
		if err := testutils.AppendStringToFile(main, filename); err != nil {
			logger.Exercise.Printf("internal error: %v", err)
			return Exercise.InternalError(err.Error())
		}
		if err := CompileWithRustc(filename); err != nil {
			return Exercise.CompilationError(err.Error())
		}
		executablePath := testutils.ExecutablePath(filename, ".rs")
		output, err := testutils.RunExecutable(executablePath, testutils.WithTimeout(500*time.Millisecond))
		if err != nil {
			return Exercise.RuntimeError(err.Error())
		}
		if expectedOutput := doPrintBytes(testString); output != expectedOutput {
			return Exercise.AssertionError(expectedOutput, output)
		}
		if err := testutils.DeleteStringFromFile(main, filename); err != nil {
			logger.Exercise.Printf("internal error: %v", err)
			return Exercise.InternalError(err.Error())
		}
	}
	return Exercise.Passed("OK")
}

func printBytes(filename string) Exercise.Result {
	return printBytesAssertionTest(filename)
}

func clippyCheck02(exercise *Exercise.Exercise) Exercise.Result {
	workingDirectory := filepath.Join(exercise.CloneDirectory, exercise.TurnInDirectory)
	if _, err := testutils.RunCommandLine(workingDirectory, "cargo", []string{"init", "--lib"}); err != nil {
		return Exercise.InternalError("cargo init failed")
	}
	concat := ""
	for _, file := range exercise.TurnInFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return Exercise.InternalError(fmt.Sprintf("Error copying %s to src/ folder: %v", file, err))
		}
		concat += string(content)
	}
	if err := os.WriteFile(filepath.Join(workingDirectory, "src/lib.rs"), []byte(concat), 0644); err != nil {
		return Exercise.InternalError(fmt.Sprintf("Error writing content to src/lib.rs: %v", err))
	}
	tmp := Exercise.Exercise{
		CloneDirectory:  exercise.CloneDirectory,
		TurnInDirectory: exercise.TurnInDirectory,
		TurnInFiles:     []string{filepath.Join(workingDirectory, "src/lib.rs")},
	}
	if err := alloweditems.Check(tmp, "", map[string]int{"unsafe": 0, "for": 1, "loop": 1, "while": 1}); err != nil {
		return Exercise.CompilationError(err.Error())
	}
	for _, loopKind := range []string{"for", "loop", "while"} {
		if err := alloweditems.Check(tmp, "", map[string]int{loopKind: 0}); err == nil {
			return Exercise.CompilationError(fmt.Sprintf("Loop kind '%s' not used", loopKind))
		}
	}
	return Exercise.Passed("OK")
}

func ex02Test(exercise *Exercise.Exercise) Exercise.Result {
	if result := clippyCheck02(exercise); !result.Passed {
		return result
	}
	if result := yes(exercise.TurnInFiles[2]); !result.Passed {
		return result
	}
	if result := collatz(exercise.TurnInFiles[0]); !result.Passed {
		return result
	}
	if result := printBytes(exercise.TurnInFiles[1]); !result.Passed {
		return result
	}
	return Exercise.Passed("OK")
}

func ex02() Exercise.Exercise {
	return Exercise.NewExercise("02", "ex02", []string{"collatz.rs", "print_bytes.rs", "yes.rs"}, 10, ex02Test)
}
