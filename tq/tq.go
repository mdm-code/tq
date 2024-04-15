package tq

import (
	"fmt"
	"io"
	"strings"

	"github.com/mdm-code/scanner"
	"github.com/mdm-code/tq/internal/interpreter"
	"github.com/mdm-code/tq/internal/lexer"
	"github.com/mdm-code/tq/internal/parser"
	"github.com/mdm-code/tq/internal/toml"
)

type Tq struct {
	input   io.Reader
	output  io.Writer
	adapter toml.Adapter
}

func New(input io.Reader, output io.Writer, adapter toml.Adapter) *Tq {
	return &Tq{
		input:   input,
		output:  output,
		adapter: adapter,
	}
}

func (t *Tq) Run(query string) error {
	reader := strings.NewReader(query)
	scanner, err := scanner.New(reader)
	if err != nil {
		return err
	}
	lexer, err := lexer.New(scanner)
	if err != nil {
		return err
	}
	parser, err := parser.New(lexer)
	if err != nil {
		return err
	}
	ast, err := parser.Parse()
	if err != nil {
		return err
	}
	interpreter := interpreter.New()
	exec := interpreter.Interpret(ast)
	var data any
	err = t.adapter.Unmarshal(t.input, &data)
	if err != nil {
		return err
	}
	filteredData, err := exec(data)
	if err != nil {
		return err
	}
	for _, d := range filteredData {
		bytes, err := t.adapter.Marshal(d)
		if err != nil {
			return err
		}
		if len(bytes) == 0 {
			continue
		}
		fmt.Fprintln(t.output, string(bytes))
	}
	return nil
}
