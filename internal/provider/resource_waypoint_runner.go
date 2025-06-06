package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceWaypointRunner() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Waypoint runner on OVH infrastructure for application deployment",

		CreateContext: resourceWaypointRunnerCreate,
		ReadContext:   resourceWaypointRunnerRead,
		UpdateContext: resourceWaypointRunnerUpdate,
		DeleteContext: resourceWaypointRunnerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Waypoint runner",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region for the runner",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH instance type for the runner",
			},
			"runner_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "static",
				Description: "Type of runner",
				ValidateFunc: validation.StringInSlice([]string{
					"static", "on-demand", "kubernetes",
				}, false),
			},
			"capacity": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      10,
				Description:  "Maximum concurrent jobs",
				ValidateFunc: validation.IntBetween(1, 100),
			},
			"docker_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Enable Docker support",
			},
			"kubernetes_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Kubernetes support",
			},
			"nomad_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Nomad support",
			},
			"web3_deployments": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable Web3 application deployments",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to apply to runner resources",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"runner_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Waypoint runner ID",
			},
			"token": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "Runner authentication token",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Runner endpoint URL",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Runner status",
			},
		},
	}
}

func resourceWaypointRunnerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	runnerConfig := map[string]interface{}{
		"name":              d.Get("name").(string),
		"region":            d.Get("region").(string),
		"instanceType":      d.Get("instance_type").(string),
		"runnerType":        d.Get("runner_type").(string),
		"capacity":          d.Get("capacity").(int),
		"dockerEnabled":     d.Get("docker_enabled").(bool),
		"kubernetesEnabled": d.Get("kubernetes_enabled").(bool),
		"nomadEnabled":      d.Get("nomad_enabled").(bool),
		"web3Deployments":   d.Get("web3_deployments").(bool),
		"tags":              d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/waypoint/runner", runnerConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Waypoint runner: %w", err))
	}

	runnerId := result["id"].(string)
	d.SetId(runnerId)

	return resourceWaypointRunnerRead(ctx, d, meta)
}

func resourceWaypointRunnerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	runnerId := d.Id()

	var runner map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/waypoint/runner/%s", runnerId), &runner)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Waypoint runner: %w", err))
	}

	d.Set("name", runner["name"])
	d.Set("region", runner["region"])
	d.Set("instance_type", runner["instanceType"])
	d.Set("runner_type", runner["runnerType"])
	d.Set("capacity", runner["capacity"])
	d.Set("docker_enabled", runner["dockerEnabled"])
	d.Set("kubernetes_enabled", runner["kubernetesEnabled"])
	d.Set("nomad_enabled", runner["nomadEnabled"])
	d.Set("web3_deployments", runner["web3Deployments"])
	d.Set("runner_id", runner["runnerId"])
	d.Set("token", runner["token"])
	d.Set("endpoint", runner["endpoint"])
	d.Set("status", runner["status"])

	if tags, ok := runner["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourceWaypointRunnerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	runnerId := d.Id()

	if d.HasChanges("capacity", "tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("capacity") {
			updateConfig["capacity"] = d.Get("capacity").(int)
		}
		if d.HasChange("tags") {
			updateConfig["tags"] = d.Get("tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/waypoint/runner/%s", runnerId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Waypoint runner: %w", err))
		}
	}

	return resourceWaypointRunnerRead(ctx, d, meta)
}

func resourceWaypointRunnerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	runnerId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/waypoint/runner/%s", runnerId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Waypoint runner: %w", err))
	}

	d.SetId("")
	return nil
}
