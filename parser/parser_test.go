package parser

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/lexer"
	"github.com/pelletier/go-toml/v2"
)

func TestXxx(t *testing.T) {
	// q := ". ['foo'][ 'bar' ][][0][:10][:][][2 : 12][\"foo\"] "
	// q := "['foo'][:1'bar']"
	q := ".['nestedDict'][][:2]"
	r := strings.NewReader(q)
	s, err := scanner.New(r)
	if err != nil {
		t.Fatal(err)
	}
	l, err := lexer.New(s)
	if err != nil {
		t.Fatal(err)

	}
	p, err := New(l)
	if err != nil {
		fmt.Println(err)
		t.Fatal()
	}
	e, err := p.Parse()
	qc := &QueryConstructor{}
	qc.Run(e)
	var data interface{}
	val := `[nestedDict]
foo = [1, 2, 3]
bar = [1, 2, 3]
`
	toml.Unmarshal([]byte(val), &data)
	d := []interface{}{data}
	log.Println(qc.Filters, err)
	for _, fn := range qc.Filters {
		d, _ = fn(d...)
	}
	log.Println(d)
}
