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

	printVersion    bool
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
		_, _ = fmt.Fprint(w, usage)
	}

	fs.BoolVar(&printVersion, "v", false, "")
	fs.BoolVar(&printVersion, "version", false, "")

	fs.StringVar(&query, "q", ".", "")
	fs.StringVar(&query, "query", ".", "")

	fs.BoolVar(&tablesInline, "tables-inline", false, "")
	fs.BoolVar(&tablesInline, "t", false, "")

	fs.BoolVar(&arraysMultiline, "arrays-multiline", false, "")
	fs.BoolVar(&arraysMultiline, "m", false, "")

	fs.StringVar(&indentSymbol, "indent-symbol", "  ", "")
	fs.StringVar(&indentSymbol, "s", "  ", "")

	fs.BoolVar(&indentTables, "indent-tables", false, "")
	fs.BoolVar(&indentTables, "i", false, "")

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
	if printVersion {
		w := flag.CommandLine.Output()
		_, _ = fmt.Fprintln(w, tq.GetVersion())
		return exitSuccess, nil
	}
	if len(args) > 0 {
		f, err := os.Open(args[0])
		input = f
		defer func() { _ = f.Close() }()
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
