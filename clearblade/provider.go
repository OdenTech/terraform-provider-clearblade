package clearblade

import (
	"context"

	"github.com/clearblade/go-iot"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
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
		/* Attributes: map[string]schema.Attribute{
			"project": schema.StringAttribute{
				Required: true,
			},
			"registry": schema.StringAttribute{
				Required: true,
			},
			"region": schema.StringAttribute{
				Required: true,
			},
			"credentials": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		}, */
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

	/* ctx = tflog.SetField(ctx, "clearblade_project", config.Project.ValueString())
	ctx = tflog.SetField(ctx, "clearblade_registry", config.Registry.ValueString())
	ctx = tflog.SetField(ctx, "clearblade_region", config.Region.ValueString())

	tflog.Debug(ctx, "Creating Clearblade IoT Core client") */

	// If practitioner provided a configuration value for any of the
	// attributes, it must be a known value.

	/* if config.Project.IsUnknown() {
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

	if config.Registry.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("registry"),
			"Unknown Clearblade IoT Core Registry",
			"The provider cannot create the Clearblade IoT Core client as there is an unknown configuration value for the Clearblade IoT Core Registry. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the CLEARBLADE_REGISTRY environment variable.",
		)
	}

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
	} */

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	/* project := os.Getenv("CLEARBLADE_PROJECT")
	registry := os.Getenv("CLEARBLADE_REGISTRY")
	region := os.Getenv("CLEARBLADE_REGION")
	credentials := os.Getenv("CLEARBLADE_CONFIGURATION")

	if !config.Project.IsNull() {
		project = config.Project.ValueString()
	}

	if !config.Registry.IsNull() {
		registry = config.Registry.ValueString()
	}

	if !config.Region.IsNull() {
		region = config.Region.ValueString()
	}

	if !config.Credentials.IsNull() {
		credentials = config.Credentials.ValueString()
	} */

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	/* if project == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("project"),
			"Missing Clearblade IoT Core Project",
			"The provider cannot create the Clearblade IoT Core client as there is a missing or empty value for the Clearblade IoT Core API project. "+
				"Set the host value in the configuration or use the CLEARBLADE_PROJECT environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if registry == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("registry"),
			"Missing Clearblade IoT Core Registry",
			"The provider cannot create the Clearblade IoT Core client as there is a missing or empty value for the Clearblade IoT Core registry. "+
				"Set the host value in the configuration or use the CLEARBLADE_REGISTRY environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if region == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("region"),
			"Missing Clearblade IoT Core Region",
			"The provider cannot create the Clearblade IoT Core client as there is a missing or empty value for the Clearblade IoT Core Region. "+
				"Set the host value in the configuration or use the CLEARBLADE_REGION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if credentials == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("credentials"),
			"Missing Clearblade IoT Core Credentials",
			"The provider cannot create the Clearblade IoT Core API client as there is a missing or empty value for the Clearblade IoT Core API credentials. "+
				"Set the host value in the configuration or use the CLEARBLADE_CONFIGURATION environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	} */

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
		NewDeviceRegistriesDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *clearbladeProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewDeviceResource,
	}
}
