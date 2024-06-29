package errors

import (
	"fmt"
)

type SubmissionError struct {
	Err     error
	Details string
}

func (e *SubmissionError) Error() string {
	return fmt.Sprintf("submission error: %v, details: %s", e.Err, e.Details)
}

func NewSubmissionError(err error, details string) *SubmissionError {
	return &SubmissionError{
		Err:     err,
		Details: details,
	}
}

var (
	ErrEmptyRepo          = fmt.Errorf("empty repository")
	ErrForbiddenItem      = fmt.Errorf("forbidden item(s) used")
	ErrInvalidOutput      = fmt.Errorf("invalid output")
	ErrInvalidCompilation = fmt.Errorf("could not compile code")
	ErrRuntime            = fmt.Errorf("code did not execute as expected")
	ErrFailedTests        = fmt.Errorf("test failed")
)
