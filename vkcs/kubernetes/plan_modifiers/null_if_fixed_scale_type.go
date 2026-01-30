package plan_modifiers

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func NullIfFixedScaleType() planmodifier.Int64 {
	return nullIfFixedScaleTypeModifier{}
}

type nullIfFixedScaleTypeModifier struct{}

func (m nullIfFixedScaleTypeModifier) Description(ctx context.Context) string {
	return m.MarkdownDescription(ctx)
}

func (m nullIfFixedScaleTypeModifier) MarkdownDescription(_ context.Context) string {
	return "If scale_type is 'fixed_scale', the plan value will be set to null."
}

func (m nullIfFixedScaleTypeModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	// Find scale_type in the config
	var scaleType types.String
	diags := req.Config.GetAttribute(ctx, req.Path.ParentPath().AtName("scale_type"), &scaleType)
	if diags.HasError() {
		return
	}

	// If scale_type is "fixed_scale", set the plan value to null
	if scaleType.ValueString() == "fixed_scale" {
		resp.PlanValue = types.Int64Null()
	}
}
