GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=gofmt
BINARY_NAME=certbot-pdns-proxy

.DEFAULT_GOAL := all
.PHONY: all build build-linux-amd64 test check-fmt fmt clean run deps

all: check-fmt fmt test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME)_linux_amd64 -v

test:
	$(GOTEST) -v ./...

check-fmt:
	$(GOFMT) -d ./

fmt:
	$(GOFMT) -w ./

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)_linux_amd64

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) -t -v ./...
