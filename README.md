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
      source  = "spectrumwebco/hashicorp-ovh"
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

## Development

```bash
# Build the provider
make build

# Run tests
make test

# Run acceptance tests
make testacc

# Install locally
make install
```

## Requirements

- Terraform >= 1.0
- Go >= 1.18
- OVH Cloud account with API access

## License

This project is licensed under the Mozilla Public License 2.0 - see the [LICENSE](LICENSE) file for details.

## Support

For issues and feature requests, please use the [GitHub issue tracker](https://github.com/spectrumwebco/terraform-provider-hashicorp-ovh/issues).
