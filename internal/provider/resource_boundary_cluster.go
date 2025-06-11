package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBoundaryCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Boundary cluster on OVH infrastructure for secure access management",

		CreateContext: resourceBoundaryClusterCreate,
		ReadContext:   resourceBoundaryClusterRead,
		UpdateContext: resourceBoundaryClusterUpdate,
		DeleteContext: resourceBoundaryClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Boundary cluster",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region for the cluster",
			},
			"controller_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of Boundary controller nodes",
				ValidateFunc: validation.IntBetween(1, 5),
			},
			"worker_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of Boundary worker nodes",
				ValidateFunc: validation.IntBetween(1, 20),
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH instance type for Boundary nodes",
			},
			"database_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "postgresql",
				Description: "Database backend type",
				ValidateFunc: validation.StringInSlice([]string{
					"postgresql", "mysql",
				}, false),
			},
			"vault_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Vault integration for credential brokering",
			},
			"ldap_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable LDAP authentication",
			},
			"oidc_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable OIDC authentication",
			},
			"session_recording": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable session recording",
			},
			"multi_hop_sessions": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable multi-hop sessions",
			},
			"web3_targets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Web3 target management",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to apply to cluster resources",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"controller_endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Boundary controller endpoints",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ui_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Boundary UI URL",
			},
			"auth_method_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default auth method ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster status",
			},
		},
	}
}

func resourceBoundaryClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterConfig := map[string]interface{}{
		"name":             d.Get("name").(string),
		"region":           d.Get("region").(string),
		"controllerCount":  d.Get("controller_count").(int),
		"workerCount":      d.Get("worker_count").(int),
		"instanceType":     d.Get("instance_type").(string),
		"databaseType":     d.Get("database_type").(string),
		"vaultIntegration": d.Get("vault_integration").(bool),
		"ldapAuth":         d.Get("ldap_auth").(bool),
		"oidcAuth":         d.Get("oidc_auth").(bool),
		"sessionRecording": d.Get("session_recording").(bool),
		"multiHopSessions": d.Get("multi_hop_sessions").(bool),
		"web3Targets":      d.Get("web3_targets").(bool),
		"tags":             d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/boundary/cluster", clusterConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Boundary cluster: %w", err))
	}

	clusterId := result["id"].(string)
	d.SetId(clusterId)

	return resourceBoundaryClusterRead(ctx, d, meta)
}

func resourceBoundaryClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	var cluster map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/boundary/cluster/%s", clusterId), &cluster)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Boundary cluster: %w", err))
	}

	d.Set("name", cluster["name"])
	d.Set("region", cluster["region"])
	d.Set("controller_count", cluster["controllerCount"])
	d.Set("worker_count", cluster["workerCount"])
	d.Set("instance_type", cluster["instanceType"])
	d.Set("database_type", cluster["databaseType"])
	d.Set("vault_integration", cluster["vaultIntegration"])
	d.Set("ldap_auth", cluster["ldapAuth"])
	d.Set("oidc_auth", cluster["oidcAuth"])
	d.Set("session_recording", cluster["sessionRecording"])
	d.Set("multi_hop_sessions", cluster["multiHopSessions"])
	d.Set("web3_targets", cluster["web3Targets"])
	d.Set("controller_endpoints", cluster["controllerEndpoints"])
	d.Set("ui_url", cluster["uiUrl"])
	d.Set("auth_method_id", cluster["authMethodId"])
	d.Set("status", cluster["status"])

	if tags, ok := cluster["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourceBoundaryClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	if d.HasChanges("controller_count", "worker_count", "tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("controller_count") {
			updateConfig["controllerCount"] = d.Get("controller_count").(int)
		}
		if d.HasChange("worker_count") {
			updateConfig["workerCount"] = d.Get("worker_count").(int)
		}
		if d.HasChange("tags") {
			updateConfig["tags"] = d.Get("tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/boundary/cluster/%s", clusterId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Boundary cluster: %w", err))
		}
	}

	return resourceBoundaryClusterRead(ctx, d, meta)
}

func resourceBoundaryClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/boundary/cluster/%s", clusterId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Boundary cluster: %w", err))
	}

	d.SetId("")
	return nil
}
