// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"

	"terraform-provider-teradata-clearscape/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure TeradataClearScapeProvider satisfies various provider interfaces.
var _ provider.Provider = &TeradataClearScapeProvider{}

// TeradataClearScapeProvider defines the provider implementation.
type TeradataClearScapeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// TeradataClearScapeProviderModel describes the provider data model.
type TeradataClearScapeProviderModel struct {
	Token types.String `tfsdk:"token"`
}

func (p *TeradataClearScapeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "teradata-clearscape"
	resp.Version = p.version
}

func (p *TeradataClearScapeProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *TeradataClearScapeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring ClearScape client")
	var config TeradataClearScapeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown ClearScape API Token",
			"The provider cannot create the ClearScape API client as there is an unknown configuration value for the ClearScape API client. ",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	token := os.Getenv("CLEARCAPE_API_TOKEN")
	if !config.Token.IsNull() {
		token = config.Token.ValueString()
	}

	if token == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing ClearScape API Token",
			"The provider cannot create the ClearScape API client as there is an unknown configuration value for the ClearScape API client.  "+
				"Set the token value in the configuration or use the CLEARCAPE_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "clearscape_token", token)

	tflog.Debug(ctx, "Creating ClearScape client")

	client, err := client.NewClient("https://api.clearscape.teradata.com/", token)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create ClearScape API client", err.Error())
		return
	}
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured ClearScape client", map[string]any{"success": true})

}

func (p *TeradataClearScapeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		EnvironmentResource,
	}
}

func (p *TeradataClearScapeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		EnvironmentDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TeradataClearScapeProvider{
			version: version,
		}
	}
}
