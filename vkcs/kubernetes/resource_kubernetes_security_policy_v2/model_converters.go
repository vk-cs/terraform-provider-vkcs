package resource_kubernetes_security_policy_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/secpolicies"
)

func (m *KubernetesSecurityPolicyV2Model) UpdateFromClusterSecPolicy(ctx context.Context, data *secpolicies.ClusterSecPolicyResponse) (diags diag.Diagnostics) {
	if data == nil {
		return nil
	}

	m.Id = types.StringValue(data.ClusterSecPolicy.ID)
	m.ClusterId = types.StringValue(data.ClusterSecPolicy.ClusterID)
	m.Namespace = types.StringValue(data.ClusterSecPolicy.Namespace)
	m.Enabled = types.BoolValue(data.ClusterSecPolicy.Enabled)
	m.PolicySettings = types.StringValue(data.ClusterSecPolicy.PolicySettings)
	m.SecurityPolicyTemplateId = types.StringValue(data.ClusterSecPolicy.SecurityPolicyID)

	return
}
