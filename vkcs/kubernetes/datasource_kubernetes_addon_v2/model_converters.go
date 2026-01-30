package datasource_kubernetes_addon_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
	mshared "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/models_shared"
)

func (m *KubernetesAddonV2Model) UpdateFromAddon(ctx context.Context, addon *addons.AddonVersion) (diags diag.Diagnostics) {
	m.AddonId = types.StringValue(addon.AddonID)
	m.Id = types.StringValue(addon.ID)
	m.Name = types.StringValue(addon.Name)
	m.Version = types.StringValue(addon.Version)

	kubeVersions, diags := mshared.FlattenStringSet(addon.SupportedKubeVersions)
	if diags.HasError() {
		return diags
	}
	m.SupportedKubeVersions = kubeVersions

	if addon.ValuesTemplate != nil {
		m.ValuesTemplate = types.StringValue(*addon.ValuesTemplate)
	}

	return diags
}
