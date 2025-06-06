package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourcePackerTemplate() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a Packer template on OVH infrastructure for image building",

		CreateContext: resourcePackerTemplateCreate,
		ReadContext:   resourcePackerTemplateRead,
		UpdateContext: resourcePackerTemplateUpdate,
		DeleteContext: resourcePackerTemplateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the Packer template",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "OVH region for image building",
			},
			"source_image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Source image for building",
			},
			"instance_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "OVH instance type for building",
			},
			"builders": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Packer builders configuration",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"provisioners": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Packer provisioners configuration",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"post_processors": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Packer post-processors configuration",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"variables": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Template variables",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"auto_build": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable automatic builds on changes",
			},
			"build_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      3600,
				Description:  "Build timeout in seconds",
				ValidateFunc: validation.IntBetween(300, 7200),
			},
			"web3_tools": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Include Web3 development tools",
			},
			"kata_support": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Include Kata containers support",
			},
			"tags": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Tags to apply to template resources",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"template_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Packer template ID",
			},
			"last_build_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Last successful build ID",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Generated image ID",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Template status",
			},
		},
	}
}

func resourcePackerTemplateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	templateConfig := map[string]interface{}{
		"name":           d.Get("name").(string),
		"region":         d.Get("region").(string),
		"sourceImage":    d.Get("source_image").(string),
		"instanceType":   d.Get("instance_type").(string),
		"builders":       d.Get("builders").([]interface{}),
		"provisioners":   d.Get("provisioners").([]interface{}),
		"postProcessors": d.Get("post_processors").([]interface{}),
		"variables":      d.Get("variables"),
		"autoBuild":      d.Get("auto_build").(bool),
		"buildTimeout":   d.Get("build_timeout").(int),
		"web3Tools":      d.Get("web3_tools").(bool),
		"kataSupport":    d.Get("kata_support").(bool),
		"tags":           d.Get("tags"),
	}

	var result map[string]interface{}
	err := config.OVHClient.Post("/cloud/project/packer/template", templateConfig, &result)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to create Packer template: %w", err))
	}

	templateId := result["id"].(string)
	d.SetId(templateId)

	return resourcePackerTemplateRead(ctx, d, meta)
}

func resourcePackerTemplateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	templateId := d.Id()

	var template map[string]interface{}
	err := config.OVHClient.Get(fmt.Sprintf("/cloud/project/packer/template/%s", templateId), &template)
	if err != nil {
		d.SetId("")
		return diag.FromErr(fmt.Errorf("failed to read Packer template: %w", err))
	}

	d.Set("name", template["name"])
	d.Set("region", template["region"])
	d.Set("source_image", template["sourceImage"])
	d.Set("instance_type", template["instanceType"])
	d.Set("builders", template["builders"])
	d.Set("provisioners", template["provisioners"])
	d.Set("post_processors", template["postProcessors"])
	d.Set("variables", template["variables"])
	d.Set("auto_build", template["autoBuild"])
	d.Set("build_timeout", template["buildTimeout"])
	d.Set("web3_tools", template["web3Tools"])
	d.Set("kata_support", template["kataSupport"])
	d.Set("template_id", template["templateId"])
	d.Set("last_build_id", template["lastBuildId"])
	d.Set("image_id", template["imageId"])
	d.Set("status", template["status"])

	if tags, ok := template["tags"].(map[string]interface{}); ok {
		d.Set("tags", tags)
	}

	return nil
}

func resourcePackerTemplateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	templateId := d.Id()

	if d.HasChanges("source_image", "builders", "provisioners", "post_processors", "variables", "auto_build", "build_timeout", "tags") {
		updateConfig := map[string]interface{}{}

		if d.HasChange("source_image") {
			updateConfig["sourceImage"] = d.Get("source_image").(string)
		}
		if d.HasChange("builders") {
			updateConfig["builders"] = d.Get("builders").([]interface{})
		}
		if d.HasChange("provisioners") {
			updateConfig["provisioners"] = d.Get("provisioners").([]interface{})
		}
		if d.HasChange("post_processors") {
			updateConfig["postProcessors"] = d.Get("post_processors").([]interface{})
		}
		if d.HasChange("variables") {
			updateConfig["variables"] = d.Get("variables")
		}
		if d.HasChange("auto_build") {
			updateConfig["autoBuild"] = d.Get("auto_build").(bool)
		}
		if d.HasChange("build_timeout") {
			updateConfig["buildTimeout"] = d.Get("build_timeout").(int)
		}
		if d.HasChange("tags") {
			updateConfig["tags"] = d.Get("tags")
		}

		err := config.OVHClient.Put(fmt.Sprintf("/cloud/project/packer/template/%s", templateId), updateConfig, nil)
		if err != nil {
			return diag.FromErr(fmt.Errorf("failed to update Packer template: %w", err))
		}
	}

	return resourcePackerTemplateRead(ctx, d, meta)
}

func resourcePackerTemplateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	_ = diag.Diagnostics{}

	templateId := d.Id()

	err := config.OVHClient.Delete(fmt.Sprintf("/cloud/project/packer/template/%s", templateId), nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to delete Packer template: %w", err))
	}

	d.SetId("")
	return nil
}
