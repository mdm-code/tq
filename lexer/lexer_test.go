package lexer

import (
	"log"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
)

func TestXxx(t *testing.T) {
	r := strings.NewReader(".[123 : 1]['\"foo.bar\"'][\"\"][0]")
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal("failed to initialize the scanner")
	}
	l, err := New(s)
	if err != nil {
		t.Fatal("failed to initialize the lexer")
	}
	for l.Next() {
		log.Println(l.Token(), l.Token().Lexeme())
	}
	log.Println(l.Errors)
}
