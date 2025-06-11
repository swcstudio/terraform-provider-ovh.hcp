terraform {
  required_version = ">= 1.0"
  
  required_providers {
    hashicorp-ovh = {
      source  = "swcstudio/hashicorp-ovh"
      version = "~> 0.1.0"
    }
  }

  # Use local backend for development
  backend "local" {
    path = "terraform.tfstate"
  }
}

# Provider configuration
# Environment variables are preferred for sensitive values:
# export OVH_ENDPOINT="ovh-eu"
# export OVH_APPLICATION_KEY="your-app-key"
# export OVH_APPLICATION_SECRET="your-app-secret"
# export OVH_CONSUMER_KEY="your-consumer-key"
# export OVH_PROJECT_ID="your-project-id"
provider "hashicorp-ovh" {
  # Configuration will be loaded from environment variables
  # Alternatively, you can specify them here (not recommended for production):
  # ovh_endpoint           = "ovh-eu"
  # ovh_application_key    = var.ovh_application_key
  # ovh_application_secret = var.ovh_application_secret
  # ovh_consumer_key       = var.ovh_consumer_key
  # ovh_project_id         = var.ovh_project_id
}

# Variables for development (optional)
variable "cluster_name" {
  description = "Name prefix for HashiCorp clusters"
  type        = string
  default     = "dev"
}

variable "region" {
  description = "OVH region for deployments"
  type        = string
  default     = "eu-west-1"
}

variable "environment" {
  description = "Environment tag"
  type        = string
  default     = "development"
}

# Example: Nomad Cluster
# Uncomment and modify as needed once the provider is fully implemented
/*
resource "hashicorp_ovh_nomad_cluster" "example" {
  name         = "${var.cluster_name}-nomad"
  region       = var.region
  server_count = 1
  client_count = 2
  
  server_type = "small"
  client_type = "small"
  
  vault_integration  = true
  consul_integration = true
  
  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
    Provider    = "hashicorp-ovh"
  }
}
*/

# Example: Vault Cluster
# Uncomment and modify as needed
/*
resource "hashicorp_ovh_vault_cluster" "example" {
  name         = "${var.cluster_name}-vault"
  region       = var.region
  node_count   = 3
  
  auto_unseal_enabled = true
  ha_enabled         = true
  
  storage_backend = "consul"
  consul_cluster_id = hashicorp_ovh_consul_cluster.example.id
  
  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
    Provider    = "hashicorp-ovh"
  }
}
*/

# Example: Consul Cluster
# Uncomment and modify as needed
/*
resource "hashicorp_ovh_consul_cluster" "example" {
  name         = "${var.cluster_name}-consul"
  region       = var.region
  node_count   = 3
  
  connect_enabled = true
  encryption_enabled = true
  
  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
    Provider    = "hashicorp-ovh"
  }
}
*/

# Example: Boundary Cluster
# Uncomment and modify as needed
/*
resource "hashicorp_ovh_boundary_cluster" "example" {
  name           = "${var.cluster_name}-boundary"
  region         = var.region
  controller_count = 1
  worker_count     = 2
  
  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
    Provider    = "hashicorp-ovh"
  }
}
*/

# Example: Waypoint Runner
# Uncomment and modify as needed
/*
resource "hashicorp_ovh_waypoint_runner" "example" {
  name   = "${var.cluster_name}-waypoint-runner"
  region = var.region
  
  runner_type = "static"
  
  nomad_cluster_id = hashicorp_ovh_nomad_cluster.example.id
  
  tags = {
    Environment = var.environment
    ManagedBy   = "terraform"
    Provider    = "hashicorp-ovh"
  }
}
*/

# Data sources for testing
# These can be used to verify the provider is working
/*
data "hashicorp_ovh_nomad_clusters" "all" {}

data "hashicorp_ovh_vault_clusters" "all" {}

data "hashicorp_ovh_consul_clusters" "all" {}
*/

# Outputs for development
# Uncomment as resources are implemented
/*
output "nomad_cluster" {
  description = "Nomad cluster information"
  value = {
    id       = hashicorp_ovh_nomad_cluster.example.id
    name     = hashicorp_ovh_nomad_cluster.example.name
    endpoint = hashicorp_ovh_nomad_cluster.example.endpoint
  }
}

output "vault_cluster" {
  description = "Vault cluster information"
  value = {
    id       = hashicorp_ovh_vault_cluster.example.id
    name     = hashicorp_ovh_vault_cluster.example.name
    endpoint = hashicorp_ovh_vault_cluster.example.endpoint
  }
}

output "consul_cluster" {
  description = "Consul cluster information"
  value = {
    id       = hashicorp_ovh_consul_cluster.example.id
    name     = hashicorp_ovh_consul_cluster.example.name
    endpoint = hashicorp_ovh_consul_cluster.example.endpoint
  }
}
*/

# Simple output to verify provider is loaded
output "provider_info" {
  description = "Information about the HashiCorp-OVH provider"
  value = {
    provider_version = "0.1.0"
    terraform_version = ">= 1.0"
    environment = var.environment
    region = var.region
  }
}