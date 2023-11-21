package lexer

import (
	"strings"

	"github.com/mdm-code/scanner"
)

const (
	// String ...
	String TokenType = iota

	// Integer ...
	Integer

	// Dot ...
	Dot

	// Colon ...
	Colon

	// ArrayOpen ...
	ArrayOpen

	// ArrayClose ...
	ArrayClose

	// Whitespace ...
	Whitespace

	// Undefined ...
	Undefined
)

// KeyCharMap ...
var KeyCharMap = map[rune]TokenType{
	'.': Dot,
	':': Colon,
	'[': ArrayOpen,
	']': ArrayClose,
}

// TokenType ...
type TokenType uint8

// Token ...
type Token struct {
	Buffer     *[]scanner.Token
	Type       TokenType
	Start, End int
}

// Lexeme ...
func (t Token) Lexeme() string {
	if t.Buffer == nil {
		return ""
	}
	chars := make([]string, t.End-t.Start)
	for _, t := range (*t.Buffer)[t.Start:t.End] {
		chars = append(chars, string(t.Rune))
	}
	return strings.Join(chars, "")
}
