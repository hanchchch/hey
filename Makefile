##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Development

.PHONY: test
TEST_ARGS ?= -v
TEST_TARGETS ?= ./...
test: ## Test the Go modules within this package.
	@ echo ▶️ go test $(TEST_ARGS) $(TEST_TARGETS)
	go test $(TEST_ARGS) $(TEST_TARGETS)
	@ echo ✅ success!

.PHONY: lint
LINT_TARGETS ?= ./...
lint: ## Lint Go code with the installed golangci-lint
	@ echo "▶️ golangci-lint run"
	golangci-lint run $(LINT_TARGETS)
	@ echo "✅ golangci-lint run"


##@ Build

BINARY_NAME=hey
CMD_HEY=cmd/hey.go
BUILD_OS=darwin linux windows
AMD64_BINS=$(BUILD_OS:%=bin/amd64/%/${BINARY_NAME})
ARM64_BINS=$(BUILD_OS:%=bin/arm64/%/${BINARY_NAME})

all: $(AMD64_BINS) $(ARM64_BINS) ## Build all binaries for all platforms

$(ARM64_BINS):
	GOARCH=$(word 2,$(subst /, ,$@)) GOOS=$(word 3,$(subst /, ,$@)) go build -o $@ $(CMD_HEY)

$(AMD64_BINS):
	GOARCH=$(word 2,$(subst /, ,$@)) GOOS=$(word 3,$(subst /, ,$@)) go build -o $@ $(CMD_HEY)
