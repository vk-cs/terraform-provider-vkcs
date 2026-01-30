package kubernetes_v2

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/nodegroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	rkubengv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/resource_kubernetes_node_group_v2"
)

const (
	createNgDelayV2        = 3
	createNgPollIntervalV2 = 30
	updateNgDelayV2        = 3
	updateNgPollIntervalV2 = 30
	deleteNgDelayV2        = 5
	deleteNgPollIntervalV2 = 30
)

var (
	_ resource.Resource                = (*kubernetesNodeGroupV2Resource)(nil)
	_ resource.ResourceWithConfigure   = (*kubernetesNodeGroupV2Resource)(nil)
	_ resource.ResourceWithImportState = (*kubernetesNodeGroupV2Resource)(nil)
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

func (r *kubernetesNodeGroupV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	// Validate node group configuration
	if !r.validateNodeGroup(ctx, region, &data, &resp.Diagnostics) {
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// Build create options
	createOpts, diags := rkubengv2.ToCreateOpts(ctx, &data, data.ClusterId.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to create node group", map[string]interface{}{
		"cluster_id": fmt.Sprintf("%#v", data.ClusterId.ValueString()),
		"createOpts": fmt.Sprintf("%#v", createOpts),
	})

	nodeGroupID, err := nodegroups.Create(client, *createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes node group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to create node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", nodeGroupID),
		"cluster_id":    fmt.Sprintf("%#v", data.ClusterId.ValueString()),
	})

	data.Id = types.StringValue(nodeGroupID)

	// Wait for node group to become active
	stateConf := r.getStateConfForNodeGroupCreate(client, data.ClusterId.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes node group to become active", err.Error())
		return
	}

	// Read the created node group
	readNgDiags := r.readNodeGroup(ctx, client, &data)
	resp.Diagnostics.Append(readNgDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesNodeGroupV2Resource) validateNodeGroup(ctx context.Context, region string, data *rkubengv2.KubernetesNodeGroupV2Model, diags *diag.Diagnostics) bool {
	scaleType := data.ScaleType.ValueString()

	switch scaleType {
	case "fixed_scale":
		if !data.AutoScaleMinSize.IsNull() && !data.AutoScaleMinSize.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"auto_scale_min_size cannot be set when scale_type is 'fixed_scale'",
			)
			return false
		}
		if !data.AutoScaleMaxSize.IsNull() && !data.AutoScaleMaxSize.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"auto_scale_max_size cannot be set when scale_type is 'fixed_scale'",
			)
			return false
		}
		if !data.AutoScaleNodeCount.IsNull() && !data.AutoScaleNodeCount.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"auto_scale_node_count cannot be set when scale_type is 'fixed_scale'",
			)
			return false
		}

		if data.FixedScaleNodeCount.IsNull() || data.FixedScaleNodeCount.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"fixed_scale_node_count is required when scale_type is 'fixed_scale'",
			)
			return false
		}
	case "auto_scale":
		if !data.FixedScaleNodeCount.IsNull() && !data.FixedScaleNodeCount.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"fixed_scale_node_count cannot be set when scale_type is 'auto_scale'",
			)
			return false
		}

		// Проверяем, что все auto_scale поля заданы и не 0
		if data.AutoScaleMinSize.IsNull() || data.AutoScaleMinSize.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"auto_scale_min_size is required when scale_type is 'auto_scale'",
			)
			return false
		}
		if data.AutoScaleMaxSize.IsNull() || data.AutoScaleMaxSize.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"auto_scale_max_size is required when scale_type is 'auto_scale'",
			)
			return false
		}
		if data.AutoScaleNodeCount.IsNull() || data.AutoScaleNodeCount.IsUnknown() {
			diags.AddError(
				"Invalid Configuration",
				"auto_scale_node_count is required when scale_type is 'auto_scale'",
			)
			return false
		}

		// Проверяем условие min <= node_count <= max
		minSize := int(data.AutoScaleMinSize.ValueInt64())
		maxSize := int(data.AutoScaleMaxSize.ValueInt64())
		nodeCount := int(data.AutoScaleNodeCount.ValueInt64())

		if minSize > nodeCount || nodeCount > maxSize {
			diags.AddError(
				"Invalid Configuration",
				fmt.Sprintf("for auto_scale, condition 'auto_scale_min_size <= auto_scale_node_count <= auto_scale_max_size' must be met (min=%d, count=%d, max=%d)", minSize, nodeCount, maxSize),
			)
			return false
		}
	default:
		diags.AddError(
			"Invalid Configuration",
			fmt.Sprintf("scale_type must be either 'fixed_scale' or 'auto_scale', got %s", scaleType),
		)
		return false
	}

	// Validate node taints
	if !data.Taints.IsNull() && !data.Taints.IsUnknown() {
		taintDiags := r.validateTaints(ctx, data.Taints)
		if taintDiags.HasError() {
			diags.Append(taintDiags...)
			return false
		}
	}

	if region != "" {
		availabilityZone := data.AvailabilityZone.ValueString()
		d := r.validateAZ(region, availabilityZone)
		if d.HasError() {
			diags.Append(d...)
			return false
		}
	}

	return true
}

