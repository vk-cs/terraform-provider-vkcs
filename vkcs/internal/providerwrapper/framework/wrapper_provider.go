package frameworkwrapper

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper"
	dswrapper "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper/framework/datasource"
	rswrapper "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper/framework/resource"
)

var (
	_ provider.Provider = &ProviderWrapper{}
)

func NewProviderWrapper(provider provider.Provider) *ProviderWrapper {
	var wrapperSchemaJSON jsonschema.ProviderWrapper
	err := json.Unmarshal([]byte(providerwrapper.ProviderSchemaJSON), &wrapperSchemaJSON)
	if err != nil {
		panic(err)
	}

	return &ProviderWrapper{
		provider:           provider,
		providerSchemaJSON: wrapperSchemaJSON.ProviderSchema,
	}
}

type ProviderWrapper struct {
	provider           provider.Provider
	providerSchemaJSON *jsonschema.ProviderSchemaJSON
}

func (pw *ProviderWrapper) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	pw.provider.Metadata(ctx, req, resp)
}

// Schema defines the provider-level schema for configuration data.
func (pw *ProviderWrapper) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	pw.provider.Schema(ctx, req, resp)
}

// Configure prepares a HashiCups API client for data sources and resources.
func (pw *ProviderWrapper) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	pw.provider.Configure(ctx, req, resp)
}

// DataSources defines the data sources implemented in the provider.
func (pw *ProviderWrapper) DataSources(ctx context.Context) []func() datasource.DataSource {
	dataSourcesFactories := pw.provider.DataSources(ctx)
	wrappedDataSourcesFactories := make([]func() datasource.DataSource, len(dataSourcesFactories))
	for i, f := range dataSourcesFactories {
		wrappedDataSourcesFactories[i] = pw.dataSourceWrapperFactory(ctx, f)
	}
	return wrappedDataSourcesFactories
}

// Resources defines the resources implemented in the provider.
func (pw ProviderWrapper) Resources(ctx context.Context) []func() resource.Resource {
	resourcesFactories := pw.provider.Resources(ctx)
	wrappedResourcesFactories := make([]func() resource.Resource, len(resourcesFactories))
	for i, f := range resourcesFactories {
		wrappedResourcesFactories[i] = pw.resourceWrapperFactory(ctx, f)
	}
	return wrappedResourcesFactories
}

func (pw *ProviderWrapper) dataSourceWrapperFactory(ctx context.Context, f func() datasource.DataSource) func() datasource.DataSource {
	return func() datasource.DataSource {
		d := f()
		dMeta := datasource.MetadataResponse{}
		d.Metadata(ctx, datasource.MetadataRequest{}, &dMeta)
		dsSchemaJSON := pw.providerSchemaJSON.DataSourcesMap[dMeta.TypeName]
		return dswrapper.NewDataSourceWrapper(d, dsSchemaJSON)
	}
}

func (pw *ProviderWrapper) resourceWrapperFactory(ctx context.Context, f func() resource.Resource) func() resource.Resource {
	return func() resource.Resource {
		r := f()
		rMeta := resource.MetadataResponse{}
		r.Metadata(ctx, resource.MetadataRequest{}, &rMeta)
		rsSchemaJSON := pw.providerSchemaJSON.ResourcesMap[rMeta.TypeName]
		return rswrapper.NewResourceWrapper(r, rsSchemaJSON)
	}
}
