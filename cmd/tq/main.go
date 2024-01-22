package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/interpreter"
	"github.com/mdm-code/tq/internal/lexer"
	"github.com/mdm-code/tq/internal/parser"
	"github.com/mdm-code/tq/internal/toml"
)

const (
	exitSuccess int = iota
	exitFailure
)

var (
	usage = `tq - query TOML configuration files

Usage:

	tq [-qtmsi] [file...]

Options:

	-h, --help              show this help message and exit
	-q, --query             query to run against the input data (default: '.')
	-t, --tables-inline     emit tables inline (default: false)
	-m, --arrays-multiline  emit arrays one element per line (default: false)
	-s, --indent-symbol     provide the indentation string (default: '  ')
	-i, --indent-tables     indent tables and array tables (default: false)

Example:

	tq -q '["servers"][]["ip"]' <<EOF
	[servers]

	[servers.prod]
	ip = "10.0.0.1"
	role = "backend"

	[servers.staging]
	ip = "10.0.0.2"
	role = "backend"
	EOF

Output:

	'10.0.0.1'
	'10.0.0.2'

Tq is a tool for querying TOML configuration files with a sequence of intuitive
filters. It works as a regular Unix filter program reading input data from the
standard input and producing results to the standard output.
`

	query           string
	tablesInline    bool
	arraysMultiline bool
	indentSymbol    string
	indentTables    bool
)

func setupCLI(args []string) error {
	fs := flag.NewFlagSet("tq", flag.ExitOnError)
	fs.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprint(w, usage)
	}

	queryUsage := "query to run against the input data"
	fs.StringVar(&query, "q", ".", queryUsage)
	fs.StringVar(&query, "query", ".", queryUsage)

	tablesInlineUsage := "emit tables inline"
	fs.BoolVar(&tablesInline, "tables-inline", false, tablesInlineUsage)
	fs.BoolVar(&tablesInline, "t", false, tablesInlineUsage)

	arraysMultilineDefault := "emit arrays one element per line"
	fs.BoolVar(&arraysMultiline, "arrays-multiline", false, arraysMultilineDefault)
	fs.BoolVar(&arraysMultiline, "m", false, arraysMultilineDefault)

	indentSymbolDefault := "provide the indentation string"
	fs.StringVar(&indentSymbol, "indent-symbol", "  ", indentSymbolDefault)
	fs.StringVar(&indentSymbol, "s", "  ", indentSymbolDefault)

	indentTablesDefault := "indent tables and array tables"
	fs.BoolVar(&indentTables, "indent-tables", false, indentTablesDefault)
	fs.BoolVar(&indentTables, "i", false, indentTablesDefault)

	err := fs.Parse(args)
	return err
}

func setupTOMLAdapter() toml.Adapter {
	conf := toml.GoTOMLConf{
		Encoder: struct {
			TablesInline    bool
			ArraysMultiline bool
			IndentSymbol    string
			IndentTables    bool
		}{
			tablesInline,
			arraysMultiline,
			indentSymbol,
			indentTables,
		},
	}

	goToml := toml.NewGoTOML(conf)
	adapter := toml.NewAdapter(goToml)
	return adapter
}

func run(args []string) (int, error) {
	err := setupCLI(args)
	if err != nil {
		return exitFailure, err
	}
	reader := strings.NewReader(query)
	scanner, err := scanner.New(reader)
	if err != nil {
		return exitFailure, err
	}
	lexer, err := lexer.New(scanner)
	if err != nil {
		return exitFailure, err
	}
	parser, err := parser.New(lexer)
	if err != nil {
		return exitFailure, err
	}
	ast, err := parser.Parse()
	if err != nil {
		return exitFailure, err
	}
	interpreter := interpreter.New()
	exec := interpreter.Interpret(ast)
	var data any
	tomlAdapter := setupTOMLAdapter()
	err = tomlAdapter.Unmarshal(os.Stdin, &data)
	if err != nil {
		return exitFailure, err
	}
	filteredData, err := exec(data)
	if err != nil {
		return exitFailure, err
	}
	for _, d := range filteredData {
		bytes, err := tomlAdapter.Marshal(d)
		if err != nil {
			return exitFailure, err
		}
		if len(bytes) == 0 {
			continue
		}
		fmt.Fprintln(os.Stdout, string(bytes))
	}
	return exitSuccess, nil
}

func main() {
	exitCode, err := run(os.Args[1:])
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	os.Exit(exitCode)
}
