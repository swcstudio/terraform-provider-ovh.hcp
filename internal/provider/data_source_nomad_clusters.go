package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNomadClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about Nomad clusters on OVH infrastructure",

		ReadContext: dataSourceNomadClustersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter clusters by OVH region",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter clusters by status",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Nomad clusters",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster ID",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster name",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "OVH region",
						},
						"server_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of server nodes",
						},
						"client_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of client nodes",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type",
						},
						"datacenter": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Datacenter name",
						},
						"vault_integration": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Vault integration enabled",
						},
						"consul_integration": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Consul integration enabled",
						},
						"server_endpoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Server endpoints",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"ui_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "UI URL",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster status",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Creation timestamp",
						},
						"tags": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Cluster tags",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceNomadClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	var clusters []map[string]interface{}
	err := config.OVHClient.Get("/cloud/project/nomad/cluster", &clusters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Nomad clusters: %w", err))
	}

	region := d.Get("region").(string)
	status := d.Get("status").(string)

	var filteredClusters []map[string]interface{}
	for _, cluster := range clusters {
		if region != "" && cluster["region"].(string) != region {
			continue
		}
		if status != "" && cluster["status"].(string) != status {
			continue
		}
		filteredClusters = append(filteredClusters, cluster)
	}

	clusterList := make([]interface{}, len(filteredClusters))
	for i, cluster := range filteredClusters {
		clusterMap := map[string]interface{}{
			"id":                 cluster["id"],
			"name":               cluster["name"],
			"region":             cluster["region"],
			"server_count":       cluster["serverCount"],
			"client_count":       cluster["clientCount"],
			"instance_type":      cluster["instanceType"],
			"datacenter":         cluster["datacenter"],
			"vault_integration":  cluster["vaultIntegration"],
			"consul_integration": cluster["consulIntegration"],
			"server_endpoints":   cluster["serverEndpoints"],
			"ui_url":             cluster["uiUrl"],
			"status":             cluster["status"],
			"created_at":         cluster["createdAt"],
		}

		if tags, ok := cluster["tags"].(map[string]interface{}); ok {
			clusterMap["tags"] = tags
		}

		clusterList[i] = clusterMap
	}

	d.Set("clusters", clusterList)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}
