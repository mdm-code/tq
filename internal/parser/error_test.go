package parser

import (
	"errors"
	"testing"
)

// Check if the Error string matches the expected output.
func TestErrorError(t *testing.T) {
	cases := []struct {
		lexeme string
		err    error
		want   string
	}{
		{
			lexeme: "EOF",
			err:    nil,
			want:   "Parser error: nil",
		},
		{
			lexeme: "]",
			err:    ErrQueryElement,
			want:   "Parser error: expected '.' or '[' to parse query element but got ']'",
		},
		{
			lexeme: "persons",
			err:    ErrSelectorUnterminated,
			want:   "Parser error: expected ']' to terminate selector but got 'persons'",
		},
	}
	for _, c := range cases {
		t.Run(c.lexeme, func(t *testing.T) {
			err := Error{c.lexeme, c.err}
			have := err.Error()
			if have != c.want {
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
			err := &Error{"", c.want}
			if !errors.Is(err, c.want) {
				t.Error("expected the underlying error to match parser error")
			}
		})
	}
}
