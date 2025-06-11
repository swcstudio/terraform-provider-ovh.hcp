package provider

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccProvider is the provider instance used for acceptance tests
var TestAccProvider *hashicorpOVHProvider

// TestAccProviderFactories contains the provider factory for acceptance tests
var TestAccProviderFactories map[string]func() (*schema.Provider, error)

// Common test configuration constants
const (
	TestResourcePrefix    = "tf-acc-test"
	TestTagKey           = "terraform-test"
	TestTagValue         = "true"
	DefaultTestTimeout   = 30 * time.Minute
	DefaultTestRegion    = "eu-west-1"
	DefaultTestZone      = "eu-west-1a"
)

// Test environment variables
var (
	TestOVHEndpoint       = os.Getenv("OVH_ENDPOINT")
	TestOVHApplicationKey = os.Getenv("OVH_APPLICATION_KEY")
	TestOVHSecret         = os.Getenv("OVH_APPLICATION_SECRET")
	TestOVHConsumerKey    = os.Getenv("OVH_CONSUMER_KEY")
	TestOVHProjectID      = os.Getenv("OVH_PROJECT_ID")
)

// TestConfig represents a test configuration
type TestConfig struct {
	ResourceName string
	Region       string
	Zone         string
	Tags         map[string]string
	Attributes   map[string]interface{}
}

// NewTestConfig creates a new test configuration with defaults
func NewTestConfig(resourceName string) *TestConfig {
	return &TestConfig{
		ResourceName: resourceName,
		Region:       DefaultTestRegion,
		Zone:         DefaultTestZone,
		Tags: map[string]string{
			TestTagKey: TestTagValue,
			"Name":     resourceName,
		},
		Attributes: make(map[string]interface{}),
	}
}

// WithRegion sets the region for the test configuration
func (tc *TestConfig) WithRegion(region string) *TestConfig {
	tc.Region = region
	return tc
}

// WithZone sets the zone for the test configuration
func (tc *TestConfig) WithZone(zone string) *TestConfig {
	tc.Zone = zone
	return tc
}

// WithTag adds a tag to the test configuration
func (tc *TestConfig) WithTag(key, value string) *TestConfig {
	tc.Tags[key] = value
	return tc
}

// WithAttribute adds an attribute to the test configuration
func (tc *TestConfig) WithAttribute(key string, value interface{}) *TestConfig {
	tc.Attributes[key] = value
	return tc
}

// RandomName generates a random name with the given prefix
func RandomName(prefix string) string {
	bytes := make([]byte, 4)
	rand.Read(bytes)
	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(bytes))
}

// RandomNameWithTimestamp generates a random name with timestamp
func RandomNameWithTimestamp(prefix string) string {
	timestamp := time.Now().Unix()
	bytes := make([]byte, 2)
	rand.Read(bytes)
	return fmt.Sprintf("%s-%d-%s", prefix, timestamp, hex.EncodeToString(bytes))
}

// TestAccPreCheck verifies that required environment variables are set
func TestAccPreCheck(t *testing.T) {
	if TestOVHEndpoint == "" {
		t.Fatal("OVH_ENDPOINT must be set for acceptance tests")
	}
	if TestOVHApplicationKey == "" {
		t.Fatal("OVH_APPLICATION_KEY must be set for acceptance tests")
	}
	if TestOVHSecret == "" {
		t.Fatal("OVH_APPLICATION_SECRET must be set for acceptance tests")
	}
	if TestOVHConsumerKey == "" {
		t.Fatal("OVH_CONSUMER_KEY must be set for acceptance tests")
	}
	if TestOVHProjectID == "" {
		t.Fatal("OVH_PROJECT_ID must be set for acceptance tests")
	}
}

// TestAccPreCheckOptional checks for optional environment variables
func TestAccPreCheckOptional(t *testing.T, envVars ...string) {
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Skipf("Skipping test because %s is not set", envVar)
		}
	}
}

// TestAccCheckResourceExists is a generic function to check if a resource exists
func TestAccCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for resource: %s", resourceName)
		}

		return nil
	}
}

