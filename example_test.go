package tq_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mdm-code/tq/v2"
	"github.com/mdm-code/tq/v2/toml"
)

// ExampleTq_Run demonstrates how to use the Tq struct to run a query against
// TOML data. The example uses a TOML configuration file with two servers and
// queries for IP addresses on the first server. The output is written to the
// standard output.
func ExampleTq_Run() {
	input := strings.NewReader(`
[servers]

[servers.alpha]
ip = "10.0.0.1"
role = "frontend"

[servers.beta]
ip = "10.0.0.2"
role = "backend"
`)
	var output bytes.Buffer
	query := `
    .servers
        .alpha
            .ip
`
	config := toml.GoTOMLConf{}
	goToml := toml.NewGoTOML(config)
	adapter := toml.NewAdapter(goToml)
	tq := tq.New(adapter)
	_ = tq.Run(input, &output, query)
	fmt.Println(output.String())
	// Output:
	// 10.0.0.1
}

// ExampleTq_Validate shows how to use the Tq struct to validate whether a
// given query is syntactically correct. The example shows how the error is
// reported and represented as a string.
func ExampleTq_Validate() {
	query := "['servers'][['ip']"
	config := toml.GoTOMLConf{}
	goToml := toml.NewGoTOML(config)
	adapter := toml.NewAdapter(goToml)
	tq := tq.New(adapter)
	err := tq.Validate(query)
	fmt.Println(err)
	// Output:
	// ['servers'][['ip']
	//             ^
	// Parser error: expected ']' to terminate selector; got '['
}
