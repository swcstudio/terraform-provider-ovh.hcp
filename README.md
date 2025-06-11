# Terraform Provider for HashiCorp Stack on OVHcloud

This Terraform provider enables you to manage HashiCorp infrastructure services on OVHcloud, providing enterprise-grade orchestration, security, and networking capabilities with seamless integration between HashiCorp's cloud-native services and OVHcloud's robust infrastructure platform.

## Features

- **Nomad Clusters**: Container and workload orchestration with multi-region support
- **Vault Clusters**: Secrets management with auto-unseal and HA configurations
- **Consul Clusters**: Service mesh and service discovery with Connect integration
- **Boundary Clusters**: Secure remote access and session management
- **Waypoint Runners**: Application deployment automation
- **Packer Templates**: Infrastructure image building and management
- **Cost Optimization**: Leverage OVH's competitive pricing and resource scheduling

## Quick Start

```hcl
terraform {
  required_providers {
    hashicorp-ovh = {
      source  = "swcstudio/hashicorp-ovh"
      version = "~> 0.1.0"
    }
  }
}

provider "hashicorp-ovh" {
  ovh_endpoint           = "ovh-eu"
  ovh_application_key    = var.ovh_application_key
  ovh_application_secret = var.ovh_application_secret
  ovh_consumer_key       = var.ovh_consumer_key
  ovh_project_id         = var.ovh_project_id
}

resource "hashicorp_ovh_nomad_cluster" "main" {
  name         = "production-nomad"
  region       = "eu-west-1"
  server_count = 3
  client_count = 5
  
  vault_integration  = true
  consul_integration = true
}
```

## Resources

- `hashicorp_ovh_nomad_cluster` - Nomad orchestration clusters
- `hashicorp_ovh_vault_cluster` - Vault secrets management
- `hashicorp_ovh_consul_cluster` - Consul service mesh
- `hashicorp_ovh_boundary_cluster` - Boundary access management
- `hashicorp_ovh_waypoint_runner` - Waypoint deployment automation
- `hashicorp_ovh_packer_template` - Packer image building

## Data Sources

- `hashicorp_ovh_nomad_clusters` - List available Nomad clusters
- `hashicorp_ovh_vault_clusters` - Query Vault cluster information
- `hashicorp_ovh_consul_clusters` - Consul cluster discovery

## Authentication

The provider requires OVH API credentials:

```bash
export OVH_ENDPOINT="ovh-eu"
export OVH_APPLICATION_KEY="your-app-key"
export OVH_APPLICATION_SECRET="your-app-secret"
export OVH_CONSUMER_KEY="your-consumer-key"
export OVH_PROJECT_ID="your-project-id"
```

## Examples

See the `examples/` directory for complete configuration examples including:

- Full HashiCorp stack deployment
- Multi-region cluster configurations
- Integrated security and networking
- Advanced orchestration patterns

## Deployment

### Prerequisites

Before deploying this provider to the Terraform Registry, ensure you have:

- ✅ **Go 1.18+** installed
- ✅ **GoReleaser** installed (`brew install goreleaser`)
- ✅ **GPG key** configured for signing releases
- ✅ **GitHub Personal Access Token** with repo permissions
- ✅ **Clean git state** (no uncommitted changes)

### Automated CI/CD Deployment

This provider uses GitHub Actions for fully automated CI/CD:

1. **Automated Testing:** Every pull request triggers comprehensive testing
2. **Security Scanning:** Automated security and vulnerability scanning
3. **Automated Releases:** Tagged commits automatically trigger releases
4. **Registry Publishing:** Releases are automatically published to Terraform Registry

### Creating a Release

