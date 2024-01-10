package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/interpreter"
	"github.com/mdm-code/tq/internal/lexer"
	"github.com/mdm-code/tq/internal/parser"
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
	i := interpreter.New()
	execFn := i.Interpret(e)
	var data any
	in, _ := ioutil.ReadAll(os.Stdin)
	toml.Unmarshal(in, &data)
	d, err := execFn(data)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	for _, dd := range d {
		b, _ := toml.Marshal(dd)
		if len(b) <= 0 {
			continue
		}
		fmt.Fprintln(os.Stdout, string(b))
	}
}
