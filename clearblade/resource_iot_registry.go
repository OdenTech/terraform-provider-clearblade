package clearblade

import (
	"context"
	"fmt"
	"os"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &deviceRegistryResource{}
	_ resource.ResourceWithConfigure   = &deviceRegistryResource{}
	_ resource.ResourceWithImportState = &deviceRegistryResource{}
)

/*
type deviceRegistryResourceModel struct {
	Project     types.String   `tfsdk:"project"`
	Region      types.String   `tfsdk:"region"`
	LastUpdated types.String   `tfsdk:"last_updated"`
	Registry    *registryModel `tfsdk:"registry"`
}
*/

type deviceRegistryResourceModel struct {
	ID                       types.String                    `tfsdk:"id"`
	Name                     types.String                    `tfsdk:"name"`
	EventNotificationConfigs []eventNotificationConfigsModel `tfsdk:"event_notification_configs"`
	StateNotificationConfig  *stateNotificationConfigModel   `tfsdk:"state_notification_config"`
	MqttConfig               *mqttConfigModel                `tfsdk:"mqtt_config"`
	HttpConfig               *httpConfigModel                `tfsdk:"http_config"`
	LogLevel                 types.String                    `tfsdk:"log_level"`
	Credentials              []credentialsModel              `tfsdk:"credentials"`
	// Region                   types.String                    `tfsdk:"region"`
	// Project                  types.String                    `tfsdk:"project"`
	//LastUpdated              types.String                    `tfsdk:"last_updated"`
}

type eventNotificationConfigsModel struct {
	PubsubTopicName  types.String `tfsdk:"pubsub_topic_name"`
	SubfolderMatches types.String `tfsdk:"subfolder_matches"`
}

type stateNotificationConfigModel struct {
	PubsubTopicName types.String `tfsdk:"pubsub_topic_name"`
}

type mqttConfigModel struct {
	MqttEnabledState types.String `tfsdk:"mqtt_enabled_state"`
}
type httpConfigModel struct {
	HttpEnabledState types.String `tfsdk:"http_enabled_state"`
}

type credentialsModel struct {
	PublicKeyCertificate publicKeyCertificateModel `tfsdk:"public_key_certificate"`
}

type publicKeyCertificateModel struct {
	Format      types.String     `tfsdk:"format"`
	Certificate types.String     `tfsdk:"certificate"`
	X509Details x509DetailsModel `tfsdk:"x509_details"`
}

type x509DetailsModel struct {
	X509CertificateDetail x509CertificateDetailModel `tfsdk:"x509_certificate_detail"`
}

