package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	rkubeclusterv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/resource_kubernetes_cluster_v2"
)

const (
	regionalClusterMastersMin = 3
	regionalClusterMastersMax = 5
)

const (
	clusterStatusV2Provisioning = "PROVISIONING"
	clusterStatusV2Running      = "RUNNING"
	clusterStatusV2Reconciling  = "RECONCILING"
	clusterStatusV2Deleting     = "DELETING"
	clusterStatusV2Deleted      = "DELETED"
)

const (
	clusterV2DefaultTimeout = time.Minute * 60
)

var (
	_ resource.Resource                   = (*kubernetesClusterV2Resource)(nil)
	_ resource.ResourceWithConfigure      = (*kubernetesClusterV2Resource)(nil)
	_ resource.ResourceWithImportState    = (*kubernetesClusterV2Resource)(nil)
	_ resource.ResourceWithValidateConfig = (*kubernetesClusterV2Resource)(nil)
)

func NewKubernetesClusterV2Resource() resource.Resource {
	return &kubernetesClusterV2Resource{}
}

type kubernetesClusterV2Resource struct {
	config clients.Config
}

func (r *kubernetesClusterV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_v2"
}

func (r *kubernetesClusterV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rkubeclusterv2.KubernetesClusterV2ResourceSchema(ctx)
}

func (r *kubernetesClusterV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.config = req.ProviderData.(clients.Config)
}

func (r *kubernetesClusterV2Resource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	if req.Config.Raw.IsNull() {
		return
	}

	// Read Terraform configuration data into the model
	var data rkubeclusterv2.KubernetesClusterV2Model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.validateControlPlaneConfig(ctx, &data, &resp.Diagnostics)
	r.validateNetworkConfig(ctx, &data, &resp.Diagnostics)
}

// validateControlPlaneConfig verifies number of master nodes for standard and regional cluster, checks azs count match
func (r *kubernetesClusterV2Resource) validateControlPlaneConfig(_ context.Context, plan *rkubeclusterv2.KubernetesClusterV2Model, diags *diag.Diagnostics) {
	clusterType := plan.ClusterType.ValueString()
	masterCount := plan.MasterCount.ValueInt64()
	numberOfAZs := len(plan.AvailabilityZones.Elements())

	errSummary := "Invalid control plane configuration"

	switch clusterType {
	case clusterTypeRegional:
		if masterCount != regionalClusterMastersMin && masterCount != regionalClusterMastersMax {
			diags.AddAttributeError(
				path.Root("master_count"),
				errSummary,
				fmt.Sprintf("Regional cluster requires %d or %d master nodes, got %d", regionalClusterMastersMin, regionalClusterMastersMax, masterCount),
			)
		}
		if numberOfAZs != 3 {
			diags.AddAttributeError(
				path.Root("availability_zones"),
				errSummary,
				fmt.Sprintf("Regional cluster requires 3 availability zones, got %d", numberOfAZs),
			)
		}
	case clusterTypeStandard:
		if masterCount != 1 && masterCount != 3 && masterCount != 5 {
			diags.AddAttributeError(
				path.Root("master_count"),
				errSummary,
				fmt.Sprintf("Standard cluster requires 1, 3, or 5 master nodes, got %d", masterCount),
			)
		}
		if numberOfAZs != 1 {
			diags.AddAttributeError(
				path.Root("availability_zones"),
				errSummary,
				fmt.Sprintf("Standard cluster requires 1 availability zone, got %d", numberOfAZs),
			)
		}
	}
}

// validateNetworkConfig checks whether 'external_network_id' exists in the config when 'public_ip' is True.
func (r *kubernetesClusterV2Resource) validateNetworkConfig(_ context.Context, plan *rkubeclusterv2.KubernetesClusterV2Model, diags *diag.Diagnostics) {
	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if plan.ExternalNetworkId.IsUnknown() {
		return
	}

	if plan.PublicIp.ValueBool() && plan.ExternalNetworkId.ValueString() == "" {
		diags.AddAttributeError(
			path.Root("external_network_id"),
			"Missing attribute 'external_network_id'",
			"Attribute 'external_network_id' must be specified when 'public_ip' is true",
		)
	}
}

