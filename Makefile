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
	@INSTALL_DIR=""; \
	if [ -n "$$GOPATH" ]; then \
		INSTALL_DIR="$$GOPATH/bin"; \
	elif echo "$$PATH" | tr ':' '\n' | grep -qx "$$HOME/.local/bin"; then \
		INSTALL_DIR="$$HOME/.local/bin"; \
	elif [ -w /usr/local/bin ]; then \
		INSTALL_DIR="/usr/local/bin"; \
	else \
		echo "Error: No suitable install directory found"; \
		echo "  - GOPATH is not set"; \
		echo "  - $$HOME/.local/bin is not in PATH"; \
		echo "  - /usr/local/bin is not writable"; \
		exit 1; \
	fi; \
	mkdir -p "$$INSTALL_DIR" && \
	cp $(BINARY) "$$INSTALL_DIR/" && \
	echo "Installed to: $$INSTALL_DIR/$(BINARY)"
