package toml

import (
	"bytes"
	"io"

	"github.com/pelletier/go-toml/v2"
)

// Decoder ...
type Decoder interface {
	Decode(io.Reader, any) error
}

// Encoder ...
type Encoder interface {
	Encode(any) ([]byte, error)
}

// DecodeEncoder ...
type DecodeEncoder interface {
	Decoder
	Encoder
}

// Adapter ...
type Adapter struct {
	adapted DecodeEncoder
}

// NewAdapter ...
func NewAdapter(adapted DecodeEncoder) Adapter {
	return Adapter{adapted: adapted}
}

// Unmarshal ...
func (a *Adapter) Unmarshal(r io.Reader, v any) error {
	return a.adapted.Decode(r, v)
}

// Marshal ...
func (a *Adapter) Marshal(v any) ([]byte, error) {
	return a.adapted.Encode(v)
}

// GoTOML ...
type GoTOML struct {
	conf GoTOMLConf
}

// NewGoTOML ...
func NewGoTOML(c GoTOMLConf) GoTOML {
	return GoTOML{conf: c}
}

// Decode ...
func (t GoTOML) Decode(r io.Reader, v any) error {
	d := toml.NewDecoder(r)
	if t.conf.Decoder.Strict {
		d.DisallowUnknownFields()
	}
	return d.Decode(v)
}

// Encode ...
func (t GoTOML) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer
	e := toml.NewEncoder(&buf)
	e.SetTablesInline(t.conf.Encoder.TablesInline)
	e.SetArraysMultiline(t.conf.Encoder.ArraysMultiline)
	e.SetIndentSymbol(t.conf.Encoder.IndentSymbol)
	e.SetIndentTables(t.conf.Encoder.IndentTables)
	err := e.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// GoTOMLConf ...
type GoTOMLConf struct {
	Decoder struct {
		Strict bool
	}
	Encoder struct {
		TablesInline    bool
		ArraysMultiline bool
		IndentSymbol    string
		IndentTables    bool
	}
}
