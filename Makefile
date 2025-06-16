.PHONY: check build clean client server fmt vet test lint

# Default target
check: fmt vet lint test

# Format code
fmt:
	go fmt ./...

# Vet code for issues
vet:
	go vet ./...

# Run golangci-lint
lint:
	golangci-lint run

# Run tests
test:
	go test ./...

# Build both client and server
build: client server

# Build client
client:
	go build -o bin/client ./client

# Build server  
server:
	go build -o bin/server ./server

# Clean build artifacts
clean:
	rm -rf bin/

# Tidy dependencies
tidy:
	go mod tidy