func (r *kubernetesNodeGroupV2Resource) validateTaints(ctx context.Context, taintsSet types.Set) (diags diag.Diagnostics) {
	if taintsSet.IsNull() || taintsSet.IsUnknown() {
		return
	}

	elements := taintsSet.Elements()
	if len(elements) == 0 {
		return
	}

	seen := make(map[string]bool)
	for _, elem := range elements {
		objValuable, ok := elem.(basetypes.ObjectValuable)
		if !ok {
			continue
		}

		objValue, objDiags := objValuable.ToObjectValue(ctx)
		if objDiags.HasError() {
			diags.Append(objDiags...)
			continue
		}

		attrs := objValue.Attributes()

		keyStr, keyOk := attrs["key"].(basetypes.StringValue)
		effectStr, effectOk := attrs["effect"].(basetypes.StringValue)
		if !keyOk || !effectOk {
			continue
		}

		if keyStr.IsNull() || keyStr.IsUnknown() || effectStr.IsNull() || effectStr.IsUnknown() {
			continue
		}

		key := keyStr.ValueString()
		effect := effectStr.ValueString()
		keyEffect := key + ":" + effect

		if seen[keyEffect] {
			diags.AddError(
				"Invalid Taints Configuration",
				fmt.Sprintf("duplicate taint with key '%s' and effect '%s' found; each combination of key and effect must be unique.", key, effect),
			)
			return
		}
		seen[keyEffect] = true

	}

	return
}

func (r *kubernetesNodeGroupV2Resource) validateAZ(region string, ngAZ string) (diags diag.Diagnostics) {
	if region != "" {
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		diags.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	res := clusters.GetClusterAZs(client)
	if res.Err != nil {
		diags.AddError("Error reading Kubernetes cluster V2 availability zones", res.Err.Error())
		return
	}

	listAZs, err := res.Extract()
	if err != nil {
		diags.AddError("Error extracting Kubernetes cluster V2 availability zones from response", err.Error())
		return
	}

	// Проверяем, что указанная зона доступна для кластеров
	found := false
	for _, az := range listAZs.AZs {
		if az == ngAZ {
			found = true
			break
		}
	}

	if !found {
		diags.AddError(
			"Invalid Availability Zone",
			fmt.Sprintf("availability zone '%s' is not available for Kubernetes clusters; available zones: %#v", ngAZ, listAZs),
		)
		return
	}

	return
}

func (r *kubernetesNodeGroupV2Resource) getStateConfForNodeGroupCreate(client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Reconciling,
		},
		Target: []string{
			clusterStatusV2Running,
		},
		Refresh: func() (interface{}, string, error) {
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, "", err
			}
			return cluster, cluster.Status, nil
		},
		Timeout:      30 * time.Minute,
		Delay:        createNgDelayV2 * time.Second,
		PollInterval: createNgPollIntervalV2 * time.Second,
	}
}

