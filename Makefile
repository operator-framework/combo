ORG := github.com/operator-framework
PKG := $(ORG)/combo
VERSION_PATH := $(PKG)/pkg/version
GIT_COMMIT := $(shell git rev-parse HEAD)

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

KUBERNETES_VERSION=v0.22.2
VERSION_FLAGS=-ldflags "-X $(VERSION_PATH).GitCommit=$(GIT_COMMIT) -X $(VERSION_PATH).ComboVersion=`cat COMBO_VERSION` -X $(VERSION_PATH).KubernetesVersion=$(KUBERNETES_VERSION)"
build-cli: ## Build the CLI binary. Speciy VERSION_PATH, GIT_COMMIT, or KUBERNETES_VERSION to change the binary version.
	$(Q)go build -a $(VERSION_FLAGS) -o ./bin/combo

IMAGE_REPO=quay.io/operator-framework/combo
IMAGE_TAG=latest
build-container: ## Build the Combo container. Accepts IMAGE_REPO and IMAGE_TAG overrides.
	docker build . -f Dockerfile -t $(IMAGE_REPO):$(IMAGE_TAG)

CONTROLLER_GEN=$(Q)go run sigs.k8s.io/controller-tools/cmd/controller-gen

# Static tests.
.PHONY: test test-unit verify build

test: test-unit ## Run the tests

test-unit: ## Run the unit tests
	$(Q)go test -count=1 -short ./...

deploy: generate ## Deploy the Combo operator to the current cluster
	kubectl apply --recursive -f manifests

teardown: ## Teardown the Combo operator to the current cluster
	kubectl delete --recursive -f manifests

# Binary builds
GO_BUILD := $(Q)go build
