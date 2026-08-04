GO=go
GOFLAGS=-mod=vendor -trimpath
HTMX_VERSION?=2.0.8
# NOTE
# ====
# #1 The master version of the Javascript code fom the go-wasm-http-server has
# been updated since the version v2.2.1. If for any reason it might be necessary
# to revert to the version used in Go code, here is a quick shell script to
# pull the version from the go.mod file.
#
# GO_WASM_HTTP_SERVER_VERSION?=master
#
# #2 The master version seems not to work right with Go 1.26.
#
GO_WASM_HTTP_SERVER_VERSION?=$(shell\
  grep go-wasm-http-server go.mod |\
  cut -d' ' -f2\
)
# ====

export CGO_ENABLED=0
export GOOS=js
export GOARCH=wasm

ifndef VERBOSE
.SILENT:
endif

.PHONY: assets build serve clean

.PHONY: .ONESHELL
.ONESHELL:

assets:
	cp $$(go env GOROOT)/lib/wasm/wasm_exec.js wasm_exec.js
	curl https://cdn.jsdelivr.net/npm/htmx.org@$(HTMX_VERSION)/dist/htmx.min.js > htmx.min.js
	curl https://cdn.jsdelivr.net/gh/nlepage/go-wasm-http-server@$(GO_WASM_HTTP_SERVER_VERSION)/sw.js > sw.js

build: clean
	$(GO) build $(GOFLAGS) -ldflags '-s -w' -o api.wasm .
	if command -v sha256sum >/dev/null 2>&1; then \
		VERSION=$$(sha256sum api.wasm | cut -c1-8); \
	else \
		VERSION=$$(shasum -a 256 api.wasm | cut -c1-8); \
	fi; \
	mv api.wasm api.$${VERSION}.wasm; \
	sed -e "s|__WASM__|api.$${VERSION}.wasm|g" \
	    -e "s|__CACHE__|api-$${VERSION}|g" \
	    worker.js.in > worker.js

clean:
	rm -f api.wasm api.*.wasm worker.js

serve:
	python3 -m http.server 8000
