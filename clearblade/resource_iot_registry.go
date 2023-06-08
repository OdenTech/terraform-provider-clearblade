package clearblade

import (
	"context"
	"fmt"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &deviceRegistryResource{}
	_ resource.ResourceWithConfigure = &deviceRegistryResource{}
	//_ resource.ResourceWithImportState = &deviceRegistryResource{}
)

type deviceRegistryResourceModel struct {
	Project     types.String   `tfsdk:"project"`
	Region      types.String   `tfsdk:"region"`
	LastUpdated types.String   `tfsdk:"last_updated"`
	Registry    *registryModel `tfsdk:"registry"`
}

type registryModel struct {
	ID                       types.String                    `tfsdk:"id"`
	Name                     types.String                    `tfsdk:"name"`
	EventNotificationConfigs []eventNotificationConfigsModel `tfsdk:"event_notification_configs"`
	StateNotificationConfig  *stateNotificationConfigModel   `tfsdk:"state_notification_config"`
	MqttConfig               *mqttConfigModel                `tfsdk:"mqtt_config"`
	HttpConfig               *httpConfigModel                `tfsdk:"http_config"`
	LogLevel                 types.String                    `tfsdk:"log_level"`
}

type eventNotificationConfigsModel struct {
	PubsubTopicName  types.String `tfsdk:"pubsub_topic_name"`
	SubfolderMatches types.String `tfsdk:"subfolder_matches"`
}

type stateNotificationConfigModel struct {
	PubsubTopicName types.String `tfsdk:"pubsub_topic_name"`
}

type mqttConfigModel struct {
	MqttConfig types.String `tfsdk:"mqtt_config"`
}
type httpConfigModel struct {
	HttpConfig types.String `tfsdk:"http_config"`
}

func NewDeviceRegistryResource() resource.Resource {
	return &deviceRegistryResource{}
}

// deviceRegistryResource is the resource implementation.
type deviceRegistryResource struct {
	client *iot.Service
}

// Metadata returns the data source type name.
func (r *deviceRegistryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_iot_registry"
}

// Schema defines the schema for the resource.
func (r *deviceRegistryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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
			"registry": schema.SingleNestedAttribute{
				MarkdownDescription: "A container for a group of devices.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						MarkdownDescription: "The identifier of this device registry. For example, myRegistry.",
						Required:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "The resource path name. For example, projects/example-project/locations/us-central1/registries/my-registry.",
						Optional:            true,
					},
					"event_notification_configs": schema.ListNestedAttribute{
						MarkdownDescription: "The configuration for notification of telemetry events received from the device. " +
							"All telemetry events that were successfully published by the device and acknowledged by Clearblade IoT Core are guaranteed to be delivered to Cloud Pub/Sub.",
						Required: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"pubsub_topic_name": schema.StringAttribute{
									MarkdownDescription: "A Cloud Pub/Sub topic name. For example, projects/myProject/topics/deviceEvents.",
									Required:            true,
								},
								"subfolder_matches": schema.StringAttribute{
									MarkdownDescription: "If the subfolder name matches this string exactly, this configuration will be used. The string must not include the leading '/' character. If empty, all strings are matched. " +
										"This field is used only for telemetry events; subfolders are not supported for state changes.",
									Required: true,
								},
							},
						},
					},
					"state_notification_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The configuration for notification of new states received from the device.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"pubsub_topic_name": schema.StringAttribute{
								MarkdownDescription: "A Cloud Pub/Sub topic name. For example, projects/myProject/topics/deviceEvents.",
								Optional:            true,
							},
						},
					},
					"mqtt_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The configuration of MQTT for a device registry.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"mqtt_config": schema.StringAttribute{
								MarkdownDescription: "If enabled, allows connections using the MQTT protocol. Otherwise, MQTT connections to this registry will fail.",
								Optional:            true,
							},
						},
					},
					"http_config": schema.SingleNestedAttribute{
						MarkdownDescription: "The configuration of the HTTP bridge for a device registry.",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"http_config": schema.StringAttribute{
								MarkdownDescription: "If enabled, allows devices to use DeviceService via the HTTP protocol. Otherwise, any requests to DeviceService will fail for this registry.",
								Optional:            true,
							},
						},
					},
					"log_level": schema.StringAttribute{
						MarkdownDescription: "The default logging verbosity for activity from devices in this registry. The verbosity level can be overridden by Device.log_level.",
						Required:            true,
					},
				},
			},
		},
	}
}

// Configure adds the provider configured client to the data source.
func (r *deviceRegistryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *deviceRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Creating iot device registry resource")
}

// Read refreshes the Terraform state with the latest data.
func (r *deviceRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
}
