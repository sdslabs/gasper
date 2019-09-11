PROJECTNAME := $(shell basename "$(PWD)")
PACKAGES := $(shell go list ./... | grep -v vendor)

# Go related variables.
GOPATH := $(shell pwd)
GOBIN := $(GOPATH)/bin
GOFILES := $(GOPATH)/cmd/*.go

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## install: Install missing dependencies.
install:
	@echo "  >  Installing project dependencies to vendor..."
	@GOBIN=$(GOBIN) go get ./...
	@go mod vendor

## build: Build the project binary.
build:
	@echo "  >  Building binary..."
	@go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)
	@echo "  >  Path of the generated binary is $(GOBIN)/$(PROJECTNAME)"

## tools: Install development tools.
tools:
	@./scripts/install_fresh.sh

## start: Start in development mode. Auto-reloads when code changes.
start: tools
	@$(GOPATH)/bin/fresh

## clean: Clean build files.
clean:
	@echo "  >  Deleting project binary"
	@-rm $(GOBIN)/$(PROJECTNAME) 2> /dev/null
	@echo "  >  Cleaning build cache"
	@go clean

# lint: Lint code using gofmt and govet.
lint:
	@echo "  >  Formatting..."
	@go fmt $(PACKAGES)
	@echo "  >  Vetting..."
	@go vet $(PACKAGES)

help: Makefile
	@echo
	@echo " Choose a command to run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
