package interpreter

import "errors"

var (
	// ErrTOMLDataType indicates unexpected data type passed to the function.
	ErrTOMLDataType = errors.New("wrong type error")
)