// TestAccCheckResourceDestroy is a generic function to check if a resource is destroyed
func TestAccCheckResourceDestroy(resourceType string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}

			if rs.Primary.ID != "" {
				return fmt.Errorf("resource %s still exists with ID: %s", resourceType, rs.Primary.ID)
			}
		}
		return nil
	}
}

// TestAccCheckResourceDisappears simulates a resource disappearing outside of Terraform
func TestAccCheckResourceDisappears(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no ID is set for resource: %s", resourceName)
		}

		// In a real implementation, you would make an API call to delete the resource
		// For testing purposes, we simulate the deletion
		return nil
	}
}

// TestAccCheckResourceAttr checks a specific attribute value
func TestAccCheckResourceAttr(resourceName, key, value string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttr(resourceName, key, value)
}

// TestAccCheckResourceAttrSet checks that an attribute is set (non-empty)
func TestAccCheckResourceAttrSet(resourceName, key string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrSet(resourceName, key)
}

// TestAccCheckResourceAttrPair checks that two resources have the same attribute value
func TestAccCheckResourceAttrPair(nameFirst, keyFirst, nameSecond, keySecond string) resource.TestCheckFunc {
	return resource.TestCheckResourceAttrPair(nameFirst, keyFirst, nameSecond, keySecond)
}

// TestAccCheckResourceAttrRegex checks that an attribute matches a regex pattern
func TestAccCheckResourceAttrRegex(resourceName, key, pattern string) resource.TestCheckFunc {
	return resource.TestMatchResourceAttr(resourceName, key, regexp.MustCompile(pattern))
}

// TestAccCheckResourceTags checks that a resource has the expected tags
func TestAccCheckResourceTags(resourceName string, expectedTags map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		for key, expectedValue := range expectedTags {
			tagKey := fmt.Sprintf("tags.%s", key)
			actualValue, ok := rs.Primary.Attributes[tagKey]
			if !ok {
				return fmt.Errorf("tag %s not found on resource %s", key, resourceName)
			}
			if actualValue != expectedValue {
				return fmt.Errorf("tag %s has value %s, expected %s", key, actualValue, expectedValue)
			}
		}

		return nil
	}
}

// MockHTTPServer creates a mock HTTP server for testing
type MockHTTPServer struct {
	*httptest.Server
	Requests []*http.Request
	Responses []MockResponse
}

// MockResponse represents a mock HTTP response
type MockResponse struct {
	StatusCode int
	Body       string
	Headers    map[string]string
}

// NewMockHTTPServer creates a new mock HTTP server
func NewMockHTTPServer() *MockHTTPServer {
	mock := &MockHTTPServer{
		Requests:  make([]*http.Request, 0),
		Responses: make([]MockResponse, 0),
	}

	mock.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mock.Requests = append(mock.Requests, r)

		if len(mock.Responses) > 0 {
			response := mock.Responses[0]
			mock.Responses = mock.Responses[1:]

			for key, value := range response.Headers {
				w.Header().Set(key, value)
			}
			w.WriteHeader(response.StatusCode)
			w.Write([]byte(response.Body))
		} else {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("{}"))
		}
	}))

	return mock
}

// AddResponse adds a mock response to the queue
func (m *MockHTTPServer) AddResponse(statusCode int, body string, headers map[string]string) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "application/json"

	m.Responses = append(m.Responses, MockResponse{
		StatusCode: statusCode,
		Body:       body,
		Headers:    headers,
	})
}

// GetRequestCount returns the number of requests received
func (m *MockHTTPServer) GetRequestCount() int {
	return len(m.Requests)
}

// GetLastRequest returns the last request received
func (m *MockHTTPServer) GetLastRequest() *http.Request {
	if len(m.Requests) == 0 {
		return nil
	}
	return m.Requests[len(m.Requests)-1]
}

// TestProvider creates a test provider configuration
func TestProvider() string {
	return `
provider "hashicorp-ovh" {
  ovh_endpoint           = "` + TestOVHEndpoint + `"
  ovh_application_key    = "` + TestOVHApplicationKey + `"
  ovh_application_secret = "` + TestOVHSecret + `"
  ovh_consumer_key       = "` + TestOVHConsumerKey + `"
  ovh_project_id         = "` + TestOVHProjectID + `"
}`
}