type x509CertificateDetailModel struct {
	Issuer             types.String `tfsdk:"issuer"`
	Subject            types.String `tfsdk:"subject"`
	StartTime          types.String `tfsdk:"start_time"`
	ExpiryTime         types.String `tfsdk:"expiry_time"`
	SignatureAlgorithm types.String `tfsdk:"signature_algorithm"`
	PublicKeyType      types.String `tfsdk:"public_key_type"`
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
			"id": schema.StringAttribute{
				MarkdownDescription: "The identifier of this device registry. For example, myRegistry.",
				Optional:            true,
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
					"mqtt_enabled_state": schema.StringAttribute{
						MarkdownDescription: "If enabled, allows connections using the MQTT protocol. Otherwise, MQTT connections to this registry will fail.",
						Optional:            true,
					},
				},
			},
			"http_config": schema.SingleNestedAttribute{
				MarkdownDescription: "The configuration of the HTTP bridge for a device registry.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"http_enabled_state": schema.StringAttribute{
						MarkdownDescription: "If enabled, allows devices to use DeviceService via the HTTP protocol. Otherwise, any requests to DeviceService will fail for this registry.",
						Optional:            true,
					},
				},
			},
			"log_level": schema.StringAttribute{
				MarkdownDescription: "The default logging verbosity for activity from devices in this registry. The verbosity level can be overridden by Device.log_level.",
				Required:            true,
			},
			"credentials": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					// Attributes: map[string]schema.Attribute{
					// 	"credential": schema.SingleNestedAttribute{
					// 		Required:    true,
					// 		Description: "A server-stored registry credential used to validate device credentials.",
					Attributes: map[string]schema.Attribute{
						"public_key_certificate": schema.SingleNestedAttribute{
							Optional:    true,
							Description: "A public key certificate format and data.",
							Attributes: map[string]schema.Attribute{
								"format": schema.StringAttribute{
									Description: "The certificate format.",
									Computed:    true,
								},
								"certificate": schema.StringAttribute{
									Description: "The certificate data.",
									Computed:    true,
								},
								//
								"x509_details": schema.SingleNestedAttribute{
									Optional:    true,
									Description: "Details of an X.509 certificate.",
									Attributes: map[string]schema.Attribute{
										"x509_certificate_detail": schema.SingleNestedAttribute{
											Optional:    true,
											Description: "The certificate details. Used only for X.509 certificates.",
											Attributes: map[string]schema.Attribute{
												"issuer": schema.StringAttribute{
													Description: "The entity that signed the certificate.",
													Computed:    true,
												},
												"subject": schema.StringAttribute{
													Description: "The entity the certificate and public key belong to.",
													Computed:    true,
												},
												"start_time": schema.StringAttribute{
													Description: "The time the certificate becomes valid.",
													Computed:    true,
												},
												"expiry_time": schema.StringAttribute{
													Description: "The time the certificate becomes invalid.",
													Computed:    true,
												},
												"signature_algorithm": schema.StringAttribute{
													Description: "The algorithm used to sign the certificate.",
													Computed:    true,
												},
												"public_key_type": schema.StringAttribute{
													Description: "The type of public key in the certificate.",
													Computed:    true,
												},
											},
										},
									},
								}, //
							},
						},
					},
					//}, //
					//},//
				},
			},
			// "project": schema.StringAttribute{
			// 	MarkdownDescription: "The id of the project.",
			// 	Optional:            true,
			// },
			// "region": schema.StringAttribute{
			// 	MarkdownDescription: "The name of the cloud region.",
			// 	Optional:            true,
			// },
			// "last_updated": schema.StringAttribute{
			// 	Optional: true,
			// },
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

	var plan deviceRegistryResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	event_notification_configs := []*iot.EventNotificationConfig{}
	for _, v := range plan.EventNotificationConfigs {
		event_notification_configs = append(event_notification_configs, &iot.EventNotificationConfig{
			PubsubTopicName:  v.PubsubTopicName.ValueString(),
			SubfolderMatches: v.SubfolderMatches.ValueString(),
		})
	}

	credentials := []*iot.RegistryCredential{}
	for _, v := range plan.Credentials {
		credentials = append(credentials, &iot.RegistryCredential{
			PublicKeyCertificate: &iot.PublicKeyCertificate{
				Certificate: v.PublicKeyCertificate.Certificate.ValueString(),
				Format:      v.PublicKeyCertificate.Format.ValueString(),
				X509Details: &iot.X509CertificateDetails{
					//X509CertificateDetail: &iot. {
					ExpiryTime:         v.PublicKeyCertificate.X509Details.X509CertificateDetail.ExpiryTime.ValueString(),
					Issuer:             v.PublicKeyCertificate.X509Details.X509CertificateDetail.Issuer.ValueString(),
					PublicKeyType:      v.PublicKeyCertificate.X509Details.X509CertificateDetail.PublicKeyType.ValueString(),
					SignatureAlgorithm: v.PublicKeyCertificate.X509Details.X509CertificateDetail.SignatureAlgorithm.ValueString(),
					StartTime:          v.PublicKeyCertificate.X509Details.X509CertificateDetail.StartTime.ValueString(),
					Subject:            v.PublicKeyCertificate.X509Details.X509CertificateDetail.Subject.ValueString(),
					//},
				},
			},
		})
	}

	stateNotificationConfig := &iot.StateNotificationConfig{
		PubsubTopicName: plan.StateNotificationConfig.PubsubTopicName.ValueString(),
	}

	mqttConfig := &iot.MqttConfig{
		MqttEnabledState: plan.MqttConfig.MqttEnabledState.ValueString(),
	}

	httpConfig := &iot.HttpConfig{
		HttpEnabledState: plan.HttpConfig.HttpEnabledState.ValueString(),
	}

	createRequestPayload := iot.DeviceRegistry{
		EventNotificationConfigs: event_notification_configs,
		Credentials:              credentials,
		Id:                       plan.ID.ValueString(),
		Name:                     plan.Name.ValueString(),
		StateNotificationConfig:  stateNotificationConfig,
		MqttConfig:               mqttConfig,
		HttpConfig:               httpConfig,
		LogLevel:                 plan.LogLevel.ValueString(),
	}

	//parent := fmt.Sprintf("projects/%s/locations/%s", plan.Project.ValueString(), plan.Region.ValueString())
	parent := fmt.Sprintf("projects/%s/locations/%s", os.Getenv("CLEARBLADE_PROJECT"), os.Getenv("CLEARBLADE_REGION"))

	// Create new registry
	_, err := r.client.Projects.Locations.Registries.Create(parent, &createRequestPayload).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating a device registry",
			"Could not create a new device registry, unexpected error: "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Created a device registry")

	// Map response body to schema and populate Computed attribute values
	//plan.ID = types.StringValue(registry.Id)
	//plan.Project = types.StringValue(plan.Project.ValueString())
	//plan.Region = types.StringValue(plan.Region.ValueString())

	//plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "In Create method: Start PLAN logging")
	ctx = tflog.SetField(ctx, "clearblade_provider_plan", plan)
	tflog.Info(ctx, "In Create method: Stop PLAN logging")
}

