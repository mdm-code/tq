package main

import (
	_ "embed"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/mdm-code/tq/v2"
	"github.com/mdm-code/tq/v2/toml"
)

const (
	exitSuccess int = iota
	exitFailure
)

var (
	//go:embed usage.txt
	usage string

	query           string
	tablesInline    bool
	arraysMultiline bool
	indentSymbol    string
	indentTables    bool
)

func setupCLI(args []string) ([]string, error) {
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
	return fs.Args(), err
}

func setupTOMLAdapter() *toml.Adapter {
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

func run(args []string, input io.Reader, output io.Writer) (int, error) {
	args, err := setupCLI(args)
	if err != nil {
		return exitFailure, err
	}
	if len(args) > 0 {
		f, err := os.Open(args[0])
		input = f
		defer func() { f.Close() }()
		if err != nil {
			return exitFailure, err
		}
	}
	adapter := setupTOMLAdapter()
	tq := tq.New(adapter)
	err = tq.Run(input, output, query)
	if err != nil {
		return exitFailure, err
	}
	return exitSuccess, nil
}

func main() {
	exitCode, err := run(os.Args[1:], os.Stdin, os.Stdout)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	os.Exit(exitCode)
}
