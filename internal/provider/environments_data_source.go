package provider

import (
	"context"
	"fmt"
	"terraform-provider-teradata-clearscape/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &environmentDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentDataSource{}
)

func EnvironmentDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

type environmentDataSource struct {
	client *client.Client
}

// Metadata returns the data source type name.
func (d *environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environments"
}

type environmentDataSourceModel struct {
	Environments []environmentModel `tfsdk:"environments"`
}

type environmentModel struct {
	Name        types.String   `tfsdk:"name"`
	Region      types.String   `tfsdk:"region"`
	State       types.String   `tfsdk:"state"`
	IP          types.String   `tfsdk:"ip"`
	DNSName     types.String   `tfsdk:"dnsname"`
	Owner       types.String   `tfsdk:"owner"`
	Type        types.String   `tfsdk:"type"`
	Services    []serviceModel `tfsdk:"services"`
	LastUpdated types.String   `tfsdk:"last_updated"`
	Operation   types.String   `tfsdk:"operation"`
}

type serviceModel struct {
	Name        types.String      `tfsdk:"name"`
	URL         types.String      `tfsdk:"url"`
	Credentials []credentialModel `tfsdk:"credentials"`
}

type credentialModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// Schema defines the schema for the data source.
func (d *environmentDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"environments": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed: true,
						},
						"region": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
						"ip": schema.StringAttribute{
							Computed: true,
						},
						"dnsname": schema.StringAttribute{
							Computed: true,
						},
						"owner": schema.StringAttribute{
							Computed: true,
						},
						"type": schema.StringAttribute{
							Computed: true,
						},
						"services": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed: true,
									},
									"url": schema.StringAttribute{
										Computed: true,
									},
									"credentials": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"name": schema.StringAttribute{
													Computed: true,
												},
												"value": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state environmentDataSourceModel

	environments, err := d.client.GetEnvironments()
	if err != nil {
		resp.Diagnostics.AddError("Failed to get environments", err.Error())
		return
	}

	for _, env := range *environments {
		environment := environmentModel{
			Name:    types.StringValue(env.Name),
			Region:  types.StringValue(env.Region),
			State:   types.StringValue(env.State),
			IP:      types.StringValue(env.IP),
			DNSName: types.StringValue(env.DNSName),
			Owner:   types.StringValue(env.Owner),
			Type:    types.StringValue(env.Type),
		}

		for _, service := range env.Services {
			s := serviceModel{
				Name: types.StringValue(service.Name),
				URL:  types.StringValue(service.URL),
			}

			for _, cred := range service.Credentials {
				s.Credentials = append(s.Credentials, credentialModel{
					Name:  types.StringValue(cred.Name),
					Value: types.StringValue(cred.Value),
				})
			}

			environment.Services = append(environment.Services, s)
		}

		state.Environments = append(state.Environments, environment)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Configure adds the provider configured client to the data source.
func (d *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}
