package interpreter

import (
	"errors"
	"fmt"
)

var (
	// ErrTOMLDataType indicates unexpected data type passed to the function.
	ErrTOMLDataType = errors.New("wrong type error")
)

// Error wraps an interpreter error to show how a given data type and value
// cannot be queried with the user-specified filter.
type Error struct {
	data   any
	filter string
	err    error
}

// Is allows to check if Error.err matches the target error.
func (e *Error) Is(target error) bool {
	return e.err == target
}

// Error reports the Interpreter error with the data type and value followed
// by the name of the data filter that was to be applied to this data.
func (e *Error) Error() string {
	return fmt.Sprintf(
		"Interpreter error: cannot query [ %T ] ( %v ) with ( %s )",
		e.data,
		e.data,
		e.filter,
	)
}
