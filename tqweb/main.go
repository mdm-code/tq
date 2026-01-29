//go:build js && wasm

package main

import (
	_ "embed"
	"html/template"
	"net/http"

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
	http.HandleFunc("/process", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{"query": r.FormValue("query")}
		if err := indexT.ExecuteTemplate(w, "form", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Instantiate TQ.
	// var err error
	conf := toml.GoTOMLConf{}
	goToml := toml.NewGoTOML(conf)
	adapter := toml.NewAdapter(goToml)
	_ = tq.New(adapter)
	// err = tq.Run(input, output, query)
	// The error from the parser, lexer, scanner, toml adapter etc. -- all of these can be passed to the template in place of the output.
	// This way the user is able to see that, first, he gets an error without any flashy effects on the front end, and second, the output
	// is placed in a visible place.

	wasmhttp.Serve(nil)

	select {}
}
