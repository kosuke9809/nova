# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=nova
BINARY_UNIX=$(BINARY_NAME)_unix

# Build targets
all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/nova

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/nova
	./$(BINARY_NAME)

deps:
	$(GOGET) ./...
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Linting
GOLINT=golangci-lint
lint:
	$(GOLINT) run