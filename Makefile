.PHONY: build test lint clean install

BINARY := run_silent
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)

build:
	go build -ldflags="$(LDFLAGS)" -o $(BINARY) .

test:
	go test -v -race ./...

lint:
	golangci-lint run ./...

clean:
	rm -f $(BINARY)

install: build
	cp $(BINARY) $(HOME)/.local/bin/
