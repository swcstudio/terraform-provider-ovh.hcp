package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVaultClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about Vault clusters on OVH infrastructure",

		ReadContext: dataSourceVaultClustersRead,

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
				Description: "List of Vault clusters",
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
						"node_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of nodes",
						},
						"instance_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Instance type",
						},
						"storage_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Storage backend type",
						},
						"auto_unseal": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Auto-unseal enabled",
						},
						"audit_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Audit logging enabled",
						},
						"cluster_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Cluster URL",
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

func dataSourceVaultClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	var clusters []map[string]interface{}
	err := config.OVHClient.Get("/cloud/project/vault/cluster", &clusters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Vault clusters: %w", err))
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
			"id":            cluster["id"],
			"name":          cluster["name"],
			"region":        cluster["region"],
			"node_count":    cluster["nodeCount"],
			"instance_type": cluster["instanceType"],
			"storage_type":  cluster["storageType"],
			"auto_unseal":   cluster["autoUnseal"],
			"audit_enabled": cluster["auditEnabled"],
			"cluster_url":   cluster["clusterUrl"],
			"ui_url":        cluster["uiUrl"],
			"status":        cluster["status"],
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