// Read resource information and refreshes the Terraform state with the latest data.
func (r *deviceRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state deviceRegistryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "In Read method: Start STATE logging")
	ctx = tflog.SetField(ctx, "clearblade_provider_state", state)
	tflog.Info(ctx, "In Read method: Stop STATE logging")

	// Get refreshed registry value from ClearBlade IoT Core
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", os.Getenv("CLEARBLADE_PROJECT"), os.Getenv("CLEARBLADE_REGION"), state.ID.ValueString())
	registry, err := r.client.Projects.Locations.Registries.Get(parent).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading ClearBlade IoT Core Registry",
			"Could not read ClearBlade IoT Core registry ID "+state.Name.ValueString()+": "+err.Error(),
		)
		return
	}

	//state.Registry.EventNotificationConfigs = types.StringValue(registry.EventNotificationConfigs)
	//state.StateNotificationConfig.PubsubTopicName = types.StringValue(registry.StateNotificationConfig.PubsubTopicName)
	//state.MqttConfig.MqttConfig = types.StringValue(registry.MqttConfig.MqttEnabledState)
	//state.HttpConfig.HttpConfig = types.StringValue(registry.HttpConfig.HttpEnabledState)
	state.LogLevel = types.StringValue(registry.LogLevel)
	//state.Region = types.StringValue(os.Getenv("CLEARBLADE_REGION"))
	//state.Project = types.StringValue(os.Getenv("CLEARBLADE_PROJECT"))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deviceRegistryResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing registry
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", os.Getenv("CLEARBLADE_PROJECT"), os.Getenv("CLEARBLADE_REGION"), state.ID.ValueString())
	_, err := r.client.Projects.Locations.Registries.Delete(parent).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting ClearBlade IoT Core Registry",
			"Could not delete Registry, unexpected error: "+err.Error(),
		)
		return
	}
}

func (r *deviceRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
