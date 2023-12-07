package planmodifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
)

func GetRegionPlanModifier(resp *resource.SchemaResponse) stringplanmodifier.RequiresReplaceIfFunc {
	return func(ctx context.Context, sr planmodifier.StringRequest, rrifr *stringplanmodifier.RequiresReplaceIfFuncResponse) {
		var configValue, stateValue types.String
		resp.Diagnostics.Append(sr.Config.GetAttribute(ctx, sr.Path, &configValue)...)
		resp.Diagnostics.Append(sr.State.GetAttribute(ctx, sr.Path, &stateValue)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !configValue.IsNull() && !configValue.Equal(stateValue) {
			rrifr.RequiresReplace = true
		}
	}
}
