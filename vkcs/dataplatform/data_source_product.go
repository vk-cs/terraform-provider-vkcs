package dataplatform

import (
	"context"
	"fmt"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/datasource_product"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/products"

	"strconv"
	"time"
)

var (
	_ datasource.DataSource              = (*productDataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*productDataSource)(nil)
)

func NewProductDataSource() datasource.DataSource {
	return &productDataSource{}
}

type productDataSource struct {
	config clients.Config
}

func (d *productDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_product"
}

func (d *productDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasource_product.ProductDataSourceSchema(ctx)
}

func (d *productDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *productDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data datasource_product.ProductModel

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

	var products []products.Product
	productName := data.ProductName.ValueString()
	productVersion := data.ProductVersion.ValueString()

	for _, product := range productsResp.Products {
		if product.ProductName == productName {
			if productVersion == "" || product.ProductVersion == productVersion {
				products = append(products, product)
			}
		}
	}

	if len(products) < 1 {
		resp.Diagnostics.AddError("Your query returned no results", "Please change your search criteria and try again.")
		return
	}

	if len(products) > 1 {
		resp.Diagnostics.AddError("Your query returned more than one result", "Please try a more specific search criteria")
		return
	}

	diags := data.UpdateFromProduct(ctx, &products[0])
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
