package monitoring

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	icapabilities "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/imagecapabilities"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/templater"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var (
	_ resource.Resource              = &сloudMonitoringResource{}
	_ resource.ResourceWithConfigure = &сloudMonitoringResource{}
)

func NewResource() resource.Resource {
	return &сloudMonitoringResource{}
}

func (r *сloudMonitoringResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cloud_monitoring"
}

type сloudMonitoringResource struct {
	config clients.Config
}

type CloudMonitoringResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	ImageID       types.String `tfsdk:"image_id"`
	Script        types.String `tfsdk:"script"`
	ServiceUserID types.String `tfsdk:"service_user_id"`
}

func (r *сloudMonitoringResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
						"require replacement if configuration value changes", "require replacement if configuration value changes"),
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"image_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIfConfigured(),
				},
				Description: "ID of the image to create cloud monitoring for.",
			},

			"script": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "Shell script of the cloud monitoring.",
			},

			"service_user_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the service monitoring user.",
			},
		},
		Description: "Receives settings for cloud monitoring for the `vkcs_compute_instance'.",
	}
}

func (r *сloudMonitoringResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *сloudMonitoringResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data CloudMonitoringResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	projectID := r.config.GetTenantID()

	icsClient, err := r.config.ICSV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ics client", err.Error())
		return
	}

	templaterClient, err := r.config.TemplaterV2Client(region, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS templater client", err.Error())
		return
	}

	imageID := data.ImageID.ValueString()
	capabilities, err := icapabilities.GetImageCapabilities(icsClient, imageID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error when get image capabilities", err.Error())
		return
	}

	if !isMonitoringSupported(&capabilities) {
		resp.Diagnostics.AddError("Cloud monitoring is not supported on this image", "")
		return
	}

	opts := templater.CreateUserOpts{
		ImageID: imageID,
		Capabilities: []string{
			"telegraf",
		},
	}

	settings, err := templater.CreateUser(templaterClient, projectID, opts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error when getting cloud monitoring settings", err.Error())
		return
	}

	data.ID = types.StringValue(settings.UserID)
	data.Region = types.StringValue(region)
	data.Script = types.StringValue(settings.Script)
	data.ServiceUserID = types.StringValue(settings.UserID)

	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *сloudMonitoringResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data CloudMonitoringResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *сloudMonitoringResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *CloudMonitoringResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Unable to update the cloud monitoring",
		"Not implemented. Please report this issue to the provider developers.")
}

func (r *сloudMonitoringResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data CloudMonitoringResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	projectID := r.config.GetTenantID()

	templaterClient, err := r.config.TemplaterV2Client(region, projectID)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS templater client", err.Error())
		return
	}

	if data.ServiceUserID.IsNull() {
		return
	}

	err = templater.DeleteUser(templaterClient, projectID, data.ServiceUserID.ValueString()).ExtractErr()
	if err != nil {
		if errutil.IsNotFound(err) {
			return
		}

		resp.Diagnostics.AddError("Error when deleting cloud monitoring service user", err.Error())
		return
	}
}

func isMonitoringSupported(capabilities *icapabilities.ImageCapabilities) bool {
	for _, capability := range capabilities.CapabilityVersions {
		if capability.Capability.Name == "telegraf" {
			return true
		}
	}

	return false
}
