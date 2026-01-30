package kubernetes

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	configvalidators "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/config_validators"
	rkubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/resource_kubernetes_node_group_v2"
)

const (
	createNgDelayV2        = 3 * time.Minute
	createNgPollIntervalV2 = 15 * time.Second

	deleteNgDelayV2        = 2 * time.Minute
	deleteNgPollIntervalV2 = 15 * time.Second

	updateNgDelayV2        = 30 * time.Second
	updateNgPollIntervalV2 = 15 * time.Second
)

var (
	_ resource.Resource                     = (*kubernetesNodeGroupV2Resource)(nil)
	_ resource.ResourceWithConfigure        = (*kubernetesNodeGroupV2Resource)(nil)
	_ resource.ResourceWithImportState      = (*kubernetesNodeGroupV2Resource)(nil)
	_ resource.ResourceWithConfigValidators = (*kubernetesNodeGroupV2Resource)(nil)
)

func NewKubernetesNodeGroupV2Resource() resource.Resource {
	return &kubernetesNodeGroupV2Resource{}
}

type kubernetesNodeGroupV2Resource struct {
	config clients.Config
}

func (r *kubernetesNodeGroupV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_node_group_v2"
}

func (r *kubernetesNodeGroupV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rkubengv2.KubernetesNodeGroupV2ResourceSchema(ctx)
}

func (r *kubernetesNodeGroupV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *kubernetesNodeGroupV2Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		&configvalidators.ScaleTypeConfigValidator{},
	}
}

func (r *kubernetesNodeGroupV2Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan rkubengv2.KubernetesNodeGroupV2Model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set default timeout if needed
	if util.IsNullOrUnknown(plan.Timeouts) {
		plan.Timeouts = rkubengv2.GetDefaultNgV2Timeouts(ctx)
	}
	if util.IsNullOrUnknown(plan.Timeouts.Create) {
		plan.Timeouts.Create = rkubengv2.GetDefaultNgV2CreateTimeout()
	}
	if util.IsNullOrUnknown(plan.Timeouts.Delete) {
		plan.Timeouts.Delete = rkubengv2.GetDefaultNgV2DeleteTimeout()
	}
	if util.IsNullOrUnknown(plan.Timeouts.Update) {
		plan.Timeouts.Update = rkubengv2.GetDefaultNgV2UpdateTimeout()
	}
	resp.Plan.Set(ctx, &plan)
}

