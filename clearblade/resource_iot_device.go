package clearblade

import (
	"context"
	"fmt"
	"os"

	//"time"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deviceResource{}
	_ resource.ResourceWithConfigure   = &deviceResource{}
	_ resource.ResourceWithImportState = &deviceResource{}
)

type deviceResourceModel struct {
	ID                 types.String             `tfsdk:"id"`
	Name               types.String             `tfsdk:"name"`
	NumID              types.String             `tfsdk:"num_id"`
	Credentials        []deviceCredentialsModel `tfsdk:"credentials"`
	LastHeartbeatTime  types.String             `tfsdk:"last_heartbeat_time"`
	LastEventTime      types.String             `tfsdk:"last_event_time"`
	LastStateTime      types.String             `tfsdk:"last_state_time"`
	LastConfigAckTime  types.String             `tfsdk:"last_config_ack_time"`
	LastConfigSendTime types.String             `tfsdk:"last_config_send_time"`
	blocked            types.Bool               `tfsdk:"blocked"`
	LastErrorTime      types.String             `tfsdk:"last_error_time"`
	LastErrorStatus    *lastErrorStatusModel    `tfsdk:"last_error_status"`
	Config             *configModel             `tfsdk:"config"`
	State              *stateModel              `tfsdk:"state"`
	LogLevel           types.String             `tfsdk:"log_level"`
	Metadata           types.String             `tfsdk:"metadata"`
	GatewayConfig      *gatewayConfigModel      `tfsdk:"gateway_config"`
	Registry           types.String             `tfsdk:"id"`
}

type deviceCredentialsModel struct {
	PublicKeyCertificate *devicePublicKeyCertificateModel `tfsdk:"public_key_certificate"`
}

type devicePublicKeyCertificateModel struct {
	ExpirationTime types.String    `tfsdk:"expiration_time"`
	PublicKey      *publicKeyModel `tfsdk:"public_key"`
}

type publicKeyModel struct {
	format types.String `tfsdk:"format"`
	Key    types.String `tfsdk:"key"`
}

type lastErrorStatusModel struct {
	Code    types.Number   `tfsdk:"code"`
	Details []detailsModel `tfsdk:"details"`
	Message types.String   `tfsdk:"message"`
}

type detailsModel struct{}

type configModel struct {
	Version         types.String `tfsdk:"version"`
	CloudUpdateTime types.String `tfsdk:"cloud_update_time"`
	DeviceAckTime   types.String `tfsdk:"device_ack_time"`
	BinaryData      types.String `tfsdk:"binary_data"`
}

type stateModel struct {
	UpdateTime types.String `tfsdk:"update_time"`
	BinaryData types.String `tfsdk:"binary_data"`
}

