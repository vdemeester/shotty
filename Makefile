BINARY = shotty
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS = -w -s -X main.version=$(VERSION)

.PHONY: build test lint clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

test:
	go test -v ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY)
