package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	dskubevolumetypesv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_volume_types_v2"
)

var (
	_ datasource.DataSource = (*kubernetesVolumeTypesV2DataSource)(nil)
)

func NewKubernetesVolumeTypesV2DataSource() datasource.DataSource {
	return &kubernetesVolumeTypesV2DataSource{}
}

type kubernetesVolumeTypesV2DataSource struct {
	config clients.Config
}

func (d *kubernetesVolumeTypesV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_volume_types_v2"
}

func (d *kubernetesVolumeTypesV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubevolumetypesv2.KubernetesVolumeTypesV2DataSourceSchema(ctx)
}

func (d *kubernetesVolumeTypesV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesVolumeTypesV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubevolumetypesv2.KubernetesVolumeTypesV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_volume_types_v2", err.Error())
		return
	}

	// API call
	volumeTypes, err := clusters.GetVolumeTypes(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling API to get root volume types", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromVolumeType(ctx, volumeTypes)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
