package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNomadCluster() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Nomad cluster on OVH infrastructure with enterprise features",

		CreateContext: resourceNomadClusterCreate,
		ReadContext:   resourceNomadClusterRead,
		UpdateContext: resourceNomadClusterUpdate,
		DeleteContext: resourceNomadClusterDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Nomad cluster",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region for the cluster",
				ValidateFunc: validation.StringInSlice([]string{
					"GRA", "SBG", "RBX", "BHS", "WAW", "DE", "UK", "SGP", "SYD", "US-EAST", "US-WEST",
				}, false),
			},
			"server_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of Nomad server nodes",
				ValidateFunc: validation.IntBetween(1, 5),
			},
			"client_count": {
				Type:         schema.TypeInt,
				Required:     true,
				Description:  "Number of Nomad client nodes",
				ValidateFunc: validation.IntBetween(0, 100),
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH instance type for cluster nodes",
				ValidateFunc: validation.StringInSlice([]string{
					"s1-2", "s1-4", "s1-8", "c2-7", "c2-15", "c2-30", "c2-60", "c2-120",
					"r2-15", "r2-30", "r2-60", "r2-120", "t1-45", "t1-90", "t1-180",
				}, false),
			},
			"datacenter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Nomad datacenter name",
			},
			"vault_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Vault integration for secrets management",
			},
			"consul_integration": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Consul integration for service discovery",
			},
			"acl_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Nomad ACL system",
			},
			"tls_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable TLS encryption",
			},
			"web3_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Web3 blockchain integration",
			},
			"kata_containers": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Kata containers for secure workloads",
			},
			"gpu_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable GPU support for ML workloads",
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
				Description: "Nomad server endpoints",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ui_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Nomad UI URL",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster status",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Cluster creation timestamp",
			},
		},
	}
}

func resourceNomadClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterName := d.Get("name").(string)
	region := d.Get("region").(string)
	serverCount := d.Get("server_count").(int)
	clientCount := d.Get("client_count").(int)
	instanceType := d.Get("instance_type").(string)
	datacenter := d.Get("datacenter").(string)

	clusterConfig := map[string]interface{}{
		"name":               clusterName,
		"region":             region,
		"serverCount":        serverCount,
		"clientCount":        clientCount,
		"instanceType":       instanceType,
		"datacenter":         datacenter,
		"vaultIntegration":   d.Get("vault_integration").(bool),
		"consulIntegration":  d.Get("consul_integration").(bool),
		"aclEnabled":         d.Get("acl_enabled").(bool),
		"tlsEnabled":         d.Get("tls_enabled").(bool),
		"web3Enabled":        d.Get("web3_enabled").(bool),
		"kataContainers":     d.Get("kata_containers").(bool),
		"gpuSupport":         d.Get("gpu_support").(bool),
		"tags":               d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/nomad/cluster", clusterConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Nomad cluster: %w", err))
	}

	clusterId := result["id"].(string)
	d.SetId(clusterId)

	if err := waitForClusterReady(ctx, config, clusterId); err != nil {
		return diag.FromErr(fmt.Errorf("cluster creation timeout: %w", err))
	}

	return resourceNomadClusterRead(ctx, d, meta)
}

func resourceNomadClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	var cluster map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/nomad/cluster/%s", clusterId), &cluster)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Nomad cluster: %w", err))
	}

	d.Set("name", cluster["name"])
	d.Set("region", cluster["region"])
	d.Set("server_count", cluster["serverCount"])
	d.Set("client_count", cluster["clientCount"])
	d.Set("instance_type", cluster["instanceType"])
	d.Set("datacenter", cluster["datacenter"])
	d.Set("vault_integration", cluster["vaultIntegration"])
	d.Set("consul_integration", cluster["consulIntegration"])
	d.Set("acl_enabled", cluster["aclEnabled"])
	d.Set("tls_enabled", cluster["tlsEnabled"])
	d.Set("web3_enabled", cluster["web3Enabled"])
	d.Set("kata_containers", cluster["kataContainers"])
	d.Set("gpu_support", cluster["gpuSupport"])
	d.Set("server_endpoints", cluster["serverEndpoints"])
	d.Set("ui_url", cluster["uiUrl"])
	d.Set("status", cluster["status"])
	d.Set("created_at", cluster["createdAt"])

	if tags, ok := cluster["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourceNomadClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/nomad/cluster/%s", clusterId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Nomad cluster: %w", err))
		}

		if err := waitForClusterReady(ctx, config, clusterId); err != nil {
			return diag.FromErr(fmt.Errorf("cluster update timeout: %w", err))
		}
	}

	return resourceNomadClusterRead(ctx, d, meta)
}

func resourceNomadClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	clusterId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/nomad/cluster/%s", clusterId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Nomad cluster: %w", err))
	}

	d.SetId("")
	return nil
}

func waitForClusterReady(ctx context.Context, config *Config, clusterId string) error {
	timeout := time.After(30 * time.Minute)
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for cluster to be ready")
		case <-ticker.C:
			var cluster map[string]interface{}
			err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/nomad/cluster/%s", clusterId), &cluster)
			if err != nil {
				continue
			}

			if status, ok := cluster["status"].(string); ok && status == "READY" {
				return nil
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
