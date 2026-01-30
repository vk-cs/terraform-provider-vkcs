package datasource_kubernetes_volume_types_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
)

func (m *KubernetesVolumeTypesV2Model) UpdateFromVolumeType(ctx context.Context, volumeTypes clusters.ListVolumeTypes) (diags diag.Diagnostics) {
	// It's a synthetic identifier
	m.Id = types.StringValue("volume_types")

	if len(volumeTypes.StorageClasses) == 0 {
		m.VolumeTypes = types.SetNull(VolumeTypesValue{}.Type(ctx))
		return diags
	}

	attrTypes := VolumeTypesValue{}.AttributeTypes(ctx)
	elements := make([]attr.Value, 0, len(volumeTypes.StorageClasses))

	for _, sc := range volumeTypes.StorageClasses {
		// handle volume type azs
		zoneVals := make([]attr.Value, 0, len(sc.Zones))
		for _, z := range sc.Zones {
			zoneVals = append(zoneVals, types.StringValue(z))
		}
		zonesSet, d := types.SetValue(types.StringType, zoneVals)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		// handle volume type name
		objVal, d := types.ObjectValue(
			attrTypes,
			map[string]attr.Value{
				"name":  types.StringValue(sc.Name),
				"zones": zonesSet,
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		volumeTypeVal, d := VolumeTypesType{}.ValueFromObject(ctx, objVal)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		elements = append(elements, volumeTypeVal)
	}

	setVal, d := types.SetValue(
		VolumeTypesType{}.ValueType(ctx).Type(ctx),
		elements,
	)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.VolumeTypes = setVal
	return diags
}
