package resource_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"
)

func ExpandClusterConfigs(ctx context.Context, v ConfigsValue) (*clusters.ClusterCreateConfig, diag.Diagnostics) {
	result := &clusters.ClusterCreateConfig{}

	if o := v.Maintenance; !o.IsUnknown() && !o.IsNull() {
		maintenance, diags := ExpandClusterConfigsMaintenance(ctx, o)
		if diags != nil && diags.HasError() {
			return nil, diags
		}
		result.Maintenance = maintenance
	}

	if o := v.Settings; !o.IsUnknown() && !o.IsNull() {
		settings, diags := ExpandClusterConfigsSettings(ctx, o)
		if diags != nil && diags.HasError() {
			return nil, diags
		}
		result.Settings = settings
	}

	if o := v.Users; !o.IsUnknown() && !o.IsNull() {
		users, diags := ExpandClusterConfigsUsers(ctx, o)
		if diags != nil && diags.HasError() {
			return nil, diags
		}
		result.Users = users
	}

	if o := v.Warehouses; !o.IsUnknown() && !o.IsNull() {
		warehouses, diags := ExpandClusterConfigsWarehouses(ctx, o)
		if diags != nil && diags.HasError() {
			return nil, diags
		}
		result.Warehouses = warehouses
	}

	return result, nil
}

