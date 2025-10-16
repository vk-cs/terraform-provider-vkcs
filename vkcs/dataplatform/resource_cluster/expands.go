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

func BuildUpdateClusterConfigsMaintenance(ctx context.Context, currentMaintenance basetypes.ObjectValue, planMaintenance basetypes.ObjectValue) (*clusters.ClusterUpdateConfigsMaintenance, diag.Diagnostics) {
	var diags diag.Diagnostics

	var current, plan *clusters.ClusterConfigMaintenance
	current, diags = ReadExpandClusterConfigsMaintenance(ctx, currentMaintenance)
	if diags.HasError() {
		return nil, diags
	}
	plan, diags = ReadExpandClusterConfigsMaintenance(ctx, planMaintenance)
	if diags.HasError() {
		return nil, diags
	}

	var changed bool
	update := &clusters.ClusterUpdateConfigsMaintenance{}

	update.Start = plan.Start
	if current.Start != plan.Start {
		changed = true
	}

	if !backupEqual(current.Backup, plan.Backup) {
		if plan.Backup == nil {
			update.Backup = clusters.ClusterConfigMaintenanceBackup{}
		} else {
			update.Backup = *plan.Backup
		}
		changed = true
	}

	crontabs := buildCrontabChanges(current.CronTabs, plan.CronTabs)
	if len(crontabs.Create) > 0 || len(crontabs.Update) > 0 || len(crontabs.Delete) > 0 {
		update.Crontabs = crontabs
		changed = true
	}

	if changed {
		return update, nil
	}

	return nil, nil
}

func backupEqual(a, b *clusters.ClusterConfigMaintenanceBackup) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return backupObjEqual(a.Full, b.Full) &&
		backupObjEqual(a.Incremental, b.Incremental) &&
		backupObjEqual(a.Differential, b.Differential)
}

func backupObjEqual(x, y *clusters.ClusterConfigMaintenanceBackupObj) bool {
	if x == nil && y == nil {
		return true
	}
	if x == nil || y == nil {
		return false
	}

	if x.Start != y.Start {
		return false
	}

	if (x.Enabled == nil) != (y.Enabled == nil) || (x.Enabled != nil && *x.Enabled != *y.Enabled) {
		return false
	}

	if (x.KeepCount == nil) != (y.KeepCount == nil) || (x.KeepCount != nil && *x.KeepCount != *y.KeepCount) {
		return false
	}

	if (x.KeepTime == nil) != (y.KeepTime == nil) || (x.KeepTime != nil && *x.KeepTime != *y.KeepTime) {
		return false
	}

	return true
}

func buildCrontabChanges(current, plan []clusters.ClusterConfigMaintenanceCronTabs) clusters.ClusterUpdateConfigsMaintenanceCrontabs {
	result := clusters.ClusterUpdateConfigsMaintenanceCrontabs{}

	currentByName := make(map[string]clusters.ClusterConfigMaintenanceCronTabs)
	for _, c := range current {
		currentByName[c.Name] = c
	}

	planNames := make(map[string]struct{})

	for _, p := range plan {
		planNames[p.Name] = struct{}{}

		c, exists := currentByName[p.Name]
		if !exists {
			result.Create = append(result.Create, clusters.ClusterUpdateConfigsMaintenanceCrontabsCreate{
				Name:     p.Name,
				Start:    p.Start,
				Settings: toUpdateSettings(p.Settings),
			})
			continue
		}

		if c.Start != p.Start || !settingsEqual(c.Settings, p.Settings) {
			result.Update = append(result.Update, clusters.ClusterUpdateConfigsMaintenanceCrontabsUpdate{
				ID:       c.ID,
				Start:    p.Start,
				Settings: toUpdateSettings(p.Settings),
			})
		}
	}

	for _, c := range current {
		if _, exists := planNames[c.Name]; !exists {
			result.Delete = append(result.Delete, clusters.ClusterUpdateConfigsMaintenanceCrontabsDelete{ID: c.ID})
		}
	}

	return result
}

func settingsEqual(a, b []clusters.ClusterConfigSetting) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]string, len(a))
	for _, s := range a {
		m[s.Alias] = s.Value
	}
	for _, s := range b {
		if v, ok := m[s.Alias]; !ok || v != s.Value {
			return false
		}
	}
	return true
}

func toUpdateSettings(settings []clusters.ClusterConfigSetting) []clusters.ClusterCreateConfigSetting {
	out := make([]clusters.ClusterCreateConfigSetting, 0, len(settings))
	for _, s := range settings {
		out = append(out, clusters.ClusterCreateConfigSetting{
			Alias: s.Alias,
			Value: s.Value,
		})
	}
	return out
}

