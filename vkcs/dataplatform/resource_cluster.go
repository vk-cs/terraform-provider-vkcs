package dataplatform

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/dataplatform/resource_cluster"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/dataplatform/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"

	"fmt"
	"time"
)

const (
	clusterCreateDelay      = 5 * time.Second
	clusterCreateMinTimeout = 5 * time.Second
	clusterCreateTimeout    = 60 * time.Minute
	clusterUpdateDelay      = 5 * time.Second
	clusterUpdateMinTimeout = 5 * time.Second
	clusterUpdateTimeout    = 60 * time.Minute
	clusterDeleteDelay      = 5 * time.Second
	clusterDeleteMinTimeout = 5 * time.Second
	clusterDeleteTimeout    = 60 * time.Minute
)

type clusterStatus string

const (
	clusterStatusCreating        clusterStatus = "InfraUpdating"
	clusterStatusConfiguring     clusterStatus = "Configuring"
	clusterStatusScaling         clusterStatus = "Scaling"
	clusterStatusUpdating        clusterStatus = "Updating"
	clusterStatusActive          clusterStatus = "Active"
	clusterStatusWaitingDeleting clusterStatus = "Waiting deleting"
	clusterStatusDeleting        clusterStatus = "Deleting"
	clusterStatusDeleted         clusterStatus = "Deleted"
)

var (
	_ resource.Resource = (*clusterResource)(nil)
)

func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

type clusterResource struct {
	config clients.Config
}

func (r *clusterResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dataplatform_cluster"
}

func (r *clusterResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_cluster.ClusterResourceSchema(ctx)
}

func (r *clusterResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_cluster.ClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.DataPlatformClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS dataplatform client", err.Error())
		return
	}

	configOpts, diags := resource_cluster.ExpandClusterConfigs(ctx, data.Configs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterTemplate, diags := data.GetClusterTemplate(client)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	podGroupOpts, diags := resource_cluster.ExpandClusterPodGroups(ctx, clusterTemplate, data.PodGroups)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createOpts := clusters.ClusterCreate{
		Name:              data.Name.ValueString(),
		Description:       data.Description.ValueString(),
		ClusterTemplateID: clusterTemplate.ID,
		NetworkID:         data.NetworkId.ValueString(),
		SubnetID:          data.SubnetId.ValueString(),
		ProductName:       data.ProductName.ValueString(),
		ProductVersion:    data.ProductVersion.ValueString(),
		AvailabilityZone:  data.AvailabilityZone.ValueString(),
		Configs:           configOpts,
		PodGroups:         podGroupOpts,
		StackID:           data.StackId.ValueString(),
		FloatingIPPool:    data.FloatingIpPool.ValueString(),
	}

	tflog.Trace(ctx, "Calling Data Platform API to create cluster", map[string]interface{}{"opts": fmt.Sprintf("%#v", createOpts)})

	clusterShort, err := clusters.Create(client, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Data Platform API to create cluster", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Data Platform API to create cluster", map[string]interface{}{"cluster": fmt.Sprintf("%#v", clusterShort)})

	id := types.StringValue(clusterShort.ID)
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{string(clusterStatusCreating), string(clusterStatusConfiguring)},
		Target:     []string{string(clusterStatusActive)},
		Refresh:    clusterStateRefreshFunc(client, clusterShort.ID),
		Timeout:    clusterCreateTimeout,
		Delay:      clusterCreateDelay,
		MinTimeout: clusterCreateMinTimeout,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for cluster creation", err.Error())
	}

	cluster, err := clusters.Get(client, clusterShort.ID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error getting cluster", err.Error())
	}

	data.ClusterTemplateId = types.StringValue(clusterTemplate.ID)

	resp.Diagnostics.Append(data.UpdateState(ctx, cluster, data.Configs, &resp.State)...)
}

func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_cluster.ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.DataPlatformClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS dataplatform client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "cluster_id", id)

	tflog.Trace(ctx, "Calling Data Platform API to retrieve cluster")

	cluster, err := clusters.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling Data Platform to retrieve cluster", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Data Platform API to retrieve cluster", map[string]interface{}{"cluster": fmt.Sprintf("%#v", cluster)})

	resp.Diagnostics.Append(data.UpdateState(ctx, cluster, data.Configs, &resp.State)...)
}

