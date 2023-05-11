package clearblade

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &clearbladeProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &clearbladeProvider{}
}

// clearbladeProvider is the provider implementation.
type clearbladeProvider struct{}

// Metadata returns the provider type name.
func (p *clearbladeProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "clearblade"
}

// Schema defines the provider-level schema for configuration data.
func (p *clearbladeProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure prepares a ClearBlade IoT API client for data sources and resources.
func (p *clearbladeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// DataSources defines the data sources implemented in the provider.
func (p *clearbladeProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return nil
}

// Resources defines the resources implemented in the provider.
func (p *clearbladeProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}
