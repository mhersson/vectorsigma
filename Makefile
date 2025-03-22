##
# vectorsigma Makefile
#

shell=bash

VERSION = $(shell git describe --tags --always)
COMMIT = $(shell git rev-parse --short HEAD)
BUILDTIME = $(shell date -u '+%Y-%m-%dT%H:%M:%SZ'.1.0)

# make will interpret non-option arguments in the command line as targets.
# This turns them into do-nothing targets, so make won't complain:
# If the first argument is "run"...
ifeq (run,$(firstword $(MAKECMDGOALS)))
# use the rest as arguments for "run"
	RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
# ...and turn them into do-nothing targets
	$(eval $(RUN_ARGS):;@:)
endif

LDFLAGS="-s -w \
	-X github.com/mhersson/vectorsigma/cmd.Version=$(VERSION) \
	-X github.com/mhersson/vectorsigma/cmd.CommitSHA=$(COMMIT) \
	-X github.com/mhersson/vectorsigma/cmd.BuildTime=$(BUILDTIME)"

all: build

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

.PHONY: version
version:
	@echo $(VERSION)

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

vet: ## Run go vet
	go vet ./...

build: fmt vet ## Build the binary.
	@go build -ldflags $(LDFLAGS) -o vectorsigma

install: ## Install the binary.
	@go install -ldflags $(LDFLAGS)

test: ## Run tests.
	@go test ./pkgs/generator ./pkgs/uml ./internal/statemachine ./cmd --coverprofile=cover.out

run: ## Run main.go with arguments.
	@go run -ldflags $(LDFLAGS) ./main.go $(RUN_ARGS)

golden: ## Update golden files.
	go test ./cmd -coverprofile cover.out -args -update

docker-image: ## Build Docker image.
	@docker buildx build --build-arg VERSION=$(VERSION) -t vectorsigma:$(VERSION) .
