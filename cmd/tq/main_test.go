package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	args := []string{
		"-q",
		".fruit[0].color[].name",
		"-i",
		"-m",
		"-t",
		"-s",
		"  ",
	}
	tomlInput := `
[[fruit]]
  name = "apple"

  [fruit.geometry]
    shape = "round"

  [[fruit.color]]
    name = "red"

  [[fruit.color]]
    name = "green"

[[fruit]]
  name = "banana"

  [[fruit.color]]
    name = "yellow"
`
	input := strings.NewReader(tomlInput)
	var output bytes.Buffer
	exitCode, err := run(args, input, &output)
	if exitCode != exitSuccess {
		t.Errorf("have: %d; want: %d", exitCode, exitSuccess)
	}
	if err != nil {
		t.Errorf("should not return an error: %s", err)
	}
	want := `'red'
'green'
`
	if have := output.String(); have != want {
		t.Errorf("have: %s\nwant: %s", have, want)
	}
}
