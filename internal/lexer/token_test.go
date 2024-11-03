package lexer

import (
	"testing"

	"github.com/mdm-code/scanner"
)

// Test that given Token state Lexeme() value receiver function return value
// matches the expected output.
func TestLexeme(t *testing.T) {
	cases := []struct {
		name  string
		token Token
		want  string
	}{
		{
			name: "default",
			token: Token{
				Buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '['}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '"'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 't'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'o'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'o'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'l'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 's'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '"'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: ']'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '.'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '['}, Buffer: nil},
					{Pos: scanner.Pos{Rune: ']'}, Buffer: nil},
				},
				Type:  String,
				Start: 2,
				End:   9,
			},
			want: "\"tools\"",
		},
		{
			name: "nil-buffer",
			token: Token{
				Type:  Undefined,
				Start: 0,
				End:   10,
			},
			want: "",
		},
		{
			name: "empty-buffer",
			token: Token{
				Buffer: &[]scanner.Token{},
				Type:   Undefined,
				Start:  0,
				End:    10,
			},
			want: "",
		},
		{
			name: "start-gt",
			token: Token{
				Buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.'}, Buffer: nil},
				},
				Type:  Dot,
				Start: 2,
				End:   1,
			},
			want: "",
		},
		{
			name:  "bare-token",
			token: Token{},
			want:  "",
		},
		{
			name: "shorter-buffer",
			token: Token{
				Buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '['}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '8'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '0'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '2'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '4'}, Buffer: nil},
				},
				Type:  Integer,
				Start: 2,
				End:   8,
			},
			want: "8024",
		},
		{
			name: "bare-string",
			token: Token{
				Buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: 'n'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'a'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'm'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: 'e'}, Buffer: nil},
				},
				Type:  String,
				Start: 0,
				End:   4,
			},
			want: "name",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if have := c.token.Lexeme(); have != c.want {
				t.Errorf("want: %s; have %s", c.want, have)
			}
		})
	}
}