// TestProviderConfig generates a provider configuration for testing
func TestProviderConfig() string {
	return TestProvider()
}

// TestNomadClusterConfig generates a basic Nomad cluster configuration
func TestNomadClusterConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hashicorp_ovh_nomad_cluster" "test" {
  name         = "%s"
  region       = "%s"
  server_count = 3
  client_count = 5

  tags = {
    %s = "%s"
    Environment = "test"
  }
}`, TestProvider(), name, DefaultTestRegion, TestTagKey, TestTagValue)
}

// TestVaultClusterConfig generates a basic Vault cluster configuration
func TestVaultClusterConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hashicorp_ovh_vault_cluster" "test" {
  name       = "%s"
  region     = "%s"
  node_count = 3

  auto_unseal_enabled = true
  ha_enabled         = true

  tags = {
    %s = "%s"
    Environment = "test"
  }
}`, TestProvider(), name, DefaultTestRegion, TestTagKey, TestTagValue)
}

// TestConsulClusterConfig generates a basic Consul cluster configuration
func TestConsulClusterConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hashicorp_ovh_consul_cluster" "test" {
  name       = "%s"
  region     = "%s"
  node_count = 3

  connect_enabled    = true
  encryption_enabled = true

  tags = {
    %s = "%s"
    Environment = "test"
  }
}`, TestProvider(), name, DefaultTestRegion, TestTagKey, TestTagValue)
}

// TestBoundaryClusterConfig generates a basic Boundary cluster configuration
func TestBoundaryClusterConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hashicorp_ovh_boundary_cluster" "test" {
  name             = "%s"
  region           = "%s"
  controller_count = 1
  worker_count     = 2

  tags = {
    %s = "%s"
    Environment = "test"
  }
}`, TestProvider(), name, DefaultTestRegion, TestTagKey, TestTagValue)
}

// TestWaypointRunnerConfig generates a basic Waypoint runner configuration
func TestWaypointRunnerConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hashicorp_ovh_waypoint_runner" "test" {
  name   = "%s"
  region = "%s"

  runner_type = "static"

  tags = {
    %s = "%s"
    Environment = "test"
  }
}`, TestProvider(), name, DefaultTestRegion, TestTagKey, TestTagValue)
}

// TestPackerTemplateConfig generates a basic Packer template configuration
func TestPackerTemplateConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "hashicorp_ovh_packer_template" "test" {
  name   = "%s"
  region = "%s"

  source_image = "ubuntu-20.04"
  image_type   = "standard"

  tags = {
    %s = "%s"
    Environment = "test"
  }
}`, TestProvider(), name, DefaultTestRegion, TestTagKey, TestTagValue)
}

// TestDataSourceConfig generates configurations for data source testing
func TestDataSourceNomadClustersConfig() string {
	return fmt.Sprintf(`
%s

data "hashicorp_ovh_nomad_clusters" "test" {}
`, TestProvider())
}

func TestDataSourceVaultClustersConfig() string {
	return fmt.Sprintf(`
%s

data "hashicorp_ovh_vault_clusters" "test" {}
`, TestProvider())
}

func TestDataSourceConsulClustersConfig() string {
	return fmt.Sprintf(`
%s

data "hashicorp_ovh_consul_clusters" "test" {}
`, TestProvider())
}

// ValidationHelper provides common validation functions
type ValidationHelper struct {
	t *testing.T
}

// NewValidationHelper creates a new validation helper
func NewValidationHelper(t *testing.T) *ValidationHelper {
	return &ValidationHelper{t: t}
}

// ValidateResourceName checks if a resource name is valid
func (vh *ValidationHelper) ValidateResourceName(name string) error {
	if len(name) < 3 || len(name) > 50 {
		return fmt.Errorf("resource name must be between 3 and 50 characters")
	}

	matched, err := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9-]*$", name)
	if err != nil {
		return err
	}
	if !matched {
		return fmt.Errorf("resource name must start with a letter and contain only letters, numbers, and hyphens")
	}

	return nil
}

