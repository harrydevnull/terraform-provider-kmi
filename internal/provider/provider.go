package provider

import (
	"context"
	"os"
	"terraform-provider-kmi/internal/kmi"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &kmiProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &kmiProvider{
			version: version,
		}
	}
}

// kmiProvider is the provider implementation.
type kmiProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Metadata returns the provider type name.
func (p *kmiProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "kmi"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *kmiProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"api_crt": schema.StringAttribute{
				Optional: true,
			},
			"akamai_ca": schema.StringAttribute{
				Optional: true,
			},
			"api_key_path": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"api_crt_path": schema.StringAttribute{
				Optional: true,
			},
			"akamai_ca_path": schema.StringAttribute{
				Optional: true,
			},
			"proxy_host": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

// kmiProviderModel maps provider schema data to a Go type.
type kmiProviderModel struct {
	Host         types.String `tfsdk:"host"`
	ApiKey       types.String `tfsdk:"api_key"`
	ApiCrt       types.String `tfsdk:"api_crt"`
	AkamaiCA     types.String `tfsdk:"akamai_ca"`
	ApiKeyPath   types.String `tfsdk:"api_key_path"`
	ApiCrtPath   types.String `tfsdk:"api_crt_path"`
	AkamaiCAPath types.String `tfsdk:"akamai_ca_path"`
	ProxyHost    types.String `tfsdk:"proxy_host"`
}

// Configure prepares a kmi API client for data sources and resources.
func (p *kmiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config kmiProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown KMI API Host",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_HOST environment variable.",
		)
	}

	if config.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Unknown KMI api_key",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI API key. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_API_KEY environment variable.",
		)
	}

	if config.ApiCrt.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_crt"),
			"Unknown KMI API CRT",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI API certificate. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_API_CRT environment variable.",
		)
	}

	if config.AkamaiCA.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("akamai_ca"),
			"Unknown KMI akamai_ca",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI_AKAMAI_CA. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_AKAMAI_CA environment variable.",
		)
	}

	if config.AkamaiCAPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("akamai_ca_path"),
			"Unknown KMI akamai_ca_path",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI_AKAMAI_CA_PATH. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_AKAMAI_CA_PATH environment variable.",
		)
	}

	if config.AkamaiCAPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_crt_path"),
			"Unknown KMI akamai_ca_path",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI_API_CRT_PATH. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_API_CRT_PATH environment variable.",
		)
	}

	if config.AkamaiCAPath.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key_path"),
			"Unknown KMI akamai_ca_path",
			"The provider cannot create the KMI API client as there is an unknown configuration value for the KMI_API_KEY_PATH. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the KMI_API_KEY_PATH environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("KMI_HOST")
	apikey := os.Getenv("KMI_API_KEY")
	apicrt := os.Getenv("KMI_API_CRT")
	akamaica := os.Getenv("KMI_AKAMAI_CA")
	apikeyPath := os.Getenv("KMI_API_KEY_PATH")
	apicrtPath := os.Getenv("KMI_API_CRT_PATH")
	akamaicaPath := os.Getenv("KMI_AKAMAI_CA_PATH")
	proxyHost := os.Getenv("KMI_PROXY_HOST")

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.ApiKey.IsNull() {
		apikey = config.ApiKey.ValueString()
	}

	if !config.ApiCrt.IsNull() {
		apicrt = config.ApiCrt.ValueString()
	}

	if !config.AkamaiCA.IsNull() {
		akamaica = config.AkamaiCA.ValueString()
	}

	if !config.AkamaiCAPath.IsNull() {
		akamaicaPath = config.AkamaiCAPath.ValueString()
	}

	if !config.ApiKeyPath.IsNull() {
		apikeyPath = config.ApiKeyPath.ValueString()
	}

	if !config.ApiCrtPath.IsNull() {
		apicrtPath = config.ApiCrtPath.ValueString()
	}

	if !config.ProxyHost.IsNull() {
		proxyHost = config.ProxyHost.ValueString()
	}

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing KMI API Host",
			"The provider cannot create the KMI API client as there is a missing or empty value for the kmi API host. "+
				"Set the host value in the configuration or use the KMI_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apikey == "" && apikeyPath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"Missing KMI API Key",
			"The provider cannot create the KMI API client as there is a missing or empty value for the kmi API host. "+
				"Set the host value in the configuration or use the KMI_API_KEY/KMI_API_KEY_PATH environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apicrt == "" && apicrtPath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("api_crt"),
			"Missing KMI API Certificate",
			"The provider cannot create the KMI API client as there is a missing or empty value for the kmi API host. "+
				"Set the host value in the configuration or use the KMI_API_CRT/KMI_API_CRT_PATH environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if akamaica == "" && akamaicaPath == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("akamai_ca"),
			"Missing KMI Certificate Authority",
			"The provider cannot create the KMI API client as there is a missing or empty value for the kmi API host. "+
				"Set the host value in the configuration or use the KMI_AKAMAI_CA/KMI_AKAMAI_CA_PATH environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "kmi_host", host)
	ctx = tflog.SetField(ctx, "kmi_key", apikey)
	ctx = tflog.SetField(ctx, "kmi_crt", apicrt)
	ctx = tflog.SetField(ctx, "kmi_ca", akamaica)
	ctx = tflog.SetField(ctx, "kmi_key_path", apikeyPath)
	ctx = tflog.SetField(ctx, "kmi_crt_path", apicrtPath)
	ctx = tflog.SetField(ctx, "kmi_ca_path", akamaicaPath)
	ctx = tflog.SetField(ctx, "proxy_host", proxyHost)

	tflog.Debug(ctx, "Creating KMI client")

	// setting a path variable takes a precedence over having it configured with the string
	if apikeyPath != "" {
		client, err := kmi.NewKMIRestClientPath(host, apikeyPath, apicrtPath, akamaicaPath, proxyHost)
		tflog.Info(ctx, "Configured KMI client", map[string]any{"success": true})

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create kmi API Client",
				"An unexpected error occurred when creating the kmi API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"kmi Client Error: "+err.Error(),
			)
			return
		}

		resp.DataSourceData = client
		resp.ResourceData = client
	} else {
		client, err := kmi.NewKMIRestClient(host, apikey, apicrt, akamaica, proxyHost)
		tflog.Info(ctx, "Configured KMI client", map[string]any{"success": true})

		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Create kmi API Client",
				"An unexpected error occurred when creating the kmi API client. "+
					"If the error is not clear, please contact the provider developers.\n\n"+
					"kmi Client Error: "+err.Error(),
			)
			return
		}

		resp.DataSourceData = client
		resp.ResourceData = client
	}
}

// DataSources defines the data sources implemented in the provider.
func (p *kmiProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAccountDataSource,
		NewCollectionsDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *kmiProvider) Resources(_ context.Context) []func() resource.Resource {

	return []func() resource.Resource{
		NewEngineResource,
		NewCollectionsResource,
		NewGroupsResource,
		NewDefinitionsResource,
		NewGroupsMembershipResource,
		NewTemplateResource,
		NewWorkloadResource,
	}
}
