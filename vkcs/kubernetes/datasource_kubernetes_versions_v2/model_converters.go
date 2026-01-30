package datasource_kubernetes_versions_v2

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	mshared "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/models_shared"
)

func (m *KubernetesVersionsV2Model) UpdateFromClusterVersion(k8sVersions clusters.ListK8SVersion) (diags diag.Diagnostics) {
	// It's a synthetic identifier
	m.Id = types.StringValue("kubernetes_versions")

	versions := make([]string, 0, len(k8sVersions.Versions))
	for _, v := range k8sVersions.Versions {
		versions = append(versions, v.Version)
	}

	clusterAZs, diags := mshared.FlattenStringSet(versions)
	if diags.HasError() {
		return diags
	}

	m.K8sVersions = clusterAZs
	return
}