func (r *kubernetesNodeGroupV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	readNgDiags := r.readNodeGroup(ctx, client, &data)
	resp.Diagnostics.Append(readNgDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesNodeGroupV2Resource) readNodeGroup(ctx context.Context, client *gophercloud.ServiceClient, data *rkubengv2.KubernetesNodeGroupV2Model) (diags diag.Diagnostics) {
	tflog.Trace(ctx, "Calling Container Infra V2 API to get node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", data.Id.ValueString()),
		"cluster_id":    fmt.Sprintf("%#v", data.ClusterId.ValueString()),
	})

	nodeGroup, err := nodegroups.Get(client, data.Id.ValueString()).Extract()
	if err != nil {
		if !util.CheckDeletedStatus(err) {
			tflog.Error(ctx, "Error reading Kubernetes node group V2", map[string]interface{}{"error": err.Error()})
			diags.AddError("Error reading Kubernetes node group V2", err.Error())
			return
		}
		data.Id = types.StringNull()
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to get node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", data.Id.ValueString()),
		"cluster_id":    fmt.Sprintf("%#v", data.ClusterId.ValueString()),
		"node_group":    fmt.Sprintf("%#v", nodeGroup),
	})

	// Use converter function to convert node group to model
	convertedModel, d := rkubengv2.ToNodeGroupModel(ctx, nodeGroup)
	diags.Append(d...)
	if diags.HasError() {
		return
	}

	*data = convertedModel

	return
}

func (r *kubernetesNodeGroupV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rkubengv2.KubernetesNodeGroupV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	// Validate node group configuration
	if !r.validateNodeGroup(ctx, region, &plan, &resp.Diagnostics) {
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	// Check what needs to be updated
	updateOpts := nodegroups.UpdateOpts{}

	// Update node flavor if changed
	if !plan.NodeFlavor.Equal(state.NodeFlavor) {
		updateOpts.VMEngine = &nodegroups.VMEngine{
			NovaEngine: nodegroups.NovaEngine{
				FlavorID: plan.NodeFlavor.ValueString(),
			},
		}
	}

	// Update scale specification if changed
	if !plan.ScaleType.Equal(state.ScaleType) ||
		!plan.FixedScaleNodeCount.Equal(state.FixedScaleNodeCount) ||
		!plan.AutoScaleNodeCount.Equal(state.AutoScaleNodeCount) ||
		!plan.AutoScaleMinSize.Equal(state.AutoScaleMinSize) ||
		!plan.AutoScaleMaxSize.Equal(state.AutoScaleMaxSize) {

		scaleType := plan.ScaleType.ValueString()
		scaleSpec := nodegroups.ScaleSpec{}

		switch scaleType {
		case "fixed_scale":
			scaleSpec.FixedScale = &nodegroups.FixedScale{
				Size: int(plan.FixedScaleNodeCount.ValueInt64()),
			}
		case "auto_scale":
			scaleSpec.AutoScale = &nodegroups.AutoScale{
				MinSize: int(plan.AutoScaleMinSize.ValueInt64()),
				MaxSize: int(plan.AutoScaleMaxSize.ValueInt64()),
				Size:    int(plan.AutoScaleNodeCount.ValueInt64()),
			}
		}

		updateOpts.ScaleSpec = &scaleSpec
	}

	// Update labels if changed
	if !plan.Labels.Equal(state.Labels) {
		labels := make(map[string]string)
		if !plan.Labels.IsNull() {
			elements := make(map[string]types.String, len(plan.Labels.Elements()))
			resp.Diagnostics.Append(plan.Labels.ElementsAs(ctx, &elements, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			for k, v := range elements {
				labels[k] = v.ValueString()
			}
		}
		updateOpts.Labels = labels
	}

	// Update taints if changed
	if !plan.Taints.Equal(state.Taints) {
		var taints []nodegroups.Taint
		if !plan.Taints.IsNull() {
			var taintValues []rkubengv2.TaintsValue
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
		updateOpts.Taints = taints
	}

	// Update parallel upgrade chunk if changed
	if !plan.ParallelUpgradeChunk.Equal(state.ParallelUpgradeChunk) {
		if plan.ParallelUpgradeChunk.ValueInt64() != 0 {
			parallelUpgradeChunk := int(plan.ParallelUpgradeChunk.ValueInt64())
			updateOpts.ParallelUpgradeChunk = &parallelUpgradeChunk
		} else {
			// Set to null/zero
			zero := 0
			updateOpts.ParallelUpgradeChunk = &zero
		}
	}

	// Only update if there are changes
	if updateOpts.VMEngine == nil && updateOpts.ScaleSpec == nil &&
		updateOpts.Labels == nil && updateOpts.Taints == nil &&
		updateOpts.ParallelUpgradeChunk == nil {

		readNgDiags := r.readNodeGroup(ctx, client, &plan)
		resp.Diagnostics.Append(readNgDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.Region = types.StringValue(region)

		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to update node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", plan.Id.ValueString()),
		"cluster_id":    fmt.Sprintf("%#v", plan.ClusterId.ValueString()),
		"updateOpts":    fmt.Sprintf("%#v", updateOpts),
	})

	err = nodegroups.Update(client, plan.Id.ValueString(), updateOpts)
	if err != nil {
		resp.Diagnostics.AddError("Error updating Kubernetes node group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to update node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", plan.Id.ValueString()),
		"cluster_id":    fmt.Sprintf("%#v", plan.ClusterId.ValueString()),
	})

	stateConf := r.getStateConfForNodeGroupUpdate(client, plan.ClusterId.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes node group to become active", err.Error())
		return
	}

	// Read the updated node group
	readNgDiags := r.readNodeGroup(ctx, client, &plan)
	if readNgDiags.HasError() {
		resp.Diagnostics.Append(readNgDiags...)
		return
	}

	plan.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *kubernetesNodeGroupV2Resource) getStateConfForNodeGroupUpdate(client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
	// Wait for node group to become active
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Reconciling,
		},
		Target: []string{
			clusterStatusV2Running,
		},
		Refresh: func() (interface{}, string, error) {
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				return nil, "", err
			}
			return cluster, cluster.Status, nil
		},
		Timeout:      30 * time.Minute,
		Delay:        updateNgDelayV2 * time.Second,
		PollInterval: updateNgPollIntervalV2 * time.Second,
	}
}