func ExpandClusterConfigsMaintenance(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterCreateConfigMaintenance, diag.Diagnostics) {
	maintenanceObjV, diags := ConfigsMaintenanceType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	maintenance := maintenanceObjV.(ConfigsMaintenanceValue)
	result := clusters.ClusterCreateConfigMaintenance{
		Start: maintenance.Start.ValueString(),
	}

	if o := maintenance.Crontabs; !o.IsUnknown() && !o.IsNull() {
		crontabs, diags := ExpandClusterCrontabs(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.CronTabs = crontabs
	}

	if o := maintenance.Backup; !o.IsUnknown() && !o.IsNull() {
		backup, diags := ExpandClusterBackup(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Backup = backup
	}

	return &result, nil
}

func ExpandClusterCrontabs(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigMaintenanceCronTabs, diag.Diagnostics) {
	crontabs := make([]ConfigsMaintenanceCrontabsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &crontabs, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigMaintenanceCronTabs, len(crontabs))
	for i, v := range crontabs {
		result[i] = clusters.ClusterCreateConfigMaintenanceCronTabs{
			Name:  v.Name.ValueString(),
			Start: v.Start.ValueString(),
		}

		if o := v.Settings; !o.IsUnknown() && !o.IsNull() {
			var settings []clusters.ClusterCreateConfigSetting
			settings, diags = ExpandClusterCrontabSettings(ctx, &v.Settings)
			if diags.HasError() {
				return nil, diags
			}

			result[i].Settings = settings
		}
	}

	return result, nil
}

func ExpandClusterCrontabSettings(ctx context.Context, o *basetypes.ListValue) ([]clusters.ClusterCreateConfigSetting, diag.Diagnostics) {
	settings := make([]ConfigsMaintenanceCrontabsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settings, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigSetting, len(settings))
	for i, s := range settings {
		result[i] = clusters.ClusterCreateConfigSetting{
			Alias: s.Alias.ValueString(),
			Value: s.Value.ValueString(),
		}
	}

	return result, nil
}

func ExpandClusterConfigsSettings(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigSetting, diag.Diagnostics) {
	settingsV := make([]ConfigsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settingsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigSetting, len(settingsV))
	for i, s := range settingsV {
		result[i] = clusters.ClusterCreateConfigSetting{
			Alias: s.Alias.ValueString(),
			Value: s.Value.ValueString(),
		}
	}
	return result, nil
}

func ExpandClusterBackup(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterCreateConfigMaintenanceBackup, diag.Diagnostics) {
	backupObjV, diags := ConfigsMaintenanceBackupType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backup := backupObjV.(ConfigsMaintenanceBackupValue)

	var result clusters.ClusterCreateConfigMaintenanceBackup

	if o := backup.Full; !o.IsUnknown() && !o.IsNull() {
		full, diags := ExpandClusterBackupFull(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Full = full
	}

	if o := backup.Incremental; !o.IsUnknown() && !o.IsNull() {
		incremental, diags := ExpandClusterBackupIncremental(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Incremental = incremental
	}

	if o := backup.Differential; !o.IsUnknown() && !o.IsNull() {
		differential, diags := ExpandClusterBackupDifferential(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Differential = differential
	}

	if result.Full == nil && result.Incremental == nil && result.Differential == nil {
		return nil, nil
	}

	return &result, nil
}

func ExpandClusterBackupFull(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterCreateConfigMaintenanceBackupObj, diag.Diagnostics) {
	fullObjV, diags := ConfigsMaintenanceBackupFullType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backupFull := fullObjV.(ConfigsMaintenanceBackupFullValue)

	result := clusters.ClusterCreateConfigMaintenanceBackupObj{
		Start:     backupFull.Start.ValueString(),
		KeepCount: int(backupFull.KeepCount.ValueInt64()),
		KeepTime:  int(backupFull.KeepTime.ValueInt64()),
	}

	return &result, nil
}

func ExpandClusterBackupIncremental(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterCreateConfigMaintenanceBackupObj, diag.Diagnostics) {
	incrementalObjV, diags := ConfigsMaintenanceBackupIncrementalType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backupIncremental := incrementalObjV.(ConfigsMaintenanceBackupIncrementalValue)

	result := clusters.ClusterCreateConfigMaintenanceBackupObj{
		Start:     backupIncremental.Start.ValueString(),
		KeepCount: int(backupIncremental.KeepCount.ValueInt64()),
		KeepTime:  int(backupIncremental.KeepTime.ValueInt64()),
	}

	return &result, nil
}

func ExpandClusterBackupDifferential(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterCreateConfigMaintenanceBackupObj, diag.Diagnostics) {
	differentialObjV, diags := ConfigsMaintenanceBackupDifferentialType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backupDifferential := differentialObjV.(ConfigsMaintenanceBackupDifferentialValue)

	result := clusters.ClusterCreateConfigMaintenanceBackupObj{
		Start:     backupDifferential.Start.ValueString(),
		KeepCount: int(backupDifferential.KeepCount.ValueInt64()),
		KeepTime:  int(backupDifferential.KeepTime.ValueInt64()),
	}

	return &result, nil
}

func ExpandClusterConfigsUsers(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigUser, diag.Diagnostics) {
	usersV, diags := ReadClusterConfigsUsers(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigUser, len(usersV))
	for i, u := range usersV {
		result[i] = clusters.ClusterCreateConfigUser{
			Username: u.Username.ValueString(),
			Password: u.Password.ValueString(),
			Role:     u.Role.ValueString(),
		}
	}
	return result, nil
}

func ReadClusterConfigsUsers(ctx context.Context, o basetypes.ListValue) ([]ConfigsUsersValue, diag.Diagnostics) {
	usersV := make([]ConfigsUsersValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &usersV, false)
	if diags.HasError() {
		return nil, diags
	}
	return usersV, nil
}

func ExpandClusterConfigsWarehouses(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigWarehouse, diag.Diagnostics) {
	warehousesV, diags := ReadClusterConfigsWarehouses(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigWarehouse, len(warehousesV))
	for i, v := range warehousesV {
		result[i] = clusters.ClusterCreateConfigWarehouse{
			Name: v.Name.ValueString(),
		}

		if o := v.Connections; !o.IsUnknown() && !o.IsNull() {
			wConnections, diags := ExpandClusterConfigsWarehousesConnections(ctx, o)
			if diags.HasError() {
				return nil, diags
			}
			result[i].Connections = wConnections
		}
	}
	return result, nil
}

func ReadClusterConfigsWarehouses(ctx context.Context, o basetypes.ListValue) ([]ConfigsWarehousesValue, diag.Diagnostics) {
	warehousesV := make([]ConfigsWarehousesValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &warehousesV, false)
	if diags.HasError() {
		return nil, diags
	}
	return warehousesV, nil
}

func ExpandClusterConfigsWarehousesConnections(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigWarehouseConnection, diag.Diagnostics) {
	connectionsV, diags := ReadClusterConfigsWarehousesConnections(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigWarehouseConnection, len(connectionsV))
	for i, v := range connectionsV {
		result[i] = clusters.ClusterCreateConfigWarehouseConnection{
			Name: v.Name.ValueString(),
			Plug: v.Plug.ValueString(),
		}

		if o := v.Settings; !o.IsUnknown() && !o.IsNull() {
			wcSettings, diags := ExpandClusterWarehousesConnectionsSettings(ctx, o)
			if diags.HasError() {
				return nil, diags
			}
			result[i].Settings = wcSettings
		}
	}
	return result, nil
}

func ReadClusterConfigsWarehousesConnections(ctx context.Context, o basetypes.ListValue) ([]ConfigsWarehousesConnectionsValue, diag.Diagnostics) {
	connectionsV := make([]ConfigsWarehousesConnectionsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &connectionsV, false)
	if diags.HasError() {
		return nil, diags
	}
	return connectionsV, nil
}

func ExpandClusterWarehousesConnectionsSettings(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigSetting, diag.Diagnostics) {
	settingsV := make([]ConfigsWarehousesConnectionsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settingsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigSetting, len(settingsV))
	for i, v := range settingsV {
		result[i] = clusters.ClusterCreateConfigSetting{
			Alias: v.Alias.ValueString(),
			Value: v.Value.ValueString(),
		}
	}
	return result, nil
}

func ReadClusterConfigsWarehousesConnectionsSettings(ctx context.Context, o basetypes.ListValue) ([]ConfigsWarehousesConnectionsSettingsValue, diag.Diagnostics) {
	settingsV := make([]ConfigsWarehousesConnectionsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settingsV, false)
	if diags.HasError() {
		return nil, diags
	}
	return settingsV, nil
}

func ExpandClusterPodGroups(ctx context.Context, template *templates.ClusterTemplate, o basetypes.ListValue) ([]clusters.ClusterCreatePodGroup, diag.Diagnostics) {
	podGroupsV := make([]PodGroupsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &podGroupsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreatePodGroup, len(podGroupsV))
	for i, v := range podGroupsV {
		name := v.Name.ValueString()
		var podGroupTemplateID string
		for _, podGroupTemplate := range template.PodGroups {
			if name == podGroupTemplate.Name {
				podGroupTemplateID = podGroupTemplate.ID
				break
			}
		}
		if podGroupTemplateID == "" {
			diags.AddError("unknown pod group name", "could not find pod group template")
		}

		count := int(v.Count.ValueInt64())
		result[i] = clusters.ClusterCreatePodGroup{
			Count:              &count,
			PodGroupTemplateID: podGroupTemplateID,
		}

		if o := v.Resource; !o.IsUnknown() && !o.IsNull() {
			pgResource, diags := ExpandClusterPodGroupsResource(ctx, o)
			if diags.HasError() {
				return nil, diags
			}
			result[i].Resource = pgResource
		}

		if o := v.Volumes; !o.IsUnknown() && !o.IsNull() {
			pgVolumes, diags := ExpandClusterPodGroupsVolumes(ctx, o)
			if diags.HasError() {
				return nil, diags
			}
			result[i].Volumes = pgVolumes
		}
	}
	return result, nil
}

func ExpandClusterPodGroupsResource(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterCreatePodGroupResource, diag.Diagnostics) {
	resource, diags := ReadClusterPodGroupResources(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	result := clusters.ClusterCreatePodGroupResource{
		CPURequest: resource.CpuRequest.ValueString(),
		RAMRequest: resource.RamRequest.ValueString(),
	}
	return &result, nil
}

func ExpandClusterPodGroupsVolumes(ctx context.Context, o basetypes.MapValue) (map[string]clusters.ClusterCreatePodGroupVolume, diag.Diagnostics) {
	volumesV := make(map[string]PodGroupsVolumesValue)
	diags := o.ElementsAs(ctx, &volumesV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make(map[string]clusters.ClusterCreatePodGroupVolume)
	for i, v := range volumesV {
		result[i] = clusters.ClusterCreatePodGroupVolume{
			StorageClassName: v.StorageClassName.ValueString(),
			Storage:          v.Storage.ValueString(),
			Count:            int(v.Count.ValueInt64()),
		}
	}
	return result, nil
}

func ReadClusterPodGroups(ctx context.Context, o basetypes.ListValue) ([]PodGroupsValue, diag.Diagnostics) {
	podGroupsV := make([]PodGroupsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &podGroupsV, false)
	if diags.HasError() {
		return nil, diags
	}
	return podGroupsV, nil
}

func ReadClusterPodGroupResources(ctx context.Context, o basetypes.ObjectValue) (*PodGroupsResourceValue, diag.Diagnostics) {
	resourceV, diags := PodGroupsResourceType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}
	resource := resourceV.(PodGroupsResourceValue)
	return &resource, nil
}

func ReadClusterPodGroupVolumes(ctx context.Context, o basetypes.MapValue) (map[string]PodGroupsVolumesValue, diag.Diagnostics) {
	volumesV := make(map[string]PodGroupsVolumesValue)

	if o.IsNull() || o.IsUnknown() {
		return volumesV, nil
	}

	diags := o.ElementsAs(ctx, &volumesV, false)
	if diags.HasError() {
		return nil, diags
	}
	return volumesV, nil
}

func BuildUpdateClusterConfigsMaintenance(ctx context.Context, stateMaintenance basetypes.ObjectValue, planMaintenance basetypes.ObjectValue) (*clusters.ClusterUpdateConfigsMaintenance, diag.Diagnostics) {
	var diags diag.Diagnostics

	stateMaintenanceObj, d := ConfigsMaintenanceType{}.ValueFromObject(ctx, stateMaintenance)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	stateVal := stateMaintenanceObj.(ConfigsMaintenanceValue)

	planMaintenanceObj, d := ConfigsMaintenanceType{}.ValueFromObject(ctx, planMaintenance)
	diags.Append(d...)
	if diags.HasError() {
		return nil, diags
	}
	planVal := planMaintenanceObj.(ConfigsMaintenanceValue)

	var changed bool
	update := &clusters.ClusterUpdateConfigsMaintenance{}

	if (!planVal.Start.IsUnknown() && !planVal.Start.IsNull()) && !planVal.Start.Equal(stateVal.Start) {
		start := planVal.Start.ValueString()
		update.Start = &start
		changed = true
	}

	if !planVal.Backup.IsNull() && !planVal.Backup.IsUnknown() {
		if stateVal.Backup.IsNull() || stateVal.Backup.IsUnknown() || !stateVal.Backup.Equal(planVal.Backup) {
			backup, d := ExpandClusterBackup(ctx, planVal.Backup)
			diags.Append(d...)
			if diags.HasError() {
				return nil, diags
			}
			update.Backup = backup
			changed = true
		}
	}

	if !planVal.Crontabs.IsNull() && !planVal.Crontabs.IsUnknown() {
		crontabChanges, d := expandClusterCrontabsUpdate(ctx, stateVal.Crontabs, planVal.Crontabs)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
		if crontabChanges != nil {
			update.Crontabs = crontabChanges
			changed = true
		}
	}

	if !changed {
		return nil, diags
	}

	return update, nil
}

func expandClusterCrontabsUpdate(ctx context.Context, stateLV, planLV basetypes.ListValue) (*clusters.ClusterUpdateConfigsMaintenanceCrontabs, diag.Diagnostics) {
	var diags diag.Diagnostics
	result := clusters.ClusterUpdateConfigsMaintenanceCrontabs{}

	var stateCrontabs []ConfigsMaintenanceCrontabsValue
	if !stateLV.IsNull() && !stateLV.IsUnknown() {
		d := stateLV.ElementsAs(ctx, &stateCrontabs, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	var planCrontabs []ConfigsMaintenanceCrontabsValue
	if !planLV.IsNull() && !planLV.IsUnknown() {
		d := planLV.ElementsAs(ctx, &planCrontabs, false)
		diags.Append(d...)
		if diags.HasError() {
			return nil, diags
		}
	}

	stateCrontabsByName := make(map[string]ConfigsMaintenanceCrontabsValue)
	for _, c := range stateCrontabs {
		stateCrontabsByName[c.Name.ValueString()] = c
	}

	planNames := make(map[string]struct{})
	for _, p := range planCrontabs {
		name := p.Name.ValueString()
		planNames[name] = struct{}{}

		c, exists := stateCrontabsByName[name]
		if !exists {
			var settings []clusters.ClusterCreateConfigSetting

			if !p.Settings.IsNull() && !p.Settings.IsUnknown() {
				var d diag.Diagnostics
				settings, d = ExpandClusterCrontabSettings(ctx, &p.Settings)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
			}

			result.Create = append(result.Create, clusters.ClusterUpdateConfigsMaintenanceCrontabsCreate{
				Name:     p.Name.ValueString(),
				Start:    p.Start.ValueString(),
				Settings: settings,
			})
			continue
		}

		needUpdateSettings := false
		pEmpty := p.Settings.IsNull() || p.Settings.IsUnknown()
		cEmpty := c.Settings.IsNull() || c.Settings.IsUnknown()
		if (pEmpty && !cEmpty) || (!pEmpty && cEmpty) || (!pEmpty && !cEmpty && !p.Settings.Equal(c.Settings)) {
			needUpdateSettings = true
		}

		if p.Start.ValueString() != c.Start.ValueString() || needUpdateSettings {
			var settings []clusters.ClusterCreateConfigSetting

			if !p.Settings.IsNull() && !p.Settings.IsUnknown() {
				var d diag.Diagnostics
				settings, d = ExpandClusterCrontabSettings(ctx, &p.Settings)
				diags.Append(d...)
				if diags.HasError() {
					return nil, diags
				}
			}

			result.Update = append(result.Update, clusters.ClusterUpdateConfigsMaintenanceCrontabsUpdate{
				ID:       c.Id.ValueString(),
				Start:    p.Start.ValueString(),
				Settings: settings,
			})
		}
	}

	for _, c := range stateCrontabs {
		if _, exists := planNames[c.Name.ValueString()]; !exists {
			result.Delete = append(result.Delete, clusters.ClusterUpdateConfigsMaintenanceCrontabsDelete{ID: c.Id.ValueString()})
		}
	}

	if result.Create == nil && result.Update == nil && result.Delete == nil {
		return nil, nil
	}

	return &result, diags
}
