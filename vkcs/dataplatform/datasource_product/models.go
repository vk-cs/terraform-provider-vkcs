package datasource_product

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/products"
)

func (m *ProductModel) UpdateFromProduct(ctx context.Context, product *products.Product) diag.Diagnostics {
	var diags diag.Diagnostics

	if product == nil {
		return diags
	}

	configs, d := FlattenConfigs(ctx, product.Configs)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.Configs = configs
	m.ProductName = types.StringValue(product.ProductName)
	m.ProductVersion = types.StringValue(product.ProductVersion)

	return diags
}
