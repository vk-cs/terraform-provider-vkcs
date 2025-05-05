package dataplatform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/datasource_products"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

var (
	_ datasource.DataSource = (*productsDataSource)(nil)
)

func NewProductsDataSource() datasource.DataSource {
	return &productsDataSource{}
}

type productsDataSource struct {
	config clients.Config
}

func (r *productsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_products"
}

func (r *productsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_products.ProductsDataSourceSchema(ctx)
}

func (r *productsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *productsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

}
