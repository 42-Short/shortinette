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

func InternalError() Result {
	return Result{Passed: true, Output: "no idea if you actually passed but the software broke so you have the benefit of the doubt"}
}

func Passed(message string) Result {
	return Result{Passed: true, Output: message}
}
