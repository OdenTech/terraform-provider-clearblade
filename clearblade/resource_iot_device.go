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
	Name        types.String `tfsdk:"name"`
	Registry    types.String `tfsdk:"registry"`
	Project     types.String `tfsdk:"project"`
	Region      types.String `tfsdk:"region"`
	LastUpdated types.String `tfsdk:"last_updated"`
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the device.",
				Required:            true,
			},
			"registry": schema.StringAttribute{
				MarkdownDescription: "The name of the device registry.",
				Required:            true,
			},
			"project": schema.StringAttribute{
				MarkdownDescription: "The name of the device registry.",
				Required:            true,
			},
			"region": schema.StringAttribute{
				MarkdownDescription: "The name of the device registry.",
				Required:            true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
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
	tflog.Debug(ctx, "Creating iot device registry resource")

	var plan deviceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//var device iot.Device

	device := iot.Device{
		Id: plan.Name.ValueString(),
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", plan.Project.ValueString(), plan.Region.ValueString(), plan.Registry.ValueString())

	// Create new device
	_, err := r.client.Projects.Locations.Registries.Devices.Create(parent, &device).Do()
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
