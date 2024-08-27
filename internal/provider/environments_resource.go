package provider

import (
	"context"
	"fmt"
	"terraform-provider-teradata-clearscape/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &environmentResource{}
	_ resource.ResourceWithConfigure = &environmentResource{}
)

func EnvironmentResource() resource.Resource {
	return &environmentResource{}
}

// environmentResource implements the resource.Resource interface.
type environmentResource struct {
	client *client.Client
}

// Metadata returns the resource type name.
func (r *environmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_environment"
}

type environmentCredentialModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

type environmentServiceModel struct {
	Name        types.String `tfsdk:"name"`
	URL         types.String `tfsdk:"url"`
	Credentials types.List   `tfsdk:"credentials"`
}

type environmentResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Region      types.String `tfsdk:"region"`
	State       types.String `tfsdk:"state"`
	IP          types.String `tfsdk:"ip"`
	DNSName     types.String `tfsdk:"dnsname"`
	Owner       types.String `tfsdk:"owner"`
	Type        types.String `tfsdk:"type"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Operation   types.String `tfsdk:"operation"`
	Password    types.String `tfsdk:"password"`
	Services    types.List   `tfsdk:"services"`
}

// Schema defines the schema for the resource.
func (r *environmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the environment.",
			},
			"last_updated": schema.StringAttribute{
				Computed:    true,
				Description: "The last time the environment was updated.",
			},
			"operation": schema.StringAttribute{
				Computed:    true,
				Description: "The last operation performed on the environment.",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The region of the environment.",
			},
			"password": schema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for the environment.",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "The current state of the environment.",
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "The IP address of the environment.",
			},
			"dnsname": schema.StringAttribute{
				Computed:    true,
				Description: "The DNS name of the environment.",
			},
			"owner": schema.StringAttribute{
				Computed:    true,
				Description: "The owner of the environment.",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "The type of the environment.",
			},
			"services": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The name of the service.",
						},
						"url": schema.StringAttribute{
							Computed:    true,
							Description: "The URL of the service.",
						},
						"credentials": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:    true,
										Description: "The name of the credential.",
									},
									"value": schema.StringAttribute{
										Computed:    true,
										Description: "The value of the credential.",
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

// Configure adds the provider configured client to the resource.
func (r *environmentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *environmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	tflog.Info(ctx, "Creating New ClearScape Environment")

	var plan environmentResourceModel

	tflog.Info(ctx, "Before Mapping")
	diags := req.Plan.Get(ctx, &plan)
	tflog.Info(ctx, "After Mapping")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "After Error", map[string]interface{}{"name": plan.Name.ValueString(), "region": plan.Region.ValueString(), "state": plan.State.ValueString(), "ip": plan.IP.ValueString(), "dnsname": plan.DNSName.ValueString(), "owner": plan.Owner.ValueString(), "type": plan.Type.ValueString(), "password": plan.Password.ValueString()})
	// Generate API request body from plan

	var envRequest client.EnvironmentCreateRequest
	envRequest.Name = plan.Name.ValueString()
	envRequest.Region = plan.Region.ValueString()
	envRequest.Password = plan.Password.ValueString()

	tflog.Info(ctx, "Environment Details %s", map[string]interface{}{"name": plan.Name.ValueString(), "region": plan.Region.ValueString(), "password": plan.Password.ValueString()})

	env, err := r.client.CreateEnvironment(envRequest)
	tflog.Info(ctx, "Environment Created %s", map[string]interface{}{"name": env.Name, "region": env.Region, "state": env.State, "ip": env.IP, "dnsname": env.DNSName, "owner": env.Owner, "type": env.Type})

	if err != nil {
		resp.Diagnostics.AddError("Failed to create environment", err.Error())
		return
	}

	for _, service := range env.Services {
		tflog.Info(ctx, "Service Details %s", map[string]interface{}{"name": service.Name, "url": service.URL})
		for _, cred := range service.Credentials {
			tflog.Info(ctx, "Service Credentials %s", map[string]interface{}{"name": cred.Name, "value": cred.Value})
		}

	}

	environment := environmentResourceModel{
		Name:     types.StringValue(env.Name),
		Region:   types.StringValue(env.Region),
		State:    types.StringValue(env.State),
		IP:       types.StringValue(env.IP),
		DNSName:  types.StringValue(env.DNSName),
		Owner:    types.StringValue(env.Owner),
		Type:     types.StringValue(env.Type),
		Password: types.StringValue(envRequest.Password),
	}

	services, _ := types.ListValueFrom(ctx, plan.Services.ElementType(ctx), env.Services)
	environment.Services = services
	tflog.Info(ctx, "Services %s", map[string]interface{}{"services": services})

	/* for _, service := range env.Services {
		var creds []environmentCredentialModel
		for _, cred := range service.Credentials {
			creds = append(creds, environmentCredentialModel{
				Name:  types.StringValue(cred.Name),
				Value: types.StringValue(cred.Value),
			})
		}

		environment.Services = append(environment.Services, environmentServiceModel{
			Name:        types.StringValue(service.Name),
			URL:         types.StringValue(service.URL),
			Credentials: creds,
		})
	} */

	plan = environment

	// Set state to fully populated data
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *environmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	var state environmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading %s ClearScape Environment", map[string]interface{}{"name": state.Name.ValueString()})
	env, err := r.client.GetEnvironment(state.Name.ValueString())
	if err != nil {
		return
	}

	// Overwrite items with refreshed state

	environment := environmentResourceModel{
		Name:     types.StringValue(env.Name),
		Region:   types.StringValue(env.Region),
		State:    types.StringValue(env.State),
		IP:       types.StringValue(env.IP),
		DNSName:  types.StringValue(env.DNSName),
		Owner:    types.StringValue(env.Owner),
		Type:     types.StringValue(env.Type),
		Password: types.StringValue(state.Password.String()),
	}

	services, d := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: environmentServiceModel{}.AttributeTypes()}, env.Services[1])
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}
	environment.Services = services

	state = environment

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (m environmentCredentialModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}
}

func (m environmentServiceModel) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"url":  types.StringType,
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *environmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan environmentResourceModel
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	env, err := r.client.UpdateEnvironment(plan.Name.String(), plan.Operation.String())
	if err != nil {
		resp.Diagnostics.AddError("Failed to update environment", err.Error())
		return
	}

	// Overwrite items with refreshed state
	environment := environmentResourceModel{
		Name:     types.StringValue(env.Name),
		Region:   types.StringValue(env.Region),
		State:    types.StringValue(env.State),
		IP:       types.StringValue(env.IP),
		DNSName:  types.StringValue(env.DNSName),
		Owner:    types.StringValue(env.Owner),
		Type:     types.StringValue(env.Type),
		Password: types.StringValue(plan.Password.String()),
	}

	services, _ := types.ListValueFrom(ctx, plan.Services.ElementType(ctx), env.Services)
	environment.Services = services
	tflog.Info(ctx, "Services %s", map[string]interface{}{"services": services})
	plan = environment

	// Set refreshed state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *environmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state environmentResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEnvironment(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Failed to Delete %s ClearScape Environment", state.Name.ValueString()), err.Error())
		return
	}
	tflog.Info(ctx, "Deleted %s ClearScape Environment", map[string]interface{}{"name": state.Name.ValueString()})
}
