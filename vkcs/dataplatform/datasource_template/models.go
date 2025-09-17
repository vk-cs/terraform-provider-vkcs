package datasource_template

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"
)

func (m *TemplateModel) UpdateFromTemplate(ctx context.Context, template *templates.ClusterTemplate) diag.Diagnostics {
	var diags diag.Diagnostics

	if template == nil {
		return diags
	}

	podGroups, d := FlattenPodGroups(ctx, template.PodGroups)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	m.Id = types.StringValue(template.ID)
	m.Name = types.StringValue(template.Name)
	m.ProductName = types.StringValue(template.ProductName)
	m.ProductVersion = types.StringValue(template.ProductVersion)
	m.PodGroups = podGroups

	return diags
}
