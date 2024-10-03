package datasource_origin_group

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/origingroups"
)

func FlattenOrigins(ctx context.Context, origins []origingroups.Origin) (types.List, diag.Diagnostics) {
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
