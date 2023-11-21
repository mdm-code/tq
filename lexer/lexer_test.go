package lexer

import (
	"errors"
	"fmt"
	"log"
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

func TestXxx(t *testing.T) {
	r := strings.NewReader(".[123 : 1]['\"foo.bar\"][\"\"][0]")
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal("failed to initialize the scanner")
	}
	l, err := New(s)
	if err != nil {
		t.Fatal("failed to initialize the lexer")
	}
	tt, ok := l.ScanAll(false)
	if !ok {
		fmt.Println(errors.Join(l.Errors...))
	}
	log.Println(tt)
}
