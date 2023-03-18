SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec
GO := GOPROXY=https://proxy.golang.org go
GOR := goreleaser
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

all: clean fmt vet test build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: ## Builds the binary file
	$(GOR) build --clean --snapshot

.PHONY: run
run: ## Runs main help
	$(GO) run main.go

.PHONY: test
test: ## Placeholder to run unit tests
	@echo "Running unit tests"
	@mkdir -p bin
	$(GO) test -cover -coverprofile=bin/c.out $$( go list ./... | egrep -v 'mocks|qqq|vendor|exec|cmd' )
	$(GO) tool cover -html=bin/c.out -o bin/coverage.html
	@echo

.PHONY: check
check: ## Runs linting
	@echo "Linting"
	@for d in $$(go list ./... | egrep -v 'vendor|asdf'); do staticcheck $${d}; done
	@echo

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet ./...

.PHONY: mod
mod: ## Run go mod tidy/vendor
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: clean
clean: ## Removes build artifacts from source code
	@echo "Cleaning"
	@rm -fr bin
	@rm -fr vendor
	@echo

.PHONY: update-here
update-here: ## Helper target to start editing all occurances with UPDATE_HERE.
	@echo "Update following files for release:"
	@grep --color -nHR UPDATE_HERE .