func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resource_cluster.ClusterModel
	var data resource_cluster.ClusterModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	client, err := r.config.DataPlatformClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS dataplatform client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "cluster_id", id)

	updateOpts := clusters.ClusterUpdate{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}

	tflog.Trace(ctx, "Calling Data Platform API to update cluster", map[string]interface{}{"opts": fmt.Sprintf("%#v", updateOpts)})

	_, err = clusters.Update(client, id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Data Platform API to update cluster", err.Error())
		return
	}

	if !plan.Configs.Maintenance.IsUnknown() && !plan.Configs.Maintenance.IsNull() {
		diags := clusterUpdateConfigsMaintenance(ctx, client, id, data.Configs.Maintenance, plan.Configs.Maintenance)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !plan.Configs.Settings.IsUnknown() && !plan.Configs.Settings.IsNull() {
		diags := clusterUpdateConfigsSettings(ctx, client, id, plan.Configs.Settings)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !plan.Configs.Users.IsUnknown() && !plan.Configs.Users.IsNull() {
		diags := clusterUpdateConfigsUsers(ctx, client, id, data.Configs.Users, plan.Configs.Users)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !plan.Configs.Warehouses.IsUnknown() && !plan.Configs.Warehouses.IsNull() {
		diags := clusterUpdateConfigsWarehouses(ctx, client, id, data.Configs.Warehouses, plan.Configs.Warehouses)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	if !plan.PodGroups.IsUnknown() && !plan.PodGroups.IsNull() {
		diags := clusterUpdatePodGroups(ctx, client, id, data.PodGroups, plan.PodGroups)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	cluster, err := clusters.Get(client, id).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error getting cluster", err.Error())
		return
	}
	tflog.Trace(ctx, "Called Data Platform API to update cluster", map[string]interface{}{"cluster": fmt.Sprintf("%#v", cluster)})

	resp.Diagnostics.Append(data.UpdateState(ctx, cluster, plan.Configs, &resp.State)...)
}

func clusterUpdateConfigsMaintenance(ctx context.Context, client *gophercloud.ServiceClient, id string, currentMaintenance basetypes.ObjectValue, planMaintenance basetypes.ObjectValue) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	updateMaintenanceOpts, d := resource_cluster.BuildUpdateClusterConfigsMaintenance(ctx, currentMaintenance, planMaintenance)
	if d.HasError() {
		return d
	}

	if updateMaintenanceOpts != nil {
		tflog.Trace(ctx, "Calling Data Platform API to update cluster maintenance", map[string]interface{}{"opts": fmt.Sprintf("%#v", planMaintenance)})

		resp, err := clusters.UpdateMaintenance(client, id, updateMaintenanceOpts).Extract()
		if err != nil {
			diags.AddError("Error calling Data Platform API to update cluster maintenance", err.Error())
			return diags
		}

		tflog.Info(ctx, "Cluster maintenance update request sent", map[string]interface{}{"resp": resp})

		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, id),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update configs.maintenance", err.Error())
			return diags
		}
	}

	return nil
}

func clusterUpdateConfigsSettings(ctx context.Context, client *gophercloud.ServiceClient, id string, settings basetypes.ListValue) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	updateSettingsOpts := make([]clusters.ClusterUpdateSetting, 0)
	planSettings, d := resource_cluster.ExpandClusterConfigsSettings(ctx, settings)
	if d.HasError() {
		return d
	}
	for _, planSetting := range planSettings {
		updateSettingsOpts = append(updateSettingsOpts, clusters.ClusterUpdateSetting(planSetting))
	}
	if len(updateSettingsOpts) > 0 {
		tflog.Trace(ctx, "Calling Data Platform API to update cluster settings", map[string]interface{}{"opts": fmt.Sprintf("%#v", updateSettingsOpts)})
		_, err := clusters.UpdateSettings(client, id, &clusters.ClusterUpdateSettings{Settings: updateSettingsOpts}).Extract()
		if err != nil {
			diags.AddError("Error calling Data Platform API to update cluster settings", err.Error())
			return diags
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, id),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update configs.settings", err.Error())
			return diags
		}
	}
	return nil
}

