package tq_test

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mdm-code/tq"
	"github.com/mdm-code/tq/toml"
)

// ExampleTq_Run demonstrates how to use the Tq struct to run a query against
// TOML data. The example uses a TOML configuration file with two servers and
// queries the IP address of the beta server. The output is written to the
// output.
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
	query := "['servers']['beta']['ip']"
	config := toml.GoTOMLConf{}
	goToml := toml.NewGoTOML(config)
	adapter := toml.NewAdapter(goToml)
	tq := tq.New(adapter)
	_ = tq.Run(input, &output, query)
	fmt.Println(output.String())
	// Output:
	// '10.0.0.2'
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
	// Parser error: expected ']' to terminate selector but got '['
}
