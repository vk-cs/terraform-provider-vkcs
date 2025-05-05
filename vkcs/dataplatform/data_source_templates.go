package dataplatform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/datasource_templates"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

var (
	_ datasource.DataSource = (*templatesDataSource)(nil)
)

func NewTemplatesDataSource() datasource.DataSource {
	return &templatesDataSource{}
}

type templatesDataSource struct {
	config clients.Config
}

func (r *templatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_templates"
}

func (r *templatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_templates.TemplatesDataSourceSchema(ctx)
}

func (r *templatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *templatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

}
