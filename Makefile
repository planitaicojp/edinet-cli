BINARY := edinet
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X github.com/planitaicojp/edinet-cli/cmd.version=$(VERSION)

.PHONY: build test lint clean install all

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test ./... -v

lint:
	golangci-lint run ./...

fmt:
	gofmt -w .
	goimports -w .

clean:
	rm -f $(BINARY)
	rm -rf dist/

install:
	go install -ldflags "$(LDFLAGS)" .

all: lint test build
