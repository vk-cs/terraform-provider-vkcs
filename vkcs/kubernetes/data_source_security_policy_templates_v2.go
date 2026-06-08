package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	dskubesecpolicytmpltsv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_security_policy_templates_v2"
)

var (
	_ datasource.DataSource = (*kubernetesSecurityPolicyTemplatesV2DataSource)(nil)
)

func NewKubernetesSecurityPolicyTemplatesV2DataSource() datasource.DataSource {
	return &kubernetesSecurityPolicyTemplatesV2DataSource{}
}

type kubernetesSecurityPolicyTemplatesV2DataSource struct {
	config clients.Config
}

func (d *kubernetesSecurityPolicyTemplatesV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_security_policy_templates_v2"
}

func (d *kubernetesSecurityPolicyTemplatesV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubesecpolicytmpltsv2.KubernetesSecurityPolicyTemplatesV2DataSourceSchema(ctx)
}

func (d *kubernetesSecurityPolicyTemplatesV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesSecurityPolicyTemplatesV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubesecpolicytmpltsv2.KubernetesSecurityPolicyTemplatesV2Model

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the region in which to obtain the Managed K8S client
	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	// Init Managed K8S client
	client, err := d.config.ManagedK8SClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_security_policy_templates_v2", err.Error())
		return
	}

	// API call
	listSecPolicyTemplates, err := clusters.GetListSecPolicyTemplates(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling API to list Kubernetes security policy templates", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromListSecPolicyTemplates(ctx, listSecPolicyTemplates)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
