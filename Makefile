PROJECTNAME := $(shell basename "$(PWD)")
PACKAGES := $(shell go list ./... | grep -v vendor)

# Go related variables.
PROJECTROOT := $(shell pwd)
GOBIN := $(PROJECTROOT)/bin
GOFILES := $(PROJECTROOT)/*.go

# Shell script related variables.
UTILDIR := $(PROJECTROOT)/scripts/utils
SPINNER := $(UTILDIR)/spinner.sh
BUILDIR := $(PROJECTROOT)/scripts/build

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## install: Install missing dependencies
install:
	@$(SPINNER) "Installing project dependencies to vendor" "GOBIN=$(GOBIN) go get ./... && go mod vendor"
	@printf "\n👍 Done\n"

## build: Build the project binary
build:
	@$(SPINNER) "Building binary $(GOBIN)/$(PROJECTNAME)" "go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)"
	@printf "\n👍 Done\n"

## tools: Install development tools
tools:
	@$(SPINNER) "Installing fresh" $(BUILDIR)/install_fresh.sh
	@printf "\n👍 Done\n"

## start: Start in development mode with hot-reload enabled
start: tools
	@$(PROJECTROOT)/bin/fresh

## clean: Clean build files
clean:
	@$(SPINNER) "Cleaning build cache" "go clean $(PACKAGES)"
	@printf "\n👍 Done\n"
	@-rm $(GOBIN)/$(PROJECTNAME) 2>/dev/null

## fmt: Format entire codebase
fmt:
	@$(SPINNER) "Formatting" "go fmt $(PACKAGES)"
	@printf "\n👍 Done\n"

## vet: Vet entire codebase
vet:
	@$(SPINNER) "Vetting" "go vet $(PACKAGES)"
	@printf "\n👍 Done\n"

## lint: Check codebase for style mistakes
lint:
	@$(SPINNER) "Linting" "golint $(PACKAGES)"
	@printf "\n👍 Done\n"

## help: Display this help
help: Makefile
	@printf "\n Gasper: Your cloud in a binary\n\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@printf ""
