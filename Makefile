default: test

NAME=bf
BUILD_PKG=github.com/thesoulless/bf/cmd/bf

.PHONY: all
all: build test lint

.PHONY: build
build:
	@go build ./...

.PHONY: test
test:
	@go test -v ./...

.PHONY: lint
lint:
	@go install honnef.co/go/tools/cmd/staticcheck@HEAD
	@staticcheck ./...
