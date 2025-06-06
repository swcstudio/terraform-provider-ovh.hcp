package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceConsulCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Consul cluster on OVH infrastructure with service mesh capabilities",

		CreateContext: resourceConsulClusterCreate,
		ReadContext:   resourceConsulClusterRead,
		UpdateContext: resourceConsulClusterUpdate,
		DeleteContext: resourceConsulClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Consul cluster",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region for the cluster",
			},
			"server_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of Consul server nodes",
				ValidateFunc: validation.IntBetween(1, 7),
			},
			"client_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3,
				Description:  "Number of Consul client nodes",
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH instance type for Consul nodes",
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Consul datacenter name",
			},
			"connect_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Consul Connect service mesh",
			},
			"acl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Consul ACL system",
			},
			"encryption_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable gossip encryption",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable TLS encryption",
			},
			"ui_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Consul UI",
			},
			"monitoring_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable monitoring and metrics",
			},
			"backup_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable automated backups",
			},
			"web3_services": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Web3 service discovery",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to apply to cluster resources",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"server_endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Consul server endpoints",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ui_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Consul UI URL",
			},
			"gossip_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Gossip encryption key",
			},
			"master_token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "ACL master token",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster status",
			},
		},
	}
}

func resourceConsulClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterConfig := map[string]interface{}{
		"name":               d.Get("name").(string),
		"region":             d.Get("region").(string),
		"serverCount":        d.Get("server_count").(int),
		"clientCount":        d.Get("client_count").(int),
		"instanceType":       d.Get("instance_type").(string),
		"datacenter":         d.Get("datacenter").(string),
		"connectEnabled":     d.Get("connect_enabled").(bool),
		"aclEnabled":         d.Get("acl_enabled").(bool),
		"encryptionEnabled":  d.Get("encryption_enabled").(bool),
		"tlsEnabled":         d.Get("tls_enabled").(bool),
		"uiEnabled":          d.Get("ui_enabled").(bool),
		"monitoringEnabled":  d.Get("monitoring_enabled").(bool),
		"backupEnabled":      d.Get("backup_enabled").(bool),
		"web3Services":       d.Get("web3_services").(bool),
		"tags":               d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/consul/cluster", clusterConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Consul cluster: %w", err))
	}

	clusterId := result["id"].(string)
	d.SetId(clusterId)

	return resourceConsulClusterRead(ctx, d, meta)
}

func resourceConsulClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	var cluster map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/consul/cluster/%s", clusterId), &cluster)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Consul cluster: %w", err))
	}

	d.Set("name", cluster["name"])
	d.Set("region", cluster["region"])
	d.Set("server_count", cluster["serverCount"])
	d.Set("client_count", cluster["clientCount"])
	d.Set("instance_type", cluster["instanceType"])
	d.Set("datacenter", cluster["datacenter"])
	d.Set("connect_enabled", cluster["connectEnabled"])
	d.Set("acl_enabled", cluster["aclEnabled"])
	d.Set("encryption_enabled", cluster["encryptionEnabled"])
	d.Set("tls_enabled", cluster["tlsEnabled"])
	d.Set("ui_enabled", cluster["uiEnabled"])
	d.Set("monitoring_enabled", cluster["monitoringEnabled"])
	d.Set("backup_enabled", cluster["backupEnabled"])
	d.Set("web3_services", cluster["web3Services"])
	d.Set("server_endpoints", cluster["serverEndpoints"])
	d.Set("ui_url", cluster["uiUrl"])
	d.Set("status", cluster["status"])

	if gossipKey, ok := cluster["gossipKey"].(string); ok {
		d.Set("gossip_key", gossipKey)
	}

	if masterToken, ok := cluster["masterToken"].(string); ok {
		d.Set("master_token", masterToken)
	}

	if tags, ok := cluster["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourceConsulClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	if d.HasChanges("server_count", "client_count", "tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("server_count") {
			updateConfig["serverCount"] = d.Get("server_count").(int)
		}
		if d.HasChange("client_count") {
			updateConfig["clientCount"] = d.Get("client_count").(int)
		}
		if d.HasChange("tags") {
			updateConfig["tags"] = d.Get("tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/consul/cluster/%s", clusterId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Consul cluster: %w", err))
		}
	}

	return resourceConsulClusterRead(ctx, d, meta)
}

func resourceConsulClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/consul/cluster/%s", clusterId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Consul cluster: %w", err))
	}

	d.SetId("")
	return nil
}
