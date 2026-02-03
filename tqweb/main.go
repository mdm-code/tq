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

/*
TODO:
1. Index has to be rewritten imediately using htmx and the index api endpoint.
2. I can use structs with templates pretty much the same way it's done with templ.
3. Change in query or input triggers processing action with short debounce.
4. The output and ALL errors are displayed in the output box.
5. If there's an error, it's printed out the same way as in the terminal.
-------------------------
|          |            |
|  query   |            |
|          |            |
|          |            |
------------   output   |
|          |            |
|  input   |            |
|          |            |
|          |            |
-------------------------
6. It seems the only thing I need to replace is the output box; query and input can stay as is.
7. That said, I don't need to add any hightlight to query and input or output for that matter.
8, Output stays read-only.
*/

var (
	//go:embed views/index.html
	index  string
	indexT = template.Must(template.New("index").Parse(index))
)

func main() {
	conf := toml.GoTOMLConf{}
	goToml := toml.NewGoTOML(conf)
	adapter := toml.NewAdapter(goToml)
	tq := tq.New(adapter)

	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		var output bytes.Buffer
		var outputStr string
		input := r.FormValue("input")
		query := r.FormValue("query")
		err := tq.Run(strings.NewReader(input), &output, query)
		if err != nil {
			outputStr = err.Error()
		} else {
			outputStr = output.String()
		}
		data := struct {
			Query, Input, Output string
		}{
			Query:  query,
			Output: outputStr,
			Input:  input,
		}
		if err := indexT.ExecuteTemplate(w, "form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	wasmhttp.Serve(nil)

	select {}
}
