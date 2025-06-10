package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func TestAccProvider(t *testing.T) {
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"hashicorp-ovh": providerserver.NewProtocol6WithError(New("test")()),
	}
}

var testAccProtoV6ProviderFactories map[string]func() (tfprotov6.ProviderServer, error)
