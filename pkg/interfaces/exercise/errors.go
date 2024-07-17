package Exercise

import (
	"fmt"
	"strings"
)

func RuntimeError(errorMessage string) Result {
	return Result{Passed: false, Output: fmt.Sprintf("runtime error: %s", errorMessage)}
}

func CompilationError(errorMessage string) Result {
	return Result{Passed: false, Output: fmt.Sprintf("compilation error: %s", errorMessage)}
}

func InvalidFileError() Result {
	return Result{Passed: false, Output: "invalid file(s) found in turn in directory"}
}

func AssertionError(expected string, got string) Result {
	expectedReplaced := strings.ReplaceAll(expected, "\n", "\\n")
	gotReplaced := strings.ReplaceAll(got, "\n", "\\n")
	return Result{Passed: false, Output: fmt.Sprintf("invalid output: expected '%s', got '%s'", expectedReplaced, gotReplaced)}
}

func InternalError(errorMessage string) Result {
	return Result{Passed: false, Output: fmt.Sprintf("internal error: %v", errorMessage)}
}

func Passed(message string) Result {
	return Result{Passed: true, Output: message}
}
