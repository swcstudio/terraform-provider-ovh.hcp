package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccNomadCluster_basic tests basic Nomad cluster creation
func TestAccNomadCluster_basic(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-cluster"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "region", "eu-west-1"),
					resource.TestCheckResourceAttr(resourceName, "server_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "vault_integration", "true"),
					resource.TestCheckResourceAttr(resourceName, "consul_integration", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "endpoint"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

// TestAccNomadCluster_update tests Nomad cluster updates
func TestAccNomadCluster_update(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-cluster-update"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "server_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_count", "5"),
				),
			},
			{
				Config: testAccNomadClusterConfig_updated(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "server_count", "5"),
					resource.TestCheckResourceAttr(resourceName, "client_count", "8"),
					resource.TestCheckResourceAttr(resourceName, "vault_integration", "false"),
				),
			},
		},
	})
}

// TestAccNomadCluster_withTags tests Nomad cluster with tags
func TestAccNomadCluster_withTags(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-cluster-tags"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_withTags(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "test"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "platform"),
					resource.TestCheckResourceAttr(resourceName, "tags.ManagedBy", "terraform"),
				),
			},
			{
				Config: testAccNomadClusterConfig_withUpdatedTags(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.Environment", "staging"),
					resource.TestCheckResourceAttr(resourceName, "tags.Team", "devops"),
					resource.TestCheckResourceAttr(resourceName, "tags.Owner", "john.doe"),
					resource.TestCheckNoResourceAttr(resourceName, "tags.ManagedBy"),
				),
			},
		},
	})
}

// TestAccNomadCluster_minimalConfig tests minimal configuration
func TestAccNomadCluster_minimalConfig(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-minimal"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_minimal(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", clusterName),
					resource.TestCheckResourceAttr(resourceName, "region", "eu-west-1"),
					// Check default values
					resource.TestCheckResourceAttr(resourceName, "server_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "client_count", "3"),
					resource.TestCheckResourceAttr(resourceName, "vault_integration", "false"),
					resource.TestCheckResourceAttr(resourceName, "consul_integration", "false"),
				),
			},
		},
	})
}

// TestAccNomadCluster_invalidConfiguration tests error handling
func TestAccNomadCluster_invalidConfiguration(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccNomadClusterConfig_invalidServerCount(),
				ExpectError: regexp.MustCompile("server_count must be between 1 and 10"),
			},
			{
				Config:      testAccNomadClusterConfig_invalidClientCount(),
				ExpectError: regexp.MustCompile("client_count must be between 0 and 100"),
			},
			{
				Config:      testAccNomadClusterConfig_invalidRegion(),
				ExpectError: regexp.MustCompile("invalid region"),
			},
			{
				Config:      testAccNomadClusterConfig_invalidName(),
				ExpectError: regexp.MustCompile("name must be between 3 and 50 characters"),
			},
		},
	})
}

// TestAccNomadCluster_disappears tests resource disappears scenario
func TestAccNomadCluster_disappears(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-disappears"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					testAccCheckNomadClusterDisappears(resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

// TestAccNomadCluster_import tests resource import functionality
func TestAccNomadCluster_import(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-import"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_basic(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"vault_integration", // May not be returned in read
					"consul_integration",
				},
			},
		},
	})
}

// TestAccNomadCluster_withSecurityGroups tests security group configuration
func TestAccNomadCluster_withSecurityGroups(t *testing.T) {
	resourceName := "hashicorp_ovh_nomad_cluster.test"
	clusterName := "test-nomad-sg"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckNomadClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNomadClusterConfig_withSecurityGroups(clusterName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNomadClusterExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "security_groups.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "security_groups.*", "sg-nomad-servers"),
					resource.TestCheckTypeSetElemAttr(resourceName, "security_groups.*", "sg-nomad-clients"),
				),
			},
		},
	})
}

// Helper functions for test checks

func testAccCheckNomadClusterExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Nomad cluster ID is set")
		}

		// Here you would typically make an API call to verify the resource exists
		// For now, we'll assume the resource exists if it has an ID
		return nil
	}
}

func testAccCheckNomadClusterDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "hashicorp_ovh_nomad_cluster" {
			continue
		}

		// Here you would typically make an API call to verify the resource is destroyed
		// For now, we'll check that the ID is cleared
		if rs.Primary.ID != "" {
			return fmt.Errorf("Nomad cluster %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckNomadClusterDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Nomad cluster ID is set")
		}

		// Here you would typically make an API call to delete the resource
		// simulating it disappearing outside of Terraform
		return nil
	}
}

// Test configuration functions

func testAccNomadClusterConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "%s"
  region       = "eu-west-1"
  server_count = 3
  client_count = 5
  
  vault_integration  = true
  consul_integration = true
}
`, name)
}

func testAccNomadClusterConfig_updated(name string) string {
	return fmt.Sprintf(`
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "%s"
  region       = "eu-west-1"
  server_count = 5
  client_count = 8
  
  vault_integration  = false
  consul_integration = true
}
`, name)
}

func testAccNomadClusterConfig_minimal(name string) string {
	return fmt.Sprintf(`
resource "hashicorp_ovh_nomad_cluster" "test" {
  name   = "%s"
  region = "eu-west-1"
}
`, name)
}

func testAccNomadClusterConfig_withTags(name string) string {
	return fmt.Sprintf(`
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "%s"
  region       = "eu-west-1"
  server_count = 3
  client_count = 5

  tags = {
    Environment = "test"
    Team        = "platform"
    ManagedBy   = "terraform"
  }
}
`, name)
}

func testAccNomadClusterConfig_withUpdatedTags(name string) string {
	return fmt.Sprintf(`
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "%s"
  region       = "eu-west-1"
  server_count = 3
  client_count = 5

  tags = {
    Environment = "staging"
    Team        = "devops"
    Owner       = "john.doe"
  }
}
`, name)
}

func testAccNomadClusterConfig_withSecurityGroups(name string) string {
	return fmt.Sprintf(`
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "%s"
  region       = "eu-west-1"
  server_count = 3
  client_count = 5

  security_groups = [
    "sg-nomad-servers",
    "sg-nomad-clients"
  ]
}
`, name)
}

func testAccNomadClusterConfig_invalidServerCount() string {
	return `
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "test-invalid"
  region       = "eu-west-1"
  server_count = 15  # Invalid: exceeds maximum
  client_count = 5
}
`
}

func testAccNomadClusterConfig_invalidClientCount() string {
	return `
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "test-invalid"
  region       = "eu-west-1"
  server_count = 3
  client_count = 150  # Invalid: exceeds maximum
}
`
}

func testAccNomadClusterConfig_invalidRegion() string {
	return `
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "test-invalid"
  region       = "invalid-region"
  server_count = 3
  client_count = 5
}
`
}

func testAccNomadClusterConfig_invalidName() string {
	return `
resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "ab"  # Invalid: too short
  region       = "eu-west-1"
  server_count = 3
  client_count = 5
}
`
}

// Unit tests for resource logic will be added when resources are implemented

// TODO: Add resource schema tests when nomadClusterResource is implemented
// TODO: Add validation tests when resource validation is implemented  
// TODO: Add benchmark tests when resource operations are implemented