func (r *kubernetesNodeGroupV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data rkubengv2.KubernetesNodeGroupV2Model

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_node_group_v2", err.Error())
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to delete node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", data.Id.ValueString()),
		"cluster_id":    fmt.Sprintf("%#v", data.ClusterId.ValueString()),
	})

	err = nodegroups.Delete(client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Kubernetes node group", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to delete node group", map[string]interface{}{
		"node_group_id": fmt.Sprintf("%#v", data.Id.ValueString()),
		"cluster_id":    fmt.Sprintf("%#v", data.ClusterId.ValueString()),
	})

	stateConf := r.getStateConfForNodeGroupDelete(client, data.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes node group to be deleted", err.Error())
		return
	}
}

func (r *kubernetesNodeGroupV2Resource) getStateConfForNodeGroupDelete(client *gophercloud.ServiceClient, nodeGroupID string) *retry.StateChangeConf {
	// Wait for node group to be deleted
	return &retry.StateChangeConf{
		Pending: []string{
			"EXISTS",
			"ERROR",
		},
		Target: []string{
			"DELETED",
		},
		Refresh: func() (interface{}, string, error) {
			nodeGroup, err := nodegroups.Get(client, nodeGroupID).Extract()
			if err != nil {
				if util.CheckDeletedStatus(err) {
					return nodeGroup, "DELETED", nil
				}
				return nil, "ERROR", err
			}
			return nodeGroup, "EXISTS", nil
		},
		Timeout:      30 * time.Minute,
		Delay:        deleteNgDelayV2 * time.Second,
		PollInterval: deleteNgPollIntervalV2 * time.Second,
	}
}

func (r *kubernetesNodeGroupV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
