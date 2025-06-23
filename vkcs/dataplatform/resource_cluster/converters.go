package resource_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/templates"

	"strings"
)

func FlattenClusterConfigsSettings(ctx context.Context, o []clusters.ClusterConfigSetting) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics
	if o == nil {
		return types.ListNull(ConfigsSettingsValue{}.Type(ctx)), nil
	}

	settingsV := make([]attr.Value, len(o))
	for i, s := range o {
		settingsV[i] = ConfigsSettingsValue{
			Alias: types.StringValue(s.Alias),
			Value: types.StringValue(s.Value),
			state: attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsSettingsValue{}.Type(ctx), settingsV)
	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsSettingsValue{}.Type(ctx)), diags
	}
	return result, nil
}

func FlattenClusterConfigsMaintenance(ctx context.Context, o *clusters.ClusterConfigMaintenance) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ObjectNull(ConfigsMaintenanceValue{}.AttributeTypes(ctx)), nil
	}

	backup, d := FlattenClusterConfigsMaintenanceBackup(ctx, o.Backup)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceValue{}.AttributeTypes(ctx)), diags
	}

	cronTabs, d := FlattenClusterConfigsMaintenanceCronTabs(ctx, o.CronTabs)

	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceValue{}.AttributeTypes(ctx)), diags
	}

	maintenanceV := ConfigsMaintenanceValue{
		Backup:   backup,
		Crontabs: cronTabs,
		Start:    types.StringValue(o.Start),
		state:    attr.ValueStateKnown,
	}

	result, d := maintenanceV.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceValue{}.AttributeTypes(ctx)), nil
	}

	return result, nil
}

func FlattenClusterConfigsMaintenanceBackup(ctx context.Context, o *clusters.ClusterConfigMaintenanceBackup) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ObjectNull(ConfigsMaintenanceBackupValue{}.AttributeTypes(ctx)), nil
	}

	differential, d := FlattenClusterConfigsMaintenanceBackupObj(ctx, o.Differential)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceBackupValue{}.AttributeTypes(ctx)), diags
	}

	full, d := FlattenClusterConfigsMaintenanceBackupObj(ctx, o.Full)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceBackupValue{}.AttributeTypes(ctx)), diags
	}

	incremental, d := FlattenClusterConfigsMaintenanceBackupObj(ctx, o.Incremental)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceBackupValue{}.AttributeTypes(ctx)), diags
	}

	backupV := ConfigsMaintenanceBackupValue{
		Differential: differential,
		Full:         full,
		Incremental:  incremental,
		state:        attr.ValueStateKnown,
	}
	result, d := backupV.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceBackupValue{}.AttributeTypes(ctx)), nil
	}

	return result, nil
}

func FlattenClusterConfigsMaintenanceBackupObj(ctx context.Context, o *clusters.ClusterConfigMaintenanceBackupObj) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ObjectNull(ConfigsMaintenanceBackupDifferentialValue{}.AttributeTypes(ctx)), nil
	}

	objV := ConfigsMaintenanceBackupDifferentialValue{
		Enabled:   types.BoolValue(o.Enabled),
		KeepCount: types.Int64Value(int64(o.KeepCount)),
		KeepTime:  types.Int64Value(int64(o.KeepTime)),
		Start:     types.StringValue(o.Start),
		state:     attr.ValueStateKnown,
	}
	result, d := objV.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(ConfigsMaintenanceBackupDifferentialValue{}.AttributeTypes(ctx)), diags
	}

	return result, nil
}

func FlattenClusterConfigsMaintenanceCronTabs(ctx context.Context, o []clusters.ClusterConfigMaintenanceCronTabs) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ListNull(ConfigsMaintenanceCrontabsValue{}.Type(ctx)), nil
	}

	cronTabsV := make([]attr.Value, len(o))
	for i, s := range o {
		settings, d := FlattenClusterConfigsSettings(ctx, s.Settings)
		diags.Append(d...)
		if diags.HasError() {
			return types.ListUnknown(ConfigsMaintenanceCrontabsValue{}.Type(ctx)), diags
		}

		cronTabsV[i] = ConfigsMaintenanceCrontabsValue{
			Name:     types.StringValue(s.Name),
			Required: types.BoolValue(s.Required),
			Settings: settings,
			Start:    types.StringValue(s.Start),
			state:    attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(ConfigsMaintenanceCrontabsValue{}.Type(ctx), cronTabsV)

	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(ConfigsMaintenanceCrontabsValue{}.Type(ctx)), diags
	}

	return result, nil
}

func FlattenClusterPodGroupsNodeProcesses(ctx context.Context, o []string) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ListNull(types.StringType), nil
	}

	nodeProcessesV := make([]attr.Value, len(o))
	for i, n := range o {
		nodeProcessesV[i] = types.StringValue(n)
	}
	result, d := types.ListValue(types.StringType, nodeProcessesV)

	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(types.StringType), diags
	}

	return result, nil
}

func FlattenClusterPodGroupsResource(ctx context.Context, o *clusters.ClusterPodGroupResource) (basetypes.ObjectValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ObjectNull(PodGroupsResourceValue{}.AttributeTypes(ctx)), nil
	}

	ram_request := strings.ReplaceAll(o.RAMRequest, "Gi", "")
	ram_limit := strings.ReplaceAll(o.RAMLimit, "Gi", "")
	resourceV := PodGroupsResourceValue{
		CpuLimit:   types.StringValue(o.CPULimit),
		CpuRequest: types.StringValue(o.CPURequest),
		RamLimit:   types.StringValue(ram_limit),
		RamRequest: types.StringValue(ram_request),
		state:      attr.ValueStateKnown,
	}
	result, d := resourceV.ToObjectValue(ctx)
	diags.Append(d...)
	if diags.HasError() {
		return types.ObjectUnknown(PodGroupsResourceValue{}.AttributeTypes(ctx)), diags
	}

	return result, nil
}

