package datasource_shielding_pops

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/cdn/v1/shieldingpop"
)

func FlattenShieldingePops(ctx context.Context, shieldingPops []shieldingpop.ShieldingPop) (types.List, diag.Diagnostics) {
	shieldingPopsVType := ShieldingPopsValue{}.Type(ctx)

	if len(shieldingPops) == 0 {
		return types.ListNull(shieldingPopsVType), nil
	}

	shieldingPopsV := make([]attr.Value, len(shieldingPops))
	for i, p := range shieldingPops {
		shieldingPopsV[i] = ShieldingPopsValue{
			City:       types.StringValue(p.City),
			Country:    types.StringValue(p.Country),
			Datacenter: types.StringValue(p.Datacenter),
			Id:         types.Int64Value(int64(p.ID)),
			state:      attr.ValueStateKnown,
		}
	}

	return types.ListValue(ShieldingPopsValue{}.Type(ctx), shieldingPopsV)
}
