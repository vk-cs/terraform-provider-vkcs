package resource_kubernetes_security_policy_v2

import (
	"context"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/secpolicies"
)

func ToCreateOpts(ctx context.Context, model KubernetesSecurityPolicyV2Model) secpolicies.CreateOpts {
	return secpolicies.CreateOpts{
		ClusterID:        model.ClusterId.ValueString(),
		SecurityPolicyID: model.SecurityPolicyTemplateId.ValueString(),
		PolicySettings:   model.PolicySettings.ValueString(),
		Namespace:        model.Namespace.ValueString(),
		Enabled:          model.Enabled.ValueBool(),
	}
}
