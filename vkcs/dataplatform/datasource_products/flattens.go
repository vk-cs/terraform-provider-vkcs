package datasource_products

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/products"
)

func FlattenProducts(ctx context.Context, o []products.Product) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ProductsValue{}.Type(ctx)), nil
	}

	productsV := make([]attr.Value, len(o))
	for i, p := range o {
		productsV[i] = ProductsValue{
			ProductName:    types.StringValue(p.ProductName),
			ProductVersion: types.StringValue(p.ProductVersion),
			state:          attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ProductsValue{}.Type(ctx), productsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ProductsValue{}.Type(ctx)), diags
	}
	return result, nil
}
