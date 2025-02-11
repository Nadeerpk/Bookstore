# Go environment variables
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOCLEAN := $(GOCMD) clean
BINARY_NAME := bookstore

# Default target: Build the binary
all: build

# Build the application
build:
	$(GOBUILD) -o $(BINARY_NAME) .

# Run the application
run: build
	./$(BINARY_NAME)

# Test the application
test:
	$(GOTEST) ./...

# Clean up generated files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Format the code
fmt:
	$(GOCMD) fmt ./...

# Install dependencies
deps:
	$(GOCMD) mod tidy

# Lint the code (requires golangci-lint)
lint:
	golangci-lint run ./...

.PHONY: all build run test clean fmt deps lint