func (r *kubernetesClusterV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rkubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// Build create options
	createOpts, diags := rkubeclusterv2.ToCreateOpts(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Make API call
	clusterID, err := clusters.Create(client, createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes cluster", err.Error())
		return
	}

	// Set the ID immediately
	data.Id = types.StringValue(clusterID)

	// Parse timeout for operation
	createTimeout, diags := data.Timeouts.Create(ctx, clusterV2DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wait for cluster to become active
	stateConf := r.getStateConfForClusterCreate(client, createTimeout, clusterID)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster to become active", err.Error())
		return
	}

	// Read the cluster from Managed K8S API to populate fields
	apiCluster, kubeconfig, diags := r.readCluster(client, clusterID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// It means that cluster was not found
	if apiCluster == nil && kubeconfig == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.UpdateFromCluster(ctx, apiCluster, kubeconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterV2Resource) getStateConfForClusterCreate(client *gophercloud.ServiceClient, createTimeout time.Duration, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Provisioning,
			clusterStatusV2Reconciling,
		},
		Target: []string{
			clusterStatusV2Running,
		},
		Refresh:      kubernetesStateRefreshFuncV2(client, clusterID),
		Timeout:      createTimeout,
		Delay:        time.Minute * 5,
		PollInterval: time.Second * 20,
	}
}

func (r *kubernetesClusterV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rkubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// Read the cluster from Managed K8S API to populate fields
	apiCluster, kubeconfig, diags := r.readCluster(client, data.Id.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// It means that cluster was not found
	if apiCluster == nil && kubeconfig == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(data.UpdateFromCluster(ctx, apiCluster, kubeconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterV2Resource) readCluster(client *gophercloud.ServiceClient, clusterID string) (cluster *clusters.Cluster, clusterKubeconfig *string, diags diag.Diagnostics) {
	// Get cluster data
	cluster, err := clusters.Get(client, clusterID).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			return nil, nil, diags
		}
		diags.AddError("Error reading Kubernetes cluster", err.Error())
		return
	}

	// Get cluster kubeconfig
	kubeconfig, err := clusters.GetKubeconfig(client, clusterID)
	if err != nil {
		diags.AddError("Error retrieving cluster kubeconfig", err.Error())
		return
	}

	clusterKubeconfig = &kubeconfig
	return
}

func (r *kubernetesClusterV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rkubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// Check what needs to be updated
	upgradeOpts := clusters.UpgradeOpts{}
	scaleOpts := clusters.ScaleOpts{}
	hasChanged := false

	// Update version if changed
	if !plan.Version.Equal(state.Version) {
		hasChanged = true
		upgradeOpts.Version = plan.Version.ValueString()
	}

	// Update master count if changed
	if !plan.MasterFlavor.Equal(state.MasterFlavor) {
		hasChanged = true
		masterSpec := clusters.MasterSpecOpts{
			Engine: clusters.MasterEngineOpts{
				NovaEngine: clusters.NovaEngineOpts{
					FlavorID: plan.MasterFlavor.ValueString(),
				},
			},
			Replicas: int(plan.MasterCount.ValueInt64()),
		}
		scaleOpts.MasterSpec = masterSpec
	}

	if !hasChanged {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	clusterID := plan.Id.ValueString()

	if upgradeOpts.Version != "" {
		// API call
		err = clusters.Upgrade(client, clusterID, upgradeOpts)
		if err != nil {
			resp.Diagnostics.AddError("Error upgrading resource vkcs_kubernetes_cluster_v2", err.Error())
			return
		}
	}

	if scaleOpts.MasterSpec.Engine.NovaEngine.FlavorID != "" {
		// API call
		err = clusters.Scale(client, clusterID, scaleOpts)
		if err != nil {
			resp.Diagnostics.AddError("Error scaling resource vkcs_kubernetes_cluster_v2", err.Error())
			return
		}
	}

	// Parse timeout for operation
	updateTimeout, diags := plan.Timeouts.Update(ctx, clusterV2DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wait for update to complete
	stateConf := r.getStateConfForClusterUpdate(client, updateTimeout, clusterID)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster update to complete", err.Error())
		return
	}

	// Read the cluster from Managed K8S API to populate fields
	apiCluster, kubeconfig, diags := r.readCluster(client, clusterID)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// It means that cluster was not found
	if apiCluster == nil && kubeconfig == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(plan.UpdateFromCluster(ctx, apiCluster, kubeconfig)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *kubernetesClusterV2Resource) getStateConfForClusterUpdate(client *gophercloud.ServiceClient, updateTimeout time.Duration, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Reconciling,
		},
		Target: []string{
			clusterStatusV2Running,
		},
		Refresh:      kubernetesStateRefreshFuncV2(client, clusterID),
		Timeout:      updateTimeout,
		Delay:        time.Minute * 3,
		PollInterval: time.Second * 20,
	}
}

func (r *kubernetesClusterV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data rkubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// API call
	err = clusters.Delete(client, data.Id.ValueString())
	if err != nil {
		if errutil.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting Kubernetes cluster", err.Error())
		return
	}

	// Parse timeout for operation
	deleteTimeout, diags := data.Timeouts.Delete(ctx, clusterV2DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wait for deletion to complete
	stateConf := r.getStateConfForClusterDelete(client, deleteTimeout, data.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster deletion to complete", err.Error())
		return
	}
}

func (r *kubernetesClusterV2Resource) getStateConfForClusterDelete(client *gophercloud.ServiceClient, deleteTimeout time.Duration, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Deleting,
		},
		Target: []string{
			clusterStatusV2Deleted,
		},
		Refresh:      kubernetesStateRefreshFuncV2(client, clusterID),
		Timeout:      deleteTimeout,
		Delay:        time.Minute * 3,
		PollInterval: time.Second * 20,
	}
}

func (r *kubernetesClusterV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
