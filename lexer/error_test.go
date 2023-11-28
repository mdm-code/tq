package lexer

import (
	"fmt"
	"testing"

	"github.com/mdm-code/scanner"
)

// Check if the nil buffer results in a bare error line.
func TestErrorNilBuffer(t *testing.T) {
	err := Error{
		buffer: nil,
		offset: 0,
		err:    ErrDisallowedChar,
	}
	want := fmt.Sprintf("Lexer error: %s", err.err)
	have := err.Error()
	if have != want {
		t.Errorf("have: %s; want: %s", have, want)
	}
}

// Check if the empty buffer results in a bare error line.
func TestErrorEmptyBuffer(t *testing.T) {
	err := Error{
		buffer: &[]scanner.Token{},
		offset: 0,
		err:    ErrUnterminatedString,
	}
	want := fmt.Sprintf("Lexer error: %s", err.err)
	have := err.Error()
	if have != want {
		t.Errorf("have: %s; want: %s", have, want)
	}
}

// Test error string constructed upon calling Error().
func TestErrorError(t *testing.T) {
	cases := []struct {
		name string
		want string
		err  error
	}{
		{
			"nil error",
			".['foo]\n  ^\nLexer error: nil",
			&Error{
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.', Start: 0, End: 1}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '[', Start: 1, End: 2}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '\'', Start: 2, End: 3}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'f', Start: 3, End: 4}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'o', Start: 4, End: 5}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'o', Start: 5, End: 6}, Buffer: nil},
					{Pos: scanner.Pos{Rune: ']', Start: 6, End: 7}, Buffer: nil},
				},
				offset: 3,
				err:    nil,
			},
		},
		{
			"negative offset",
			".\n^\nLexer error: nil",
			&Error{
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.', Start: 0, End: 1}, Buffer: nil},
				},
				offset: -1,
				err:    nil,
			},
		},
		{
			"empty buffer",
			"Lexer error: nil",
			&Error{
				buffer: &[]scanner.Token{},
				offset: 0,
				err:    nil,
			},
		},
		{
			"nil buffer",
			"Lexer error: nil",
			&Error{
				buffer: nil,
				offset: 0,
				err:    nil,
			},
		},
		{
			"disallowed character error",
			".?\n ^\nLexer error: disallowed character",
			&Error{
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.', Start: 0, End: 1}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '?', Start: 1, End: 2}, Buffer: nil},
				},
				offset: 2,
				err:    ErrDisallowedChar,
			},
		},
		{
			"unterminated string error",
			".['g\n  ^\nLexer error: unterminated string literal",
			&Error{
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.', Start: 0, End: 1}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '[', Start: 1, End: 2}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '\'', Start: 2, End: 3}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'g', Start: 3, End: 4}, Buffer: nil},
				},
				offset: 3,
				err:    ErrUnterminatedString,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := c.err.Error(); have != c.want {
				t.Errorf("\nWant:\n%s\nHave:\n%s", have, c.want)
			}
		})
	}
}
