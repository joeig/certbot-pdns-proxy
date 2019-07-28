GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
GODEP=dep
BINARY_NAME=certbot-pdns-proxy
GOFILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.DEFAULT_GOAL := all
.PHONY: all build build-linux-amd64 test check-fmt fmt clean run deps

all: check-fmt fmt test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/certbot-pdns-proxy

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux_amd64 -v ./cmd/certbot-pdns-proxy

test:
	$(GOTEST) -v ./...

check-fmt:
	$(GOFMT) -d ${GOFILES}

fmt:
	$(GOFMT) -w ${GOFILES}

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_linux_amd64

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GODEP) ensure
