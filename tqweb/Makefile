HTMX_VERSION?=2.0.8

ifndef VERBOSE
.SILENT:
endif

.PHONY: assets build

.PHONY: .ONESHELL
.ONESHELL:

assets:
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js ./assets/js/
	curl https://cdn.jsdelivr.net/npm/htmx.org@$(HTMX_VERSION)/dist/htmx.min.js > ./assets/js/htmx.min.js

build: assets
	echo "building..."
