package datasource_kubernetes_cluster_addon_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
)

func (m *KubernetesClusterAddonV2Model) UpdateFromAddon(ctx context.Context, addon *addons.ClusterAddon) (diags diag.Diagnostics) {
	if addon == nil {
		return
	}

	m.AddonId = types.StringValue(addon.AddonID)
	m.AddonName = types.StringValue(addon.AddonName)
	m.BaseAddonName = types.StringValue(addon.BaseAddonName)
	m.AddonVersionId = types.StringValue(addon.AddonVersionID)
	m.ClusterId = types.StringValue(addon.ClusterID)
	m.Id = types.StringValue(addon.ID)
	m.Namespace = types.StringValue(addon.Namespace)
	m.Values = types.StringValue(addon.Values)
	m.Status = types.StringValue(addon.Status)
	m.CreatedAt = types.StringValue(addon.CreatedAt)
	m.UpdatedAt = types.StringValue(addon.UpdatedAt)

	return diags
}
