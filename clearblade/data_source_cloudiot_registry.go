package clearblade

import (
	"context"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &deviceRegistriesDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceRegistriesDataSource{}
)

// deviceRegistriesDataSourceModel maps the data source schema data.
type deviceRegistriesDataSourceModel struct {
	Project          types.String            `tfsdk:"project"`
	Region           types.String            `tfsdk:"region"`
	DeviceRegistries []deviceRegistriesModel `tfsdk:"device_registries"`
}

// deviceRegistriesModel maps deviceRegistry schema data.
type deviceRegistriesModel struct {
	ID                       types.String                   `tfsdk:"id"`
	Name                     types.String                   `tfsdk:"name"`
	EventNotificationConfigs []eventNotificationConfigModel `tfsdk:"event_notification_configs"`
	StateNotificationConfig  stateNotificationConfigModel   `tfsdk:"state_notification_config"`
	HttpConfig               httpConfigModel                `tfsdk:"http_config"`
	MqttConfig               mqttConfigModel                `tfsdk:"mqtt_config"`
	LogLevel                 types.String                   `tfsdk:"log_level"`
	Credentials              []credentialsModel             `tfsdk:"credentials"`
}

type eventNotificationConfigModel struct {
	SubFolderMatches types.String `tfsdk:"subfolder_matches"`
	PubsubTopicName  types.String `tfsdk:"pubsub_topic_name"`
}

type stateNotificationConfigModel struct {
	PubsubTopicName types.String `tfsdk:"pubsub_topic_name"`
}

type httpConfigModel struct {
	HttpEnabledState types.String `tfsdk:"http_enabled_state"`
}
type mqttConfigModel struct {
	MqttEnabledState types.String `tfsdk:"mqtt_enabled_state"`
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

func NewDeviceRegistriesDataSource() datasource.DataSource {
	return &deviceRegistriesDataSource{}
}

type deviceRegistriesDataSource struct {
	client *iot.Service
}

func (d *deviceRegistriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registries"
}

// Configure adds the provider configured client to the data source.
func (d *deviceRegistriesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*iot.Service)
}

func (d *deviceRegistriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List of device registries in a project.",
		Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Required:    true,
				Description: "The name of the project to list device registries for.",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The name of the region to list device registries for.",
			},
			"device_registries": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"event_notification_configs": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"event_notification_config": schema.SingleNestedAttribute{
										Required:    true,
										Description: "The configuration for forwarding telemetry events.",
										Attributes: map[string]schema.Attribute{
											"subfolder_matches": schema.StringAttribute{
												Description: "This field is used only for telemetry events; subfolders are not supported for state changes.",
												Computed:    true,
											},
											"pubsub_topic_name": schema.StringAttribute{
												Description: "A Cloud Pub/Sub topic name. For example, projects/myProject/topics/deviceEvents.",
												Computed:    true,
											},
										},
									},
								},
							},
						},
						"state_notification_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "The configuration for notification of new states received from the device.",
							Attributes: map[string]schema.Attribute{
								"pubsub_topic_name": schema.StringAttribute{
									Description: "A Cloud Pub/Sub topic name. For example, projects/myProject/topics/deviceEvents.",
									Computed:    true,
								},
							},
						},
						"http_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "The configuration of the HTTP bridge for a device registry.",
							Attributes: map[string]schema.Attribute{
								"http_enabled_state": schema.StringAttribute{
									Description: "If enabled, allows devices to use DeviceService via the HTTP protocol. Otherwise, any requests to DeviceService will fail for this registry.",
									Computed:    true,
								},
							},
						},
						"mqtt_config": schema.SingleNestedAttribute{
							Required:    true,
							Description: "The configuration of MQTT for a device registry.",
							Attributes: map[string]schema.Attribute{
								"mqtt_enabled_state": schema.StringAttribute{
									Description: "If enabled, allows connections using the MQTT protocol. Otherwise, MQTT connections to this registry will fail.",
									Computed:    true,
								},
							},
						},
						"log_level": schema.StringAttribute{
							Description: "The logging verbosity for device activity. Specifies which events should be written to logs.",
							Computed:    true,
						},
						"credentials": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"credential": schema.SingleNestedAttribute{
										Required:    true,
										Description: "A server-stored registry credential used to validate device credentials.",
										Attributes: map[string]schema.Attribute{
											"public_key_certificate": schema.SingleNestedAttribute{
												Required:    true,
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
													"x509_certificate_detail": schema.SingleNestedAttribute{
														Required:    true,
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

func (d *deviceRegistriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state deviceRegistriesDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)

	if state.Project.ValueString() != "" {
		tflog.SetField(ctx, "project", state.Project.ValueString())
	}

	if state.Region.ValueString() != "" {
		tflog.SetField(ctx, "region", state.Region.ValueString())
	}

	tflog.Info(ctx, "requesting device registry listing from Clearblade IoT Core")
	// To-Do: Seems like the Go client for IoT Core is returning a single value instead of an array of device registries
	/* parent := fmt.Sprintf("projects/%s/locations/%s", state.Project.ValueString(), state.Region.ValueString())
	device_registries, err := d.client.Projects.Locations.Registries.List(parent).Do()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to read Clearblade IoT Core device registries. Make sure your credentials are correct and you have access "+
				"to the project, or that you have the correct permissions.",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "device registry")
	tflog.Info(ctx, device_registries)
	
	for _, device_registry := range device_registries {
		drState := deviceRegistriesModel{
			Name:  types.StringValue(device_registry.Name),
			LogLevel: types.StringValue(device_registry.LogLevel),
		}
		state.DeviceRegistries = append(state.DeviceRegistries, drState)
	} */

}
