include version.mk

ORGANIZATION := christophercampbell
PROJECT := bridge-connector

ARCH := $(shell uname -m)

ifeq ($(ARCH),x86_64)
	ARCH = amd64
else
	ifeq ($(ARCH),aarch64)
		ARCH = arm64
	endif
endif

GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/dist
GOOS := $(shell uname -s  | tr '[:upper:]' '[:lower:]')
GOENVVARS := GOBIN=$(GOBIN) CGO_ENABLED=1 GOOS=$(GOOS) GOARCH=$(ARCH)

GOBINARY := bridge-connector
GOCMD := $(GOBASE)/cmd

LDFLAGS += -X 'github.com/$(ORGANIZATION)/$(PROJECT).Version=$(VERSION)'
LDFLAGS += -X 'github.com/$(ORGANIZATION)/$(PROJECT).GitRev=$(GITREV)'
LDFLAGS += -X 'github.com/$(ORGANIZATION)/$(PROJECT).GitBranch=$(GITBRANCH)'
LDFLAGS += -X 'github.com/$(ORGANIZATION)/$(PROJECT).BuildDate=$(DATE)'

.PHONY: build
build: ## Builds the binary locally into ./dist
	$(GOENVVARS) go build -ldflags "all=$(LDFLAGS)" -o $(GOBIN)/$(GOBINARY) $(GOCMD)

.PHONY: build-docker
build-docker: ## Builds a docker image with the node binary
	docker build -t $(PROJECT) -f ./Dockerfile .

.PHONY: lint
lint: ## Runs the linter
	export "GOROOT=$$(go env GOROOT)" && $$(go env GOPATH)/bin/golangci-lint run

.PHONY:
clean: ## clean build artifacts
	rm -rf $(GOBIN)

## Help display.
## Pulls comments from beside commands and prints a nicely formatted
## display with the commands and their usage information.
.DEFAULT_GOAL := help

.PHONY: help
help: ## Prints this help
		@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
