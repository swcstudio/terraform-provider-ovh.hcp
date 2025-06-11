terraform {
  required_providers {
    hashicorp-ovh = {
      source  = "swcstudio/hashicorp-ovh"
      version = "~> 0.1.0"
    }
  }
}

provider "hashicorp-ovh" {
  ovh_endpoint           = var.ovh_endpoint
  ovh_application_key    = var.ovh_application_key
  ovh_application_secret = var.ovh_application_secret
  ovh_consumer_key       = var.ovh_consumer_key
  ovh_project_id         = var.ovh_project_id
}

resource "hashicorp_ovh_nomad_cluster" "example" {
  name         = "example-nomad"
  region       = "eu-west-1"
  version      = "1.6.2"
  server_count = 3
  client_count = 3
  
  instance_type = "s1-4"
  datacenter    = "dc1"
  
  encrypt             = true
  acl                 = true
  tls                 = true
  vault_integration   = true
  consul_integration  = true
  monitoring          = true
  backup              = true
  
  tags = {
    Environment = "production"
    Team        = "platform"
  }
}

resource "hashicorp_ovh_vault_cluster" "example" {
  name       = "example-vault"
  region     = "eu-west-1"
  version    = "1.14.2"
  node_count = 3
  
  instance_type           = "s1-4"
  storage_type           = "integrated"
  auto_unseal            = true
  audit_logging          = true
  performance_replication = false
  disaster_recovery      = false
  monitoring             = true
  backup                 = true
  
  tags = {
    Environment = "production"
    Team        = "security"
  }
}

resource "hashicorp_ovh_consul_cluster" "example" {
  name         = "example-consul"
  region       = "eu-west-1"
  version      = "1.16.1"
  server_count = 3
  client_count = 3
  
  instance_type = "s1-4"
  datacenter    = "dc1"
  
  encrypt      = true
  acl          = true
  tls          = true
  connect      = true
  mesh_gateway = true
  monitoring   = true
  backup       = true
  
  tags = {
    Environment = "production"
    Team        = "platform"
  }
}

resource "hashicorp_ovh_boundary_cluster" "example" {
  name             = "example-boundary"
  region           = "eu-west-1"
  version          = "0.13.2"
  controller_count = 3
  worker_count     = 3
  
  instance_type = "s1-4"
  database_type = "postgres"
  kms_type      = "aead"
  tls           = true
  monitoring    = true
  backup        = true
  
  tags = {
    Environment = "production"
    Team        = "security"
  }
}

output "nomad_endpoint" {
  value = hashicorp_ovh_nomad_cluster.example.endpoint
}

output "vault_endpoint" {
  value = hashicorp_ovh_vault_cluster.example.endpoint
}

output "consul_endpoint" {
  value = hashicorp_ovh_consul_cluster.example.endpoint
}

output "boundary_endpoint" {
  value = hashicorp_ovh_boundary_cluster.example.controller_endpoint
}
