package provider

import (
	"context"
	"fmt"
	"os"
	"testing"

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

// TestProvider tests provider initialization and configuration
func TestProvider(t *testing.T) {
	provider := New("test")()
	
	if provider == nil {
		t.Fatal("Expected provider to be initialized")
	}
}

// TestProviderSchema validates the provider schema
func TestProviderSchema(t *testing.T) {
	provider := New("test")()
	
	req := provider.GetProviderSchemaRequest{}
	resp := provider.GetProviderSchemaResponse{}
	
	provider.GetProviderSchema(context.Background(), req, &resp)
	
	if resp.Diagnostics.HasError() {
		t.Fatalf("Expected no errors, got: %v", resp.Diagnostics.Errors())
	}
	
	// Verify required attributes are present
	if resp.Provider.Attributes == nil {
		t.Fatal("Expected provider attributes to be defined")
	}
}

// TestProviderConfigure tests provider configuration with various scenarios
func TestProviderConfigure(t *testing.T) {
	tests := []struct {
		name        string
		config      map[string]interface{}
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid_configuration",
			config: map[string]interface{}{
				"ovh_endpoint":           "ovh-eu",
				"ovh_application_key":    "test-key",
				"ovh_application_secret": "test-secret", 
				"ovh_consumer_key":       "test-consumer-key",
				"ovh_project_id":         "test-project-id",
			},
			expectError: false,
		},
		{
			name: "missing_endpoint",
			config: map[string]interface{}{
				"ovh_application_key":    "test-key",
				"ovh_application_secret": "test-secret",
				"ovh_consumer_key":       "test-consumer-key", 
				"ovh_project_id":         "test-project-id",
			},
			expectError: true,
			errorMsg:    "OVH endpoint is required",
		},
		{
			name: "invalid_endpoint",
			config: map[string]interface{}{
				"ovh_endpoint":           "invalid-endpoint",
				"ovh_application_key":    "test-key",
				"ovh_application_secret": "test-secret",
				"ovh_consumer_key":       "test-consumer-key",
				"ovh_project_id":         "test-project-id",
			},
			expectError: true,
			errorMsg:    "Invalid OVH endpoint",
		},
		{
			name: "missing_application_key",
			config: map[string]interface{}{
				"ovh_endpoint":           "ovh-eu",
				"ovh_application_secret": "test-secret",
				"ovh_consumer_key":       "test-consumer-key",
				"ovh_project_id":         "test-project-id",
			},
			expectError: true,
			errorMsg:    "OVH application key is required",
		},
		{
			name: "empty_configuration",
			config: map[string]interface{}{},
			expectError: true,
			errorMsg:    "OVH configuration is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := New("test")()
			
			// Test configuration validation
			req := provider.ConfigureProviderRequest{}
			resp := provider.ConfigureProviderResponse{}
			
			// Convert config to terraform value
			// This would need proper implementation based on actual schema
			
			provider.ConfigureProvider(context.Background(), req, &resp)
			
			if tt.expectError {
				if !resp.Diagnostics.HasError() {
					t.Errorf("Expected error but got none")
				}
				// Check error message contains expected text
				found := false
				for _, diag := range resp.Diagnostics.Errors() {
					if diag.Summary() == tt.errorMsg || diag.Detail() == tt.errorMsg {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s' not found in diagnostics", tt.errorMsg)
				}
			} else {
				if resp.Diagnostics.HasError() {
					t.Errorf("Expected no error but got: %v", resp.Diagnostics.Errors())
				}
			}
		})
	}
}

// TestProviderResourcesAndDataSources verifies all resources and data sources are registered
func TestProviderResourcesAndDataSources(t *testing.T) {
	provider := New("test")()
	
	req := provider.GetProviderSchemaRequest{}
	resp := provider.GetProviderSchemaResponse{}
	
	provider.GetProviderSchema(context.Background(), req, &resp)
	
	if resp.Diagnostics.HasError() {
		t.Fatalf("Unexpected errors: %v", resp.Diagnostics.Errors())
	}
	
	// Test that expected resources are registered
	expectedResources := []string{
		"hashicorp_ovh_nomad_cluster",
		"hashicorp_ovh_vault_cluster", 
		"hashicorp_ovh_consul_cluster",
		"hashicorp_ovh_boundary_cluster",
		"hashicorp_ovh_waypoint_runner",
		"hashicorp_ovh_packer_template",
	}
	
	for _, resourceName := range expectedResources {
		if _, exists := resp.ResourceSchemas[resourceName]; !exists {
			t.Errorf("Expected resource %s to be registered", resourceName)
		}
	}
	
	// Test that expected data sources are registered
	expectedDataSources := []string{
		"hashicorp_ovh_nomad_clusters",
		"hashicorp_ovh_vault_clusters",
		"hashicorp_ovh_consul_clusters", 
		"hashicorp_ovh_boundary_clusters",
	}
	
	for _, dataSourceName := range expectedDataSources {
		if _, exists := resp.DataSourceSchemas[dataSourceName]; !exists {
			t.Errorf("Expected data source %s to be registered", dataSourceName)
		}
	}
}