func clusterUpdateConfigsUsers(ctx context.Context, client *gophercloud.ServiceClient, id string, dataUsersRaw basetypes.ListValue, planUsersRaw basetypes.ListValue) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	dataUsers, d := resource_cluster.ReadClusterConfigsUsers(ctx, dataUsersRaw)
	if d.HasError() {
		return d
	}

	planUsers, d := resource_cluster.ReadClusterConfigsUsers(ctx, planUsersRaw)
	if d.HasError() {
		return d
	}

	var usersToAdd []clusters.ClusterUpdateUser
	var usersToDelete []string
	remainingUsers := make(map[string]bool)
	userIDs := make(map[string]string)

	for _, user := range dataUsers {
		remainingUsers[user.Username.ValueString()] = false
		userIDs[user.Username.ValueString()] = user.Id.ValueString()
	}

	for _, planUser := range planUsers {
		if _, ok := remainingUsers[planUser.Username.ValueString()]; !ok {
			usersToAdd = append(usersToAdd, clusters.ClusterUpdateUser{
				Username: planUser.Username.ValueString(),
				Password: planUser.Password.ValueString(),
				Role:     planUser.Role.ValueString(),
			})
		} else {
			remainingUsers[planUser.Username.ValueString()] = true
		}
	}

	for userName, isRemaining := range remainingUsers {
		if !isRemaining {
			usersToDelete = append(usersToDelete, userIDs[userName])
		}
	}

	if len(usersToAdd) > 0 {
		tflog.Trace(ctx, "Calling Data Platform API to add cluster users", map[string]interface{}{"opts": fmt.Sprintf("%#v", usersToAdd)})
		_, err := clusters.AddClusterUsers(client, id, &clusters.ClusterUpdateUsers{Users: usersToAdd}).Extract()
		if err != nil {
			diags.AddError("Error calling Data Platform API to add cluster users", err.Error())
			return diags
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, id),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update configs.users", err.Error())
			return diags
		}
	}

	if len(usersToDelete) > 0 {
		tflog.Trace(ctx, "Calling Data Platform API to delete cluster users", map[string]interface{}{"opts": fmt.Sprintf("%#v", usersToDelete)})
		err := clusters.DeleteClusterUsers(client, id, &clusters.ClusterDeleteUsers{ClusterUsersIDs: usersToDelete}).ExtractErr()
		if err != nil {
			diags.AddError("Error calling Data Platform API to delete cluster users", err.Error())
			return diags
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, id),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update configs.users", err.Error())
			return diags
		}
	}
	return nil
}

func clusterUpdateConfigsWarehouses(ctx context.Context, client *gophercloud.ServiceClient, clusterID string, dataWarehousesRaw basetypes.ListValue, planWarehousesRaw basetypes.ListValue) diag.Diagnostics {
	var d diag.Diagnostics

	dataWarehouses, d := resource_cluster.ReadClusterConfigsWarehouses(ctx, dataWarehousesRaw)
	if d.HasError() {
		return d
	}

	planWarehouses, d := resource_cluster.ReadClusterConfigsWarehouses(ctx, planWarehousesRaw)
	if d.HasError() {
		return d
	}

	warehouses := make(map[string]resource_cluster.ConfigsWarehousesValue)
	for _, dataWarehouse := range dataWarehouses {
		warehouses[dataWarehouse.Name.ValueString()] = dataWarehouse
	}

	for _, planWarehouse := range planWarehouses {
		dataConnections := warehouses[planWarehouse.Name.ValueString()].Connections
		if !dataConnections.Equal(planWarehouse.Connections) {
			d := clusterUpdateConfigsWarehousesConnections(ctx, client, clusterID, warehouses[planWarehouse.Name.ValueString()].Id.ValueString(), dataConnections, planWarehouse.Connections)
			if d.HasError() {
				return d
			}
		}
	}

	return nil
}

