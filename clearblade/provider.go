package clearblade

import (
	"context"
	"net"
	"net/http"
	"os"
	"time"

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
	Credentials     types.String `tfsdk:"credentials"`
	CredentialsFile types.String `tfsdk:"credentials_file"`
	Project         types.String `tfsdk:"project"`
	Region          types.String `tfsdk:"region"`
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
			"credentials_file": schema.StringAttribute{
				Optional:  true,
				Sensitive: false,
			},
			"project": schema.StringAttribute{
				Optional: true,
			},
			"region": schema.StringAttribute{
				Optional: true,
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

	tflog.Debug(ctx, "Creating Clearblade IoT Core client")

	if config.Project.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Unknown Clearblade IoT Core Project",
			"The provider cannot create the Clearblade IoT Core client as there is an unknown configuration value for the Clearblade IoT Core Project. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLEARBLADE_CONFIGURATION environment variable.",
		)
	}

	if config.Region.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Unknown Clearblade IoT Core Region",
			"The provider cannot create the Clearblade IoT Core client as there is an unknown configuration value for the Clearblade IoT Core Region. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLEARBLADE_REGION environment variable.",
		)
	}

	if config.Credentials.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("credentials"),
			"Unknown Credentials",
			"The provider cannot create the Clearblade IoT Core client as there is an unknown configuration value for the Clearblade IoT Core API credentials.",
		)
	}

	if config.Credentials.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("credentials_file"),
			"Unknown Credentials File",
			"The provider cannot create the Clearblade IoT Core client as there is an unknown configuration value for the Clearblade IoT Core API credentials. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLEARBLADE_CONFIGURATION environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	os.Setenv("CLEARBLADE_PROJECT", config.Project.ValueString())
	os.Setenv("CLEARBLADE_REGION", config.Region.ValueString())

	var credentialsOption iot.ServiceOption
	switch {
	case os.Getenv("CLEARBLADE_CONFIGURATION") != "":
		credentialsOption = iot.WithFileCredentials()
	case !config.CredentialsFile.IsNull():
		os.Setenv("CLEARBLADE_CONFIGURATION", config.CredentialsFile.ValueString())
		credentialsOption = iot.WithFileCredentials()
	case !config.Credentials.IsNull():
		credentialsOption = iot.WithServiceAccountCredentials(config.Credentials.ValueString())
	default:
		resp.Diagnostics.AddError(
			"Missing Credentials",
			"The provider cannot create the Clearblade IoT Core client as there is no configuration value for the Clearblade IoT Core API credentials. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLEARBLADE_CONFIGURATION environment variable.",
		)
		return
	}

	// Create a new Clearblade IoT Core client using the configuration values
	client, err := iot.NewService(
		ctx,
		iot.WithHTTPClient(p.newHTTPClient()),
		credentialsOption,
	)
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
		NewDeviceRegistriesDataSource,
		NewDevicesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *clearbladeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
		NewDeviceRegistryResource,
	}
}

func (p *clearbladeProvider) newHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: (&net.Dialer{
				Timeout:   2 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			IdleConnTimeout: 60 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}