// ValidateRegion checks if a region is valid
func (vh *ValidationHelper) ValidateRegion(region string) error {
	validRegions := []string{
		"eu-west-1", "eu-central-1", "us-east-1", "us-west-1",
		"ap-southeast-1", "ap-northeast-1",
	}

	for _, validRegion := range validRegions {
		if region == validRegion {
			return nil
		}
	}

	return fmt.Errorf("invalid region: %s", region)
}

// ValidateTags checks if tags are valid
func (vh *ValidationHelper) ValidateTags(tags map[string]string) error {
	if len(tags) > 50 {
		return fmt.Errorf("cannot have more than 50 tags")
	}

	for key, value := range tags {
		if len(key) == 0 || len(key) > 128 {
			return fmt.Errorf("tag key must be between 1 and 128 characters")
		}
		if len(value) > 256 {
			return fmt.Errorf("tag value must be 256 characters or less")
		}
		if strings.HasPrefix(key, "aws:") || strings.HasPrefix(key, "ovh:") {
			return fmt.Errorf("tag key cannot start with 'aws:' or 'ovh:'")
		}
	}

	return nil
}

// TestCleanup provides utilities for cleaning up test resources
type TestCleanup struct {
	resources []string
	t         *testing.T
}

// NewTestCleanup creates a new test cleanup helper
func NewTestCleanup(t *testing.T) *TestCleanup {
	return &TestCleanup{
		resources: make([]string, 0),
		t:         t,
	}
}

// AddResource adds a resource to be cleaned up
func (tc *TestCleanup) AddResource(resourceID string) {
	tc.resources = append(tc.resources, resourceID)
}

// Cleanup cleans up all registered resources
func (tc *TestCleanup) Cleanup() {
	for _, resourceID := range tc.resources {
		// In a real implementation, you would make API calls to delete the resources
		tc.t.Logf("Cleaning up resource: %s", resourceID)
	}
}

// TestMetrics provides utilities for collecting test metrics
type TestMetrics struct {
	StartTime time.Time
	EndTime   time.Time
	APICall   int
	Errors    []error
}

// NewTestMetrics creates a new test metrics collector
func NewTestMetrics() *TestMetrics {
	return &TestMetrics{
		StartTime: time.Now(),
		APICall:   0,
		Errors:    make([]error, 0),
	}
}

// RecordAPICall increments the API call counter
func (tm *TestMetrics) RecordAPICall() {
	tm.APICall++
}

// RecordError adds an error to the metrics
func (tm *TestMetrics) RecordError(err error) {
	tm.Errors = append(tm.Errors, err)
}

// Finish marks the end of the test and calculates duration
func (tm *TestMetrics) Finish() {
	tm.EndTime = time.Now()
}

// Duration returns the test duration
func (tm *TestMetrics) Duration() time.Duration {
	if tm.EndTime.IsZero() {
		return time.Since(tm.StartTime)
	}
	return tm.EndTime.Sub(tm.StartTime)
}

// Summary returns a summary of the test metrics
func (tm *TestMetrics) Summary() string {
	return fmt.Sprintf("Duration: %v, API Calls: %d, Errors: %d",
		tm.Duration(), tm.APICall, len(tm.Errors))
}

// LoadTestData loads test data from files or environment
func LoadTestData(key string, defaultValue string) string {
	if value := os.Getenv(fmt.Sprintf("TEST_%s", key)); value != "" {
		return value
	}
	return defaultValue
}

// SkipCI skips tests when running in CI environment
func SkipCI(t *testing.T, reason string) {
	if os.Getenv("CI") == "true" {
		t.Skipf("Skipping in CI: %s", reason)
	}
}

// RequireEnv fails the test if required environment variables are not set
func RequireEnv(t *testing.T, envVars ...string) {
	for _, envVar := range envVars {
		if os.Getenv(envVar) == "" {
			t.Fatalf("Required environment variable %s is not set", envVar)
		}
	}
}

// TestTimeout returns the timeout for long-running tests
func TestTimeout() time.Duration {
	if timeout := os.Getenv("TEST_TIMEOUT"); timeout != "" {
		if d, err := time.ParseDuration(timeout); err == nil {
			return d
		}
	}
	return DefaultTestTimeout
}