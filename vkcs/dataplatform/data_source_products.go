package dataplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/datasource_products"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/products"

	"strconv"
	"time"
)

var (
	_ datasource.DataSource              = (*productsDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*productsDataSource)(nil)
)

func NewProductsDataSource() datasource.DataSource {
	return &productsDataSource{}
}

type productsDataSource struct {
	config clients.Config
}

func (d *productsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_products"
}

func (d *productsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_products.ProductsDataSourceSchema(ctx)
}

func (d *productsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *productsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_products.ProductsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := d.config.DataPlatformClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS dataplatform client", err.Error())
		return
	}

	tflog.Trace(ctx, "Calling Data Platform API to list products")

	productsResp, err := products.List(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Data Platform API to list products", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Data Platform API to list products", map[string]interface{}{"products": fmt.Sprintf("%#v", productsResp.Products)})

	resp.State.SetAttribute(ctx, path.Root("id"), strconv.FormatInt(time.Now().Unix(), 10))

	productsModel, diags := datasource_products.FlattenProducts(ctx, productsResp.Products)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Products = productsModel

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
