package lexer

import (
	"strings"

	"github.com/mdm-code/scanner"
)

const (
	// Undefined represents an undefined token type.
	Undefined TokenType = iota

	// String represents a string token type.
	String

	// Integer represents an integer token type.
	Integer

	// Dot represents a full stop token type.
	Dot

	// Colon represents a colon token type.
	Colon

	// ArrayOpen represents an opening bracket token type.
	ArrayOpen

	// ArrayClose represents a closing bracket token type.
	ArrayClose

	// Whitespace represents a white space token type.
	Whitespace
)

// keyCharMap maps runes onto TokenTypes.
var keyCharMap = map[rune]TokenType{
	'.': Dot,
	':': Colon,
	'[': ArrayOpen,
	']': ArrayClose,
}

// escapeSequenceMap maps popular escape sequence characters onto its Go string
// Unicode representation.
var escapeSequenceMap = map[rune]string{
	'b':  "\b",
	't':  "\t",
	'n':  "\n",
	'f':  "\f",
	'r':  "\r",
	'"':  "\"",
	'\'': "'",
	'\\': "\\",
}

// TokenType indicates the type of the lexer Token.
type TokenType uint8

// Token represents a single lexeme read from the Scanner token buffer.
type Token struct {
	Type                   TokenType
	Buffer                 *[]scanner.Token
	Start, End, LineOffset int
}

// Lexeme returns the string representation of the Token.
func (t Token) Lexeme() string {
	var result string
	if t.Buffer == nil || len(*t.Buffer) < 1 || t.Start > t.End {
		return result
	}
	switch t.Type {
	case String:
		result = t.reprString()
	default:
		result = t.reprDefault()
	}
	return result
}

func (t Token) reprString() string {
	end := t.End
	head := t.Start
	size := t.End - t.Start
	if end > len(*t.Buffer) {
		end = len(*t.Buffer)
	}
	chars := make([]string, size)
	for head != end {
		token := (*t.Buffer)[head]
		// NOTE: For quoted strings, check if the current token initiates an
		// escape sequence and there is at least a single token left to look up
		// followed by the terminating quote character. Bare strings may not
		// contain escape sequence characters.
		if token.Rune == '\\' && head+2 != end {
			v, ok := escapeSequenceMap[(*t.Buffer)[head+1].Rune]
			if ok {
				token = (*t.Buffer)[head]
				head += 2
				chars = append(chars, v)
				continue
			}
		}
		chars = append(chars, string(token.Rune))
		head++
	}
	// NOTE: If a given string startswith ' or ", then trim the prefix and
	// trim the suffix: chars = chars[1:len(chars)-1].
	// TODO: Then remove Trim on the String from AST and use the regular Value.
	// TODO: Lexer should be able to track whether a given quote is escaped.
	return strings.Join(chars, "")
}

func (t Token) reprDefault() string {
	end := t.End
	size := t.End - t.Start
	if end > len(*t.Buffer) {
		end = len(*t.Buffer)
	}
	chars := make([]string, size)
	for _, t := range (*t.Buffer)[t.Start:end] {
		chars = append(chars, string(t.Rune))
	}
	return strings.Join(chars, "")
}
