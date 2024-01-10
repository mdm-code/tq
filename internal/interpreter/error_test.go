package interpreter

import (
	"errors"
	"testing"
)

// Test if interpreter Error has the expected string representation.
func TestErrorError(t *testing.T) {
	cases := []struct {
		name   string
		data   any
		filter string
		err    error
		want   string
	}{
		{
			name:   "interface",
			data:   map[string]any{"x": "y"},
			filter: "string \"persons\"",
			err:    ErrTOMLDataType,
			want: "Interpreter error: cannot query " +
				"[ map[string]interface {} ] " +
				"( map[x:y] ) " +
				"with ( string \"persons\" )",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := &Error{
				data:   c.data,
				filter: c.filter,
				err:    c.err,
			}
			if have := err.Error(); have != c.want {
				t.Errorf("have: %s; want: %s", have, c.want)
			}
		})
	}
}

// Test if the Error matches the embedded error with errors.Is().
func TestErrorIs(t *testing.T) {
	cases := []struct {
		want error
	}{
		{want: ErrTOMLDataType},
	}
	for _, c := range cases {
		t.Run(c.want.Error(), func(t *testing.T) {
			err := &Error{err: c.want}
			if !errors.Is(err, c.want) {
				t.Errorf("the errors should match: %s : %s", err, c.want)
			}
		})
	}
}
