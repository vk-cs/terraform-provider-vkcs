package resource_kubernetes_cluster_addon_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
)

func (m *KubernetesClusterAddonV2Model) UpdateFromClusterAddon(ctx context.Context, data *addons.ClusterAddon) (diags diag.Diagnostics) {
	if data == nil {
		return nil
	}

	m.AddonId = types.StringValue(data.AddonID)
	m.AddonName = types.StringValue(data.AddonName)
	m.AddonVersionId = types.StringValue(data.AddonVersionID)
	m.ClusterId = types.StringValue(data.ClusterID)
	m.Id = types.StringValue(data.ID)
	m.Namespace = types.StringValue(data.Namespace)
	m.Values = types.StringValue(data.Values)
	m.Status = types.StringValue(data.Status)
	m.CreatedAt = types.StringValue(data.CreatedAt)
	m.UpdatedAt = types.StringValue(data.UpdatedAt)

	return
}
