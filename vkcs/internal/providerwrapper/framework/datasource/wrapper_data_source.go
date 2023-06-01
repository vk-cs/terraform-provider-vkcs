package frameworkwrapper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	jsonschema "github.com/vk-cs/terraform-provider-vkcs/helpers/providerjson/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/providerwrapper/framework/datasource/customschema"
)

var (
	_ datasource.DataSource = &DataSourceWrapper{}
)

func NewDataSourceWrapper(dataSource datasource.DataSource, dataSourceJSON jsonschema.ResourceJSON) *DataSourceWrapper {
	return &DataSourceWrapper{
		dataSource:     dataSource,
		dataSourceJSON: dataSourceJSON,
	}
}

type DataSourceWrapper struct {
	dataSource     datasource.DataSource
	dataSourceJSON jsonschema.ResourceJSON
}

func (dw *DataSourceWrapper) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	dw.dataSource.Metadata(ctx, req, resp)
}

func (dw *DataSourceWrapper) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	dw.dataSource.Schema(ctx, req, resp)
	resp.Schema = customschema.CustomizeSchema(dw.dataSourceJSON, resp.Schema)
}

func (dw *DataSourceWrapper) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	dw.dataSource.Read(ctx, req, resp)
}

func (dw *DataSourceWrapper) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if ds, ok := dw.dataSource.(datasource.DataSourceWithConfigure); ok {
		ds.Configure(ctx, req, resp)
	}
}

func (dw *DataSourceWrapper) ConfigValidators(ctx context.Context) []datasource.ConfigValidator {
	if ds, ok := dw.dataSource.(datasource.DataSourceWithConfigValidators); ok {
		return ds.ConfigValidators(ctx)
	}
	return nil
}

func (dw *DataSourceWrapper) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if ds, ok := dw.dataSource.(datasource.DataSourceWithValidateConfig); ok {
		ds.ValidateConfig(ctx, req, resp)
	}
}