type gatewayConfigModel struct {
	GatewayType             types.String `tfsdk:"gateway_type"`
	GatewayAuthMethod       types.String `tfsdk:"gateway_auth_method"`
	LastAccessedGatewayID   types.String `tfsdk:"last_accessed_gateway_id"`
	LastAccessedGatewayTime types.String `tfsdk:"last_accessed_gateway_time"`
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
			"id": schema.StringAttribute{
				Description: "The user-defined device identifier. The device ID must be unique within a device registry.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "The resource path name. For example, projects/p1/locations/us-central1/registries/registry0/devices/dev0 or projects/p1/locations/us-central1/registries/registry0/devices/{numId}.",
				Optional:    true,
			},
			"num_id": schema.StringAttribute{
				Description: "A server-defined unique numeric ID for the device. This is a more compact way to identify devices, and it is globally unique.",
				Computed:    true,
			},
			"credentials": schema.ListNestedAttribute{
				Optional:    true,
				Description: "The credentials used to authenticate this device.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"expiration_time": schema.StringAttribute{
							Computed:    true,
							Optional:    true,
							Description: "The time at which this credential becomes invalid.",
						},
						"public_key": schema.SingleNestedAttribute{
							Required:            true,
							MarkdownDescription: "A public key used to verify the signature of JSON Web Tokens (JWTs).",
							Attributes: map[string]schema.Attribute{
								"format": schema.StringAttribute{
									Required: true,
									Validators: []validator.String{
										stringvalidator.OneOf(
											"RSA_PEM",
											"RSA_X509_PEM",
											"ES256_PEM",
											"ES256_X509_PEM",
										),
									},
									Description: `The format of the key. Possible values: ["RSA_PEM", "RSA_X509_PEM", "ES256_PEM", "ES256_X509_PEM"]`,
								},
								"key": schema.StringAttribute{
									Required:    true,
									Description: "The key data.",
								},
							},
						},
					},
				},
			},
			"last_heartbeat_time": schema.StringAttribute{
				MarkdownDescription: "The last time an MQTT PINGREQ was received.",
				Computed:            true,
			},
			"last_event_time": schema.StringAttribute{
				Description: "The last time a telemetry event was received.",
				Computed:    true,
			},
			"last_state_time": schema.StringAttribute{
				MarkdownDescription: "The last time a state event was received.",
				Computed:            true,
			},
			"last_config_ack_time": schema.StringAttribute{
				Description: "The last time a cloud-to-device config version acknowledgment was received from the device.",
				Computed:    true,
			},
			"last_config_send_time": schema.StringAttribute{
				Description: "The last time a cloud-to-device config version was sent to the device.",
				Computed:    true,
			},
			"blocked": schema.BoolAttribute{
				Description: "If a device is blocked, connections or requests from this device will fail.",
				Optional:    true,
			},
			"last_error_time": schema.StringAttribute{
				Description: "The time the most recent error occurred, such as a failure to publish to Cloud Pub/Sub.",
				Computed:    true,
			},
			"last_error_status": schema.SingleNestedAttribute{
				Description: "The error message of the most recent error, such as a failure to publish to Cloud Pub/Sub.",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"code": schema.NumberAttribute{
						Optional:    true,
						Description: `The status code, which should be an enum value of google.rpc.Code.`,
					},
					"details": schema.ListNestedAttribute{
						Optional:    true,
						Description: `A list of messages that carry the error details.`,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								// "@type": schema.MapAttribute{
								// 	/* ... */
								// },
							},
						},
					},
					"message": schema.StringAttribute{
						Optional:    true,
						Description: `A developer-facing error message, which should be in English.`,
					},
				},
			},
			"config": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The most recent device configuration, which is eventually sent from Cloud IoT Core to the device.",
				Attributes: map[string]schema.Attribute{
					"version": schema.StringAttribute{
						Computed:    true,
						Description: `The version of this update.`,
					},
					"cloud_update_time": schema.StringAttribute{
						Computed:    true,
						Description: `The time at which this configuration version was updated in Cloud IoT Core.`,
					},
					"device_ack_time": schema.StringAttribute{
						Computed:    true,
						Description: `The time at which Cloud IoT Core received the acknowledgment from the device, indicating that the device has received this configuration version.`,
					},
					"binary_data": schema.StringAttribute{
						Optional:    true,
						Description: `The device configuration data.`,
					},
				},
			},
			"state": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "The state most recently received from the device.",
				Attributes: map[string]schema.Attribute{
					"update_time": schema.StringAttribute{
						Optional:    true,
						Description: `The time at which this state version was updated in Cloud IoT Core.`,
					},
					"binary_data": schema.StringAttribute{
						Optional:    true,
						Description: `The device state data.`,
					},
				},
			},
			"log_level": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(
						"NONE",
						"ERROR",
						"INFO",
						"DEBUG",
						"",
					),
				},
				Description: `The logging verbosity for device activity. Possible values: ["NONE", "ERROR", "INFO", "DEBUG"]`,
			},
			"metadata": schema.SingleNestedAttribute{
				Optional:    true,
				Description: `The metadata key-value pairs assigned to the device.`,
				Attributes:  map[string]schema.Attribute{},
			},
			"gateway_config": schema.SingleNestedAttribute{
				Optional:    true,
				Description: `Gateway-related configuration and state.`,
				Attributes: map[string]schema.Attribute{
					"gateway_type": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"GATEWAY",
								"NON_GATEWAY",
								"",
							),
						},
						Description: `Indicates whether the device is a gateway. Default value: "NON_GATEWAY" Possible values: ["GATEWAY", "NON_GATEWAY"]`,
					},
					"gateway_auth_method": schema.StringAttribute{
						Optional: true,
						Validators: []validator.String{
							stringvalidator.OneOf(
								"ASSOCIATION_ONLY",
								"DEVICE_AUTH_TOKEN_ONLY",
								"ASSOCIATION_AND_DEVICE_AUTH_TOKEN",
							),
						},
						Description: `Indicates whether the device is a gateway. Possible values: ["ASSOCIATION_ONLY", "DEVICE_AUTH_TOKEN_ONLY", "ASSOCIATION_AND_DEVICE_AUTH_TOKEN"]`,
					},
					"last_accessed_gateway_id": schema.StringAttribute{
						Computed:    true,
						Description: `The ID of the gateway the device accessed most recently.`,
					},
					"last_accessed_gateway_time": schema.StringAttribute{
						Computed:    true,
						Description: `The most recent time at which the device accessed the gateway specified in last_accessed_gateway.`,
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
	payload.Id = plan.ID.ValueString()
	if plan.GatewayConfig.GatewayType.ValueString() != "" {
		payload.GatewayConfig = &iot.GatewayConfig{
			GatewayType:       plan.GatewayConfig.GatewayType.ValueString(),
			GatewayAuthMethod: plan.GatewayConfig.GatewayAuthMethod.ValueString(),
		}
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", os.Getenv("CLEARBLADE_PROJECT"), os.Getenv("CLEARBLADE_REGION"), plan.Registry.ValueString())

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

	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

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

func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
