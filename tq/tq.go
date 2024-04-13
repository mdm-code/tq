package tq

import (
	"io"

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

func (t *Tq) Run() {
}
