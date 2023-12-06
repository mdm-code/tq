package lexer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
)

// Verify if the lexer gets initialized with a valid scanner.
func TestNewLexer(t *testing.T) {
	r := strings.NewReader("")
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal("failed to initialize the scanner")
	}
	_, err = New(s)
	if err != nil {
		t.Fatal("failed to initialize the lexer with a valid scanner")
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

func TestLexerScanError(t *testing.T) {
}

// Check if happy-path lexer ScanAll returns the expected output.
func TestLexerScanAll(t *testing.T) {
	r := strings.NewReader(".['package'][][9]")
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal("failed to initialize the scanner")
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
	have, ok := l.ScanAll(true)
	if !ok {
		t.Errorf("")
	}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have: %v; want %v", have, want)
	}
}
