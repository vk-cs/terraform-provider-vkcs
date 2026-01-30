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
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	rkubeclusterv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/v2/resource_kubernetes_cluster_v2"
)

// timeouts of CRUD operations
const (
	operationCreateV2 = 60
	operationUpdateV2 = 60
	operationDeleteV2 = 30

	createClusterDelayV2        = 5
	createClusterPollIntervalV2 = 30
	updateClusterDelayV2        = 3
	updateClusterPollIntervalV2 = 20
	deleteClusterDelayV2        = 2
	deleteClusterPollIntervalV2 = 20
)

// cluster statuses from new API
const (
	clusterStatusV2Provisioning = "PROVISIONING"
	clusterStatusV2Starting     = "STARTING"
	clusterStatusV2Running      = "RUNNING"
	clusterStatusV2Reconciling  = "RECONCILING"
	clusterStatusV2Deleting     = "DELETING"
	clusterStatusV2Deleted      = "DELETED"
)

var (
	_ resource.Resource                = (*kubernetesClusterV2Resource)(nil)
	_ resource.ResourceWithConfigure   = (*kubernetesClusterV2Resource)(nil)
	_ resource.ResourceWithImportState = (*kubernetesClusterV2Resource)(nil)
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

func (r *kubernetesClusterV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rkubeclusterv2.KubernetesClusterV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	// Validate cluster configuration
	if !r.validateCluster(ctx, region, &data, &resp.Diagnostics) {
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// Build create options
	createOpts, diags := rkubeclusterv2.ToCreateOpts(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to create cluster", map[string]interface{}{"createOpts": fmt.Sprintf("%#v", createOpts)})

	clusterID, err := clusters.Create(client, *createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes cluster V2", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to create cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", clusterID)})

	// Set the ID immediately
	data.Id = types.StringValue(clusterID)

	// Wait for cluster to become active
	stateConf := r.getStateConfForClusterCreate(client, data.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster V2 to become active", err.Error())
		return
	}

	// Read the cluster to populate computed fields
	readClusterDiags := r.readCluster(ctx, client, &data)
	resp.Diagnostics.Append(readClusterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterV2Resource) getStateConfForClusterCreate(client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Provisioning,
			clusterStatusV2Starting,
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
		Timeout:    operationCreateV2 * time.Minute,
		Delay:      createClusterDelayV2 * time.Second,
		MinTimeout: createClusterPollIntervalV2 * time.Second,
	}
}

func (r *kubernetesClusterV2Resource) validateCluster(ctx context.Context, region string, data *rkubeclusterv2.KubernetesClusterV2Model, diags *diag.Diagnostics) bool {
	// Validate cluster_type and master_count
	clusterType := data.ClusterType.ValueString()
	masterCount := int(data.MasterCount.ValueInt64())

	switch clusterType {
	case "standard":
		if masterCount != 1 {
			diags.AddError(
				"Invalid Configuration",
				fmt.Sprintf("standard cluster requires exactly 1 master node, got %d", masterCount),
			)
			return false
		}
	case "regional":
		if masterCount != 3 && masterCount != 5 {
			diags.AddError(
				"Invalid Configuration",
				fmt.Sprintf("regional cluster requires 3 or 5 master nodes, got %d", masterCount),
			)
			return false
		}
	default:
		diags.AddError(
			"Invalid Configuration",
			fmt.Sprintf("cluster_type must be either 'standard' or 'regional', got %s", clusterType),
		)
		return false
	}

	// Validate availability_zones
	var availabilityZones []string
	if !data.AvailabilityZones.IsNull() {
		d := data.AvailabilityZones.ElementsAs(ctx, &availabilityZones, false)
		if d.HasError() {
			diags.Append(d...)
			return false
		}
	}

	if clusterType == "regional" && len(availabilityZones) != 3 {
		diags.AddError(
			"Invalid Configuration",
			fmt.Sprintf("regional cluster requires 3 availability zones, got %d", len(availabilityZones)),
		)
		return false
	}

	if clusterType == "standard" && len(availabilityZones) != 1 {
		diags.AddError(
			"Invalid Configuration",
			fmt.Sprintf("standard cluster requires exactly 1 availability zone, got %d", len(availabilityZones)),
		)
		return false
	}

	// Validate external network is specified when public IP is enabled
	if data.EnablePublicIp.ValueBool() && data.ExternalNetworkId.ValueString() == "" {
		diags.AddError(
			"Invalid Configuration",
			"external_network_id must be specified when enable_public_ip is true",
		)
		return false
	}

	if region != "" {
		// Validate specified values
		dAZs := r.validateAZs(region, availabilityZones)
		if dAZs.HasError() {
			diags.Append(dAZs...)
			return false
		}

		// Validate Kubernetes version of cluster
		dKubeVersion := r.validateKubeVersion(region, data.Version.ValueString())
		if dKubeVersion.HasError() {
			diags.Append(dKubeVersion...)
			return false
		}
	}

	return true
}

func (r *kubernetesClusterV2Resource) validateAZs(region string, clusterAZs []string) (diags diag.Diagnostics) {
	if region == "" {
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		diags.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_cluster_v2", err.Error())
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

	mk8sAZs := make(map[string]struct{}, len(listAZs.AZs))
	for _, az := range listAZs.AZs {
		mk8sAZs[az] = struct{}{}
	}

	for _, az := range clusterAZs {
		if _, ok := mk8sAZs[az]; !ok {
			allowedAZs := r.extractKeysFromMap(mk8sAZs)
			diags.AddError("Configuration of resource vkcs_kubernetes_cluster_v2 contains invalid availability zone", fmt.Sprintf("allowed %#v, got %#v", allowedAZs, clusterAZs))
		}
	}

	return
}

func (r *kubernetesClusterV2Resource) validateKubeVersion(region string, kubeVersion string) (diags diag.Diagnostics) {
	if region == "" {
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		diags.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	res := clusters.GetListK8SVersion(client)
	if res.Err != nil {
		diags.AddError("Error reading Kubernetes cluster V2 available versions", res.Err.Error())
		return
	}

	listKubeVersions, err := res.Extract()
	if err != nil {
		diags.AddError("Error extracting Kubernetes cluster V2 available versions from response", err.Error())
		return
	}

	mk8sKubeVersions := make(map[string]struct{}, len(listKubeVersions.Versions))
	for _, version := range listKubeVersions.Versions {
		mk8sKubeVersions[version.Version] = struct{}{}
	}

	if _, ok := mk8sKubeVersions[kubeVersion]; !ok {
		allowedVersions := r.extractKeysFromMap(mk8sKubeVersions)
		diags.AddError("Configuration of resource vkcs_kubernetes_cluster_v2 contains invalid version of Kubernetes", fmt.Sprintf("allowed %#v, got %#v", allowedVersions, kubeVersion))
	}

	return
}

func (r *kubernetesClusterV2Resource) extractKeysFromMap(m map[string]struct{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (r *kubernetesClusterV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rkubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	readClusterDiags := r.readCluster(ctx, client, &data)
	resp.Diagnostics.Append(readClusterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterV2Resource) readCluster(ctx context.Context, client *gophercloud.ServiceClient, data *rkubeclusterv2.KubernetesClusterV2Model) (diags diag.Diagnostics) {
	tflog.Trace(ctx, "Calling Container Infra V2 API to get cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", data.Id.ValueString())})

	clusterID := data.Id.ValueString()

	cluster, err := clusters.Get(client, clusterID).Extract()
	if err != nil {
		// Handle deleted cluster
		if !util.CheckDeletedStatus(err) {
			tflog.Error(ctx, "Error reading Kubernetes cluster V2", map[string]interface{}{"error": err.Error()})
			diags.AddError("Error reading Kubernetes cluster V2", err.Error())
			return
		}
		data.Id = types.StringNull()
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to get cluster", map[string]interface{}{
		"cluster_id": fmt.Sprintf("%#v", data.Id.ValueString()),
		"cluster":    fmt.Sprintf("%#v", cluster),
	})

	// Use converter function to convert cluster to model
	convertedModel, d := rkubeclusterv2.ToClusterModel(ctx, cluster)
	diags.Append(d...)
	if diags.HasError() {
		return
	}
	*data = convertedModel

	// Get cluster kubeconfig
	kubeconfig, err := clusters.GetKubeconfig(client, clusterID)
	if err != nil {
		diags.AddError("Error retrieving cluster kubeconfig", err.Error())
		return
	}
	data.K8sConfig = types.StringValue(kubeconfig)

	return
}

func (r *kubernetesClusterV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rkubeclusterv2.KubernetesClusterV2Model

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	// Validate cluster configuration
	if !r.validateCluster(ctx, region, &plan, &resp.Diagnostics) {
		return
	}

	// Validate simultaneous scaling and upgrade
	versionChanged := !plan.Version.Equal(state.Version)
	masterFlavorChanged := !plan.MasterFlavor.Equal(state.MasterFlavor)

	if versionChanged && masterFlavorChanged {
		resp.Diagnostics.AddError(
			"Invalid Update",
			"simultaneous scaling and cluster upgrade is not available",
		)
		return
	}

	client, err := r.config.ContainerInfraV2Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	// Check what needs to be updated
	upgradeOpts := clusters.UpgradeOpts{}
	scaleOpts := clusters.ScaleOpts{}

	// Update version if changed
	if !plan.Version.Equal(state.Version) {
		upgradeOpts.Version = plan.Version.ValueString()
	}

	// Update master count if changed
	if !plan.MasterFlavor.Equal(state.MasterFlavor) {
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

	// Only update if there are changes
	if upgradeOpts.Version == "" && scaleOpts.MasterSpec.Engine.NovaEngine.FlavorID == "" {
		tflog.Trace(ctx, "No allowed changes detected for vkcs_kubernetes_cluster_v2")

		diags := r.readCluster(ctx, client, &plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		plan.Region = types.StringValue(region)

		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return

	}

	if upgradeOpts.Version != "" {
		// ugrade cluster
		tflog.Trace(ctx, "Calling Container Infra V2 API to upgrade cluster", map[string]interface{}{
			"cluster_id":  fmt.Sprintf("%#v", plan.Id.ValueString()),
			"upgradeOpts": fmt.Sprintf("%#v", upgradeOpts),
		})

		err = clusters.Upgrade(client, plan.Id.ValueString(), upgradeOpts)
		if err != nil {
			resp.Diagnostics.AddError("error upgrading vkcs_kubernetes_cluster_v2", err.Error())
			return
		}

		tflog.Trace(ctx, "Called Container Infra V2 API to upgrade cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", plan.Id.ValueString())})
	}

	if scaleOpts.MasterSpec.Engine.NovaEngine.FlavorID != "" {
		// scale cluster
		tflog.Trace(ctx, "Calling Container Infra V2 API to scale cluster", map[string]interface{}{
			"cluster_id": fmt.Sprintf("%#v", plan.Id.ValueString()),
			"scaleOpts":  fmt.Sprintf("%#v", scaleOpts),
		})

		err = clusters.Scale(client, plan.Id.ValueString(), scaleOpts)
		if err != nil {
			resp.Diagnostics.AddError("error scaling vkcs_kubernetes_cluster_v2", err.Error())
			return
		}

		tflog.Trace(ctx, "Called Container Infra V2 API to scale cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", plan.Id.ValueString())})
	}

	// Wait for update to complete
	stateConf := r.getStateConfForClusterUpdate(client, plan.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster V2 update to complete", err.Error())
		return
	}

	// Read the updated cluster
	readClusterDiags := r.readCluster(ctx, client, &plan)
	resp.Diagnostics.Append(readClusterDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Region = types.StringValue(region)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *kubernetesClusterV2Resource) getStateConfForClusterUpdate(client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
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
		Timeout:    operationUpdateV2 * time.Minute,
		Delay:      updateClusterDelayV2 * time.Second,
		MinTimeout: updateClusterPollIntervalV2 * time.Second,
	}
}

func (r *kubernetesClusterV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data rkubeclusterv2.KubernetesClusterV2Model

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
		resp.Diagnostics.AddError("Error creating Container Infra V2 API client for resource vkcs_kubernetes_cluster_v2", err.Error())
		return
	}

	tflog.Trace(ctx, "Calling Container Infra V2 API to delete cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", data.Id.ValueString())})

	err = clusters.Delete(client, data.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting Kubernetes cluster V2", err.Error())
		return
	}

	tflog.Trace(ctx, "Called Container Infra V2 API to delete cluster", map[string]interface{}{"cluster_id": fmt.Sprintf("%#v", data.Id.ValueString())})

	// Wait for deletion to complete
	stateConf := r.getStateConfForClusterDelete(client, data.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster V2 deletion to complete", err.Error())
		return
	}
}

func (r *kubernetesClusterV2Resource) getStateConfForClusterDelete(client *gophercloud.ServiceClient, clusterID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterStatusV2Deleting,
			clusterStatusV2Running,
		},
		Target: []string{
			clusterStatusV2Deleted,
		},
		Refresh: func() (interface{}, string, error) {
			cluster, err := clusters.Get(client, clusterID).Extract()
			if err != nil {
				if util.CheckDeletedStatus(err) {
					return nil, clusterStatusV2Deleted, nil
				}
				return nil, "", err
			}
			return cluster, cluster.Status, nil
		},
		Timeout:    operationDeleteV2 * time.Minute,
		Delay:      deleteClusterDelayV2 * time.Second,
		MinTimeout: deleteClusterPollIntervalV2 * time.Second,
	}
}

func (r *kubernetesClusterV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
