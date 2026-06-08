package resource_kubernetes_cluster_addon_v2

import (
	"context"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
)

func ToCreateOpts(ctx context.Context, model KubernetesClusterAddonV2Model) addons.CreateOpts {
	return addons.CreateOpts{
		ClusterID:      model.ClusterId.ValueString(),
		AddonID:        model.AddonId.ValueString(),
		AddonVersionID: model.AddonVersionId.ValueString(),
		Namespace:      model.Namespace.ValueString(),
		Values:         model.Values.ValueString(),
		AddonName:      model.AddonName.ValueString(),
	}
}

func ToUpdateOpts(ctx context.Context, model KubernetesClusterAddonV2Model) addons.UpdateOpts {
	return addons.UpdateOpts{
		ClusterID:      model.ClusterId.ValueString(),
		ClusterAddonID: model.Id.ValueString(),
		AddonVersionID: model.AddonVersionId.ValueString(),
		Values:         model.Values.ValueString(),
	}
}
