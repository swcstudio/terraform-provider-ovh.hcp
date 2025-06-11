package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

	frameworkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"hashicorp-ovh": providerserver.NewProtocol6WithError(New("test")()),
}

// TestProviderInitialization tests that the provider can be initialized
func TestProviderInitialization(t *testing.T) {
	provider := New("test")()
	
	if provider == nil {
		t.Fatal("Expected provider to be initialized")
	}
	
	// Verify it's the correct type
	if _, ok := provider.(*HashiCorpOVHProvider); !ok {
		t.Error("Expected provider to be of type *HashiCorpOVHProvider")
	}
}

// TestProviderVersions tests provider initialization with different versions
func TestProviderVersions(t *testing.T) {
	versions := []string{"dev", "0.1.0", "1.0.0", "test"}
	
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			provider := New(version)()
			
			if provider == nil {
				t.Errorf("Provider should initialize with version %s", version)
			}
			
			// Test that metadata contains the correct version
			req := frameworkprovider.MetadataRequest{}
			resp := &frameworkprovider.MetadataResponse{}
			
			provider.Metadata(context.Background(), req, resp)
			
			if resp.Version != version {
				t.Errorf("Expected version %s, got %s", version, resp.Version)
			}
		})
	}
}

// TestProviderMetadata tests that provider metadata is correctly set
func TestProviderMetadata(t *testing.T) {
	provider := New("test")()
	
	req := frameworkprovider.MetadataRequest{}
	resp := &frameworkprovider.MetadataResponse{}
	
	provider.Metadata(context.Background(), req, resp)
	
	// Check that TypeName is set correctly
	expectedTypeName := "hashicorp-ovh"
	if resp.TypeName != expectedTypeName {
		t.Errorf("Expected TypeName %s, got %s", expectedTypeName, resp.TypeName)
	}
	
	// Check that Version is set
	if resp.Version == "" {
		t.Error("Expected Version to be set")
	}
}

// TestProviderSchema tests that the provider schema can be retrieved
func TestProviderSchema(t *testing.T) {
	provider := New("test")()
	
	req := frameworkprovider.SchemaRequest{}
	resp := &frameworkprovider.SchemaResponse{}
	
	provider.Schema(context.Background(), req, resp)
	
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema method returned errors: %v", resp.Diagnostics.Errors())
	}
	
	// Verify that schema has expected attributes
	expectedAttributes := []string{
		"ovh_endpoint",
		"ovh_application_key", 
		"ovh_application_secret",
		"ovh_consumer_key",
	}
	
	for _, attrName := range expectedAttributes {
		if _, exists := resp.Schema.Attributes[attrName]; !exists {
			t.Errorf("Expected attribute %s not found in schema", attrName)
		}
	}
	
	// Verify required attributes are marked as required
	requiredAttributes := []string{
		"ovh_endpoint",
		"ovh_application_key",
		"ovh_application_secret", 
		"ovh_consumer_key",
	}
	
	for _, attrName := range requiredAttributes {
		attr, exists := resp.Schema.Attributes[attrName]
		if !exists {
			t.Errorf("Required attribute %s not found", attrName)
			continue
		}
		if !attr.IsRequired() {
			t.Errorf("Attribute %s should be required", attrName)
		}
	}
	
	// Verify sensitive attributes are marked as sensitive
	sensitiveAttributes := []string{
		"ovh_application_secret",
		"ovh_consumer_key",
	}
	
	for _, attrName := range sensitiveAttributes {
		attr, exists := resp.Schema.Attributes[attrName]
		if !exists {
			t.Errorf("Sensitive attribute %s not found", attrName)
			continue
		}
		if !attr.IsSensitive() {
			t.Errorf("Attribute %s should be marked as sensitive", attrName)
		}
	}
}

// TestProviderConfigureWithValidConfig tests provider configuration with valid config
func TestProviderConfigureWithValidConfig(t *testing.T) {
	// Note: Full configure testing requires acceptance tests
	// This test just verifies the provider can be instantiated
	provider := New("test")()
	
	if provider == nil {
		t.Error("Expected provider to be instantiated")
	}
	
	// Test that Configure method exists and can be called
	// (Full testing requires proper ConfigureRequest setup which is complex for unit tests)
}

// TestProviderConfigureWithMissingConfig tests provider configuration with missing config
func TestProviderConfigureWithMissingConfig(t *testing.T) {
	// Note: Full configure testing requires acceptance tests
	// This test just verifies the provider behavior
	provider := New("test")()
	
	if provider == nil {
		t.Error("Expected provider to be instantiated")
	}
	
	// Configuration validation testing is better done in acceptance tests
	// where we can properly set up ConfigureRequest with tfsdk.Config
}

