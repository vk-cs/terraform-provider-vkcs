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
		warehouses, diags := ExpandClusterWarehouses(ctx, o)
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

func ExpandClusterWarehouses(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigWarehouse, diag.Diagnostics) {
	warehousesV := make([]ConfigsWarehousesValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &warehousesV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterCreateConfigWarehouse, len(warehousesV))
	for i, v := range warehousesV {
		result[i] = clusters.ClusterCreateConfigWarehouse{
			Name: v.Name.ValueString(),
		}

		if o := v.Connections; !o.IsUnknown() && !o.IsNull() {
			wConnections, diags := ExpandClusterWarehousesConnections(ctx, o)
			if diags.HasError() {
				return nil, diags
			}
			result[i].Connections = wConnections
		}
	}
	return result, nil
}

func ExpandClusterWarehousesConnections(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterCreateConfigWarehouseConnection, diag.Diagnostics) {
	connectionsV := make([]ConfigsWarehousesConnectionsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &connectionsV, false)
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
	resourceV, diags := PodGroupsResourceType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	resource := resourceV.(PodGroupsResourceValue)
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
