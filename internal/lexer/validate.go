package lexer

import "unicode"

// isWhitespace verifies if the rune r is a whitespace character.
func isWhitespace(r rune) bool {
	return unicode.IsSpace(r)
}

// isDigit verifies if the rune r is a digit character.
func isDigit(r rune) bool {
	return unicode.IsDigit(r)
}

// isKeyChar verifies if the rune r is a key character.
func isKeyChar(r rune) bool {
	_, ok := keyCharMap[r]
	return ok
}

// isQuote verifies if the rune r is a quote character.
func isQuote(r rune) bool {
	return r == '"' || r == '\''
}

// isBareChar checks if the rune r is an accepted TOML bare key character.
func isBareChar(r rune) bool {
	if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '-' && r != '_' {
		return false
	}
	return true
}
