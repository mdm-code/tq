package lexer

import "unicode"

// IsNewline ...
func IsNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

// IsWhitespace ...
func IsWhitespace(r rune) bool {
	return r == ' ' || r == '\t'
}

// IsDigit ...
func IsDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// IsKeyChar ...
func IsKeyChar(r rune) bool {
	_, ok := KeyCharMap[r]
	return ok
}

// IsQuote ...
func IsQuote(r rune) bool {
	return r == '"' || r == '\''
}