func (r *kubernetesNodeGroupV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the region in which to obtain the Managed K8S client
	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	// Init Managed K8S client
	client, err := r.config.ManagedK8SClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// Build create options
	createOpts, diags := rkubengv2.ToCreateOpts(ctx, &data, data.ClusterId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API call
	nodeGroupID, err := nodegroups.Create(client, *createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes node group", err.Error())
		return
	}

	// Set the ID immediately
	data.Id = types.StringValue(nodeGroupID)

	// Parse timeout for operation
	createTimeout, err := time.ParseDuration(data.Timeouts.Create.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid duration format for create timeout. Expected format like '60m', '2h', '1h30m', '30s'", err.Error())
		return
	}

	// Wait for node group to become active
	stateConf := r.getStateConfForNodeGroupCreate(createTimeout, client, data.ClusterId.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes node group to become active", err.Error())
		return
	}

	// Read the node group from Managed K8S API to populate fields
	apiNodeGroup, diags := r.readNodeGroup(client, nodeGroupID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// It means that node group was not found
	if apiNodeGroup == nil {
		resp.State.RemoveResource(ctx)
	}

	resp.Diagnostics.Append(data.UpdateFromNodeGroup(ctx, apiNodeGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesNodeGroupV2Resource) getStateConfForNodeGroupCreate(createNgTimeoutV2 time.Duration, client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Reconciling,
		},
		Target: []string{
			clusterStatusV2Running,
		},
		Refresh:      kubernetesStateRefreshFuncV2(client, clusterID),
		Timeout:      createNgTimeoutV2,
		Delay:        createNgDelayV2,
		PollInterval: createNgPollIntervalV2,
	}
}

func (r *kubernetesNodeGroupV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the region in which to obtain the Managed K8S client
	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	data.Region = types.StringValue(region)

	// Init Managed K8S client
	client, err := r.config.ManagedK8SClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// Read the node group from Managed K8S API to populate fields
	apiNodeGroup, diags := r.readNodeGroup(client, data.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// It means that node group was not found
	if apiNodeGroup == nil {
		resp.State.RemoveResource(ctx)
	}

	resp.Diagnostics.Append(data.UpdateFromNodeGroup(ctx, apiNodeGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesNodeGroupV2Resource) readNodeGroup(client *gophercloud.ServiceClient, nodeGroupID string) (nodeGroup *nodegroups.NodeGroup, diags diag.Diagnostics) {
	apiNodeGroup, err := nodegroups.Get(client, nodeGroupID).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			return nil, diags
		}
		diags.AddError("Error reading Kubernetes node group", err.Error())
		return
	}

	nodeGroup = apiNodeGroup
	return
}

func (r *kubernetesNodeGroupV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rkubengv2.KubernetesNodeGroupV2Model

	// Read Terraform configuration data into the models
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the region in which to obtain the Managed K8S client
	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	plan.Region = types.StringValue(region)

	// Init Managed K8S client
	client, err := r.config.ManagedK8SClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// Check what needs to be updated
	updateOpts := nodegroups.UpdateOpts{}
	hasChanged := false

	// Update node flavor if changed
	if !plan.NodeFlavor.Equal(state.NodeFlavor) {
		hasChanged = true
		updateOpts.VMEngine = &nodegroups.VMEngine{
			NovaEngine: nodegroups.NovaEngine{
				FlavorID: plan.NodeFlavor.ValueString(),
			},
		}
	}

	// Update scale specification if changed
	if !plan.ScaleType.Equal(state.ScaleType) ||
		!plan.AutoScaleMinSize.Equal(state.AutoScaleMinSize) ||
		!plan.AutoScaleMaxSize.Equal(state.AutoScaleMaxSize) ||
		!plan.FixedScaleNodeCount.Equal(state.FixedScaleNodeCount) {
		hasChanged = true

		scaleType := plan.ScaleType.ValueString()
		scaleSpec := nodegroups.ScaleSpec{}

		switch scaleType {
		case "fixed_scale":
			scaleSpec.FixedScale = &nodegroups.FixedScale{
				Size: int(plan.FixedScaleNodeCount.ValueInt64()),
			}
		case "auto_scale":
			autoScaleNodeCount := state.AutoScaleNodeCount.ValueInt64()
			autoScaleMinSize := plan.AutoScaleMinSize.ValueInt64()
			autoScaleMaxSize := plan.AutoScaleMaxSize.ValueInt64()

			if autoScaleNodeCount < autoScaleMinSize {
				autoScaleNodeCount = autoScaleMinSize
			}
			if autoScaleNodeCount > autoScaleMaxSize {
				autoScaleNodeCount = autoScaleMaxSize
			}

			scaleSpec.AutoScale = &nodegroups.AutoScale{
				MinSize: int(autoScaleMinSize),
				MaxSize: int(autoScaleMaxSize),
				Size:    int(autoScaleNodeCount),
			}
		}

		updateOpts.ScaleSpec = &scaleSpec
	}

	// Update labels if changed
	if !plan.Labels.Equal(state.Labels) {
		hasChanged = true
		labels := make(map[string]string, len(plan.Labels.Elements()))
		if !plan.Labels.IsNull() {
			resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &labels, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}
		updateOpts.Labels = &labels
	}

	// Update taints if changed
	if !plan.Taints.Equal(state.Taints) {
		hasChanged = true
		taints := make([]nodegroups.Taint, 0, len(plan.Taints.Elements()))
		if !plan.Taints.IsNull() {
			taintValues := make([]rkubengv2.TaintsValue, 0, len(plan.Taints.Elements()))
			resp.Diagnostics.Append(plan.Taints.ElementsAs(ctx, &taintValues, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			for _, taint := range taintValues {
				taints = append(taints, nodegroups.Taint{
					Key:    taint.Key.ValueString(),
					Value:  taint.Value.ValueString(),
					Effect: taint.Effect.ValueString(),
				})
			}
		}
		updateOpts.Taints = &taints
	}

	// Update parallel upgrade chunk if changed
	if !plan.ParallelUpgradeChunk.Equal(state.ParallelUpgradeChunk) {
		hasChanged = true
		if plan.ParallelUpgradeChunk.ValueInt64() != 0 {
			parallelUpgradeChunk := int(plan.ParallelUpgradeChunk.ValueInt64())
			updateOpts.ParallelUpgradeChunk = &parallelUpgradeChunk
		} else {
			// Set to null/zero
			zero := 0
			updateOpts.ParallelUpgradeChunk = &zero
		}
	}

	if hasChanged {
		// API call
		err = nodegroups.Update(ctx, client, plan.Id.ValueString(), updateOpts)
		if err != nil {
			resp.Diagnostics.AddError("Error updating Kubernetes node group", err.Error())
			return
		}

		// Parse timeout for operation
		updateTimeout, err := time.ParseDuration(plan.Timeouts.Update.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Invalid duration format for update timeout. Expected format like '60m', '2h', '1h30m', '30s'", err.Error())
			return
		}

		stateConf := r.getStateConfForNodeGroupUpdate(updateTimeout, client, plan.ClusterId.ValueString())
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			resp.Diagnostics.AddError("Error waiting for Kubernetes node group to become active", err.Error())
			return
		}
	}

	// Read the node group from Managed K8S API to populate fields
	apiNodeGroup, diags := r.readNodeGroup(client, plan.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// It means that node group was not found
	if apiNodeGroup == nil {
		resp.State.RemoveResource(ctx)
	}

	resp.Diagnostics.Append(plan.UpdateFromNodeGroup(ctx, apiNodeGroup)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *kubernetesNodeGroupV2Resource) getStateConfForNodeGroupUpdate(updateNgTimeoutV2 time.Duration, client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Reconciling,
		},
		Target: []string{
			clusterStatusV2Running,
		},
		Refresh:      kubernetesStateRefreshFuncV2(client, clusterID),
		Timeout:      updateNgTimeoutV2,
		Delay:        updateNgDelayV2,
		PollInterval: updateNgPollIntervalV2,
	}
}

func (r *kubernetesNodeGroupV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	// Read Terraform configuration data into the models
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the region in which to obtain the Managed K8S client
	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	// Init Managed K8S client
	client, err := r.config.ManagedK8SClient(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// API call
	err = nodegroups.Delete(client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Kubernetes node group", err.Error())
		return
	}

	// Parse timeout for operation
	deleteTimeout, err := time.ParseDuration(data.Timeouts.Delete.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid duration format for delete timeout. Expected format like '60m', '2h', '1h30m', '30s'", err.Error())
		return
	}

	stateConf := r.getStateConfForNodeGroupDelete(deleteTimeout, client, data.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes node group to be deleted", err.Error())
		return
	}
}

func (r *kubernetesNodeGroupV2Resource) getStateConfForNodeGroupDelete(deleteNgTimeoutV2 time.Duration, client *gophercloud.ServiceClient, nodeGroupID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			string(nodeGroupStatusRunning),
			string(nodeGroupStatusError),
		},
		Target: []string{
			string(nodeGroupStatusNotFound),
		},
		Refresh: func() (interface{}, string, error) {
			nodeGroup, err := nodegroups.Get(client, nodeGroupID).Extract()
			if err != nil {
				if errutil.IsNotFound(err) {
					return nodeGroup, string(nodeGroupStatusNotFound), nil
				}
				return nil, string(nodeGroupStatusError), err
			}
			return nodeGroup, string(nodeGroupStatusRunning), nil
		},
		Timeout:      deleteNgTimeoutV2,
		Delay:        deleteNgDelayV2,
		PollInterval: deleteNgPollIntervalV2,
	}
}

func (r *kubernetesNodeGroupV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
