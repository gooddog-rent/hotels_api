# Go parameters
GOCMD=go
GOARCH=amd64
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOTOOL=$(GOCMD) tool
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINST=$(GOCMD) install
PLATFORMS=darwin linux windows

MODULE_NAME=hotels_api

# Binary name
BINARY_NAME=hotels_api

# Build app
build:
	@$(GOBUILD) -o $(BINARY_NAME)_windows ./cmd
	@echo "📦 Build Done"

# Build app for all platforms
build_all:
	@$(foreach GOOS, $(PLATFORMS), ($(GOBUILD) -o $(BINARY_NAME)_$(GOOS) ./cmd))
	@echo "📦 Builds Done"

# Build and run
run:
	@echo "🚀 Running App"
	@./$(BINARY_NAME)

# Clean app binaries
clean:
	@echo "🧹 Clean App Binaries"
	@$(foreach GOOS, $(PLATFORMS), (shell rm $(BINARY_NAME)_$(GOOS)))

# Run unit tests
test: 
	go test -v -cover -race ./...
.PHONY: test

# Lint
lint:
	@golangci-lint run
	@echo "🔦 Code Linted"