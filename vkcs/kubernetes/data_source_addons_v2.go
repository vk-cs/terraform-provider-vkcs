package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	addons "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
	dskubeaddonsv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_addons_v2"
)

var (
	_ datasource.DataSource              = (*kubernetesAddonsV2DataSource)(nil)
	_ datasource.DataSourceWithConfigure = (*kubernetesAddonsV2DataSource)(nil)
)

func NewKubernetesAddonsV2DataSource() datasource.DataSource {
	return &kubernetesAddonsV2DataSource{}
}

type kubernetesAddonsV2DataSource struct {
	config clients.Config
}

func (d *kubernetesAddonsV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_addons_v2"
}

func (d *kubernetesAddonsV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubeaddonsv2.KubernetesAddonsV2DataSourceSchema(ctx)
}

func (d *kubernetesAddonsV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesAddonsV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubeaddonsv2.KubernetesAddonsV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_addons_v2", err.Error())
		return
	}

	// API call
	addonList, err := addons.ListAddons(client).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling API to get list of available addons", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromAddonList(ctx, addonList)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
