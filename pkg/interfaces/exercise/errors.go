// Package `Exercise` provides structures and functions for defining and handling results of exercises.
package Exercise

import (
	"fmt"
	"strings"
)

// RuntimeError returns a Result indicating a runtime error occurred with the provided
// error message.
//
//   - errorMessage: The error message describing the runtime error.
//
// Returns a Result with Passed set to false and the error message included in the Output.
func RuntimeError(errorMessage string, runTests ...string) (res Result) {
	testExplanation := ""
	if len(runTests) != 0 {
		testExplanation = "failed test:\n$ " + strings.Join(runTests, "\n$ ")

	}
	return Result{Passed: false, Output: fmt.Sprintf("%s\nruntime error: %s", testExplanation, errorMessage)}
}

// CompilationError returns a Result indicating a compilation error occurred with the provided
// error message.
//
//   - errorMessage: The error message describing the compilation error.
//
// Returns a Result with Passed set to false and the error message included in the Output.
func CompilationError(errorMessage string) (res Result) {
	return Result{Passed: false, Output: fmt.Sprintf("compilation error: %s", errorMessage)}
}

// InvalidFileError returns a Result indicating that invalid files were found in the
// submission.
//
// Returns a Result with Passed set to false and a message indicating the presence of invalid files in the Output.
func InvalidFileError() (res Result) {
	return Result{Passed: false, Output: "invalid file(s) found in turn in directory"}
}

// AssertionError returns a Result indicating that the output of the student's code did not match
// the expected output.
//
//   - expected: The expected output as a string.
//   - got: The actual output produced by the student's code.
//
// Returns a Result with Passed set to false and a message detailing the discrepancy between expected and actual output in the Output.
func AssertionError(expected string, got string, runTests ...string) (res Result) {
	testExplanation := ""
	if len(runTests) != 0 {
		testExplanation = "failed test:\n$ " + strings.Join(runTests, "\n$ ")

	}
	return Result{Passed: false, Output: fmt.Sprintf("%s\ninvalid output: expected:\n%s\ngot:\n%s\n", testExplanation, expected, got)}
}

// InternalError returns a Result indicating an internal error occurred during the execution
// of the exercise.
//
//   - errorMessage: The error message describing the internal error.
//
// Returns a Result with Passed set to false and the error message included in the Output.
func InternalError(errorMessage string) (res Result) {
	return Result{Passed: false, Output: fmt.Sprintf("internal error: %v", errorMessage)}
}

// Passed returns a Result indicating that the exercise was successfully completed.
//
//   - message: A success message to include in the Output.
//
// Returns a Result with Passed set to true and the success message included in the Output.
func Passed(message string) (res Result) {
	return Result{Passed: true, Output: message}
}
