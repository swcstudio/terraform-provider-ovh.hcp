package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/ovh/go-ovh/ovh"
)

type Config struct {
	OVHClient *ovh.Client
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"ovh_endpoint": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_ENDPOINT", nil),
					Description: "OVH API endpoint (ovh-eu, ovh-us, ovh-ca, kimsufi-eu, kimsufi-ca, soyoustart-eu, soyoustart-ca, runabove-ca)",
				},
				"ovh_application_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_KEY", nil),
					Description: "OVH API application key",
				},
				"ovh_application_secret": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_APPLICATION_SECRET", nil),
					Description: "OVH API application secret",
					Sensitive:   true,
				},
				"ovh_consumer_key": {
					Type:        schema.TypeString,
					Required:    true,
					DefaultFunc: schema.EnvDefaultFunc("OVH_CONSUMER_KEY", nil),
					Description: "OVH API consumer key",
					Sensitive:   true,
				},

			},
			ResourcesMap: map[string]*schema.Resource{
				"hashicorp_ovh_nomad_cluster":    resourceNomadCluster(),
				"hashicorp_ovh_vault_cluster":    resourceVaultCluster(),
				"hashicorp_ovh_consul_cluster":   resourceConsulCluster(),
				"hashicorp_ovh_boundary_cluster": resourceBoundaryCluster(),
				"hashicorp_ovh_waypoint_runner":  resourceWaypointRunner(),
				"hashicorp_ovh_packer_template":  resourcePackerTemplate(),
			},
			DataSourcesMap: map[string]*schema.Resource{
				"hashicorp_ovh_nomad_clusters":    dataSourceNomadClusters(),
				"hashicorp_ovh_vault_clusters":    dataSourceVaultClusters(),
				"hashicorp_ovh_consul_clusters":   dataSourceConsulClusters(),
				"hashicorp_ovh_boundary_clusters": dataSourceBoundaryClusters(),
			},
			ConfigureContextFunc: providerConfigure,
		}

		return p
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	ovhClient, err := ovh.NewClient(
		d.Get("ovh_endpoint").(string),
		d.Get("ovh_application_key").(string),
		d.Get("ovh_application_secret").(string),
		d.Get("ovh_consumer_key").(string),
	)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("failed to create OVH client: %w", err))
	}

	config := &Config{
		OVHClient: ovhClient,
	}

	return config, diags
}
