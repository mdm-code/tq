package lexer

import (
	"testing"
)

// Test if whitespace characters are correctly identified.
func TestIsWhitespace(t *testing.T) {
	cases := []struct {
		input rune
		want  bool
		name  string
	}{
		{'\t', true, "\\t"},
		{' ', true, " "},
		{'\r', true, "\\r"},
		{'\n', true, "\\n"},
		{'a', false, "a"},
		{'0', false, "0"},
		{'"', false, "\""},
		{'\'', false, "'"},
		{'.', false, "."},
		{':', false, ":"},
		{'[', false, "["},
		{']', false, "]"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := isWhitespace(c.input); have != c.want {
				t.Errorf("want: %t; have: %t", c.want, have)
			}
		})
	}
}

// Verify if digit characters are correctly identified.
func TestIsDigit(t *testing.T) {
	cases := []struct {
		input rune
		want  bool
		name  string
	}{
		{'0', true, "0"},
		{'9', true, "9"},
		{'\t', false, "\\t"},
		{' ', false, " "},
		{'\r', false, "\\r"},
		{'\n', false, "\\n"},
		{'a', false, "a"},
		{'"', false, "\""},
		{'\'', false, "'"},
		{'.', false, "."},
		{':', false, ":"},
		{'[', false, "["},
		{']', false, "]"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := isDigit(c.input); have != c.want {
				t.Errorf("want: %t; have: %t", c.want, have)
			}
		})
	}
}

// Test if key characters are correctly identified.
func TestIsKeyChar(t *testing.T) {
	cases := []struct {
		input rune
		want  bool
		name  string
	}{
		{'.', true, "."},
		{':', true, ":"},
		{'[', true, "["},
		{']', true, "]"},
		{'\t', false, "\\t"},
		{' ', false, " "},
		{'\r', false, "\\r"},
		{'\n', false, "\\n"},
		{'a', false, "a"},
		{'0', false, "0"},
		{'"', false, "\""},
		{'\'', false, "'"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := isKeyChar(c.input); have != c.want {
				t.Errorf("want: %t; have: %t", c.want, have)
			}
		})
	}
}

// Check if quote characters are correctly identified.
func TestIsQuote(t *testing.T) {
	cases := []struct {
		input rune
		want  bool
		name  string
	}{
		{'"', true, "\""},
		{'\'', true, "'"},
		{'\t', false, "\\t"},
		{' ', false, " "},
		{'\r', false, "\\r"},
		{'\n', false, "\\n"},
		{'a', false, "a"},
		{'0', false, "0"},
		{'.', false, "."},
		{':', false, ":"},
		{'[', false, "["},
		{']', false, "]"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := isQuote(c.input); have != c.want {
				t.Errorf("want: %t; have: %t", c.want, have)
			}
		})
	}
}
