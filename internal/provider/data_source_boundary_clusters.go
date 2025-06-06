package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBoundaryClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about Boundary clusters on OVH infrastructure",

		ReadContext: dataSourceBoundaryClustersRead,

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
				Description: "List of Boundary clusters",
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
						"controller_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of controller nodes",
						},
						"worker_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of worker nodes",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type",
						},
						"vault_integration": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Vault integration enabled",
						},
						"controller_endpoints": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Controller endpoints",
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

func dataSourceBoundaryClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	var clusters []map[string]interface{}
	err := config.OVHClient.Get("/cloud/project/boundary/cluster", &clusters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Boundary clusters: %w", err))
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
			"id":                   cluster["id"],
			"name":                 cluster["name"],
			"region":               cluster["region"],
			"controller_count":     cluster["controllerCount"],
			"worker_count":         cluster["workerCount"],
			"instance_type":        cluster["instanceType"],
			"vault_integration":    cluster["vaultIntegration"],
			"controller_endpoints": cluster["controllerEndpoints"],
			"ui_url":               cluster["uiUrl"],
			"status":               cluster["status"],
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
