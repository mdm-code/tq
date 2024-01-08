package parser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/lexer"
)

// Check if the nil buffer results in a bare error line.
func TestErrorNilBuffer(t *testing.T) {
	err := Error{
		buffer: nil,
		err:    ErrQueryElement,
	}
	want := fmt.Sprintf("Parser error: %s but got ''", err.err)
	have := err.Error()
	if have != want {
		t.Errorf("have: %s; want: %s", have, want)
	}
}

// Check if the empty buffer results in a bare error line.
func TestErrorEmptyBuffer(t *testing.T) {
	err := Error{
		buffer: &[]scanner.Token{},
		err:    ErrSelectorUnterminated,
	}
	want := fmt.Sprintf("Parser error: %s but got ''", err.err)
	have := err.Error()
	if have != want {
		t.Errorf("have: %s; want: %s", have, want)
	}
}

// Check if the Error string matches the expected output.
func TestErrorError(t *testing.T) {
	cases := []struct {
		name string
		want string
		err  error
	}{
		{
			name: "EOL",
			want: "Parser error: nil",
			err: &Error{
				lexeme: "EOL",
				err:    nil,
			},
		},
		{
			name: "]",
			want: ".['foo']]\n        ^\nParser error: expected '.' or '[' to parse query element but got ']'",
			err: &Error{
				lexeme: "]",
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.', Start: 0, End: 1}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '[', Start: 1, End: 2}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '\'', Start: 2, End: 3}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'f', Start: 3, End: 4}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'o', Start: 4, End: 5}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'o', Start: 5, End: 6}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '\'', Start: 6, End: 7}, Buffer: nil},
					{Pos: scanner.Pos{Rune: ']', Start: 7, End: 8}, Buffer: nil},
					{Pos: scanner.Pos{Rune: ']', Start: 8, End: 9}, Buffer: nil},
				},
				offset: 8,
				err:    ErrQueryElement,
			},
		},
		{
			name: "persons",
			want: "Parser error: expected ']' to terminate selector but got 'persons'",
			err: &Error{
				lexeme: "persons",
				err:    ErrSelectorUnterminated,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := c.err.Error(); have != c.want {
				t.Errorf("want: %s; have: %s", c.want, have)
			}
		})
	}
}

// Verify the error comparison using the errors.Is function.
func TestErrorIs(t *testing.T) {
	cases := []struct {
		name string
		want error
	}{
		{name: "ErrQueryElement", want: ErrQueryElement},
		{name: "ErrSelectorUnterminated", want: ErrSelectorUnterminated},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			token := lexer.Token{}
			err := &Error{"", token.Buffer, token.Start, c.want}
			if !errors.Is(err, c.want) {
				t.Error("expected the underlying error to match parser error")
			}
		})
	}
}
