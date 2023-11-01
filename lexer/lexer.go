// Package lexer ...
package lexer

import (
	"errors"

	"github.com/mdm-code/scanner"
)

const (
	// IDENT ...
	IDENT TokenType = iota
)

// ErrNilScanner ...
var ErrNilScanner = errors.New("provided Scanner is nil")

// Zero ...
var Zero = Pos{Tokens: []scanner.Token{}, Start: 0, End: 0}

// TokenType ...
type TokenType uint8

// Pos ...
type Pos struct {
	Tokens     []scanner.Token
	Start, End int
}

// Token ...
type Token struct {
	Pos
	Buffer *[]scanner.Token
	Type   TokenType
	Lexeme string
}

// Lexer ...
type Lexer struct {
	Buffer []scanner.Token
	Cursor Pos
}

// New ...
func New(s *scanner.Scanner) (*Lexer, error) {
	if s == nil {
		return nil, ErrNilScanner
	}
	buf := []scanner.Token{}
	for s.Scan() {
		t := s.Token()
		buf = append(buf, t)
	}
	l := Lexer{
		Cursor: Zero,
		Buffer: buf,
	}
	return &l, nil
}
