package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/ovh/go-ovh/ovh"
)

type HashiCorpOVHProvider struct {
	version string
}

type HashiCorpOVHProviderModel struct {
	OVHEndpoint          types.String `tfsdk:"ovh_endpoint"`
	OVHApplicationKey    types.String `tfsdk:"ovh_application_key"`
	OVHApplicationSecret types.String `tfsdk:"ovh_application_secret"`
	OVHConsumerKey       types.String `tfsdk:"ovh_consumer_key"`
}

type Config struct {
	OVHClient *ovh.Client
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &HashiCorpOVHProvider{
			version: version,
		}
	}
}

func (p *HashiCorpOVHProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "hashicorp-ovh"
	resp.Version = p.version
}

func (p *HashiCorpOVHProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"ovh_endpoint": schema.StringAttribute{
				Description: "OVH API endpoint (ovh-eu, ovh-us, ovh-ca, kimsufi-eu, kimsufi-ca, soyoustart-eu, soyoustart-ca, runabove-ca)",
				Required:    true,
			},
			"ovh_application_key": schema.StringAttribute{
				Description: "OVH API application key",
				Required:    true,
			},
			"ovh_application_secret": schema.StringAttribute{
				Description: "OVH API application secret",
				Required:    true,
				Sensitive:   true,
			},
			"ovh_consumer_key": schema.StringAttribute{
				Description: "OVH API consumer key",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *HashiCorpOVHProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring HashiCorp OVH provider")

	var config HashiCorpOVHProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ovhEndpoint := os.Getenv("OVH_ENDPOINT")
	if !config.OVHEndpoint.IsNull() {
		ovhEndpoint = config.OVHEndpoint.ValueString()
	}

	ovhApplicationKey := os.Getenv("OVH_APPLICATION_KEY")
	if !config.OVHApplicationKey.IsNull() {
		ovhApplicationKey = config.OVHApplicationKey.ValueString()
	}

	ovhApplicationSecret := os.Getenv("OVH_APPLICATION_SECRET")
	if !config.OVHApplicationSecret.IsNull() {
		ovhApplicationSecret = config.OVHApplicationSecret.ValueString()
	}

	ovhConsumerKey := os.Getenv("OVH_CONSUMER_KEY")
	if !config.OVHConsumerKey.IsNull() {
		ovhConsumerKey = config.OVHConsumerKey.ValueString()
	}

	if ovhEndpoint == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Endpoint Configuration",
			"While configuring the provider, the OVH endpoint was not found in "+
				"the OVH_ENDPOINT environment variable or provider "+
				"configuration block ovh_endpoint attribute.",
		)
	}

	if ovhApplicationKey == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Application Key Configuration",
			"While configuring the provider, the OVH application key was not found in "+
				"the OVH_APPLICATION_KEY environment variable or provider "+
				"configuration block ovh_application_key attribute.",
		)
	}

	if ovhApplicationSecret == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Application Secret Configuration",
			"While configuring the provider, the OVH application secret was not found in "+
				"the OVH_APPLICATION_SECRET environment variable or provider "+
				"configuration block ovh_application_secret attribute.",
		)
	}

	if ovhConsumerKey == "" {
		resp.Diagnostics.AddError(
			"Missing OVH Consumer Key Configuration",
			"While configuring the provider, the OVH consumer key was not found in "+
				"the OVH_CONSUMER_KEY environment variable or provider "+
				"configuration block ovh_consumer_key attribute.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "ovh_endpoint", ovhEndpoint)
	ctx = tflog.SetField(ctx, "ovh_application_key", ovhApplicationKey)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ovh_application_secret")
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "ovh_consumer_key")

	tflog.Debug(ctx, "Creating OVH client")

	ovhClient, err := ovh.NewClient(
		ovhEndpoint,
		ovhApplicationKey,
		ovhApplicationSecret,
		ovhConsumerKey,
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create OVH API Client",
			"An unexpected error occurred when creating the OVH API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"OVH Client Error: "+err.Error(),
		)
		return
	}

	providerConfig := &Config{
		OVHClient: ovhClient,
	}

	resp.DataSourceData = providerConfig
	resp.ResourceData = providerConfig

	tflog.Info(ctx, "Configured HashiCorp OVH provider", map[string]any{"success": true})
}

func (p *HashiCorpOVHProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
	}
}

func (p *HashiCorpOVHProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
	}
}
