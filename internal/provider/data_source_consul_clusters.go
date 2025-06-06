package provider

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceConsulClusters() *schema.Resource {
	return &schema.Resource{
		Description: "Retrieves information about Consul clusters on OVH infrastructure",

		ReadContext: dataSourceConsulClustersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter clusters by OVH region",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter clusters by datacenter",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter clusters by status",
			},
			"clusters": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of Consul clusters",
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
						"connect_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Connect service mesh enabled",
						},
						"acl_enabled": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "ACL system enabled",
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

func dataSourceConsulClustersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	var diags diag.Diagnostics

	var clusters []map[string]interface{}
	err := config.OVHClient.Get("/cloud/project/consul/cluster", &clusters)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to read Consul clusters: %w", err))
	}

	region := d.Get("region").(string)
	datacenter := d.Get("datacenter").(string)
	status := d.Get("status").(string)

	var filteredClusters []map[string]interface{}
	for _, cluster := range clusters {
		if region != "" && cluster["region"].(string) != region {
			continue
		}
		if datacenter != "" && cluster["datacenter"].(string) != datacenter {
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
			"id":               cluster["id"],
			"name":             cluster["name"],
			"region":           cluster["region"],
			"server_count":     cluster["serverCount"],
			"client_count":     cluster["clientCount"],
			"instance_type":    cluster["instanceType"],
			"datacenter":       cluster["datacenter"],
			"connect_enabled":  cluster["connectEnabled"],
			"acl_enabled":      cluster["aclEnabled"],
			"server_endpoints": cluster["serverEndpoints"],
			"ui_url":           cluster["uiUrl"],
			"status":           cluster["status"],
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