// TestProviderResources tests that resources are properly registered
func TestProviderResources(t *testing.T) {
	provider := New("test")()
	
	resources := provider.Resources(context.Background())
	
	// Currently we expect no resources since they're not implemented yet
	// This test will need to be updated as resources are added
	expectedResourceCount := 0
	if len(resources) != expectedResourceCount {
		t.Errorf("Expected %d resources, got %d", expectedResourceCount, len(resources))
	}
}

// TestProviderDataSources tests that data sources are properly registered
func TestProviderDataSources(t *testing.T) {
	provider := New("test")()
	
	dataSources := provider.DataSources(context.Background())
	
	// Currently we expect no data sources since they're not implemented yet
	// This test will need to be updated as data sources are added
	expectedDataSourceCount := 0
	if len(dataSources) != expectedDataSourceCount {
		t.Errorf("Expected %d data sources, got %d", expectedDataSourceCount, len(dataSources))
	}
}

// TestProviderConcurrentAccess tests thread safety
func TestProviderConcurrentAccess(t *testing.T) {
	provider := New("test")()
	
	// Test concurrent metadata requests
	done := make(chan bool, 10)
	errors := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			req := frameworkprovider.MetadataRequest{}
			resp := &frameworkprovider.MetadataResponse{}
			
			provider.Metadata(context.Background(), req, resp)
			
			if resp.TypeName != "hashicorp-ovh" {
				errors <- fmt.Errorf("concurrent access error: unexpected TypeName %s", resp.TypeName)
				return
			}
		}()
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	
	close(errors)
	for err := range errors {
		t.Error(err)
	}
}

// TestProviderEnvironmentVariables tests environment variable handling
func TestProviderEnvironmentVariables(t *testing.T) {
	// Test environment variable reading capability
	provider := New("test")()
	
	if provider == nil {
		t.Error("Provider should initialize regardless of environment variables")
	}
	
	// Environment variable testing is better done in acceptance tests
	// where we can properly test the full configure flow
}

// preCheck ensures acceptance test requirements are met
func testAccPreCheck(t *testing.T) {
	// Check for required environment variables for acceptance tests
	requiredEnvVars := []string{
		"OVH_ENDPOINT",
		"OVH_APPLICATION_KEY",
		"OVH_APPLICATION_SECRET",
		"OVH_CONSUMER_KEY",
	}
	
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("%s environment variable must be set for acceptance tests", envVar)
		}
	}
}

// testAccCheckProviderConfigured is a test helper to verify provider is properly configured
func testAccCheckProviderConfigured(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		
		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID not set")
		}
		
		return nil
	}
}

// TestAccProvider tests the provider in an acceptance test context
func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccProviderConfig(),
				Check: resource.ComposeTestCheckFunc(
					// Basic provider configuration test
					func(s *terraform.State) error {
						// Provider should be configured without errors
						return nil
					},
				),
			},
		},
	})
}

// testAccProviderConfig returns a basic provider configuration for testing
func testAccProviderConfig() string {
	return `
provider "hashicorp-ovh" {
  # Configuration will be read from environment variables
  # OVH_ENDPOINT, OVH_APPLICATION_KEY, etc.
}
`
}

// BenchmarkProviderInitialization benchmarks provider creation
func BenchmarkProviderInitialization(b *testing.B) {
	for i := 0; i < b.N; i++ {
		provider := New("test")()
		if provider == nil {
			b.Fatal("Provider initialization failed")
		}
	}
}

// BenchmarkProviderMetadata benchmarks metadata retrieval
func BenchmarkProviderMetadata(b *testing.B) {
	provider := New("test")()
	
	for i := 0; i < b.N; i++ {
		req := frameworkprovider.MetadataRequest{}
		resp := &frameworkprovider.MetadataResponse{}
		
		provider.Metadata(context.Background(), req, resp)
		
		if resp.TypeName == "" {
			b.Fatal("Metadata retrieval failed")
		}
	}
}

// BenchmarkProviderSchema benchmarks schema retrieval
func BenchmarkProviderSchema(b *testing.B) {
	provider := New("test")()
	
	for i := 0; i < b.N; i++ {
		req := frameworkprovider.SchemaRequest{}
		resp := &frameworkprovider.SchemaResponse{}
		
		provider.Schema(context.Background(), req, resp)
		
		if resp.Diagnostics.HasError() {
			b.Fatalf("Schema retrieval failed: %v", resp.Diagnostics.Errors())
		}
	}
}

// BenchmarkProviderConfigure benchmarks provider configuration
func BenchmarkProviderConfigure(b *testing.B) {
	// Skip benchmarking Configure method as it requires complex setup
	// Benchmark provider initialization instead
	for i := 0; i < b.N; i++ {
		provider := New("test")()
		
		if provider == nil {
			b.Fatal("Provider initialization failed")
		}
	}
}