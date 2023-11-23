package lexer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/scanner"
)

var (
	// ErrNilScanner indicates that the provided Scanner is nil.
	ErrNilScanner = errors.New("provided Scanner is nil")

	// ErrKeyCharUnsupported indicates that the key character is unsupported.
	ErrKeyCharUnsupported = errors.New("unsupported key character")

	// ErrUnterminatedString indicates that the string literal is not terminated.
	ErrUnterminatedString = errors.New("unterminated string literal")

	// ErrDisallowedChar indcates the the character is disallowed.
	ErrDisallowedChar = errors.New("disallowed character")
)

// Error wraps a concrete lexer error to represent its query context in the
// error message. It stores references to the Lexer buffer context and the
// Lexer token start offset.
type Error struct {
	buffer *[]scanner.Token // Lexer buffer context pointer
	offset int              // Lexer token start offset
	err    error            // wrapped Lexer error
}

// Error reports the Lexer error wrapped inside the Lexer buffer context with
// a marker indicating the start of the Lexer token at which the occurred.
func (e *Error) Error() string {
	msg := e.getErrorMsg()
	if e.buffer == nil {
		return msg
	}
	marker := "^"
	result := e.wrapErrorMsg(msg, marker)
	return result
}

// getErrorMsg provides the wrapped error message or the nil default.
func (e *Error) getErrorMsg() string {
	if e.err != nil {
		return e.err.Error()
	}
	return "nil"
}

// wrapErrorMsg wraps the error message inside the Lexer buffer context.
func (e *Error) wrapErrorMsg(msg, marker string) string {
	var b strings.Builder
	b.Grow(e.offset*2 + 1)
	for _, t := range *e.buffer {
		b.WriteRune(t.Rune)
	}
	b.WriteString("\n")
	indent := strings.Repeat(" ", e.offset-1)
	b.WriteString(indent)
	b.WriteString(marker)
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("Lexer error: %s", msg))
	return b.String()
}
