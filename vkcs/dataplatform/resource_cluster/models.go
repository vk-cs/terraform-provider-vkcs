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
		var d diag.Diagnostics

		d = UpdateClusterConfigsMaintenance(ctx, cluster.Configs.Maintenance, oldConfigs.Maintenance, path.Root("configs").AtName("maintenance"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = UpdateClusterConfigsUsers(ctx, cluster.Configs.Users, oldConfigs.Users, path.Root("configs").AtName("users"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = UpdateClusterConfigsWarehouses(ctx, cluster.Configs.Warehouses, oldConfigs.Warehouses, path.Root("configs").AtName("warehouses"), state)
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

func UpdateClusterConfigsMaintenance(ctx context.Context, maintenance *clusters.ClusterConfigMaintenance, oldMaintenanceV basetypes.ObjectValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	oldMaintenance, d := ExpandClusterConfigsMaintenance(ctx, oldMaintenanceV)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	maintenance.CronTabs = MergeClusterConfigsMaintenanceCronTabs(maintenance.CronTabs, oldMaintenance.CronTabs)

	maintenanceV, d := FlattenClusterConfigsMaintenance(ctx, maintenance)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	d = state.SetAttribute(ctx, path, maintenanceV)
	diags.Append(d...)

	return diags
}

func MergeClusterConfigsMaintenanceCronTabs(newCrontabs []clusters.ClusterConfigMaintenanceCronTabs, oldCrontabs []clusters.ClusterCreateConfigMaintenanceCronTabs) []clusters.ClusterConfigMaintenanceCronTabs {
	var result []clusters.ClusterConfigMaintenanceCronTabs

	oldCronTabsByName := make(map[string]clusters.ClusterCreateConfigMaintenanceCronTabs, len(oldCrontabs))
	for _, old := range oldCrontabs {
		oldCronTabsByName[old.Name] = old
	}

	for _, newCron := range newCrontabs {
		oldCron, ok := oldCronTabsByName[newCron.Name]
		if !ok {
			continue
		}

		newCron.Settings = MergeClusterConfigsMaintenanceCronTabsSettings(newCron.Settings, oldCron.Settings)
		result = append(result, newCron)
	}

	return result
}

func MergeClusterConfigsMaintenanceCronTabsSettings(newSettings []clusters.ClusterConfigSetting, oldSettings []clusters.ClusterCreateConfigSetting) []clusters.ClusterConfigSetting {
	var filtered []clusters.ClusterConfigSetting

	existing := make(map[string]struct{}, len(oldSettings))
	for _, s := range oldSettings {
		existing[s.Alias] = struct{}{}
	}

	for _, s := range newSettings {
		if _, ok := existing[s.Alias]; ok {
			filtered = append(filtered, s)
		}
	}

	return filtered
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

		var roleValue types.String
		if u.Role == "" {
			roleValue = types.StringNull()
		} else {
			roleValue = types.StringValue(u.Role)
		}

		for i, user := range usersV {
			if user.Username.ValueString() == u.Username {
				usersV[i].CreatedAt = types.StringValue(u.CreatedAt)
				usersV[i].Id = types.StringValue(u.ID)
				usersV[i].Role = roleValue
				updated = true
			}
		}
		if !updated {
			usersV = append(usersV, ConfigsUsersValue{
				CreatedAt: types.StringValue(u.CreatedAt),
				Id:        types.StringValue(u.ID),
				Role:      roleValue,
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

func UpdateClusterConfigsWarehouses(ctx context.Context, warehouses []clusters.ClusterConfigWarehouse, oldWarehouses basetypes.ListValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	if warehouses == nil {
		return nil
	}

	warehousesV := make([]ConfigsWarehousesValue, 0, len(oldWarehouses.Elements()))
	if len(oldWarehouses.Elements()) > 0 {
		d := oldWarehouses.ElementsAs(ctx, &warehousesV, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	var oldWarehousesMap = make(map[string]ConfigsWarehousesValue, len(warehousesV))
	for _, w := range warehousesV {
		oldWarehousesMap[w.Name.ValueString()] = w
	}

	for i, w := range warehouses {
		d := UpdateClusterConfigsWarehousesConnections(ctx, i, w.Connections, oldWarehousesMap[w.Name].Connections, path.AtListIndex(i).AtName("connections"), state)
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

func UpdateClusterConfigsWarehousesConnections(ctx context.Context, i int, connections *[]clusters.ClusterConfigWarehouseConnection, oldConnections basetypes.ListValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	if connections == nil {
		d := state.SetAttribute(ctx, path, types.ListNull(ConfigsWarehousesConnectionsValue{}.Type(ctx)))
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		return nil
	}

	if len(*connections) == 0 {
		connectionsV, d := types.ListValue(ConfigsWarehousesConnectionsValue{}.Type(ctx), []attr.Value{})
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path, connectionsV)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		return nil
	}

	connectionsV := make([]ConfigsWarehousesConnectionsValue, 0, len(oldConnections.Elements()))
	if len(oldConnections.Elements()) > 0 {
		d := oldConnections.ElementsAs(ctx, &connectionsV, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	var oldConnectionsMap = make(map[string]ConfigsWarehousesConnectionsValue, len(connectionsV))
	for _, c := range connectionsV {
		oldConnectionsMap[c.Name.ValueString()] = c
	}
	for j, c := range *connections {
		oldConnection := oldConnectionsMap[c.Name]
		d := UpdateClusterConfigsWarehousesConnectionsSettings(ctx, c.Settings, oldConnection.Settings, path.AtListIndex(j).AtName("settings"), state)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		var stateSettings []ConfigsWarehousesConnectionsSettingsValue
		d = state.GetAttribute(ctx, path.AtListIndex(j).AtName("settings"), &stateSettings)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		d = state.SetAttribute(ctx, path.AtListIndex(j).AtName("created_at"), c.CreatedAt)
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

func UpdateClusterConfigsWarehousesConnectionsSettings(ctx context.Context, newSettings []clusters.ClusterConfigSetting, oldSettings basetypes.ListValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	oldElems := oldSettings.Elements()
	oldSettingsV := make([]ConfigsWarehousesConnectionsSettingsValue, 0, len(oldElems))
	if len(oldElems) > 0 {
		d := oldSettings.ElementsAs(ctx, &oldSettingsV, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	} else {
		return nil
	}

	newAliases := make(map[string]clusters.ClusterConfigSetting, len(newSettings))
	for _, newSetting := range newSettings {
		newAliases[newSetting.Alias] = newSetting
	}

	for i, oldSetting := range oldSettingsV {
		if _, exists := newAliases[oldSetting.Alias.ValueString()]; exists {
			oldSettingsV[i].Value = types.StringValue(newAliases[oldSetting.Alias.ValueString()].Value)
		}
	}

	d := state.SetAttribute(ctx, path, oldSettingsV)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	return nil
}

func UpdateClusterConfigsSettings(ctx context.Context, newSettings []clusters.ClusterConfigSetting, oldSettings basetypes.ListValue, path path.Path, state *tfsdk.State) diag.Diagnostics {
	var diags diag.Diagnostics

	oldElems := oldSettings.Elements()
	settingsV := make([]ConfigsSettingsValue, 0, len(oldElems))
	if len(oldElems) > 0 {
		d := oldSettings.ElementsAs(ctx, &settingsV, false)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
	}

	existingAliases := make(map[string]int, len(settingsV))
	for i, old := range settingsV {
		existingAliases[old.Alias.ValueString()] = i
	}

	for _, newSetting := range newSettings {
		if i, exists := existingAliases[newSetting.Alias]; exists {
			settingsV[i].Value = types.StringValue(newSetting.Value)
		}
	}

	var value any
	if len(settingsV) == 0 {
		value = types.ListNull(ConfigsSettingsValue{}.Type(ctx))
	} else {
		value = settingsV
	}

	d := state.SetAttribute(ctx, path, value)
	diags.Append(d...)
	if diags.HasError() {
		return diags
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
