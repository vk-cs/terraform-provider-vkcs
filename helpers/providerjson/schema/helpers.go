package schema

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	ds_schema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	rs_schema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	sdkschema "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func (p *ProviderJSON) SDKSchema() map[string]*sdkschema.Schema {
	return p.SDKProvider.Schema
}

func (p *ProviderJSON) SDKDataSourcesMap() map[string]*sdkschema.Resource {
	return p.SDKProvider.DataSourcesMap
}

func (p *ProviderJSON) SDKResourcesMap() map[string]*sdkschema.Resource {
	return p.SDKProvider.ResourcesMap
}

func (p *ProviderJSON) ResourcesMap() map[string]*rs_schema.Schema {
	resourcesMap := make(map[string]*rs_schema.Schema)
	ctx := context.Background()
	providerMeta := provider.MetadataResponse{}
	p.Provider.Metadata(ctx, provider.MetadataRequest{}, &providerMeta)

	for _, f := range p.Provider.Resources(ctx) {
		rs := f()
		rsMeta := resource.MetadataResponse{}
		rs.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: providerMeta.TypeName}, &rsMeta)
		rsSchema := resource.SchemaResponse{}
		rs.Schema(ctx, resource.SchemaRequest{}, &rsSchema)
		resourcesMap[rsMeta.TypeName] = &rsSchema.Schema
	}

	return resourcesMap
}

func (p *ProviderJSON) DataSourcesMap() map[string]*ds_schema.Schema {
	dataSourcesMap := make(map[string]*ds_schema.Schema)
	ctx := context.Background()
	providerMeta := provider.MetadataResponse{}
	p.Provider.Metadata(ctx, provider.MetadataRequest{}, &providerMeta)

	for _, f := range p.Provider.DataSources(ctx) {
		ds := f()
		dsMeta := datasource.MetadataResponse{}
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: providerMeta.TypeName}, &dsMeta)
		dsSchema := datasource.SchemaResponse{}
		ds.Schema(ctx, datasource.SchemaRequest{}, &dsSchema)
		dataSourcesMap[dsMeta.TypeName] = &dsSchema.Schema
	}

	return dataSourcesMap
}

func NodeIsBlock(input SchemaJSON) bool {
	if input.Type == SchemaTypeList || input.Type == SchemaTypeSet {
		if _, ok := input.Elem.(ResourceJSON); ok {
			return true
		}
	}

	return false
}
