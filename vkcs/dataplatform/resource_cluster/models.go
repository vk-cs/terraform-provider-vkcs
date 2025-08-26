package resource_cluster

import (
	"context"

	"github.com/google/uuid"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"
)

const floatingIPAutoMode = "auto"

var floatingIPAutoID = uuid.Nil.String()

func (m *ClusterModel) UpdateState(ctx context.Context, cluster *clusters.Cluster, state *tfsdk.State, oldSettings basetypes.ListValue) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(m.UpdateFromCluster(ctx, cluster)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(state.Set(ctx, m)...)

	diags.Append(UpdateClusterConfigs(ctx, cluster.Configs, state)...)
	diags.Append(UpdateClusterPodGroups(ctx, cluster.PodGroups, state)...)
	diags.Append(UpdateClusterConfigsSettings(ctx, oldSettings, cluster.Configs.Settings, state)...)

	return diags
}

func (m *ClusterModel) UpdateFromCluster(ctx context.Context, cluster *clusters.Cluster) diag.Diagnostics {
	var diags diag.Diagnostics

	if cluster == nil {
		return diags
	}

	m.Id = types.StringValue(cluster.ID)
	m.CreatedAt = types.StringValue(cluster.CreatedAt)
	m.Name = types.StringValue(cluster.Name)
	m.Description = types.StringValue(cluster.Description)
	m.ProductVersion = types.StringValue(cluster.ProductVersion)
	m.AvailabilityZone = types.StringValue(cluster.AvailabilityZone)
	m.Multiaz = types.BoolValue(cluster.MultiAZ)
	m.NetworkId = types.StringValue(cluster.NetworkID)
	m.SubnetId = types.StringValue(cluster.SubnetID)
	m.ProductName = types.StringValue(cluster.ProductName)
	m.ProductType = types.StringValue(cluster.ProductType)
	m.ProductVersion = types.StringValue(cluster.ProductVersion)
	m.StackId = types.StringValue(cluster.StackID)

	if cluster.FloatingIPPool == floatingIPAutoID {
		m.FloatingIpPool = types.StringValue(floatingIPAutoMode)
	} else {
		m.FloatingIpPool = types.StringValue(cluster.FloatingIPPool)
	}

	return diags
}

func UpdateClusterConfigs(ctx context.Context, config *clusters.ClusterConfig, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	if config == nil {
		return nil
	}

	maintenance, d := FlattenClusterConfigsMaintenance(ctx, config.Maintenance)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	d = state.SetAttribute(ctx, path.Root("configs").AtName("maintenance"), maintenance)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	d = UpdateClusterConfigsWarehouses(ctx, config.Warehouses, path.Root("configs").AtName("warehouses"), state)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	return nil
}

func UpdateClusterConfigsWarehouses(ctx context.Context, warehouses []clusters.ClusterConfigWarehouse, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	if warehouses == nil {
		return nil
	}

	for i, w := range warehouses {
		d := UpdateClusterConfigsWarehousesConnections(ctx, i, w.Connections, path.AtListIndex(i).AtName("connections"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.AtListIndex(i).AtName("id"), w.ID)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.AtListIndex(i).AtName("name"), w.Name)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}
	return nil
}

func UpdateClusterConfigsWarehousesConnections(ctx context.Context, i int, connections []clusters.ClusterConfigWarehouseConnection, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	if connections == nil {
		return nil
	}

	for j, c := range connections {
		d := state.SetAttribute(ctx, path.AtListIndex(j).AtName("created_at"), c.CreatedAt)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.AtListIndex(j).AtName("id"), c.ID)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.AtListIndex(j).AtName("name"), c.Name)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.AtListIndex(j).AtName("plug"), c.Plug)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}
	return nil
}

func UpdateClusterConfigsSettings(ctx context.Context, oldSettings basetypes.ListValue, o []clusters.ClusterConfigSetting, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics
	settingsV := make([]ConfigsSettingsValue, 0, len(oldSettings.Elements()))
	if len(oldSettings.Elements()) > 0 {
		d := oldSettings.ElementsAs(ctx, &settingsV, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	for _, s := range o {
		updated := false
		for i, setting := range settingsV {
			if setting.Alias.ValueString() == s.Alias {
				settingsV[i].Value = types.StringValue(s.Value)
				updated = true
			}
		}
		if !updated {
			settingsV = append(settingsV, ConfigsSettingsValue{
				Alias: types.StringValue(s.Alias),
				Value: types.StringValue(s.Value),
				state: attr.ValueStateKnown,
			})
		}
	}
	if len(settingsV) == 0 {
		d := state.SetAttribute(ctx, path.Root("configs").AtName("settings"), types.ListNull(ConfigsSettingsValue{}.Type(ctx)))
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	} else {
		d := state.SetAttribute(ctx, path.Root("configs").AtName("settings"), settingsV)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}
	return nil
}

func UpdateClusterPodGroups(ctx context.Context, o []clusters.ClusterPodGroup, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	if o == nil {
		return nil
	}

	for _, p := range o {
		var statePodGroups []PodGroupsValue
		d := state.GetAttribute(ctx, path.Root("pod_groups"), &statePodGroups)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		for i, group := range statePodGroups {
			if group.Name.ValueString() == p.Name {
				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("alias"), p.Alias)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("availability_zone"), p.AvailabilityZone)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("count"), p.Count)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("floating_ip_pool"), p.FloatingIPPool)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("id"), p.ID)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				volumes, d := FlattenClusterPodGroupsVolumes(ctx, p.Volumes)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("volumes"), volumes)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				resource, d := FlattenClusterPodGroupsResource(ctx, p.Resource)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}

				d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("resource"), resource)
				diags.Append(d...)
				if diags.HasError() {
					return diags
				}
			}
		}
	}
	return nil
}

func (m *ClusterModel) GetClusterTemplate(client *gophercloud.ServiceClient) (*templates.ClusterTemplate, diag.Diagnostics) {
	var diags diag.Diagnostics
	templatesResp, err := templates.Get(client).Extract()
	if err != nil {
		diags.AddError("unknown cluster template", "could not retrieve cluster templates")
		return nil, diags
	}
	clusterTemplateID := m.ClusterTemplateId.ValueString()
	for _, tmpl := range templatesResp.ClusterTemplates {
		if clusterTemplateID != "" && tmpl.ID == clusterTemplateID {
			return &tmpl, nil
		}
		if tmpl.ProductName == m.ProductName.ValueString() && tmpl.ProductVersion == m.ProductVersion.ValueString() {
			return &tmpl, nil
		}
	}
	diags.AddError("unknown cluster template", "could not find cluster templates")
	return nil, diags
}
