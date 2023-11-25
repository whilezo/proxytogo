# Build the proxy binary
build:
	go build -o bin/proxy cmd/main.go

# Run the proxy server
run:
	./bin/proxy

# Clean build artifacts
clean:
	rm -rf bin/

# Install project dependencies
deps:
	go mod download

# Run tests
test:
	go test ./...

# Generate code documentation
docs:
	godoc -http=:6060 -goroot=. -play

# Display help information
help:
	@echo "Available targets:"
	@echo "  build      - Build the proxy binary"
	@echo "  run        - Run the proxy server"
	@echo "  clean      - Clean build artifacts"
	@echo "  deps       - Install project dependencies"
	@echo "  test       - Run tests"
	@echo "  docs       - Generate code documentation"
	@echo "  help       - Display help information"

.PHONY: build run start clean deps test docs help