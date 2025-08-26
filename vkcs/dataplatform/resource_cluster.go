package dataplatform

import (
	"context"

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
	clusterCreateTimeout    = 90 * time.Minute
	clusterUpdateDelay      = 5 * time.Second
	clusterUpdateMinTimeout = 5 * time.Second
	clusterUpdateTimeout    = 90 * time.Minute
	clusterDeleteDelay      = 5 * time.Second
	clusterDeleteMinTimeout = 5 * time.Second
	clusterDeleteTimeout    = 60 * time.Minute
)

type clusterStatus string

const (
	clusterStatusCreating        clusterStatus = "InfraUpdating"
	clusterStatusConfiguring     clusterStatus = "Configuring"
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

	resp.Diagnostics.Append(data.UpdateState(ctx, cluster, &resp.State, data.Configs.Settings)...)
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

	resp.Diagnostics.Append(data.UpdateState(ctx, cluster, &resp.State, data.Configs.Settings)...)
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

	if !plan.Configs.Settings.IsUnknown() && !plan.Configs.Settings.IsNull() {
		updateSettingsOpts := make([]clusters.ClusterUpdateSetting, 0)
		planSettings, diags := resource_cluster.ExpandClusterConfigsSettings(ctx, plan.Configs.Settings)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		for _, planSetting := range planSettings {
			updateSettingsOpts = append(updateSettingsOpts, clusters.ClusterUpdateSetting(planSetting))
		}
		if len(updateSettingsOpts) > 0 {
			_, err = clusters.UpdateSettings(client, id, &clusters.ClusterUpdateSettings{Settings: updateSettingsOpts}).Extract()
			if err != nil {
				resp.Diagnostics.AddError("Error calling Data Platform API to update cluster settings", err.Error())
				return
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
				resp.Diagnostics.AddError("Error waiting for cluster update", err.Error())
				return
			}
		}
	}
	cluster, err := clusters.Get(client, id).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error getting cluster", err.Error())
	}
	tflog.Trace(ctx, "Called Data Platform API to update cluster", map[string]interface{}{"cluster": fmt.Sprintf("%#v", cluster)})

	resp.Diagnostics.Append(data.UpdateState(ctx, cluster, &resp.State, plan.Configs.Settings)...)
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
