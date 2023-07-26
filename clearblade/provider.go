package clearblade

import (
	"context"
	"os"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &clearbladeProvider{}
)

/* // New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &clearbladeProvider{
			version: version,
		}
	}
} */

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &clearbladeProvider{}
}

// clearbladeProviderModel maps provider schema data to a Go type.
type clearbladeProviderModel struct {
	Credentials types.String `tfsdk:"credentials"`
	/* Project types.String `tfsdk:"project"`
	Region  types.String `tfsdk:"region"`
	Registry    types.String `tfsdk:"registry"`
	Credentials types.String `tfsdk:"credentials"` */
}

// clearbladeProvider is the provider implementation.
type clearbladeProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *clearbladeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clearblade"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *clearbladeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

// Configure prepares a ClearBlade IoT Core API client for data sources and resources.
func (p *clearbladeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring ClearBlade IoT Core client")

	// Retrieve provider data from configuration
	var config clearbladeProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "clearblade_credentials", config.Credentials.ValueString())

	tflog.Debug(ctx, "Creating Clearblade IoT Core client")

	if config.Credentials.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("credentials"),
			"Unknown Credentials",
			"The provider cannot create the Clearblade IoT Core client as there is an unknown configuration value for the Clearblade IoT Core API credentials. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLEARBLADE_CONFIGURATION environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	os.Setenv("CLEARBLADE_CONFIGURATION", config.Credentials.ValueString())

	// Create a new Clearblade IoT Core client using the configuration values
	client, err := iot.NewService(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Clearblade IoT Core API Client",
			"An unexpected error occurred when creating the Clearblade IoT Core API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Clearblade IoT Core Client Error: "+err.Error(),
		)
		return
	}

	// Make the Clearblade IoT Core client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Clearblade IoT Core client", map[string]any{"success": true})
}

// DataSources defines the data sources implemented in the provider.
func (p *clearbladeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		//NewDeviceRegistriesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *clearbladeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		//NewDeviceResource,
		NewDeviceRegistryResource,
	}
}
