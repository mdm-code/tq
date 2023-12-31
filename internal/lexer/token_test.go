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
				buffer: &[]scanner.Token{
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
				start: 2,
				end:   9,
			},
			want: "\"tools\"",
		},
		{
			name: "nil-buffer",
			token: Token{
				Type:  Undefined,
				start: 0,
				end:   10,
			},
			want: "",
		},
		{
			name: "empty-buffer",
			token: Token{
				buffer: &[]scanner.Token{},
				Type:   Undefined,
				start:  0,
				end:    10,
			},
			want: "",
		},
		{
			name: "start-gt",
			token: Token{
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.'}, Buffer: nil},
				},
				Type:  Dot,
				start: 2,
				end:   1,
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
				buffer: &[]scanner.Token{
					{Pos: scanner.Pos{Rune: '.'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '['}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '8'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '0'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '2'}, Buffer: nil},
					{Pos: scanner.Pos{Rune: '4'}, Buffer: nil},
				},
				Type:  Integer,
				start: 2,
				end:   8,
			},
			want: "8024",
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
