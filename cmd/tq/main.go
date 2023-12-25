package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/lexer"
	"github.com/mdm-code/tq/parser"
	"github.com/pelletier/go-toml/v2"
)

func main() {
	query := flag.String("q", ".", "query")
	flag.Parse()
	r := strings.NewReader(*query)
	s, err := scanner.New(r)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	l, err := lexer.New(s)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	p, err := parser.New(l)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	e, err := p.Parse()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	qc := &parser.QueryConstructor{}
	qc.Interpret(e)
	var data interface{}
	in, _ := ioutil.ReadAll(os.Stdin)
	toml.Unmarshal(in, &data)
	d := []interface{}{data}
	for _, fn := range qc.Filters {
		d, _ = fn(d...)
	}
	for _, dd := range d {
		b, _ := toml.Marshal(dd)
		fmt.Fprintln(os.Stdout, string(b))
	}
}
