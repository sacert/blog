.PHONY: test test-coverage test-short benchmark run build

# Default target
all: test build

# Build the application
build:
	go build -o blog-app

# Run the application
run: build
	./blog-app

# Run all tests
test:
	go test ./... -v

# Run tests with code coverage
test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

# Run only unit tests (skip integration tests)
test-short:
	go test ./... -v -short

# Run benchmarks
benchmark:
	go test -bench=. -benchmem

# Clean build artifacts
clean:
	rm -f blog-app
	rm -f coverage.out coverage.html

# Setup - downloads dependencies
setup:
	go mod download