func FlattenClusterPodGroupsVolumes(ctx context.Context, o map[string]clusters.ClusterPodGroupVolume) (basetypes.MapValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.MapNull(PodGroupsVolumesValue{}.Type(ctx)), nil
	}

	volumesV := make(map[string]attr.Value, len(o))
	for i, v := range o {
		storage := strings.ReplaceAll(v.Storage, "Gi", "")
		volumesV[i] = PodGroupsVolumesValue{
			Count:            types.Int64Value(int64(v.Count)),
			Storage:          types.StringValue(storage),
			StorageClassName: types.StringValue(v.StorageClassName),
			state:            attr.ValueStateKnown,
		}
	}
	result, d := types.MapValue(PodGroupsVolumesValue{}.Type(ctx), volumesV)

	diags.Append(d...)
	if diags.HasError() {
		return types.MapUnknown(PodGroupsVolumesValue{}.Type(ctx)), diags
	}

	return result, nil
}

func ExpandClusterConfigs(ctx context.Context, v ConfigsValue) (*clusters.ClusterConfig, diag.Diagnostics) {
	result := &clusters.ClusterConfig{}

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

	if o := v.Warehouses; !o.IsUnknown() && !o.IsNull() {
		warehouses, diags := ExpandClusterWarehouses(ctx, o)
		if diags != nil && diags.HasError() {
			return nil, diags
		}
		result.Warehouses = warehouses
	}

	return result, nil
}

func ExpandClusterConfigsMaintenance(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterConfigMaintenance, diag.Diagnostics) {
	maintenanceObjV, diags := ConfigsMaintenanceType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	maintenance := maintenanceObjV.(ConfigsMaintenanceValue)
	result := clusters.ClusterConfigMaintenance{
		Start: maintenance.Start.ValueString(),
	}
	return &result, nil
}

func ExpandClusterConfigsSettings(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterConfigSetting, diag.Diagnostics) {
	settingsV := make([]ConfigsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settingsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterConfigSetting, len(settingsV))
	for i, s := range settingsV {
		result[i] = clusters.ClusterConfigSetting{
			Alias: s.Alias.ValueString(),
			Value: s.Value.ValueString(),
		}
	}
	return result, nil
}

func ExpandClusterWarehouses(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterConfigWarehouse, diag.Diagnostics) {
	warehousesV := make([]ConfigsWarehousesValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &warehousesV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterConfigWarehouse, len(warehousesV))
	for i, v := range warehousesV {
		result[i] = clusters.ClusterConfigWarehouse{
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

func ExpandClusterWarehousesConnections(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterConfigWarehouseConnection, diag.Diagnostics) {
	connectionsV := make([]ConfigsWarehousesConnectionsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &connectionsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterConfigWarehouseConnection, len(connectionsV))
	for i, v := range connectionsV {
		result[i] = clusters.ClusterConfigWarehouseConnection{
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

func ExpandClusterWarehousesConnectionsSettings(ctx context.Context, o basetypes.ListValue) ([]clusters.ClusterConfigSetting, diag.Diagnostics) {
	settingsV := make([]ConfigsWarehousesConnectionsSettingsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &settingsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterConfigSetting, len(settingsV))
	for i, v := range settingsV {
		result[i] = clusters.ClusterConfigSetting{
			Alias: v.Alias.ValueString(),
			Value: v.Value.ValueString(),
		}
	}
	return result, nil
}

func ExpandClusterPodGroups(ctx context.Context, template *templates.ClusterTemplate, o basetypes.ListValue) ([]clusters.ClusterPodGroup, diag.Diagnostics) {
	podGroupsV := make([]PodGroupsValue, 0, len(o.Elements()))
	diags := o.ElementsAs(ctx, &podGroupsV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make([]clusters.ClusterPodGroup, len(podGroupsV))
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

		result[i] = clusters.ClusterPodGroup{
			Count:              int(v.Count.ValueInt64()),
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

func ExpandClusterPodGroupsResource(ctx context.Context, o basetypes.ObjectValue) (*clusters.ClusterPodGroupResource, diag.Diagnostics) {
	resourceV, diags := PodGroupsResourceType{}.ValueFromObject(ctx, o)
	if diags.HasError() {
		return nil, diags
	}

	resource := resourceV.(PodGroupsResourceValue)
	result := clusters.ClusterPodGroupResource{
		CPURequest: resource.CpuRequest.ValueString(),
		RAMRequest: resource.RamRequest.ValueString(),
	}
	return &result, nil
}

func ExpandClusterPodGroupsVolumes(ctx context.Context, o basetypes.MapValue) (map[string]clusters.ClusterPodGroupVolume, diag.Diagnostics) {
	volumesV := make(map[string]PodGroupsVolumesValue)
	diags := o.ElementsAs(ctx, &volumesV, false)
	if diags.HasError() {
		return nil, diags
	}

	result := make(map[string]clusters.ClusterPodGroupVolume)
	for i, v := range volumesV {
		result[i] = clusters.ClusterPodGroupVolume{
			StorageClassName: v.StorageClassName.ValueString(),
			Storage:          v.Storage.ValueString(),
			Count:            int(v.Count.ValueInt64()),
		}
	}
	return result, nil
}
