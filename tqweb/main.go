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
1. Tq config for GoTOML is sent from the frontend selection component.
- Checkboxes: send true if checked.
- Buttons: on click, keep highlighted and add to the payload.
2. The config and TQ has to be instantiated with each call.
3. Curl tailwind css scripts.
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
		var data struct{ Output string }
		input := r.FormValue("input")
		query := r.FormValue("query")
		err := tq.Run(strings.NewReader(input), &output, query)
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