func clusterUpdateConfigsWarehousesConnections(ctx context.Context, client *gophercloud.ServiceClient, clusterID string, warehouseID string, dataConnectionsRaw basetypes.ListValue, planConnectionsRaw basetypes.ListValue) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	var dataConnections []resource_cluster.ConfigsWarehousesConnectionsValue
	var planConnections []resource_cluster.ConfigsWarehousesConnectionsValue

	if !dataConnectionsRaw.IsNull() && !dataConnectionsRaw.IsUnknown() {
		dataConnections, d = resource_cluster.ReadClusterConfigsWarehousesConnections(ctx, dataConnectionsRaw)
		if d.HasError() {
			return d
		}
	}

	if !planConnectionsRaw.IsNull() && !planConnectionsRaw.IsUnknown() {
		planConnections, d = resource_cluster.ReadClusterConfigsWarehousesConnections(ctx, planConnectionsRaw)
		if d.HasError() {
			return d
		}
	}

	var connectionsToAdd []clusters.ClusterAddConnection
	var connectionsToDelete []string
	remainingConnections := make(map[string]bool)
	connectionIDs := make(map[string]string)

	for _, connection := range dataConnections {
		remainingConnections[connection.Name.ValueString()] = false
		connectionIDs[connection.Name.ValueString()] = connection.Id.ValueString()
	}

	for _, planConnection := range planConnections {
		if _, ok := remainingConnections[planConnection.Name.ValueString()]; !ok {
			settings, d := resource_cluster.ReadClusterConfigsWarehousesConnectionsSettings(ctx, planConnection.Settings)
			if d.HasError() {
				return d
			}

			planSettings := make([]clusters.ClusterAddConnectionSetting, len(settings))
			for i, setting := range settings {
				planSettings[i] = clusters.ClusterAddConnectionSetting{
					Alias: setting.Alias.ValueString(),
					Value: setting.Value.ValueString(),
				}
			}

			connectionsToAdd = append(connectionsToAdd, clusters.ClusterAddConnection{
				Plug:     planConnection.Plug.ValueString(),
				Name:     planConnection.Name.ValueString(),
				Settings: planSettings,
			})
		} else {
			remainingConnections[planConnection.Name.ValueString()] = true
		}
	}

	for connectionName, isRemaining := range remainingConnections {
		if !isRemaining {
			connectionsToDelete = append(connectionsToDelete, connectionIDs[connectionName])
		}
	}

	if len(connectionsToAdd) > 0 {
		tflog.Trace(ctx, "Calling Data Platform API to add cluster connections", map[string]interface{}{"opts": fmt.Sprintf("%#v", connectionsToAdd)})
		_, err := clusters.AddClusterConfigsWarehouseConnections(client, clusterID, &clusters.ClusterAddWarehouseConnections{
			WarehouseID: warehouseID,
			Connections: connectionsToAdd,
		}).Extract()
		if err != nil {
			diags.AddError("Error calling Data Platform API to add cluster connections", err.Error())
			return diags
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, clusterID),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update", err.Error())
			return diags
		}
	}

	if len(connectionsToDelete) > 0 {
		tflog.Trace(ctx, "Calling Data Platform API to delete cluster connections", map[string]interface{}{"opts": fmt.Sprintf("%#v", connectionsToDelete)})
		err := clusters.DeleteClusterConnections(client, clusterID, &clusters.ClusterDeleteWarehouseConnections{
			ClusterConnections: connectionsToDelete,
			WarehouseID:        warehouseID,
		}).ExtractErr()
		if err != nil {
			diags.AddError("Error calling Data Platform API to delete cluster connections", err.Error())
			return diags
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, clusterID),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update", err.Error())
			return diags
		}
	}
	return nil
}

