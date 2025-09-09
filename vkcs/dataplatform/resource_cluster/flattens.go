package resource_cluster

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"

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
		Start: types.StringValue(o.Start),
		state: attr.ValueStateKnown,
	}
	if o.Enabled != nil {
		objV.Enabled = types.BoolValue(*o.Enabled)
	}
	if o.KeepCount != nil {
		objV.KeepCount = types.Int64Value(int64(*o.KeepCount))
	}
	if o.KeepTime != nil {
		objV.KeepTime = types.Int64Value(int64(*o.KeepTime))
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

func FlattenClusterInfo(ctx context.Context, i *clusters.ClusterInfo) (InfoValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if i == nil {
		return NewInfoValueNull(), nil
	}

	services, d := FlattenClusterInfoServices(ctx, i.Services)
	diags.Append(d...)
	if diags.HasError() {
		return NewInfoValueNull(), nil
	}

	infoV := InfoValue{
		Services: services,
		state:    attr.ValueStateKnown,
	}

	return infoV, nil
}

func FlattenClusterInfoServices(ctx context.Context, o []clusters.ClusterInfoServices) (basetypes.ListValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	if o == nil {
		return types.ListNull(InfoServicesValue{}.Type(ctx)), nil
	}

	servicesV := make([]attr.Value, len(o))
	for i, s := range o {
		servicesV[i] = InfoServicesValue{
			ConnectionString: types.StringValue(s.ConnectionString),
			Description:      types.StringValue(s.Description),
			Exposed:          types.BoolValue(s.Exposed),
			ServiceType:      types.StringValue(s.Type),
			state:            attr.ValueStateKnown,
		}
	}
	result, d := types.ListValue(InfoServicesValue{}.Type(ctx), servicesV)

	diags.Append(d...)
	if diags.HasError() {
		return types.ListUnknown(InfoServicesValue{}.Type(ctx)), diags
	}

	return result, nil
}
