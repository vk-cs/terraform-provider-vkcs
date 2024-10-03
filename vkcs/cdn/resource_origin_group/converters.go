package resource_origin_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/origingroups"
)

func (m *OriginGroupModel) UpdateFromOriginGroup(ctx context.Context, originGroup *origingroups.OriginGroup) diag.Diagnostics {
	var diags diag.Diagnostics

	if originGroup == nil {
		return diags
	}

	m.Id = types.Int64Value(int64(originGroup.ID))
	m.Name = types.StringValue(originGroup.Name)

	origins, d := flattenOrigins(ctx, originGroup.Origins)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.Origins = origins
	m.UseNext = types.BoolValue(originGroup.UseNext)

	return diags
}

func ExpandOrigins(ctx context.Context, origins types.List) ([]origingroups.Origin, diag.Diagnostics) {
	if origins.IsUnknown() || origins.IsNull() {
		return nil, nil
	}

	originsV := make([]OriginsValue, 0, len(origins.Elements()))
	diags := origins.ElementsAs(ctx, &originsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]origingroups.Origin, len(originsV))
	for i, o := range originsV {
		result[i] = origingroups.Origin{
			Backup:  o.Backup.ValueBool(),
			Enabled: o.Enabled.ValueBool(),
			Source:  o.Source.ValueString(),
		}
	}

	return result, nil
}

func flattenOrigins(ctx context.Context, origins []origingroups.Origin) (types.List, diag.Diagnostics) {
	originsVType := OriginsValue{}.Type(ctx)

	if len(origins) == 0 {
		return types.ListNull(originsVType), nil
	}

	originsV := make([]attr.Value, len(origins))
	for i, o := range origins {
		originsV[i] = OriginsValue{
			Backup:  types.BoolValue(o.Backup),
			Enabled: types.BoolValue(o.Enabled),
			Source:  types.StringValue(o.Source),
			state:   attr.ValueStateKnown,
		}
	}

	return types.ListValue(OriginsValue{}.Type(ctx), originsV)
}
