# Makefile for building a Go project

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=popci

# Directories
SRC_DIR=./src
BIN_DIR=./bin

# Build the project
all: build

build:
	$(GOBUILD) -o $(BIN_DIR)/$(BINARY_NAME) $(SRC_DIR)/...

clean:
	$(GOCLEAN)
	rm -f $(BIN_DIR)/$(BINARY_NAME)

test:
	$(GOTEST) -v ./...

get:
	$(GOGET) -v ./...

.PHONY: all build clean test get