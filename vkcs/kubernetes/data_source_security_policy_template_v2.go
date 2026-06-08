package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	dskubesecpolicytmpltv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_security_policy_template_v2"
)

var (
	_ datasource.DataSource                   = (*kubernetesSecurityPolicyTemplateV2DataSource)(nil)
	_ datasource.DataSourceWithValidateConfig = (*kubernetesSecurityPolicyTemplateV2DataSource)(nil)
)

func NewKubernetesSecurityPolicyTemplateV2DataSource() datasource.DataSource {
	return &kubernetesSecurityPolicyTemplateV2DataSource{}
}

type kubernetesSecurityPolicyTemplateV2DataSource struct {
	config clients.Config
}

func (d *kubernetesSecurityPolicyTemplateV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_security_policy_template_v2"
}

func (d *kubernetesSecurityPolicyTemplateV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubesecpolicytmpltv2.KubernetesSecurityPolicyTemplateV2DataSourceSchema(ctx)
}

func (d *kubernetesSecurityPolicyTemplateV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesSecurityPolicyTemplateV2DataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data dskubesecpolicytmpltv2.KubernetesSecurityPolicyTemplateV2Model

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	idSet := !data.Id.IsNull()
	nameSet := !data.Name.IsNull()
	versionSet := !data.Version.IsNull()

	if idSet && (nameSet || versionSet) {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Specify either 'id' OR ('name' and 'version'), not both.",
		)
		return
	}

	if !idSet && (!nameSet || !versionSet) {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"You must specify either 'id' OR both 'name' and 'version'.",
		)
	}

	if util.IsKnownValue(data.Id) && data.Id.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'id' must be not empty.",
		)
	}

	if util.IsKnownValue(data.Name) && data.Name.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'name' must be not empty.",
		)
	}

	if util.IsKnownValue(data.Version) && data.Version.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attribute 'version' must be not empty.",
		)
	}
}

func (d *kubernetesSecurityPolicyTemplateV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubesecpolicytmpltv2.KubernetesSecurityPolicyTemplateV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_security_policy_template_v2", err.Error())
		return
	}

	// Fetch security policy template by ID or by name + version
	var secPolicyTemplate clusters.SecPolicyTemplate

	if !data.Id.IsNull() {
		// Lookup by ID
		secPolicyTemplate, err = clusters.GetSecPolicyTemplateByID(client, data.Id.ValueString()).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error calling API to get Kubernetes security policy template", err.Error())
			return
		}
	} else if !data.Name.IsNull() && !data.Version.IsNull() {
		// Lookup by name + version
		secPolicyTemplate, err = clusters.GetSecPolicyTemplateByNameAndVersion(client, data.Name.ValueString(), data.Version.ValueString()).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error calling API to get Kubernetes security policy template", err.Error())
			return
		}
	}

	resp.Diagnostics.Append(data.UpdateFromSecPolicyTemplate(secPolicyTemplate)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
