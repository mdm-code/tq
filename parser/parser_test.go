package parser

import (
	"log"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/lexer"
)

func TestXxx(t *testing.T) {
	q := ". ['foo'][ 'bar' ][][0][:10][:][][2 : 12][\"foo\"] "
	r := strings.NewReader(q)
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal(err)
	}
	l, err := lexer.New(s)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: at this point whenever there's a lexer error and some tokens are
	// not returned for this reason, the query will be processed up until this
	// point, and this is wrong and has to be fixed.
	p, err := New(l, true) // the lexer errors are not handled in any way here
	if err != nil {
		t.Fatal(err)
	}
	e, err := p.Parse()
	var v AstPrinter
	log.Println(v.Print(e))
}
