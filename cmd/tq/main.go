package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/lexer"
	"github.com/mdm-code/tq/parser"
)

func main() {
	s, err := scanner.New(os.Stdin)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	l, err := lexer.New(s)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	p, err := parser.New(l)
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	e, err := p.Parse()
	var a parser.AstPrinter
	fmt.Fprintln(os.Stdout, a.Print(e))
}
