package clearblade

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// deviceRegistriesDataSourceModel maps the data source schema data.
type deviceRegistriesDataSourceModel struct {
	DeviceRegistries []deviceRegistriesModel `tfsdk:"deviceRegistries"`
}

// deviceRegistriesModel maps deviceRegistry schema data.
type deviceRegistriesModel struct {
	Name                     types.String                   `tfsdk:"name"`
	EventNotificationConfigs []eventNotificationConfigModel `tfsdk:"eventNotificationConfigs"`
	StateNotificationConfig  stateNotificationConfigModel   `tfsdk:"stateNotificationConfig"`
	MqttConfig               mqttConfigModel                `tfsdk:"mqttConfig"`
	HttpConfig               httpConfigModel                `tfsdk:"httpConfig"`
	LogLevel                 types.String                   `tfsdk:"logLevel"`
}

type eventNotificationConfigModel struct {
	SubFolderMatches types.String `tfsdk:"subfolderMatches"`
	PubsubTopicName  types.String `tfsdk:"pubsubTopicName"`
}

type stateNotificationConfigModel struct {
	PubsubTopicName types.String `tfsdk:"pubsubTopicName"`
}

type mqttConfigModel struct {
	MqttEnabledState types.String `tfsdk:"mqttEnabledState"`
}

type httpConfigModel struct {
	HttpEnabledState types.String `tfsdk:"httpEnabledState"`
}

func NewDeviceRegistriesDataSource() datasource.DataSource {
	return &deviceRegistriesDataSource{}
}

type deviceRegistriesDataSource struct{}

func (d *deviceRegistriesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_registries"
}

func (d *deviceRegistriesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

func (d *deviceRegistriesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
}