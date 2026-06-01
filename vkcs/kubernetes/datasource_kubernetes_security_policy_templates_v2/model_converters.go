package datasource_kubernetes_security_policy_templates_v2

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
)

func (m *KubernetesSecurityPolicyTemplatesV2Model) UpdateFromListSecPolicyTemplates(ctx context.Context, apiListSecPolicyTemplates clusters.ListSecPolicyTemplates) (diags diag.Diagnostics) {
	// It's a synthetic identifier
	m.Id = types.StringValue("policy_templates")

	if len(apiListSecPolicyTemplates.SecPolicyTemplates) == 0 {
		m.SecurityPolicies = types.SetNull(SecurityPoliciesValue{}.Type(ctx))
		return diags
	}

	attrTypes := SecurityPoliciesValue{}.AttributeTypes(ctx)
	elements := make([]attr.Value, 0, len(apiListSecPolicyTemplates.SecPolicyTemplates))

	for _, spt := range apiListSecPolicyTemplates.SecPolicyTemplates {
		objVal, d := types.ObjectValue(
			attrTypes,
			map[string]attr.Value{
				"id":                   types.StringValue(spt.ID),
				"name":                 types.StringValue(spt.Name),
				"description":          types.StringValue(spt.Description),
				"settings_description": types.StringValue(spt.SettingsDescription),
				"version":              types.StringValue(spt.Version),
			},
		)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		secPolicyTemplateTypeVal, d := SecurityPoliciesType{}.ValueFromObject(ctx, objVal)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		elements = append(elements, secPolicyTemplateTypeVal)
	}

	setVal, d := types.SetValue(
		SecurityPoliciesType{}.ValueType(ctx).Type(ctx),
		elements,
	)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.SecurityPolicies = setVal
	return
}
