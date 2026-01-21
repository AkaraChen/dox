.PHONY: test build clean install

# Build variables
BINARY_NAME=do
BUILD_DIR=build
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

all: test build

## test: Run all tests with coverage
test:
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

## test-short: Run tests without race detection (faster)
test-short:
	$(GOTEST) -v -short ./...

## test-ci: Run tests for CI (with coverage threshold)
test-ci:
	$(GOTEST) -v -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -func=coverage.out | grep total | awk '{if ($$3+0 < 80) {print "Coverage below 80%"; exit 1}}'

## build: Build the binary
build:
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/main.go

## build-all: Build for multiple platforms
build-all:
	mkdir -p $(BUILD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/main.go
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/main.go
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/main.go
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/main.go

## install: Install the binary locally
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/main.go

## clean: Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

## deps: Download dependencies
deps:
	$(GOMOD) download
	$(GOMOD) tidy

## lint: Run golangci-lint
lint:
	golangci-lint run ./...

## fmt: Format code
fmt:
	gofmt -w -s .