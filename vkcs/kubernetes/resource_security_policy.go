package kubernetes

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/containerinfra/v1/securitypolicies"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

var _ resource.Resource = &SecurityPolicyResource{}
var _ resource.ResourceWithConfigure = &SecurityPolicyResource{}

func NewSecurityPolicyResource() resource.Resource {
	return &SecurityPolicyResource{}
}

type SecurityPolicyResource struct {
	config clients.Config
}

type SecurityPolicyResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Region                   types.String `tfsdk:"region"`
	ClusterID                types.String `tfsdk:"cluster_id"`
	SecurityPolicyTemplateID types.String `tfsdk:"security_policy_template_id"`
	PolicySettings           types.String `tfsdk:"policy_settings"`
	Namespace                types.String `tfsdk:"namespace"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	CreatedAt                types.String `tfsdk:"created_at"`
	UpdatedAt                types.String `tfsdk:"updated_at"`
}

func (r *SecurityPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_kubernetes_security_policy"
}

func (r *SecurityPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The region in which to obtain the Container Infra client. If omitted, the `region` argument of the provider is used. Changing this creates a new security policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
						"require replacement if configuration value changes", "require replacement if configuration value changes"),
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"cluster_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the kubernetes cluster. Changing this creates a new security policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"security_policy_template_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the security policy template. Changing this creates a new security policy.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"policy_settings": schema.StringAttribute{
				Required:    true,
				Description: "Policy settings.",
			},

			"namespace": schema.StringAttribute{
				Required:    true,
				Description: "Namespace to apply security policy to.",
			},

			"enabled": schema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Controls whether the security policy is enabled. Default is true.",
			},

			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
			},

			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Update timestamp.",
			},
		},

		Description: "Provides a kubernetes cluster security policy resource. This can be used to create, modify and delete kubernetes security policies.",
	}
}

func (r *SecurityPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *SecurityPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SecurityPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if data.Region.IsUnknown() {
		region = r.config.GetRegion()
	}

	containerInfraClient, err := r.config.ContainerInfraV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes API client", err.Error())
		return
	}

	clusterID := data.ClusterID.ValueString()
	ctx = tflog.SetField(ctx, "cluster_id", clusterID)

	createOpts := securitypolicies.CreateOpts{
		ClusterID:                data.ClusterID.ValueString(),
		SecurityPolicyTemplateID: data.SecurityPolicyTemplateID.ValueString(),
		PolicySettings:           data.PolicySettings.ValueString(),
		Namespace:                data.Namespace.ValueString(),
		Enabled:                  data.Enabled.ValueBool(),
	}

	tflog.Debug(ctx, "Calling Kubernetes API to create security policy for the cluster", map[string]interface{}{"options": fmt.Sprintf("%#v", createOpts)})

	securityPolicy, err := securitypolicies.Create(containerInfraClient, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling Kubernetes API", err.Error())
		return
	}

	id := securityPolicy.UUID
	resp.State.SetAttribute(ctx, path.Root("id"), id)

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.ClusterID = types.StringValue(securityPolicy.ClusterID)
	data.SecurityPolicyTemplateID = types.StringValue(securityPolicy.SecurityPolicyTemplateID)
	data.PolicySettings = types.StringValue(securityPolicy.PolicySettings)
	data.Namespace = types.StringValue(securityPolicy.Namespace)
	data.Enabled = types.BoolValue(securityPolicy.Enabled)
	data.CreatedAt = types.StringValue(securityPolicy.CreatedAt)
	data.UpdatedAt = types.StringValue(securityPolicy.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecurityPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SecurityPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes API client", err.Error())
		return
	}

	id := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "security_policy_id", id)

	tflog.Debug(ctx, "Calling API to retrieve cluster security policy")

	securityPolicy, err := securitypolicies.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error calling Kubernetes API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called API to retrieve cluster security policy", map[string]interface{}{"security_policy": fmt.Sprintf("%#v", securityPolicy)})

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.ClusterID = types.StringValue(securityPolicy.ClusterID)
	data.SecurityPolicyTemplateID = types.StringValue(securityPolicy.SecurityPolicyTemplateID)
	data.PolicySettings = types.StringValue(securityPolicy.PolicySettings)
	data.Namespace = types.StringValue(securityPolicy.Namespace)
	data.Enabled = types.BoolValue(securityPolicy.Enabled)
	data.CreatedAt = types.StringValue(securityPolicy.CreatedAt)
	data.UpdatedAt = types.StringValue(securityPolicy.UpdatedAt)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SecurityPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SecurityPolicyResourceModel
	var data SecurityPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.ContainerInfraV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes API client", err.Error())
		return
	}

	id := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "security_policy_id", id)

	updateOpts := securitypolicies.UpdateOpts{
		PolicySettings: plan.PolicySettings.ValueString(),
		Namespace:      plan.Namespace.ValueString(),
		Enabled:        plan.Enabled.ValueBool(),
	}

	securityPolicy, err := securitypolicies.Update(client, id, &updateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_kubernetes_security_policy", err.Error())
		return
	}

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.ClusterID = types.StringValue(securityPolicy.ClusterID)
	data.SecurityPolicyTemplateID = types.StringValue(securityPolicy.SecurityPolicyTemplateID)
	data.PolicySettings = types.StringValue(securityPolicy.PolicySettings)
	data.Namespace = types.StringValue(securityPolicy.Namespace)
	data.Enabled = types.BoolValue(securityPolicy.Enabled)
	data.CreatedAt = types.StringValue(securityPolicy.CreatedAt)
	data.UpdatedAt = types.StringValue(securityPolicy.UpdatedAt)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *SecurityPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *SecurityPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	containerInfraClient, err := r.config.ContainerInfraV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating Kubernetes API client", err.Error())
		return
	}

	id := data.ID.ValueString()

	err = securitypolicies.Delete(containerInfraClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Error deleting vkcs_kubernetes_security_policy", err.Error())
		return
	}
}