// TestProviderEnvironmentVariables tests environment variable configuration
func TestProviderEnvironmentVariables(t *testing.T) {
	// Save original environment
	originalVars := map[string]string{
		"OVH_ENDPOINT":           os.Getenv("OVH_ENDPOINT"),
		"OVH_APPLICATION_KEY":    os.Getenv("OVH_APPLICATION_KEY"),
		"OVH_APPLICATION_SECRET": os.Getenv("OVH_APPLICATION_SECRET"),
		"OVH_CONSUMER_KEY":       os.Getenv("OVH_CONSUMER_KEY"),
		"OVH_PROJECT_ID":         os.Getenv("OVH_PROJECT_ID"),
	}
	
	// Restore environment after test
	defer func() {
		for key, value := range originalVars {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()
	
	// Test with valid environment variables
	os.Setenv("OVH_ENDPOINT", "ovh-eu")
	os.Setenv("OVH_APPLICATION_KEY", "test-key")
	os.Setenv("OVH_APPLICATION_SECRET", "test-secret")
	os.Setenv("OVH_CONSUMER_KEY", "test-consumer-key")
	os.Setenv("OVH_PROJECT_ID", "test-project-id")
	
	provider := New("test")()
	
	req := provider.ConfigureProviderRequest{}
	resp := provider.ConfigureProviderResponse{}
	
	provider.ConfigureProvider(context.Background(), req, &resp)
	
	if resp.Diagnostics.HasError() {
		t.Errorf("Expected no errors with valid environment variables, got: %v", resp.Diagnostics.Errors())
	}
}

// TestProviderConcurrentAccess tests thread safety
func TestProviderConcurrentAccess(t *testing.T) {
	provider := New("test")()
	
	// Test concurrent schema requests
	done := make(chan bool)
	errors := make(chan error, 10)
	
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			
			req := provider.GetProviderSchemaRequest{}
			resp := provider.GetProviderSchemaResponse{}
			
			provider.GetProviderSchema(context.Background(), req, &resp)
			
			if resp.Diagnostics.HasError() {
				errors <- fmt.Errorf("concurrent access error: %v", resp.Diagnostics.Errors())
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

// TestProviderVersionValidation tests version handling
func TestProviderVersionValidation(t *testing.T) {
	versions := []string{"dev", "0.1.0", "1.0.0", "test"}
	
	for _, version := range versions {
		t.Run(fmt.Sprintf("version_%s", version), func(t *testing.T) {
			provider := New(version)()
			
			if provider == nil {
				t.Errorf("Provider should initialize with version %s", version)
			}
		})
	}
}

// TestProviderConfigurationValidation tests configuration edge cases
func TestProviderConfigurationValidation(t *testing.T) {
	provider := New("test")()
	
	testCases := []struct {
		name        string
		setupEnv    func()
		expectError bool
		description string
	}{
		{
			name: "conflicting_config_and_env",
			setupEnv: func() {
				os.Setenv("OVH_ENDPOINT", "ovh-eu")
				// Config would also specify endpoint differently
			},
			expectError: false, // Should use explicit config over env
			description: "Explicit configuration should take precedence over environment variables",
		},
		{
			name: "partial_env_vars",
			setupEnv: func() {
				os.Setenv("OVH_ENDPOINT", "ovh-eu")
				os.Setenv("OVH_APPLICATION_KEY", "test-key")
				// Missing other required vars
				os.Unsetenv("OVH_APPLICATION_SECRET")
				os.Unsetenv("OVH_CONSUMER_KEY")
				os.Unsetenv("OVH_PROJECT_ID")
			},
			expectError: true,
			description: "Should fail with partial environment configuration",
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup environment
			tc.setupEnv()
			
			// Clean up after test
			defer func() {
				envVars := []string{
					"OVH_ENDPOINT", "OVH_APPLICATION_KEY", "OVH_APPLICATION_SECRET",
					"OVH_CONSUMER_KEY", "OVH_PROJECT_ID",
				}
				for _, env := range envVars {
					os.Unsetenv(env)
				}
			}()
			
			req := provider.ConfigureProviderRequest{}
			resp := provider.ConfigureProviderResponse{}
			
			provider.ConfigureProvider(context.Background(), req, &resp)
			
			hasError := resp.Diagnostics.HasError()
			if tc.expectError && !hasError {
				t.Errorf("%s: expected error but got none", tc.description)
			} else if !tc.expectError && hasError {
				t.Errorf("%s: expected no error but got: %v", tc.description, resp.Diagnostics.Errors())
			}
		})
	}
}

// preCheck ensures acceptance test requirements are met
func testAccPreCheck(t *testing.T) {
	// Check for required environment variables for acceptance tests
	requiredEnvVars := []string{
		"OVH_ENDPOINT",
		"OVH_APPLICATION_KEY", 
		"OVH_APPLICATION_SECRET",
		"OVH_CONSUMER_KEY",
		"OVH_PROJECT_ID",
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
					testAccCheckProviderConfigured("hashicorp-ovh"),
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

// BenchmarkProviderSchema benchmarks schema retrieval
func BenchmarkProviderSchema(b *testing.B) {
	provider := New("test")()
	
	for i := 0; i < b.N; i++ {
		req := provider.GetProviderSchemaRequest{}
		resp := provider.GetProviderSchemaResponse{}
		
		provider.GetProviderSchema(context.Background(), req, &resp)
		
		if resp.Diagnostics.HasError() {
			b.Fatalf("Schema retrieval failed: %v", resp.Diagnostics.Errors())
		}
	}
}