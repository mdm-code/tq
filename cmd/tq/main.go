/*
tq - query TOML configuration files

Usage:

	tq [-q|--query arg...] [file...]

Options:

	-h, --help   show this help message and exit
	-q, --query  specify the query to run against the input data (default: '.')

Example:

	tq -q '["servers"][]["ip"]' <<EOF
	[servers]

	[servers.alpha]
	ip = "10.0.0.1"
	role = "frontend"

	[servers.beta]
	ip = "10.0.0.2"
	role = "backend"
	EOF

Output:

	'10.0.0.1'
	'10.0.0.2'

Tq is a tool for querying TOML configuration files with a sequence of intuitive
filters. It works as a regular Unix filter program reading input data from the
standard input and producing results to the standard output.
*/
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
