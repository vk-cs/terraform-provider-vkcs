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
	return &result, nil
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

		result[i] = clusters.ClusterCreatePodGroup{
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
