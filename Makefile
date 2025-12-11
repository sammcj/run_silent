.PHONY: build test lint clean install

BINARY := run_silent
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -s -w -X main.version=$(VERSION)

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
