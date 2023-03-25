package provider

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure Device42Provider satisfies various provider interfaces.
var _ provider.Provider = &Device42Provider{}

// Device42Provider defines the provider implementation.
type Device42Provider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// Device42ProviderModel describes the provider data model.
type Device42ProviderModel struct {
	Hostname types.String `tfsdk:"hostname"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

func (p *Device42Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "device42"
	resp.Version = p.version
}

func (p *Device42Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: "hostname of the Device42 instance",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "username with access to the Device42 instance",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "password for the username with access",
				Required:    true,
				Sensitive:   true,
			},
			"insecure": schema.BoolAttribute{
				Description: "disable TLS verification - not recommended",
				Optional:    true,
			},
		},
	}
}

func (p *Device42Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	apiHostname := ""
	username := ""
	password := ""

	var data Device42ProviderModel

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if v := os.Getenv("D42_HOSTNAME"); v != "" {
		apiHostname = v
	} else {
		apiHostname = data.Hostname.ValueString()
	}

	if v := os.Getenv("D42_USERNAME"); v != "" {
		username = v
	} else {
		username = data.Username.ValueString()
	}

	if v := os.Getenv("D42_PASSWORD"); v != "" {
		password = v
	} else {
		password = data.Password.ValueString()
	}

	if apiHostname == "" {
		resp.Diagnostics.AddError(
			"Missing hostname",
			"No hostname was provided",
		)
	}
	if username == "" {
		resp.Diagnostics.AddError(
			"Missing username",
			"No username was provided",
		)
	}
	if password == "" {
		resp.Diagnostics.AddError(
			"Missing password",
			"No password was provided",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	client := http.DefaultClient
	if data.Insecure.ValueBool() == true {
		// disable TLS verification - not recommended
		resp.Diagnostics.AddWarning(
			"Insecure mode enabled",
			"TLS verification was disabled for this session.",
		)
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *Device42Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}

func (p *Device42Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		D42PasswordDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Device42Provider{
			version: version,
		}
	}
}
