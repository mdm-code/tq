package toml

import (
	"fmt"
	"strings"
	"testing"
)

type FailingReader struct{}

func (FailingReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("errored")
}

// Test Adapter Unmarshal does not return an error when given right arguments.
func TestAdapterUnmarshal(t *testing.T) {
	var conf GoTOMLConf
	goToml := NewGoTOML(conf)
	a := NewAdapter(goToml)
	r := strings.NewReader("number=13")
	var v interface{}
	err := a.Unmarshal(r, &v)
	if err != nil {
		t.Error("unmarshal should not return an error")
	}
	switch v.(type) {
	case map[string]any:
	default:
		t.Error("unmarshal was to assign a map of interface values to v")
	}
}

// Test Adapter Marshal errors when read operation fails.
func TestAdapterUnmarshalError(t *testing.T) {
	var conf GoTOMLConf
	goToml := NewGoTOML(conf)
	a := NewAdapter(goToml)
	r := FailingReader{}
	var v interface{}
	err := a.Unmarshal(r, &v)
	if err == nil {
		t.Error("unmarshal should fail when read operation fails")
	}
}

// Test Adapter Marshal does not return an error when given correct arguments.
func TestAdapterMarshal(t *testing.T) {
	var conf GoTOMLConf
	goToml := NewGoTOML(conf)
	a := NewAdapter(goToml)
	v := map[string]any{"name": "Bob", "age": 32, "education": "higher"}
	_, err := a.Marshal(v)
	if err != nil {
		t.Error("map encoding should not return an error")
	}
}

// Test Adapter Marshal errors with incorrect arguments.
func TestAdapterMarshalError(t *testing.T) {
	var conf GoTOMLConf
	goToml := NewGoTOML(conf)
	a := NewAdapter(goToml)
	var v interface{}
	_, err := a.Marshal(v)
	if err == nil {
		t.Error("encode should not pass when given an interface value")
	}
}
