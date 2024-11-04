package tester

import "errors"

const (
	InternalError = iota + 1
	EarlyGrading
	InvalidFiles
	NothingTurnedIn
	RuntimeError
	Timeout
)

type GradingError struct {
	code int
	err  string
}

func (e *GradingError) Error() string {
	return e.err
}

func TestingError(code int, err string) *GradingError {
	return &GradingError{
		code: code,
		err:  err,
	}
}

func matchesCustomError(err error, code int) bool {
	var customError *GradingError
	if errors.As(err, &customError) {
		return customError.code == code
	}
	return false
}
