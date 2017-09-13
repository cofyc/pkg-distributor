
.EXPORT_ALL_VARIABLES:

pkgs = $(shell go list ./... | grep -v /vendor/)

all:
	go build github.com/cofyc/pkg-distributor/cmd/pkg-distributor

linux: GOOS=linux
linux: GOARCH=amd64
linux: all

test:
	go test $(pkgs)
.PHONY: test

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'
.PHONY: style

fmt:
	@echo ">> formatting code"
	@go fmt $(pkgs)
.PHONY: fmt
