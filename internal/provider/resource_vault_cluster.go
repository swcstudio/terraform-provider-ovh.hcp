package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceVaultCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Vault cluster on OVH infrastructure with enterprise features",

		CreateContext: resourceVaultClusterCreate,
		ReadContext:   resourceVaultClusterRead,
		UpdateContext: resourceVaultClusterUpdate,
		DeleteContext: resourceVaultClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Vault cluster",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region for the cluster",
			},
			"node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of Vault nodes",
				ValidateFunc: validation.IntBetween(1, 7),
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH instance type for Vault nodes",
			},
			"storage_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "consul",
				Description: "Vault storage backend type",
				ValidateFunc: validation.StringInSlice([]string{
					"consul", "raft", "etcd", "dynamodb",
				}, false),
			},
			"auto_unseal": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable auto-unseal with OVH KMS",
			},
			"audit_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable audit logging",
			},
			"performance_replication": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable performance replication",
			},
			"disaster_recovery": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable disaster recovery replication",
			},
			"web3_secrets": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Web3 secrets engine",
			},
			"kubernetes_auth": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Kubernetes authentication",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to apply to cluster resources",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"cluster_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vault cluster URL",
			},
			"ui_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Vault UI URL",
			},
			"root_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Initial root token",
			},
			"unseal_keys": {
				Type:        schema.TypeList,
				Computed:    true,
				Sensitive:   true,
				Description: "Unseal keys",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster status",
			},
		},
	}
}

func resourceVaultClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterConfig := map[string]interface{}{
		"name":                   d.Get("name").(string),
		"region":                 d.Get("region").(string),
		"nodeCount":              d.Get("node_count").(int),
		"instanceType":           d.Get("instance_type").(string),
		"storageType":            d.Get("storage_type").(string),
		"autoUnseal":             d.Get("auto_unseal").(bool),
		"auditEnabled":           d.Get("audit_enabled").(bool),
		"performanceReplication": d.Get("performance_replication").(bool),
		"disasterRecovery":       d.Get("disaster_recovery").(bool),
		"web3Secrets":            d.Get("web3_secrets").(bool),
		"kubernetesAuth":         d.Get("kubernetes_auth").(bool),
		"tags":                   d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/vault/cluster", clusterConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Vault cluster: %w", err))
	}

	clusterId := result["id"].(string)
	d.SetId(clusterId)

	return resourceVaultClusterRead(ctx, d, meta)
}

func resourceVaultClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	var cluster map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/vault/cluster/%s", clusterId), &cluster)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Vault cluster: %w", err))
	}

	d.Set("name", cluster["name"])
	d.Set("region", cluster["region"])
	d.Set("node_count", cluster["nodeCount"])
	d.Set("instance_type", cluster["instanceType"])
	d.Set("storage_type", cluster["storageType"])
	d.Set("auto_unseal", cluster["autoUnseal"])
	d.Set("audit_enabled", cluster["auditEnabled"])
	d.Set("performance_replication", cluster["performanceReplication"])
	d.Set("disaster_recovery", cluster["disasterRecovery"])
	d.Set("web3_secrets", cluster["web3Secrets"])
	d.Set("kubernetes_auth", cluster["kubernetesAuth"])
	d.Set("cluster_url", cluster["clusterUrl"])
	d.Set("ui_url", cluster["uiUrl"])
	d.Set("status", cluster["status"])

	if rootToken, ok := cluster["rootToken"].(string); ok {
		d.Set("root_token", rootToken)
	}

	if unsealKeys, ok := cluster["unsealKeys"].([]interface{}); ok {
		d.Set("unseal_keys", unsealKeys)
	}

	if tags, ok := cluster["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourceVaultClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	if d.HasChanges("node_count", "tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("node_count") {
			updateConfig["nodeCount"] = d.Get("node_count").(int)
		}
		if d.HasChange("tags") {
			updateConfig["tags"] = d.Get("tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/vault/cluster/%s", clusterId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Vault cluster: %w", err))
		}
	}

	return resourceVaultClusterRead(ctx, d, meta)
}

func resourceVaultClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/vault/cluster/%s", clusterId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Vault cluster: %w", err))
	}

	d.SetId("")
	return nil
}
