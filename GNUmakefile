# Terraform Provider HashiCorp-OVH Enterprise Makefile
# This Makefile provides comprehensive build, test, and release automation

SHELL := /bin/bash
.DEFAULT_GOAL := help

# Project configuration
PROJECT_NAME := terraform-provider-hashicorp-ovh
PROVIDER_NAME := hashicorp-ovh
NAMESPACE := swcstudio
BINARY_NAME := terraform-provider-hashicorp-ovh
PKG := github.com/swcstudio/terraform-provider-hashicorp-ovh

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT ?= $(shell git rev-parse HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

# Go configuration
GO_VERSION := 1.21
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
CGO_ENABLED := 0

# Build flags
LDFLAGS := -s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(BUILD_DATE)
BUILD_FLAGS := -trimpath -ldflags="$(LDFLAGS)"

# Directories
DIST_DIR := dist
COVERAGE_DIR := coverage
DOCS_DIR := docs
EXAMPLES_DIR := examples
TOOLS_DIR := tools

# Test configuration
TEST_TIMEOUT := 120m
TEST_PARALLEL := 4
COVERAGE_THRESHOLD := 80

# Terraform configuration
TF_VERSION := 1.6.0
TF_ACC := 1
TF_LOG := DEBUG

# Tool versions
GOLANGCI_LINT_VERSION := v1.54.2
GORELEASER_VERSION := latest
TFPLUGINDOCS_VERSION := latest

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
PURPLE := \033[0;35m
CYAN := \033[0;36m
NC := \033[0m # No Color

# OS Detection
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Linux)
    OS := linux
    OPEN := xdg-open
endif
ifeq ($(UNAME_S),Darwin)
    OS := darwin
    OPEN := open
endif

##@ General

.PHONY: help
help: ## Display this help message
	@echo -e "$(BLUE)Terraform Provider HashiCorp-OVH$(NC)"
	@echo -e "$(CYAN)Enterprise Build System$(NC)"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: version
version: ## Show version information
	@echo -e "$(BLUE)Version Information:$(NC)"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"
	@echo "  Date:    $(BUILD_DATE)"
	@echo "  Go:      $(shell go version)"
	@echo "  OS/Arch: $(GOOS)/$(GOARCH)"

##@ Development

.PHONY: setup
setup: ## Set up development environment
	@echo -e "$(BLUE)Setting up development environment...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo -e "$(RED)Go is not installed$(NC)"; exit 1; }
	@go mod download
	@go mod tidy
	@$(MAKE) tools
	@echo -e "$(GREEN)Development environment ready!$(NC)"

.PHONY: tools
tools: ## Install required development tools
	@echo -e "$(BLUE)Installing development tools...$(NC)"
	@mkdir -p $(TOOLS_DIR)
	@GOBIN=$(PWD)/$(TOOLS_DIR) go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@GOBIN=$(PWD)/$(TOOLS_DIR) go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@$(TFPLUGINDOCS_VERSION)
	@GOBIN=$(PWD)/$(TOOLS_DIR) go install github.com/goreleaser/goreleaser@$(GORELEASER_VERSION)
	@GOBIN=$(PWD)/$(TOOLS_DIR) go install golang.org/x/vuln/cmd/govulncheck@latest
	@GOBIN=$(PWD)/$(TOOLS_DIR) go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
	@GOBIN=$(PWD)/$(TOOLS_DIR) go install github.com/fzipp/gocyclo/cmd/gocyclo@latest
	@echo -e "$(GREEN)Development tools installed!$(NC)"

.PHONY: clean
clean: ## Clean build artifacts and caches
	@echo -e "$(BLUE)Cleaning build artifacts...$(NC)"
	@rm -rf $(DIST_DIR)
	@rm -rf $(COVERAGE_DIR)
	@rm -rf $(TOOLS_DIR)
	@rm -f $(BINARY_NAME)
	@rm -f coverage.out coverage.html
	@rm -f *.log
	@go clean -cache -modcache -testcache
	@echo -e "$(GREEN)Clean completed!$(NC)"

##@ Building

.PHONY: build
build: ## Build the provider binary
	@echo -e "$(BLUE)Building $(BINARY_NAME)...$(NC)"
	@CGO_ENABLED=$(CGO_ENABLED) go build $(BUILD_FLAGS) -o $(BINARY_NAME) .
	@echo -e "$(GREEN)Build completed: $(BINARY_NAME)$(NC)"

.PHONY: build-all
build-all: ## Build for all supported platforms
	@echo -e "$(BLUE)Building for all platforms...$(NC)"
	@mkdir -p $(DIST_DIR)
	@for os in linux darwin windows freebsd; do \
		for arch in amd64 arm64; do \
			if [ "$$os" = "darwin" ] && [ "$$arch" = "386" ]; then continue; fi; \
			if [ "$$os" = "windows" ] && [ "$$arch" = "arm64" ]; then continue; fi; \
			echo "Building $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch CGO_ENABLED=0 go build $(BUILD_FLAGS) \
				-o $(DIST_DIR)/$(BINARY_NAME)_$$os-$$arch \
				. || exit 1; \
		done; \
	done
	@echo -e "$(GREEN)Multi-platform build completed!$(NC)"

.PHONY: install
install: build ## Install provider locally for development
	@echo -e "$(BLUE)Installing provider locally...$(NC)"
	@mkdir -p ~/.terraform.d/plugins/registry.terraform.io/$(NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS)_$(GOARCH)
	@cp $(BINARY_NAME) ~/.terraform.d/plugins/registry.terraform.io/$(NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS)_$(GOARCH)/
	@echo -e "$(GREEN)Provider installed locally!$(NC)"
	@echo "Location: ~/.terraform.d/plugins/registry.terraform.io/$(NAMESPACE)/$(PROVIDER_NAME)/$(VERSION)/$(OS)_$(GOARCH)/"

##@ Testing

.PHONY: test
test: ## Run unit tests
	@echo -e "$(BLUE)Running unit tests...$(NC)"
	@go test -v -race -timeout=30m ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	@echo -e "$(BLUE)Running tests with coverage...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@go test -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	@go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@go tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total | awk '{print "Total coverage: " $$3}'
	@echo -e "$(GREEN)Coverage report generated: $(COVERAGE_DIR)/coverage.html$(NC)"

.PHONY: test-coverage-check
test-coverage-check: test-coverage ## Check if coverage meets threshold
	@echo -e "$(BLUE)Checking coverage threshold ($(COVERAGE_THRESHOLD)%)...$(NC)"
	@COVERAGE=$$(go tool cover -func=$(COVERAGE_DIR)/coverage.out | grep total | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$COVERAGE >= $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo -e "$(GREEN)✓ Coverage ($$COVERAGE%) meets threshold ($(COVERAGE_THRESHOLD)%)$(NC)"; \
	else \
		echo -e "$(RED)✗ Coverage ($$COVERAGE%) below threshold ($(COVERAGE_THRESHOLD)%)$(NC)"; \
		exit 1; \
	fi

.PHONY: test-verbose
test-verbose: ## Run tests with verbose output
	@echo -e "$(BLUE)Running verbose tests...$(NC)"
	@go test -v -race -timeout=30m -count=1 ./...

.PHONY: test-short
test-short: ## Run short tests only
	@echo -e "$(BLUE)Running short tests...$(NC)"
	@go test -v -race -short ./...

.PHONY: test-integration
test-integration: build ## Run integration tests
	@echo -e "$(BLUE)Running integration tests...$(NC)"
	@cd $(EXAMPLES_DIR)/local-dev && \
		terraform init && \
		terraform validate && \
		terraform plan

.PHONY: testacc
testacc: ## Run acceptance tests (requires API credentials)
	@echo -e "$(BLUE)Running acceptance tests...$(NC)"
	@echo -e "$(YELLOW)Warning: This will create real resources!$(NC)"
	@TF_ACC=$(TF_ACC) go test -v -timeout $(TEST_TIMEOUT) -parallel $(TEST_PARALLEL) ./internal/provider/

.PHONY: test-benchmarks
test-benchmarks: ## Run benchmark tests
	@echo -e "$(BLUE)Running benchmark tests...$(NC)"
	@go test -bench=. -benchmem -run=^$$ ./...

##@ Quality Assurance

.PHONY: lint
lint: ## Run golangci-lint
	@echo -e "$(BLUE)Running golangci-lint...$(NC)"
	@if [ -f $(TOOLS_DIR)/golangci-lint ]; then \
		$(TOOLS_DIR)/golangci-lint run; \
	else \
		golangci-lint run; \
	fi

.PHONY: lint-fix
lint-fix: ## Run golangci-lint with auto-fix
	@echo -e "$(BLUE)Running golangci-lint with auto-fix...$(NC)"
	@if [ -f $(TOOLS_DIR)/golangci-lint ]; then \
		$(TOOLS_DIR)/golangci-lint run --fix; \
	else \
		golangci-lint run --fix; \
	fi

.PHONY: fmt
fmt: ## Format Go code
	@echo -e "$(BLUE)Formatting Go code...$(NC)"
	@gofmt -s -w .
	@goimports -w .

.PHONY: fmt-check
fmt-check: ## Check if code is formatted
	@echo -e "$(BLUE)Checking code formatting...$(NC)"
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo -e "$(RED)The following files are not formatted:$(NC)"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	else \
		echo -e "$(GREEN)All files are properly formatted$(NC)"; \
	fi

.PHONY: vet
vet: ## Run go vet
	@echo -e "$(BLUE)Running go vet...$(NC)"
	@go vet ./...

.PHONY: complexity
complexity: ## Check cyclomatic complexity
	@echo -e "$(BLUE)Checking cyclomatic complexity...$(NC)"
	@if [ -f $(TOOLS_DIR)/gocyclo ]; then \
		$(TOOLS_DIR)/gocyclo -over 15 .; \
	else \
		echo -e "$(YELLOW)gocyclo not installed, run 'make tools'$(NC)"; \
	fi

##@ Security

.PHONY: security
security: security-gosec security-vuln ## Run all security checks

.PHONY: security-gosec
security-gosec: ## Run gosec security scanner
	@echo -e "$(BLUE)Running gosec security scanner...$(NC)"
	@if [ -f $(TOOLS_DIR)/gosec ]; then \
		$(TOOLS_DIR)/gosec -severity medium -confidence medium -quiet ./...; \
	else \
		echo -e "$(YELLOW)gosec not installed, run 'make tools'$(NC)"; \
	fi

.PHONY: security-vuln
security-vuln: ## Run Go vulnerability check
	@echo -e "$(BLUE)Running vulnerability check...$(NC)"
	@if [ -f $(TOOLS_DIR)/govulncheck ]; then \
		$(TOOLS_DIR)/govulncheck ./...; \
	else \
		echo -e "$(YELLOW)govulncheck not installed, run 'make tools'$(NC)"; \
	fi

##@ Documentation

.PHONY: docs
docs: ## Generate provider documentation
	@echo -e "$(BLUE)Generating provider documentation...$(NC)"
	@if [ -f $(TOOLS_DIR)/tfplugindocs ]; then \
		$(TOOLS_DIR)/tfplugindocs generate; \
	else \
		tfplugindocs generate; \
	fi
	@echo -e "$(GREEN)Documentation generated in $(DOCS_DIR)/$(NC)"

.PHONY: docs-check
docs-check: docs ## Check if documentation is up to date
	@echo -e "$(BLUE)Checking documentation status...$(NC)"
	@if git diff --quiet $(DOCS_DIR)/; then \
		echo -e "$(GREEN)Documentation is up to date$(NC)"; \
	else \
		echo -e "$(RED)Documentation is out of date$(NC)"; \
		git diff $(DOCS_DIR)/; \
		exit 1; \
	fi

.PHONY: docs-serve
docs-serve: ## Serve documentation locally
	@echo -e "$(BLUE)Serving documentation locally...$(NC)"
	@if command -v python3 >/dev/null 2>&1; then \
		cd $(DOCS_DIR) && python3 -m http.server 8080; \
	else \
		echo -e "$(RED)Python3 not found$(NC)"; \
		exit 1; \
	fi

##@ Release

.PHONY: release-dry
release-dry: ## Dry run release process
	@echo -e "$(BLUE)Running release dry run...$(NC)"
	@if [ -f $(TOOLS_DIR)/goreleaser ]; then \
		$(TOOLS_DIR)/goreleaser release --snapshot --skip-publish --clean; \
	else \
		goreleaser release --snapshot --skip-publish --clean; \
	fi

.PHONY: release
release: ## Create a release (requires GPG and GitHub token)
	@echo -e "$(BLUE)Creating release...$(NC)"
	@if [ -z "$$GITHUB_TOKEN" ]; then \
		echo -e "$(RED)GITHUB_TOKEN environment variable is required$(NC)"; \
		exit 1; \
	fi
	@if [ -z "$$GPG_FINGERPRINT" ]; then \
		echo -e "$(RED)GPG_FINGERPRINT environment variable is required$(NC)"; \
		exit 1; \
	fi
	@if [ -f $(TOOLS_DIR)/goreleaser ]; then \
		$(TOOLS_DIR)/goreleaser release --clean; \
	else \
		goreleaser release --clean; \
	fi

.PHONY: release-notes
release-notes: ## Generate release notes
	@echo -e "$(BLUE)Generating release notes...$(NC)"
	@git log --pretty=format:"- %s" $(shell git describe --tags --abbrev=0)..HEAD

##@ Validation

.PHONY: validate
validate: fmt-check vet lint test-coverage-check security docs-check ## Run all validation checks

.PHONY: ci
ci: validate build test-integration ## Run CI pipeline locally

.PHONY: pre-commit
pre-commit: fmt lint test ## Run pre-commit checks

.PHONY: mod-tidy
mod-tidy: ## Clean up go.mod and go.sum
	@echo -e "$(BLUE)Tidying Go modules...$(NC)"
	@go mod tidy
	@go mod verify

.PHONY: mod-update
mod-update: ## Update all Go modules
	@echo -e "$(BLUE)Updating Go modules...$(NC)"
	@go get -u ./...
	@go mod tidy

##@ Examples

.PHONY: example-init
example-init: install ## Initialize example configurations
	@echo -e "$(BLUE)Initializing example configurations...$(NC)"
	@for example in $(EXAMPLES_DIR)/*/; do \
		if [ -f "$$example/main.tf" ]; then \
			echo "Initializing $$example"; \
			cd "$$example" && terraform init; \
		fi; \
	done

.PHONY: example-validate
example-validate: example-init ## Validate example configurations
	@echo -e "$(BLUE)Validating example configurations...$(NC)"
	@for example in $(EXAMPLES_DIR)/*/; do \
		if [ -f "$$example/main.tf" ]; then \
			echo "Validating $$example"; \
			cd "$$example" && terraform validate; \
		fi; \
	done

.PHONY: example-plan
example-plan: example-validate ## Plan example configurations
	@echo -e "$(BLUE)Planning example configurations...$(NC)"
	@for example in $(EXAMPLES_DIR)/*/; do \
		if [ -f "$$example/main.tf" ]; then \
			echo "Planning $$example"; \
			cd "$$example" && terraform plan; \
		fi; \
	done

##@ Debugging

.PHONY: debug-env
debug-env: ## Show debug environment information
	@echo -e "$(BLUE)Debug Environment Information:$(NC)"
	@echo "  GOOS: $(GOOS)"
	@echo "  GOARCH: $(GOARCH)"
	@echo "  GO_VERSION: $(GO_VERSION)"
	@echo "  CGO_ENABLED: $(CGO_ENABLED)"
	@echo "  BUILD_FLAGS: $(BUILD_FLAGS)"
	@echo "  LDFLAGS: $(LDFLAGS)"
	@echo "  PROJECT_NAME: $(PROJECT_NAME)"
	@echo "  VERSION: $(VERSION)"
	@echo "  COMMIT: $(COMMIT)"
	@echo "  BUILD_DATE: $(BUILD_DATE)"

.PHONY: debug-provider
debug-provider: build ## Debug provider binary
	@echo -e "$(BLUE)Provider Debug Information:$(NC)"
	@./$(BINARY_NAME) --help
	@echo ""
	@file ./$(BINARY_NAME)
	@echo ""
	@ls -la ./$(BINARY_NAME)

##@ Monitoring

.PHONY: watch-test
watch-test: ## Watch and run tests on file changes
	@echo -e "$(BLUE)Watching for changes and running tests...$(NC)"
	@echo -e "$(YELLOW)Press Ctrl+C to stop$(NC)"
	@while true; do \
		find . -name "*.go" | entr -d make test; \
	done

.PHONY: profile
profile: ## Run performance profiling
	@echo -e "$(BLUE)Running performance profiling...$(NC)"
	@go test -cpuprofile=cpu.prof -memprofile=mem.prof -bench=. ./...
	@echo -e "$(GREEN)Profiles generated: cpu.prof, mem.prof$(NC)"

##@ Docker

.PHONY: docker-build
docker-build: ## Build Docker image
	@echo -e "$(BLUE)Building Docker image...$(NC)"
	@docker build -t $(PROJECT_NAME):$(VERSION) .
	@docker build -t $(PROJECT_NAME):latest .

.PHONY: docker-run
docker-run: docker-build ## Run provider in Docker
	@echo -e "$(BLUE)Running provider in Docker...$(NC)"
	@docker run --rm -it $(PROJECT_NAME):$(VERSION)

##@ Utilities

.PHONY: open-coverage
open-coverage: test-coverage ## Open coverage report in browser
	@if [ -f $(COVERAGE_DIR)/coverage.html ]; then \
		$(OPEN) $(COVERAGE_DIR)/coverage.html; \
	else \
		echo -e "$(RED)Coverage report not found. Run 'make test-coverage' first.$(NC)"; \
	fi

.PHONY: check-deps
check-deps: ## Check for required dependencies
	@echo -e "$(BLUE)Checking dependencies...$(NC)"
	@command -v go >/dev/null 2>&1 || { echo -e "$(RED)Go is required but not installed$(NC)"; exit 1; }
	@command -v git >/dev/null 2>&1 || { echo -e "$(RED)Git is required but not installed$(NC)"; exit 1; }
	@command -v terraform >/dev/null 2>&1 || echo -e "$(YELLOW)Terraform not found (optional for development)$(NC)"
	@echo -e "$(GREEN)Dependencies check completed$(NC)"

# Include custom targets if they exist
-include Makefile.custom

# Ensure directories exist
$(DIST_DIR) $(COVERAGE_DIR) $(TOOLS_DIR):
	@mkdir -p $@

# Phony targets that don't create files
.PHONY: all
all: validate build docs ## Run all build and validation targets

# Default error handling
%:
	@echo -e "$(RED)Unknown target: $@$(NC)"
	@echo -e "$(BLUE)Run 'make help' to see available targets$(NC)"
	@exit 1