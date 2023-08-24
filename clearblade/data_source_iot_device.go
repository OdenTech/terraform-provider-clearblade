package clearblade

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &devicesDataSource{}
	_ datasource.DataSourceWithConfigure = &devicesDataSource{}
)

// devicesDataSourceModel maps the data source schema data.
type devicesDataSourceModel struct {
	Devices  []devicesModel `tfsdk:"devices"`
	Registry types.String   `tfsdk:"registry"`
}

// devicesModel maps device schema data.
type devicesModel struct {
	ID    types.String `tfsdk:"id"`
	Name  types.String `tfsdk:"name"`
	NumID types.String `tfsdk:"num_id"`
}

func NewDevicesDataSource() datasource.DataSource {
	return &devicesDataSource{}
}

type devicesDataSource struct {
	client *iot.Service
}

func (d *devicesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List devices in a device registry.",
		Attributes: map[string]schema.Attribute{
			"devices": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "The user-defined device identifier. The device ID must be unique within a device registry.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "The resource path name. For example, projects/p1/locations/us-central1/registries/registry0/devices/dev0 or projects/p1/locations/us-central1/registries/registry0/devices/{numId}.",
						},
						"num_id": schema.StringAttribute{
							Computed:    true,
							Description: "A server-defined unique numeric ID for the device. This is a more compact way to identify devices, and it is globally unique.",
						},
					},
				},
			},
			"registry": schema.StringAttribute{
				Description: "The name of the device registry where this device should be created.",
				Required:    true,
			},
		},
	}
}

func (d *devicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state devicesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	tflog.Info(ctx, "requesting device listing from Clearblade IoT Core")
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", os.Getenv("CLEARBLADE_PROJECT"), os.Getenv("CLEARBLADE_REGION"), state.Registry.ValueString())
	devices, err := d.client.Projects.Locations.Registries.Devices.List(parent).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Clearblade IoT Core devices. Make sure your credentials are correct and you have access "+
				"to the project, or that you have the correct permissions.",
			err.Error(),
		)
		return
	}

	ctx = tflog.SetField(ctx, "device length", strconv.Itoa(len(devices.Devices)))

	// Map response body to model
	for _, device := range devices.Devices {
		deviceState := devicesModel{
			ID:    types.StringValue(device.Id),
			Name:  types.StringValue(device.Name),
			NumID: types.StringValue(strconv.FormatUint(device.NumId, 10)),
		}

		state.Devices = append(state.Devices, deviceState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *devicesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

// Configure adds the provider configured client to the data source.
func (d *devicesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*iot.Service)
}
