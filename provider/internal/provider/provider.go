package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ provider.Provider = &TCPCheckProvider{}

type TCPCheckProvider struct {
	version string
}

type TCPCheckProviderModel struct{}

func (p *TCPCheckProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "tcpcheck"
	resp.Version = p.version
}

func (p *TCPCheckProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provider for TCP health checks",
	}
}

func (p *TCPCheckProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data TCPCheckProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func (p *TCPCheckProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTCPCheckResource,
	}
}

func (p *TCPCheckProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &TCPCheckProvider{
			version: version,
		}
	}
}
