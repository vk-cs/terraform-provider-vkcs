package mlplatform

import (
	"context"
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/dnsnames"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/imageinfos"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &K8SRegistryResource{}

func NewK8SRegistryResource() resource.Resource {
	return &K8SRegistryResource{}
}

type K8SRegistryResource struct {
	config clients.Config
}

func (r *K8SRegistryResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_mlplatform_k8s_registry"
}

type K8SRegistryResourceModel struct {
	ID               types.String              `tfsdk:"id"`
	Name             types.String              `tfsdk:"name"`
	DomainName       types.String              `tfsdk:"domain_name"`
	AdminName        types.String              `tfsdk:"admin_name"`
	AdminPassword    types.String              `tfsdk:"admin_password"`
	FlavorID         types.String              `tfsdk:"flavor_id"`
	AvailabilityZone types.String              `tfsdk:"availability_zone"`
	BootVolume       *MLPlatformVolumeModel    `tfsdk:"boot_volume"`
	DataVolumes      []*MLPlatformVolumeModel  `tfsdk:"data_volumes"`
	Networks         []*MLPlatformNetworkModel `tfsdk:"networks"`
	CreatedAt        types.String              `tfsdk:"created_at"`
	PrivateIP        types.String              `tfsdk:"private_ip"`
	DNSName          types.String              `tfsdk:"dns_name"`
	Region           types.String              `tfsdk:"region"`
	Timeouts         timeouts.Value            `tfsdk:"timeouts"`
}

func getSchemaK8SRegistry(ctx context.Context, resp *resource.SchemaResponse) map[string]schema.Attribute {
	k8sRegistryAttrs := map[string]schema.Attribute{
		"domain_name": schema.StringAttribute{
			Optional:    true,
			Description: "Domain name. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"admin_name": schema.StringAttribute{
			Optional:    true,
			Description: "K8SRegistry admin name. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"admin_password": schema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "K8SRegistry admin password. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"networks": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"ip_pool": schema.StringAttribute{
						Required:    true,
						Description: "ID of the ip pool",
					},
					"network_id": schema.StringAttribute{
						Required:    true,
						Description: "ID of the network",
					},
				},
			},
			Required:    true,
			Description: "Network configuration",
			Validators: []validator.List{
				listvalidator.SizeAtMost(1),
			},
		},
	}

	k8sRegistrySchema := getCommonInstanceSchema(ctx, resp)
	maps.Copy(k8sRegistrySchema, k8sRegistryAttrs)
	return k8sRegistrySchema
}

func (r *K8SRegistryResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:  getSchemaK8SRegistry(ctx, resp),
		Description: "Manages a ML Platform K8SRegistry resource.",
	}
}

