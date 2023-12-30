MODULE_PKG := github.com/jeffrydegrande/acid

.PHONY: default
default: help

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

ci: lint test security vet staticcheck ## Run all CI checks
.PHONY: ci

.PHONY: security
security: ## Run security checks
	@ echo "▶️ gg run github.com/securego/gosec/v2/cmd/gosec -fmt golint ./..."
	@ go run github.com/securego/gosec/v2/cmd/gosec -fmt golint ./...
	@ echo "✅ go run github.com/securego/gosec/v2/cmd/gosec -fmt golint ./..."

.PHONY: lint
LINT_ARGS ?= --enable-all
LINT_TARGETS ?= ./...
lint: ## Lint Go code with the installed golangci-lint
	@ echo "▶️ golangci-lint run $(LINT_ARGS) $(LINT_TARGETS)"
	golangci-lint run $(LINT_ARGS) $(LINT_TARGETS)
	@ echo "✅ golangci-lint run"

vet:
	go vet
.PHONY: vet

.PHONY: staticcheck
STATICCHECK_TARGETS ?= ./...
staticcheck: ## Run staticcheck linter
	@ echo "▶️ gstaticcheck $(STATICCHECK_TARGETS)"
	CGO_ENABLED=0 staticcheck $(STATICCHECK_TARGETS)
	@ echo "✅ staticcheck $(STATICCHECK_TARGETS)"

.PHONY: test
TEST_TARGETS ?= ./...
TEST_ARGS ?= -v -coverprofile=coverage.txt
test: ## Test the Go modules within this package.
	@ echo ▶️ go test $(TEST_ARGS) $(TEST_TARGETS)
	go test $(TEST_ARGS) $(TEST_TARGETS)
	@ echo ✅ success!

	@ echo ▶️ go tool  cover -func=coverage.txt
	go tool cover -func=coverage.txt
	@ echo ✅ success!
