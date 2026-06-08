package kubernetes

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	addons "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	dskubeaddonv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/datasource_kubernetes_addon_v2"
)

var (
	_ datasource.DataSource                   = (*kubernetesAddonV2DataSource)(nil)
	_ datasource.DataSourceWithConfigure      = (*kubernetesAddonV2DataSource)(nil)
	_ datasource.DataSourceWithValidateConfig = (*kubernetesAddonV2DataSource)(nil)
)

func NewKubernetesAddonV2DataSource() datasource.DataSource {
	return &kubernetesAddonV2DataSource{}
}

type kubernetesAddonV2DataSource struct {
	config clients.Config
}

func (d *kubernetesAddonV2DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_addon_v2"
}

func (d *kubernetesAddonV2DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = dskubeaddonv2.KubernetesAddonV2DataSourceSchema(ctx)
}

func (d *kubernetesAddonV2DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *kubernetesAddonV2DataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var data dskubeaddonv2.KubernetesAddonV2Model

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	idSet := !data.Id.IsNull()
	nameSet := !data.Name.IsNull()
	versionSet := !data.Version.IsNull()

	if nameSet != versionSet {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Attributes 'name' and 'version' must be specified together.",
		)
		return
	}

	namePair := nameSet && versionSet

	if idSet == namePair {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"Specify exactly one of: 'id' OR ('name' and 'version').",
		)
		return
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

func (d *kubernetesAddonV2DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dskubeaddonv2.KubernetesAddonV2Model

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
		resp.Diagnostics.AddError("Error creating API client for datasource vkcs_kubernetes_addon_v2", err.Error())
		return
	}

	// API call
	addon, diags := d.getAddon(client, &data)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(data.UpdateFromAddon(ctx, addon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (d *kubernetesAddonV2DataSource) getAddon(client *gophercloud.ServiceClient, data *dskubeaddonv2.KubernetesAddonV2Model) (addon *addons.AddonVersion, diags diag.Diagnostics) {
	if !data.Id.IsNull() {
		apiAddon, err := addons.GetAddonVersionByID(client, data.Id.ValueString()).Extract()
		if err != nil {
			diags.AddError("Error calling API to get addon by ID", err.Error())
			return
		}
		addon = &apiAddon
	} else {
		apiAddon, err := addons.GetAddonVersionByName(client, data.Name.ValueString(), data.Version.ValueString()).Extract()
		if err != nil {
			diags.AddError("Error calling API to get addon by name and version", err.Error())
			return
		}
		addon = &apiAddon
	}

	return
}
