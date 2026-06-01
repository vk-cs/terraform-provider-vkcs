package datasource_kubernetes_security_policy_template_v2

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
)

func (m *KubernetesSecurityPolicyTemplateV2Model) UpdateFromSecPolicyTemplate(apiSecPolicyTemplate clusters.SecPolicyTemplate) (diags diag.Diagnostics) {
	m.Id = types.StringValue(apiSecPolicyTemplate.ID)
	m.Name = types.StringValue(apiSecPolicyTemplate.Name)
	m.Description = types.StringValue(apiSecPolicyTemplate.Description)
	m.SettingsDescription = types.StringValue(apiSecPolicyTemplate.SettingsDescription)
	m.Version = types.StringValue(apiSecPolicyTemplate.Version)

	return
}
