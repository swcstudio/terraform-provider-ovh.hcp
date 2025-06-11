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

### Quick Deploy

For maintainers, use the automated release script:

```bash
# Set your GitHub token
export GITHUB_TOKEN=your_github_token_here

# Create and publish a release
./scripts/release.sh v0.1.0
```

### Manual Deployment Steps

1. **Prepare the release:**
   ```bash
   git checkout main
   git pull origin main
   make test
   ```

2. **Create and push a tag:**
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

3. **Build and publish:**
   ```bash
   export GPG_TTY=$(tty)
   export GPG_FINGERPRINT=your_gpg_fingerprint
   export GITHUB_TOKEN=your_github_token
   goreleaser release --clean
   ```

4. **Register with Terraform Registry:**
   - Go to [registry.terraform.io](https://registry.terraform.io)
   - Sign in with GitHub
   - Click "Publish" → "Provider"
   - Select `swcstudio/terraform-provider-hashicorp-ovh`
   - Add your GPG public key

## Development

### Quick Start Development

Run the automated development setup:

```bash
./scripts/dev-setup.sh
```

This will install all dependencies and configure your development environment.

### Manual Development Setup

```bash
# Clone the repository
git clone https://github.com/swcstudio/terraform-provider-hashicorp-ovh.git
cd terraform-provider-hashicorp-ovh

# Install dependencies
go mod download
go mod tidy

# Build the provider
make build

# Run tests
make test

# Run acceptance tests (requires API credentials)
make testacc

# Install locally for testing
make install
```

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
make help     # Show all available targets
make build    # Build the provider binary
make install  # Install provider locally
make test     # Run unit tests
make testacc  # Run acceptance tests
make docs     # Generate documentation
make lint     # Run linter
make fmt      # Format code
make clean    # Clean build artifacts
```

### Development Shortcuts

Use the development helper script for common tasks:

```bash
./dev.sh build    # Build provider
./dev.sh test     # Run tests
./dev.sh install  # Install locally
./dev.sh docs     # Generate docs
./dev.sh release v0.1.0  # Create release
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
├── scripts/              # Development and release scripts
│   ├── dev-setup.sh      # Development environment setup
│   └── release.sh        # Automated release script
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

#### Local Testing
```bash
# Install provider locally
make install

# Test with example configuration
cd examples/local-dev
terraform init
terraform plan
terraform apply
```

## Requirements

- Terraform >= 1.0
- Go >= 1.18
- OVH Cloud account with API access

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/swcstudio/terraform-provider-hashicorp-ovh/issues).
