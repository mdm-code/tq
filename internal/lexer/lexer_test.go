package lexer

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
)

// Verify if the lexer gets initialized with a valid scanner.
func TestNewLexer(t *testing.T) {
	cases := []struct {
		name   string
		reader io.Reader
		err    error
	}{
		{
			name:   "valid reader",
			reader: strings.NewReader(""),
			err:    nil,
		},
		{
			name:   "invalid string",
			reader: strings.NewReader("a\xc5z"),
			err:    scanner.ErrRuneError,
		},
		{
			name:   "nil scanner",
			reader: nil,
			err:    ErrNilScanner,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s, _ := scanner.New(c.reader)
			_, err := New(s)
			if !errors.Is(err, c.err) {
				t.Fatal("new lexer does not return the expected error")
			}
		})
	}
}

// Check if the happy-path lexer scan result equals the expected output.
func TestLexerScan(t *testing.T) {
	r := strings.NewReader(".['package'][][9]")
	s, err := scanner.New(r)
	if err != nil {
		t.Error("failed to initialize the scanner")
	}
	l, err := New(s)
	if err != nil {
		t.Error("failed to initialize the lexer with a valid scanner")
	}
	want := []Token{
		{Dot, &l.buffer, 0, 1},
		{ArrayOpen, &l.buffer, 1, 2},
		{String, &l.buffer, 2, 11},
		{ArrayClose, &l.buffer, 11, 12},
		{ArrayOpen, &l.buffer, 12, 13},
		{ArrayClose, &l.buffer, 13, 14},
		{ArrayOpen, &l.buffer, 14, 15},
		{Integer, &l.buffer, 15, 16},
		{ArrayClose, &l.buffer, 16, 17},
	}
	have := []Token{}
	for l.Scan() {
		have = append(have, l.Token())
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have: %v; want %v", have, want)
	}
}

// Check if declared errors are accumulated by the lexer given the query input.
func TestLexerScanError(t *testing.T) {
	cases := []struct {
		name  string
		query string
		want  error
	}{
		{
			name:  "unterminated-string",
			query: ".[495][ 'field ] ",
			want:  ErrUnterminatedString,
		},
		{
			name:  "dissallowed-char",
			query: "['texts'][  ]['chars'].[$1] 22\r",
			want:  ErrDisallowedChar,
		},
		{
			name:  "dissallowed-char-in-string",
			query: "['parent\r']",
			want:  ErrDisallowedChar,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.query)
			s, err := scanner.New(r)
			if err != nil {
				t.Fatal("failed to initialize the scanner")
			}
			l, err := New(s)
			if err != nil {
				t.Fatal("failed to initialize the lexer with a valid scanner")
			}
			_, ok := l.ScanAll(false)
			if ok {
				t.Error("lexer unexpectedly scanned all tokens without errors")
			}
			if len(l.Errors) == 0 {
				t.Error("lexer did not accumulated any errors")
			}
			have := errors.Join(l.Errors...)
			if !errors.Is(have, c.want) {
				t.Errorf("have: %v; want: %v", have, c.want)
			}
		})
	}
}

// Check if happy-path lexer ScanAll returns the expected output.
func TestLexerScanAll(t *testing.T) {
	cases := []struct {
		name             string
		query            string
		ignoreWhitespace bool
		want             []Token
	}{
		{
			name:             "whitespace ignored",
			query:            ". [ 'package' ][][ 9 ] ",
			ignoreWhitespace: true,
			want: []Token{
				{Dot, nil, 0, 1},
				{ArrayOpen, nil, 2, 3},
				{String, nil, 4, 13},
				{ArrayClose, nil, 14, 15},
				{ArrayOpen, nil, 15, 16},
				{ArrayClose, nil, 16, 17},
				{ArrayOpen, nil, 17, 18},
				{Integer, nil, 19, 20},
				{ArrayClose, nil, 21, 22},
			},
		},
		{
			name:             "whitespace included",
			query:            ". [ 'package' ][][ 9 ] ",
			ignoreWhitespace: false,
			want: []Token{
				{Dot, nil, 0, 1},
				{Whitespace, nil, 1, 2},
				{ArrayOpen, nil, 2, 3},
				{Whitespace, nil, 3, 4},
				{String, nil, 4, 13},
				{Whitespace, nil, 13, 14},
				{ArrayClose, nil, 14, 15},
				{ArrayOpen, nil, 15, 16},
				{ArrayClose, nil, 16, 17},
				{ArrayOpen, nil, 17, 18},
				{Whitespace, nil, 18, 19},
				{Integer, nil, 19, 20},
				{Whitespace, nil, 20, 21},
				{ArrayClose, nil, 21, 22},
				{Whitespace, nil, 22, 23},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := strings.NewReader(c.query)
			s, err := scanner.New(r)
			if err != nil {
				t.Fatal("failed to initialize the scanner")
			}
			l, err := New(s)
			if err != nil {
				t.Error("failed to initialize the lexer with a valid scanner")
			}
			have, ok := l.ScanAll(c.ignoreWhitespace)
			if !ok {
				t.Errorf("failed to scan valid query")
			}
			if len(have) != len(c.want) {
				t.Error("have tokens: %i; want: %i", len(have), len(c.want))
			}
			for i := 0; i < len(c.want)-1; i++ {
				have[i].Buffer = nil
				if have[i] != c.want[i] {
					t.Errorf("have: %v; want: %v", have, c.want)
				}
			}
		})
	}
}
