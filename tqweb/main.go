//go:build js && wasm

package main

import (
	"bytes"
	_ "embed"
	"html/template"
	"net/http"
	"strings"

	"github.com/mdm-code/tq/v2"
	"github.com/mdm-code/tq/v2/toml"
	wasmhttp "github.com/nlepage/go-wasm-http-server/v2"
)

var (
	//go:embed views/index.html
	index  string
	indexT = template.Must(template.New("index").Parse(index))
)

func setupTomlAdapter(r *http.Request) *toml.Adapter {
	tablesInline := r.FormValue("tables-inline") != ""
	arraysMultiline := r.FormValue("arrays-multiline") != ""
	indentSymbol := r.FormValue("indent-symbol")
	if indentSymbol == "" {
		indentSymbol = "  "
	}
	indentTables := r.FormValue("indent-tables") != ""
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

func main() {
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		var output bytes.Buffer
		var data struct{ Output string }

		adapter := setupTomlAdapter(r)
		t := tq.New(adapter)

		input := r.FormValue("input")
		query := r.FormValue("query")
		err := t.Run(strings.NewReader(input), &output, query)
		if err != nil {
			data.Output = err.Error()
		} else {
			data.Output = output.String()
		}
		if err := indexT.ExecuteTemplate(w, "output", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	wasmhttp.Serve(nil)

	select {}
}
