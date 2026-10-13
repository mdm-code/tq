VERSION      := $(shell cat VERSION)
GO_FLAGS     := -trimpath
GO_LDFLAGS   := "-s -w -X main.version=${VERSION}"
COV_PROFILE  := coverage.txt
CGO_ENABLED  := 0

.DEFAULT_GOAL := build

ifndef VERBOSE
.SILENT:
endif

.PHONY: lint test install build cover clean

.PHONY: .ONESHELL
.ONESHELL:

lint:
	golangci-lint run -v

test: lint
	go clean -testcache
	go test ./... -v

install: test
	CGO_ENABLED=${CGO_ENABLED} go install ${GO_FLAGS} \
	-ldflags ${GO_LDFLAGS} -a github.com/mdm-code/tq/...

build: test
	CGO_ENABLED=${CGO_ENABLED} go build ${GO_FLAGS} \
	-ldflags ${GO_LDFLAGS} -a github.com/mdm-code/tq/...

cover:
	go test -coverprofile=${COV_PROFILE} -covermode=atomic ./...
	go tool cover -html=${COV_PROFILE}

clean:
	go clean github.com/mdm-code/tq/...
	go mod tidy
	go clean -testcache
	rm -f ${COV_PROFILE}
