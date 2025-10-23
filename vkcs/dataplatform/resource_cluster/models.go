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
const importedPassword = "IMPORTED_PASSWORD"

var floatingIPAutoID = uuid.Nil.String()

func (m *ClusterModel) UpdateState(ctx context.Context, cluster *clusters.Cluster, oldConfigs ConfigsValue, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics
	diags.Append(m.UpdateFromCluster(ctx, cluster)...)
	if diags.HasError() {
		return diags
	}

	diags.Append(state.Set(ctx, m)...)

	if cluster.Configs != nil {
		maintenance, d := FlattenClusterConfigsMaintenance(ctx, cluster.Configs.Maintenance)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.Root("configs").AtName("maintenance"), maintenance)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = UpdateClusterConfigsUsers(ctx, cluster.Configs.Users, oldConfigs.Users, path.Root("configs").AtName("users"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = UpdateClusterConfigsWarehouses(ctx, cluster.Configs.Warehouses, path.Root("configs").AtName("warehouses"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = UpdateClusterConfigsSettings(ctx, cluster.Configs.Settings, oldConfigs.Settings, path.Root("configs").AtName("settings"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}
	diags.Append(UpdateClusterPodGroups(ctx, cluster.PodGroups, state)...)

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

	info, d := FlattenClusterInfo(ctx, cluster.Info)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	m.Info = info

	return diags
}

func UpdateClusterConfigsUsers(ctx context.Context, users []clusters.ClusterConfigUser, oldUsers basetypes.ListValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	usersV := make([]ConfigsUsersValue, 0, len(oldUsers.Elements()))
	if len(oldUsers.Elements()) > 0 {
		d := oldUsers.ElementsAs(ctx, &usersV, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	for _, u := range users {
		updated := false
		for i, user := range usersV {
			if user.Username.ValueString() == u.Username {
				usersV[i].CreatedAt = types.StringValue(u.CreatedAt)
				usersV[i].Id = types.StringValue(u.ID)
				usersV[i].Role = types.StringValue(u.Role)
				updated = true
			}
		}
		if !updated {
			usersV = append(usersV, ConfigsUsersValue{
				CreatedAt: types.StringValue(u.CreatedAt),
				Id:        types.StringValue(u.ID),
				Role:      types.StringValue(u.Role),
				Username:  types.StringValue(u.Username),
				Password:  types.StringValue(importedPassword),
				state:     attr.ValueStateKnown,
			})
		}
	}
	if len(usersV) == 0 {
		d := state.SetAttribute(ctx, path, types.ListNull(ConfigsUsersValue{}.Type(ctx)))
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	} else {
		d := state.SetAttribute(ctx, path, usersV)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
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

	if len(connections) == 0 {
		d := state.SetAttribute(ctx, path, types.ListNull(ConfigsWarehousesConnectionsValue{}.Type(ctx)))
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
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

func UpdateClusterConfigsSettings(ctx context.Context, o []clusters.ClusterConfigSetting, oldSettings basetypes.ListValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
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
		d := state.SetAttribute(ctx, path, types.ListNull(ConfigsSettingsValue{}.Type(ctx)))
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	} else {
		d := state.SetAttribute(ctx, path, settingsV)
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

	clusterPodgroupsMap := make(map[string]clusters.ClusterPodGroup)
	for _, p := range o {
		clusterPodgroupsMap[p.Name] = p
	}

	var statePodGroups []PodGroupsValue
	d := state.GetAttribute(ctx, path.Root("pod_groups"), &statePodGroups)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	for i, group := range statePodGroups {
		if clusterPodGroup, ok := clusterPodgroupsMap[group.Name.ValueString()]; ok {
			d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("alias"), clusterPodGroup.Alias)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("availability_zone"), clusterPodGroup.AvailabilityZone)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("count"), clusterPodGroup.Count)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("floating_ip_pool"), clusterPodGroup.FloatingIPPool)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("id"), clusterPodGroup.ID)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			volumes, d := FlattenClusterPodGroupsVolumes(ctx, clusterPodGroup.Volumes)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			d = state.SetAttribute(ctx, path.Root("pod_groups").AtListIndex(i).AtName("volumes"), volumes)
			diags.Append(d...)
			if diags.HasError() {
				return diags
			}

			resource, d := FlattenClusterPodGroupsResource(ctx, clusterPodGroup.Resource)
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