1. **Create and push a tag:**
   ```bash
   git checkout main
   git pull origin main
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. **GitHub Actions will automatically:**
   - Run all tests and security scans
   - Build multi-platform binaries
   - Sign releases with GPG
   - Create GitHub release
   - Publish to Terraform Registry

### Manual Release (Emergency Only)

For emergency releases when CI/CD is unavailable:

```bash
make release
```

**Prerequisites:** Set `GITHUB_TOKEN` and `GPG_FINGERPRINT` environment variables.

## Development

### Quick Start Development

Set up your development environment:

```bash
# Clone the repository
git clone https://github.com/swcstudio/terraform-provider-hashicorp-ovh.git
cd terraform-provider-hashicorp-ovh

# Set up development environment
make setup

# Build the provider
make build

# Run tests
make test

# Run acceptance tests (requires API credentials)
make testacc

# Install locally for testing
make install
```

The `make setup` command will automatically:
- Install all required development tools
- Download and verify dependencies
- Configure the development environment

### Development Workflow

1. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your API credentials
   ```

2. **Make your changes:**
   ```bash
   # Edit code in internal/
   # Add tests in internal/provider/*_test.go
   ```

3. **Test your changes:**
   ```bash
   make fmt      # Format code
   make lint     # Run linter
   make test     # Run unit tests
   make build    # Build provider
   ```

4. **Test locally:**
   ```bash
   make install  # Install to local Terraform
   cd examples/local-dev
   terraform init
   terraform plan
   ```

### Available Make Targets

```bash
make help              # Show all available targets
make setup             # Set up development environment
make build             # Build the provider binary
make build-all         # Build for all platforms
make install           # Install provider locally
make test              # Run unit tests
make test-coverage     # Run tests with coverage
make testacc           # Run acceptance tests
make lint              # Run golangci-lint
make fmt               # Format code
make security          # Run security checks
make docs              # Generate documentation
make validate          # Run all validation checks
make ci                # Run full CI pipeline locally
make clean             # Clean build artifacts
```

### Enterprise-Grade Quality Checks

This provider includes comprehensive quality assurance:

```bash
make validate          # Run all checks (formatting, linting, tests, security)
make security          # Security scanning (gosec, vulnerability check)
make test-coverage     # Test coverage with threshold checking
make complexity        # Cyclomatic complexity analysis
make ci                # Full CI pipeline simulation
```

### Contributing

1. **Fork the repository**
2. **Create a feature branch:** `git checkout -b feature/your-feature-name`
3. **Make your changes** following the development workflow above
4. **Add tests** for new functionality
5. **Update documentation** if needed
6. **Commit your changes:** `git commit -am 'Add some feature'`
7. **Push to the branch:** `git push origin feature/your-feature-name`
8. **Submit a pull request**

### Project Structure

```
├── internal/
│   └── provider/          # Provider implementation
├── examples/              # Example configurations
│   ├── main.tf           # Basic example
│   └── local-dev/        # Local development setup
├── docs/                 # Auto-generated documentation
├── .github/
│   └── workflows/        # GitHub Actions CI/CD workflows
│       ├── ci.yml        # Main CI/CD pipeline
│       └── security.yml  # Security scanning workflow
├── charts/               # Helm charts for deployment
└── .goreleaser.yml       # Release configuration
```

### Testing

#### Unit Tests
```bash
make test
```

#### Acceptance Tests
Acceptance tests require valid OVH and HCP credentials:

```bash
export OVH_ENDPOINT="ovh-eu"
export OVH_APPLICATION_KEY="your-key"
export OVH_APPLICATION_SECRET="your-secret" 
export OVH_CONSUMER_KEY="your-consumer-key"
export HCP_CLIENT_ID="your-hcp-client-id"
export HCP_CLIENT_SECRET="your-hcp-client-secret"

make testacc
```

#### Integration Testing
```bash
# Run integration tests with examples
make test-integration

# Validate all example configurations
make example-validate

# Plan example configurations
make example-plan
```

#### Security Testing
```bash
# Run comprehensive security scans
make security

# Run vulnerability checks
make security-vuln

# Run static security analysis
make security-gosec
```

## Requirements

- Terraform >= 1.0
- Go >= 1.18
- OVH Cloud account with API access

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/swcstudio/terraform-provider-hashicorp-ovh/issues).
