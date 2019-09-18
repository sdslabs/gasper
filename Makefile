PROJECTNAME := $(shell basename "$(PWD)")
PACKAGES := $(shell go list ./... | grep -v vendor)

# Go related variables.
GOPATH := $(shell pwd)
GOBIN := $(GOPATH)/bin
GOFILES := $(GOPATH)/*.go

# Shell script related variables.
UTILDIR := $(GOPATH)/scripts/utils
SPINNER := $(UTILDIR)/spinner.sh
BUILDIR := $(GOPATH)/scripts/build

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: default
default: help

## install: Install missing dependencies.
install:
	@$(SPINNER) "Installing project dependencies to vendor" "GOBIN=$(GOBIN) go get ./... && go mod vendor"
	@printf "\n👍 Done\n"

## build: Build the project binary.
build:
	@$(SPINNER) "Building binary $(GOBIN)/$(PROJECTNAME)" "go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)"
	@printf "\n👍 Done\n"

## tools: Install development tools.
tools:
	@$(SPINNER) "Installing fresh" $(BUILDIR)/install_fresh.sh
	@printf "\n👍 Done\n"

## start: Start in development mode. Auto-reloads when code changes.
start: tools
	@$(GOPATH)/bin/fresh

## clean: Clean build files.
clean:
	@$(SPINNER) "Cleaning build cache" "go clean $(PACKAGES)"
	@printf "\n👍 Done\n"
	@-rm $(GOBIN)/$(PROJECTNAME) 2>/dev/null

## lint: Lint entire codebase.
lint:
	@$(SPINNER) "Formatting" "go fmt $(PACKAGES)"
	@printf "\n👍 Done\n"
	@$(SPINNER) "Vetting" "go vet $(PACKAGES)"
	@printf "\n👍 Done\n"

## help: Display this help.
help: Makefile
	@printf "\n Choose a command to run in "$(PROJECTNAME)":\n\n"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@printf "\n"
