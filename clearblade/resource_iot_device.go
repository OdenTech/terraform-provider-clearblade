package clearblade

import (
	"context"
	"fmt"
	"time"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deviceResource{}
	_ resource.ResourceWithConfigure = &deviceResource{}
	//_ resource.ResourceWithImportState = &deviceRegistryResource{}
)

type deviceResourceModel struct {
	Registry    types.String `tfsdk:"registry"`
	Project     types.String `tfsdk:"project"`
	Region      types.String `tfsdk:"region"`
	LastUpdated types.String `tfsdk:"last_updated"`
	Device      *deviceModel `tfsdk:"device"`
}

type deviceModel struct {
	ID types.String `tfsdk:"id"`
	//Name          types.String        `tfsdk:"name"`
	GatewayConfig *gatewayConfigModel `tfsdk:"gateway_config"`
}

type gatewayConfigModel struct {
	GatewayType       types.String `tfsdk:"gateway_type"`
	GatewayAuthMethod types.String `tfsdk:"gateway_auth_method"`
}

func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceResource is the resource implementation.
type deviceResource struct {
	client *iot.Service
}

// Metadata returns the data source type name.
func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iot_device"
}

// Schema defines the schema for the resource.
func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"registry": schema.StringAttribute{
				MarkdownDescription: "The name of the device registry where this device should be created.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The id of the project.",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The name of the cloud region.",
				Required:            true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
			"device": schema.SingleNestedAttribute{
				MarkdownDescription: "A unique resource.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "A unique name for the device resource.",
						Required:            true,
					},
					"gateway_config": schema.SingleNestedAttribute{
						MarkdownDescription: "Gateway-related configuration and state.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"gateway_type": schema.StringAttribute{
								MarkdownDescription: "Indicates whether the device is a gateway.",
								Optional:            true,
							},
							"gateway_auth_method": schema.StringAttribute{
								MarkdownDescription: "Indicates how to authorize and/or authenticate devices to access the gateway.",
								Optional:            true,
							},
						},
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *deviceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*iot.Service)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *iot.Service, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating iot device resource")

	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := &iot.Device{}
	payload.Id = plan.Device.ID.ValueString()
	if plan.Device.GatewayConfig.GatewayType.ValueString() != "" {
		payload.GatewayConfig = &iot.GatewayConfig{
			GatewayType:       plan.Device.GatewayConfig.GatewayType.ValueString(),
			GatewayAuthMethod: plan.Device.GatewayConfig.GatewayAuthMethod.ValueString(),
		}
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", plan.Project.ValueString(), plan.Region.ValueString(), plan.Registry.ValueString())

	// Create new device
	_, err := r.client.Projects.Locations.Registries.Devices.Create(parent, payload).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a device",
			"Could not create a new device, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Created a device")

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
