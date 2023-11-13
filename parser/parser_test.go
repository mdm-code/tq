package parser

import (
	"log"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/lexer"
)

func TestXxx(t *testing.T) {
	q := " . ['foo'][ 'bar' ][0][:10][:][][2 : 12][\"foo\"] "
	r := strings.NewReader(q)
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal(err)
	}
	l, err := lexer.New(s)
	if err != nil {
		t.Fatal(err)
	}
	p, err := New(l, true)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(p.root())
}
