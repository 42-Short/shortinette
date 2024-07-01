package errors

import (
	"fmt"
)

type SubmissionError struct {
	Err     error
	Details string
}

type InternalError struct {
	Err     error
	Details string
}

func (e *InternalError) Error() string {
	return fmt.Sprintf("internal error: %v, details: %s", e.Err, e.Details)
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

func NewInternalError(err error, details string) *InternalError {
	return &InternalError{
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
)

var (
	ErrInternal = fmt.Errorf("internal error")
)
