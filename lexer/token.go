package lexer

import (
	"strings"

	"github.com/mdm-code/scanner"
)

const (
	// Undefined ...
	Undefined TokenType = iota

	// String ...
	String

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
	if t.Buffer == nil || len(*t.Buffer) < 1 || t.Start > t.End {
		return ""
	}
	end := t.End
	if end > len(*t.Buffer) {
		end = len(*t.Buffer)
	}
	chars := make([]string, end-t.Start)
	for _, t := range (*t.Buffer)[t.Start:end] {
		chars = append(chars, string(t.Rune))
	}
	return strings.Join(chars, "")
}
