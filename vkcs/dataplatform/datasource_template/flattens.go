package datasource_template

import (
	"context"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"
)

func FlattenPodGroups(ctx context.Context, pods []templates.ClusterTemplatePodgroup) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	if pods == nil {
		return types.ListNull(PodGroupsValue{}.Type(ctx)), diags
	}

	values := make([]attr.Value, len(pods))
	for i, p := range pods {
		resourceObj, d := FlattenPodGroupResource(ctx, p.Resource)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(PodGroupsValue{}.Type(ctx)), diags
		}

		volumesMap, d := FlattenVolumes(ctx, p.Volumes)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(PodGroupsValue{}.Type(ctx)), diags
		}

		podGroupObj, d := PodGroupsValue{
			Name:     types.StringValue(p.Name),
			Count:    types.Int64Value(int64(p.Count)),
			Resource: resourceObj,
			Volumes:  volumesMap,
			state:    attr.ValueStateKnown,
		}.ToObjectValue(ctx)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListNull(PodGroupsValue{}.Type(ctx)), diags
		}

		values[i] = podGroupObj
	}

	result, d := types.ListValue(types.ObjectType{AttrTypes: PodGroupsValue{}.AttributeTypes(ctx)}, values)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListNull(PodGroupsValue{}.Type(ctx)), diags
	}

	return result, diags
}

func FlattenPodGroupResource(ctx context.Context, r templates.ClusterTemplatePodgroupResource) (basetypes.ObjectValue, diag.Diagnostics) {
	val := PodGroupsResourceValue{
		CpuRequest: types.StringValue(r.CpuRequest),
		CpuMargin:  types.NumberValue(big.NewFloat(r.CpuMargin)),
		RamRequest: types.StringValue(r.RamRequest),
		RamMargin:  types.NumberValue(big.NewFloat(r.RamMargin)),
		state:      attr.ValueStateKnown,
	}

	result, diags := val.ToObjectValue(ctx)
	if diags.HasError() {
		return types.ObjectNull(val.AttributeTypes(ctx)), diags
	}

	return result, diags
}

func FlattenVolumes(ctx context.Context, vols map[string]templates.ClusterTemplatePodgroupVolume) (basetypes.MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if vols == nil {
		return types.MapNull(PodGroupsVolumesValue{}.Type(ctx)), diags
	}

	values := make(map[string]attr.Value, len(vols))
	for k, v := range vols {
		values[k] = PodGroupsVolumesValue{
			Count:            types.Int64Value(int64(v.Count)),
			Storage:          types.StringValue(v.Storage),
			StorageClassName: types.StringValue(v.StorageClassName),
			state:            attr.ValueStateKnown,
		}
	}

	result, d := types.MapValue(PodGroupsVolumesValue{}.Type(ctx), values)
	diags.Append(d...)
	if diags.HasError() {
		return types.MapNull(PodGroupsVolumesValue{}.Type(ctx)), diags
	}

	return result, diags
}