func ReadExpandClusterConfigsMaintenance(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterConfigMaintenance, diag.Diagnostics) {
	maintenanceObjV, diags := ConfigsMaintenanceType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	maintenance := maintenanceObjV.(ConfigsMaintenanceValue)
	result := clusters.ClusterConfigMaintenance{
		Start: maintenance.Start.ValueString(),
	}

	if o := maintenance.Crontabs; !o.IsUnknown() && !o.IsNull() {
		crontabs, diags := ReadExpandClusterCrontabs(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.CronTabs = crontabs
	}

	if o := maintenance.Backup; !o.IsUnknown() && !o.IsNull() {
		backup, diags := ReadExpandClusterBackup(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		// normalize empty backup
		if backup.Full == nil && backup.Incremental == nil && backup.Differential == nil {
			result.Backup = nil
		} else {
			result.Backup = backup
		}
	}

	return &result, nil
}

func ReadExpandClusterCrontabs(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterConfigMaintenanceCronTabs, diag.Diagnostics) {
	crontabs := make([]ConfigsMaintenanceCrontabsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &crontabs, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterConfigMaintenanceCronTabs, len(crontabs))
	for i, v := range crontabs {
		result[i] = clusters.ClusterConfigMaintenanceCronTabs{
			ID:       v.Id.ValueString(),
			Required: v.Required.ValueBool(),
			Name:     v.Name.ValueString(),
			Start:    v.Start.ValueString(),
		}

		if o := v.Settings; !o.IsUnknown() && !o.IsNull() {
			var settings []clusters.ClusterConfigSetting
			settings, diags = ReadExpandClusterCrontabSettings(ctx, &v.Settings)
			if diags.HasError() {
				return nil, diags
			}

			result[i].Settings = settings
		}
	}

	return result, nil
}

func ReadExpandClusterCrontabSettings(ctx context.Context, o *basetypes.ListValue) ([]clusters.ClusterConfigSetting, diag.Diagnostics) {
	settings := make([]ConfigsMaintenanceCrontabsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settings, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterConfigSetting, len(settings))
	for i, s := range settings {
		result[i] = clusters.ClusterConfigSetting{
			Alias: s.Alias.ValueString(),
			Value: s.Value.ValueString(),
		}
	}

	return result, nil
}

func ReadExpandClusterBackup(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterConfigMaintenanceBackup, diag.Diagnostics) {
	backupObjV, diags := ConfigsMaintenanceBackupType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backup := backupObjV.(ConfigsMaintenanceBackupValue)

	var result clusters.ClusterConfigMaintenanceBackup

	if o := backup.Full; !o.IsUnknown() && !o.IsNull() {
		full, diags := ReadExpandClusterBackupFull(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Full = full
	}

	if o := backup.Incremental; !o.IsUnknown() && !o.IsNull() {
		incremental, diags := ReadExpandClusterBackupIncremental(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Incremental = incremental
	}

	if o := backup.Differential; !o.IsUnknown() && !o.IsNull() {
		differential, diags := ReadExpandClusterBackupDifferential(ctx, o)
		if diags.HasError() {
			return nil, diags
		}
		result.Differential = differential
	}

	return &result, nil
}

func ReadExpandClusterBackupFull(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterConfigMaintenanceBackupObj, diag.Diagnostics) {
	fullObjV, diags := ConfigsMaintenanceBackupFullType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backupFull := fullObjV.(ConfigsMaintenanceBackupFullValue)

	var keepCount *int
	var keepTime *int

	if !backupFull.KeepCount.IsNull() && !backupFull.KeepCount.IsUnknown() {
		v := int(backupFull.KeepCount.ValueInt64())
		keepCount = &v
	}

	if !backupFull.KeepTime.IsNull() && !backupFull.KeepTime.IsUnknown() {
		v := int(backupFull.KeepTime.ValueInt64())
		keepTime = &v
	}

	result := clusters.ClusterConfigMaintenanceBackupObj{
		Start:     backupFull.Start.ValueString(),
		KeepCount: keepCount,
		KeepTime:  keepTime,
	}

	return &result, nil
}

func ReadExpandClusterBackupIncremental(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterConfigMaintenanceBackupObj, diag.Diagnostics) {
	incrementalObjV, diags := ConfigsMaintenanceBackupIncrementalType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backupIncremental := incrementalObjV.(ConfigsMaintenanceBackupIncrementalValue)

	keepCount := int(backupIncremental.KeepCount.ValueInt64())
	keepTime := int(backupIncremental.KeepTime.ValueInt64())

	result := clusters.ClusterConfigMaintenanceBackupObj{
		Start:     backupIncremental.Start.ValueString(),
		KeepCount: &keepCount,
		KeepTime:  &keepTime,
	}

	return &result, nil
}

func ReadExpandClusterBackupDifferential(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterConfigMaintenanceBackupObj, diag.Diagnostics) {
	differentialObjV, diags := ConfigsMaintenanceBackupDifferentialType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	backupDifferential := differentialObjV.(ConfigsMaintenanceBackupDifferentialValue)

	keepCount := int(backupDifferential.KeepCount.ValueInt64())
	keepTime := int(backupDifferential.KeepTime.ValueInt64())

	result := clusters.ClusterConfigMaintenanceBackupObj{
		Enabled:   backupDifferential.Enabled.ValueBoolPointer(),
		Start:     backupDifferential.Start.ValueString(),
		KeepCount: &keepCount,
		KeepTime:  &keepTime,
	}

	return &result, nil
}
