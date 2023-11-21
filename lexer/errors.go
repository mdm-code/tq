package lexer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mdm-code/scanner"
)

var (
	// ErrNilScanner ...
	ErrNilScanner = errors.New("provided Scanner is nil")

	// ErrKeyCharUnsupported ...
	ErrKeyCharUnsupported = errors.New("unsupported key character")

	// ErrUnterminatedString ...
	ErrUnterminatedString = errors.New("unterminated string literal")

	// ErrDisallowedChar ...
	ErrDisallowedChar = errors.New("disallowed character")
)

// Error ...
type Error struct {
	Buffer *[]scanner.Token
	Offset int
	Err    error
}

// Error ...
func (e *Error) Error() string {
	if e.Buffer == nil {
		return e.Err.Error()
	}
	var b strings.Builder
	b.Grow(e.Offset*2 + 1)
	for _, t := range *e.Buffer {
		b.WriteString(string(t.Rune))
	}
	b.WriteString("\n")
	for i := 0; i < e.Offset; i++ {
		b.WriteString(" ")
	}
	b.WriteString("^")
	b.WriteString("\n")
	var errMsg string
	if e.Err != nil {
		errMsg = e.Err.Error()
	} else {
		errMsg = "null"
	}
	b.WriteString(fmt.Sprintf("Lexer error: %s", errMsg))
	return b.String()
}