func clusterUpdatePodGroups(ctx context.Context, client *gophercloud.ServiceClient, id string, dataPodGroupsRaw basetypes.ListValue, planPodGroupsRaw basetypes.ListValue) diag.Diagnostics {
	var diags diag.Diagnostics
	var d diag.Diagnostics

	planPodGroups, d := resource_cluster.ReadClusterPodGroups(ctx, planPodGroupsRaw)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}
	dataPodGroups, d := resource_cluster.ReadClusterPodGroups(ctx, dataPodGroupsRaw)
	diags.Append(d...)
	if diags.HasError() {
		return diags
	}

	podGroupsByName := make(map[string]resource_cluster.PodGroupsValue)
	for _, dataPodGroup := range dataPodGroups {
		podGroupsByName[dataPodGroup.Name.ValueString()] = dataPodGroup
	}

	var podGroupsToUpdate []clusters.ClusterUpdatePodGroup

	for _, planPodGroup := range planPodGroups {
		dataPodGroup := podGroupsByName[planPodGroup.Name.ValueString()]

		planVolumes, d := resource_cluster.ReadClusterPodGroupVolumes(ctx, planPodGroup.Volumes)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}
		dataVolumes, d := resource_cluster.ReadClusterPodGroupVolumes(ctx, dataPodGroup.Volumes)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		oldCount := int(dataPodGroup.Count.ValueInt64())
		newCount := int(planPodGroup.Count.ValueInt64())

		changed := oldCount != newCount

		volumes := make(map[string]clusters.ClusterUpdatePodGroupVolume)
		for volumeName, planVolume := range planVolumes {
			dataVolume, exists := dataVolumes[volumeName]

			if exists && planVolume.Storage.ValueString() != dataVolume.Storage.ValueString() {
				changed = true
			}

			volumes[volumeName] = clusters.ClusterUpdatePodGroupVolume{
				StorageClassName: planVolume.StorageClassName.ValueString(),
				Storage:          planVolume.Storage.ValueString(),
				Count:            int(planVolume.Count.ValueInt64()),
			}
		}

		if !changed {
			continue
		}

		podGroupID := dataPodGroup.Id.ValueString()

		resource, d := resource_cluster.ReadClusterPodGroupResources(ctx, planPodGroup.Resource)
		diags.Append(d...)
		if diags.HasError() {
			return diags
		}

		podGroupToUpdate := clusters.ClusterUpdatePodGroup{
			ID:    podGroupID,
			Count: &newCount,
			Resource: clusters.ClusterUpdatePodGroupResource{
				CPURequest: resource.CpuRequest.ValueString(),
				RAMRequest: resource.RamRequest.ValueString(),
			},
		}

		if len(volumes) > 0 {
			podGroupToUpdate.Volumes = volumes
		}

		podGroupsToUpdate = append(podGroupsToUpdate, podGroupToUpdate)
	}

	if len(podGroupsToUpdate) > 0 {
		tflog.Trace(ctx, "Calling Data Platform API to update cluster pod groups", map[string]interface{}{"opts": fmt.Sprintf("%#v", podGroupsToUpdate)})
		_, err := clusters.UpdateClusterPodGroup(client, id, &clusters.ClusterUpdatePodGroups{PodGroups: podGroupsToUpdate}).Extract()
		if err != nil {
			diags.AddError("Error calling Data Platform API to update cluster pod groups", err.Error())
			return diags
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{string(clusterStatusConfiguring), string(clusterStatusUpdating), string(clusterStatusScaling)},
			Target:     []string{string(clusterStatusActive)},
			Refresh:    clusterStateRefreshFunc(client, id),
			Timeout:    clusterUpdateTimeout,
			Delay:      clusterUpdateDelay,
			MinTimeout: clusterUpdateMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			diags.AddError("Error waiting for cluster update podGroups", err.Error())
			return diags
		}
	}

	return nil
}

func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_cluster.ClusterModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.DataPlatformClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS dataplatform client", err.Error())
		return
	}

	id := data.Id.ValueString()
	ctx = tflog.SetField(ctx, "cluster_id", id)

	tflog.Trace(ctx, "Calling Data Platform API to delete cluster")

	err = clusters.Delete(client, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error calling Data Platform API to delete cluster", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Data Platform to delete cluster")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{string(clusterStatusCreating), string(clusterStatusActive), string(clusterStatusDeleting), string(clusterStatusWaitingDeleting), string(clusterStatusConfiguring)},
		Target:     []string{string(clusterStatusDeleted)},
		Refresh:    clusterStateRefreshFunc(client, id),
		Timeout:    clusterDeleteTimeout,
		Delay:      clusterDeleteDelay,
		MinTimeout: clusterDeleteMinTimeout,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for cluster deletion", err.Error())
	}
}

func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
