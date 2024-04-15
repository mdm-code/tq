package toml

import (
	"fmt"
	"reflect"
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
	var have interface{}
	err := a.Unmarshal(r, &have)
	if err != nil {
		t.Error("unmarshal should not return an error")
	}
	want := map[string]any{"number": int64(13)}
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have: %v, want: %v", have, want)
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
	have, err := a.Marshal(v)
	if err != nil {
		t.Error("map encoding should not return an error")
	}
	want := []byte("age = 32\neducation = 'higher'\nname = 'Bob'\n")
	if !reflect.DeepEqual(have, want) {
		t.Errorf("have: %s; want: %s", string(have), string(want))
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
