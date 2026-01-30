package datasource_kubernetes_addons_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
)

func (m *KubernetesAddonsV2Model) UpdateFromAddonList(ctx context.Context, listAddons addons.AddonList) (diags diag.Diagnostics) {
	// It's a synthetic identifier
	m.Id = types.StringValue("addons")

	if len(listAddons.Addons) == 0 {
		m.Addons = types.SetNull(AddonsValue{}.Type(ctx))
		return diags
	}

	elements := make([]attr.Value, 0, len(listAddons.Addons))
	for _, a := range listAddons.Addons {
		// Build the list of versions for this addon
		versionsList := make([]attr.Value, 0, len(a.Versions))
		for _, v := range a.Versions {
			versionObj, versionDiags := NewVersionsValue(
				VersionsValue{}.AttributeTypes(ctx),
				map[string]attr.Value{
					"id":      types.StringValue(v.ID),
					"version": types.StringValue(v.Version),
				},
			)
			diags.Append(versionDiags...)
			if diags.HasError() {
				return diags
			}
			versionsList = append(versionsList, versionObj)
		}

		versionsListValue, listDiags := types.ListValue(
			VersionsValue{}.Type(ctx),
			versionsList,
		)
		diags.Append(listDiags...)
		if diags.HasError() {
			return diags
		}

		// Build the addon object
		addonValue, addonDiags := NewAddonsValue(
			AddonsValue{}.AttributeTypes(ctx),
			map[string]attr.Value{
				"id":       types.StringValue(a.ID),
				"name":     types.StringValue(a.Name),
				"versions": versionsListValue,
			},
		)
		diags.Append(addonDiags...)
		if diags.HasError() {
			return diags
		}

		elements = append(elements, addonValue)
	}

	// Build the set of addons
	addonsSet, setDiags := types.SetValue(
		AddonsValue{}.Type(ctx),
		elements,
	)
	diags.Append(setDiags...)
	if diags.HasError() {
		return diags
	}

	m.Addons = addonsSet
	return diags
}
