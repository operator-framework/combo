# kernel-style V=1 build verbosity
ifeq ("$(origin V)", "command line")
  BUILD_VERBOSE = $(V)
endif

ifeq ($(BUILD_VERBOSE),1)
  Q =
else
  Q = @
endif


.PHONY: help
help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# Code management
.PHONY: lint format tidy generate build

PKGS = $(shell go list ./...)

tidy: ## Update dependencies
	$(Q)go mod tidy

generate: ## Generate code and manifests
	$(Q)go generate ./...

format: ## Format the source code
	$(Q)go fmt ./...

lint: ## Run golangci-lint
	$(Q)go run github.com/golangci/golangci-lint/cmd/golangci-lint run

verify: tidy generate format lint ## Verify the current code generation and lint
	git diff --exit-code

build-cli: ## Build the CLI binary
	$(Q)go build -a -o ./bin/combo

IMAGE_REPO=quay.io/operator-framework/combo
IMAGE_TAG=dev
build-container: ## Build the Combo container. Accepts IMAGE_REPO and IMAGE_TAG overrides.
	docker build . -f Dockerfile -t $(IMAGE_REPO):$(IMAGE_TAG)

CONTROLLER_GEN=$(Q)go run sigs.k8s.io/controller-tools/cmd/controller-gen

# Static tests.
.PHONY: test test-unit verify build

test: test-unit ## Run the tests

test-unit: ## Run the unit tests
	$(Q)go test -count=1 -short ./...

# Binary builds
GO_BUILD := $(Q)go build
