package kubernetes

import (
	"context"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/addons"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	rkubeclusteraddonv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/resource_kubernetes_cluster_addon_v2"
)

var (
	_ resource.Resource                = (*kubernetesClusterAddonV2Resource)(nil)
	_ resource.ResourceWithConfigure   = (*kubernetesClusterAddonV2Resource)(nil)
	_ resource.ResourceWithImportState = (*kubernetesClusterAddonV2Resource)(nil)
)

const (
	clusterAddonStatusV2Installing = "INSTALLING"
	clusterAddonStatusV2Installed  = "INSTALLED"
	clusterAddonStatusV2Updating   = "UPDATING"
	clusterAddonStatusV2Deleting   = "DELETING"
	clusterAddonStatusV2Deleted    = "DELETED"
	clusterAddonStatusV2Error      = "ERROR"
)

const (
	clusterAddonV2DefaultTimeout = time.Minute * 20
)

func NewKubernetesClusterAddonV2Resource() resource.Resource {
	return &kubernetesClusterAddonV2Resource{}
}

type kubernetesClusterAddonV2Resource struct {
	config clients.Config
}

func (r *kubernetesClusterAddonV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_cluster_addon_v2"
}

func (r *kubernetesClusterAddonV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rkubeclusteraddonv2.KubernetesClusterAddonV2ResourceSchema(ctx)
}

func (r *kubernetesClusterAddonV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.config = req.ProviderData.(clients.Config)
}

// We use ModifyPlan, because importing resource has problem with field 'updated_at'. Field is computed and it doesn't have plan_modifier = UseStateForUnknown,
// so it leads to error: After the apply operation, the provider still indicated an unknown value... All values must be known after apply, so this is always a bug in the provider
// and should be reported in the provider's own repository. Terraform will still save the other known object values in the state.
func (r *kubernetesClusterAddonV2Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.State.Raw.IsNull() || req.Plan.Raw.IsNull() {
		return
	}

	var state, plan rkubeclusteraddonv2.KubernetesClusterAddonV2Model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	valuesHaveNotChanged := state.Values.Equal(plan.Values)
	stateIsNotNull := !state.UpdatedAt.IsNull()
	planIsUnknown := plan.UpdatedAt.IsUnknown()

	if valuesHaveNotChanged && stateIsNotNull && planIsUnknown {
		plan.UpdatedAt = state.UpdatedAt
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

func (r *kubernetesClusterAddonV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rkubeclusteraddonv2.KubernetesClusterAddonV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_addon_v2", err.Error())
		return
	}

	// Build create options
	createOpts := rkubeclusteraddonv2.ToCreateOpts(ctx, data)

	// Make API call
	clusterAddonID, err := addons.CreateClusterAddon(client, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating cluster addon", err.Error())
		return
	}

	// Set the ID immediately
	data.Id = types.StringValue(clusterAddonID.ID)

	// Parse timeout for operation
	createTimeout, diags := data.Timeouts.Create(ctx, clusterAddonV2DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wait for cluster to become active
	stateConf := r.getStateConfForClusterAddonCreate(client, createTimeout, clusterAddonID.ID)
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster to become active", err.Error())
		return
	}

	// Read the cluster addon from Managed K8S API to populate fields
	clusterAddon, err := addons.GetClusterAddon(client, clusterAddonID.ID).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading cluster addons", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromClusterAddon(ctx, &clusterAddon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterAddonV2Resource) getStateConfForClusterAddonCreate(client *gophercloud.ServiceClient, createTimeout time.Duration, clusterAddonID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterAddonStatusV2Installing,
		},
		Target: []string{
			clusterAddonStatusV2Installed,
			clusterAddonStatusV2Error,
		},
		Refresh:      kubernetesAddonStateRefreshFuncV2(client, clusterAddonID),
		Timeout:      createTimeout,
		Delay:        time.Second * 30,
		PollInterval: time.Second * 10,
	}
}

func (r *kubernetesClusterAddonV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rkubeclusteraddonv2.KubernetesClusterAddonV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_addon_v2", err.Error())
		return
	}

	// Read the cluster sec policy from Managed K8S API to populate fields
	clusterAddon, err := addons.GetClusterAddon(client, data.Id.ValueString()).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading cluster addon", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromClusterAddon(ctx, &clusterAddon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesClusterAddonV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rkubeclusteraddonv2.KubernetesClusterAddonV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_addon_v2", err.Error())
		return
	}

	if plan.Values.Equal(state.Values) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Check what needs to be updated
	updateOpts := rkubeclusteraddonv2.ToUpdateOpts(ctx, plan)

	// API call
	err = addons.UpdateClusterAddon(client, &updateOpts)
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error updating vkcs_kubernetes_cluster_addon_v2", err.Error())
		return
	}

	// Parse timeout for operation
	updateTimeout, diags := plan.Timeouts.Update(ctx, clusterV2DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wait for update to complete
	stateConf := r.getStateConfForClusterAddonUpdate(client, updateTimeout, plan.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster addon update to complete", err.Error())
		return
	}

	clusterAddon, err := addons.GetClusterAddon(client, plan.Id.ValueString()).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading cluster addon", err.Error())
		return
	}

	// Read the cluster addon from Managed K8S API to populate fields
	resp.Diagnostics.Append(plan.UpdateFromClusterAddon(ctx, &clusterAddon)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *kubernetesClusterAddonV2Resource) getStateConfForClusterAddonUpdate(client *gophercloud.ServiceClient, updateTimeout time.Duration, clusterAddonID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterAddonStatusV2Updating,
		},
		Target: []string{
			clusterAddonStatusV2Installed,
			clusterAddonStatusV2Error,
		},
		Refresh:      kubernetesAddonStateRefreshFuncV2(client, clusterAddonID),
		Timeout:      updateTimeout,
		Delay:        time.Second * 30,
		PollInterval: time.Second * 10,
	}
}

func (r *kubernetesClusterAddonV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data rkubeclusteraddonv2.KubernetesClusterAddonV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_cluster_addon_v2", err.Error())
		return
	}

	// API call
	err = addons.DeleteClusterAddon(client, data.ClusterId.ValueString(), data.Id.ValueString())
	if err != nil {
		if errutil.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting cluster addon", err.Error())
		return
	}

	// Parse timeout for operation
	deleteTimeout, diags := data.Timeouts.Delete(ctx, clusterAddonV2DefaultTimeout)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	// Wait for deletion to complete
	stateConf := r.getStateConfForClusterAddonDelete(client, deleteTimeout, data.Id.ValueString())
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for Kubernetes cluster deletion to complete", err.Error())
		return
	}
}

func (r *kubernetesClusterAddonV2Resource) getStateConfForClusterAddonDelete(client *gophercloud.ServiceClient, deleteTimeout time.Duration, clusterAddonID string) *retry.StateChangeConf {
	return &retry.StateChangeConf{
		Pending: []string{
			clusterAddonStatusV2Deleting,
		},
		Target: []string{
			clusterAddonStatusV2Deleted,
			clusterAddonStatusV2Error,
		},
		Refresh:      kubernetesAddonStateRefreshFuncV2(client, clusterAddonID),
		Timeout:      deleteTimeout,
		Delay:        time.Second * 30,
		PollInterval: time.Second * 10,
	}
}

func (r *kubernetesClusterAddonV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
