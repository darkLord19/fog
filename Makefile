.PHONY: build test install clean run lint

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	@echo "Building wtx $(VERSION)..."
	@go build $(LDFLAGS) -o bin/wtx ./cmd/wtx

test:
	@echo "Running tests..."
	@go test -v ./...

install:
	@echo "Installing wtx..."
	@go install $(LDFLAGS) ./cmd/wtx

clean:
	@echo "Cleaning..."
	@rm -rf bin/ dist/

run:
	@go run ./cmd/wtx

lint:
	@echo "Running linters..."
	@go fmt ./...
	@go vet ./...

dev: build
	@./bin/wtx
