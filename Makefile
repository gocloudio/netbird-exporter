.PHONY: build test clean docker

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
BINARY_NAME=netbird-exporter
DOCKER_TAG=gocloudio/netbird-exporter:latest

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./cmd/netbird-exporter

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

docker:
	docker build -t $(DOCKER_TAG) .

run:
	./$(BINARY_NAME)

tidy:
	$(GOMOD) tidy

# Cross-compilation
build-linux:
	mkdir -p output
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o output/$(BINARY_NAME)-linux-amd64 -v ./cmd/netbird-exporter

build-arm:
	mkdir -p output	
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o output/$(BINARY_NAME)-linux-arm64 -v ./cmd/netbird-exporter 