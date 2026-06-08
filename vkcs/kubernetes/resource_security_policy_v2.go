package kubernetes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/kubernetes/containerinfra/v2/secpolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	rkubesecpolicyv2 "github.com/vk-cs/terraform-provider-vkcs/vkcs/kubernetes/resource_kubernetes_security_policy_v2"
)

var (
	_ resource.Resource                = (*kubernetesSecurityPolicyV2Resource)(nil)
	_ resource.ResourceWithConfigure   = (*kubernetesSecurityPolicyV2Resource)(nil)
	_ resource.ResourceWithImportState = (*kubernetesSecurityPolicyV2Resource)(nil)
)

func NewSecurityPolicyV2Resource() resource.Resource {
	return &kubernetesSecurityPolicyV2Resource{}
}

type kubernetesSecurityPolicyV2Resource struct {
	config clients.Config
}

func (r *kubernetesSecurityPolicyV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_kubernetes_security_policy_v2"
}

func (r *kubernetesSecurityPolicyV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = rkubesecpolicyv2.KubernetesSecurityPolicyV2ResourceSchema(ctx)
}

func (r *kubernetesSecurityPolicyV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.config = req.ProviderData.(clients.Config)
}

func (r *kubernetesSecurityPolicyV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data rkubesecpolicyv2.KubernetesSecurityPolicyV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_security_policy_v2", err.Error())
		return
	}

	// Build create options
	createOpts := rkubesecpolicyv2.ToCreateOpts(ctx, data)

	// Make API call
	clusterSecPolicyID, err := secpolicies.Create(client, createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating cluster security policy", err.Error())
		return
	}

	// Set the ID immediately
	data.Id = types.StringValue(clusterSecPolicyID)

	// Read the cluster sec policy from Managed K8S API to populate fields
	clusterSecPolicy, err := secpolicies.Get(client, clusterSecPolicyID).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading cluster security policy", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromClusterSecPolicy(ctx, clusterSecPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesSecurityPolicyV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data rkubesecpolicyv2.KubernetesSecurityPolicyV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_security_policy_v2", err.Error())
		return
	}

	// Read the cluster sec policy from Managed K8S API to populate fields
	clusterSecPolicy, err := secpolicies.Get(client, data.Id.ValueString()).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading cluster security policy", err.Error())
		return
	}

	resp.Diagnostics.Append(data.UpdateFromClusterSecPolicy(ctx, clusterSecPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *kubernetesSecurityPolicyV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state rkubesecpolicyv2.KubernetesSecurityPolicyV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_security_policy_v2", err.Error())
		return
	}

	if plan.PolicySettings.Equal(state.PolicySettings) && plan.Namespace.Equal(state.Namespace) && plan.Enabled.Equal(state.Enabled) {
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
		return
	}

	// Check what needs to be updated
	updateOpts := secpolicies.UpdateOpts{
		ClusterID:               plan.ClusterId.ValueString(),
		ClusterSecurityPolicyID: plan.Id.ValueString(),
		PolicySettings:          plan.PolicySettings.ValueString(),
		Namespace:               plan.Namespace.ValueString(),
		Enabled:                 plan.Enabled.ValueBool(),
	}

	// API call
	updatedClusterSecPolicy, err := secpolicies.Update(client, updateOpts).Extract()
	if err != nil {
		if errutil.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error updating vkcs_kubernetes_security_policy_v2", err.Error())
		return
	}

	// Read the cluster sec policy from Managed K8S API to populate fields
	resp.Diagnostics.Append(plan.UpdateFromClusterSecPolicy(ctx, updatedClusterSecPolicy)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *kubernetesSecurityPolicyV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data rkubesecpolicyv2.KubernetesSecurityPolicyV2Model

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
		resp.Diagnostics.AddError("Error creating API client for resource vkcs_kubernetes_security_policy_v2", err.Error())
		return
	}

	// API call
	err = secpolicies.Delete(client, data.ClusterId.ValueString(), data.Id.ValueString())
	if err != nil {
		if errutil.IsNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Error deleting cluster security policy", err.Error())
		return
	}
}

func (r *kubernetesSecurityPolicyV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
