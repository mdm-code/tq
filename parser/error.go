package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/tq/lexer"
)

var (
	// ErrQueryElement indicates an unprocessable query filter element.
	ErrQueryElement = errors.New("expected '.' or '[' to parse query element")

	// ErrSelectorUnterminated indicates an unterminated selector element.
	ErrSelectorUnterminated = errors.New("expected ']' to terminate selector")
)

// Error wraps a concrete Parser error to represent its context. It reports the
// token where the error has occurred.
type Error struct {
	token lexer.Token
	err   error
}

// Is allows to check if Error.err matches the target error.
func (e *Error) Is(target error) bool {
	return e.err == target
}

// Error reports the Parser error wrapped inside of the custom context.
func (e *Error) Error() string {
	return e.getErrorLine()
}

// getErrorLine provides the error line with the nil error as default.
func (e *Error) getErrorLine() string {
	var b strings.Builder
	b.WriteString("Parser error: ")
	if e.err != nil {
		b.WriteString(e.err.Error())
		b.WriteString(fmt.Sprintf(" but got '%s'", e.token.Lexeme()))
		return b.String()
	}
	return b.String()
}