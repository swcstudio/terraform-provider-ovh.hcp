default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build provider
.PHONY: build
build:
	go build -o terraform-provider-hashicorp-ovh

# Install provider locally
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/spectrumwebco/hashicorp-ovh/0.1.0/linux_amd64
	mv terraform-provider-hashicorp-ovh ~/.terraform.d/plugins/registry.terraform.io/spectrumwebco/hashicorp-ovh/0.1.0/linux_amd64

# Generate documentation
.PHONY: docs
docs:
	tfplugindocs

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Format code
.PHONY: fmt
fmt:
	gofmt -s -w .

# Test
.PHONY: test
test:
	go test -v ./...

# Clean
.PHONY: clean
clean:
	rm -f terraform-provider-hashicorp-ovh

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build     - Build the provider binary"
	@echo "  install   - Install provider locally"
	@echo "  test      - Run unit tests"
	@echo "  testacc   - Run acceptance tests"
	@echo "  docs      - Generate documentation"
	@echo "  lint      - Run linter"
	@echo "  fmt       - Format code"
	@echo "  clean     - Clean build artifacts"