func (r *K8SRegistryResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *K8SRegistryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data K8SRegistryResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, instanceCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	bootVolumeSize := data.BootVolume.Size.ValueInt64()
	if bootVolumeSize == 0 {
		imageInfoResp, err := imageinfos.Get(mlPlatformClient, k8sRegistryInstanceType).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error creating vkcs_mlplatform_k8s_registry", err.Error())
			return
		}
		bootVolumeSize = imageInfoResp.VolumeSize
	}

	volumesOpts := expandVolumeOpts(data.BootVolume, bootVolumeSize, data.DataVolumes, data.AvailabilityZone.ValueString())
	networksOpts := expandNetworkOpts(data.Networks)

	domainName := data.DomainName.ValueString()
	if domainName == "" {
		domainNameResp, err := dnsnames.Get(mlPlatformClient).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error creating vkcs_mlplatform_k8s_registry", err.Error())
			return
		}
		domainName = domainNameResp.DNS
	}

	createOpts := instances.CreateOpts{
		InstanceName:    data.Name.ValueString(),
		DomainName:      domainName,
		InstanceType:    k8sRegistryInstanceType,
		JHAdminName:     data.AdminName.ValueString(),
		JHAdminPassword: data.AdminPassword.ValueString(),
		Flavor:          data.FlavorID.ValueString(),
		Volumes:         volumesOpts,
		Networks:        networksOpts,
	}

	tflog.Debug(ctx, "Creating K8SRegistry Instance", map[string]interface{}{"options": fmt.Sprintf("%#v", createOpts)})

	response, err := instances.Create(mlPlatformClient, &createOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_mlplatform_k8s_registry", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), response.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusPrepareDBAAS, instanceStatusCreating, instanceStatusInstallScripts, instanceStatusStarting},
		Target:     []string{instanceStatusRunning},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, response.ID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the K8SRegistry instance to be created", map[string]interface{}{"timeout": timeout})

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}

	response, err = instances.Get(mlPlatformClient, response.ID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_k8s_registry", err.Error())
		return
	}

	data.ID = types.StringValue(response.ID)
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(response.Name)
	data.FlavorID = types.StringValue(response.FlavorID)
	data.CreatedAt = types.StringValue(response.CreatedAt)
	data.PrivateIP = types.StringValue(response.PrivateIP)
	data.AdminName = types.StringValue(response.JHAdminName)
	data.DNSName = types.StringValue(response.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(response.Volumes)

	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *K8SRegistryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data K8SRegistryResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS mlplatform client", err.Error())
		return
	}

	id := data.ID.ValueString()

	response, err := instances.Get(mlPlatformClient, id).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_k8s_registry", checkDeleted.Error())
		}
		return
	}
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(response.Name)
	data.FlavorID = types.StringValue(response.FlavorID)
	data.CreatedAt = types.StringValue(response.CreatedAt)
	data.PrivateIP = types.StringValue(response.PrivateIP)
	data.AdminName = types.StringValue(response.JHAdminName)
	data.DNSName = types.StringValue(response.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(response.Volumes)
	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *K8SRegistryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan K8SRegistryResourceModel
	var data K8SRegistryResourceModel

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

	timeout, diags := data.Timeouts.Update(ctx, instanceUpdateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	id := data.ID.ValueString()

	if !plan.FlavorID.Equal(data.FlavorID) {
		err := instanceUpdateFlavor(ctx, mlPlatformClient, id, plan.FlavorID.ValueString(), timeout)

		if err != nil {
			resp.Diagnostics.AddError("Error updating flavor of vkcs_mlplatform_k8s_registry", err.Error())
			return
		}

		data.FlavorID = plan.FlavorID
	}

	var changedVolumes []instances.ResizeVolumeParams
	if !plan.BootVolume.Size.IsUnknown() && !plan.BootVolume.Size.Equal(data.BootVolume.Size) {
		changedVolumes = append(changedVolumes, instances.ResizeVolumeParams{
			ID:   plan.BootVolume.VolumeID.ValueString(),
			Size: int(plan.BootVolume.Size.ValueInt64()),
		})
	}

	for _, planDataVolume := range plan.DataVolumes {
		for _, dataDataVolume := range data.DataVolumes {
			if planDataVolume.VolumeID == dataDataVolume.VolumeID && planDataVolume.Size != dataDataVolume.Size {
				changedVolumes = append(changedVolumes, instances.ResizeVolumeParams{
					ID:   planDataVolume.VolumeID.ValueString(),
					Size: int(planDataVolume.Size.ValueInt64()),
				})
			}
		}
	}

	if len(changedVolumes) > 0 {
		err := instanceUpdateVolumes(ctx, mlPlatformClient, id, changedVolumes, timeout)
		if err != nil {
			resp.Diagnostics.AddError("Error updating volumes of vkcs_mlplatform_k8s_registry", err.Error())
			return
		}

		response, err := instances.Get(mlPlatformClient, id).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_k8s_registry", err.Error())
			return
		}
		bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(response.Volumes)

		data.BootVolume = bootVolume
		data.DataVolumes = dataVolumes
		data.AvailabilityZone = types.StringValue(availabilityZone)
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *K8SRegistryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data K8SRegistryResourceModel

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, instanceDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	mlPlatformClient, err := r.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	id := data.ID.ValueString()

	err = instances.Delete(mlPlatformClient, id).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_mlplatform_k8s_registry", err.Error())
		return
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusRunning, instanceStatusDeleting},
		Target:     []string{instanceStatusDeleted},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, id),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the K8SRegistry instance to be deleted")

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}
}

func (r *K8SRegistryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
