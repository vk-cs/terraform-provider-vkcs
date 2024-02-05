package mlplatform

import (
	"context"
	"fmt"

	"golang.org/x/exp/maps"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
var _ resource.Resource = &JupyterHubResource{}

func NewJupyterHubResource() resource.Resource {
	return &JupyterHubResource{}
}

type JupyterHubResource struct {
	config clients.Config
}

func (r *JupyterHubResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_mlplatform_jupyterhub"
}

type JupyterHubResourceModel struct {
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
	S3FSBucket       types.String              `tfsdk:"s3fs_bucket"`
	CreatedAt        types.String              `tfsdk:"created_at"`
	PrivateIP        types.String              `tfsdk:"private_ip"`
	DNSName          types.String              `tfsdk:"dns_name"`
	Region           types.String              `tfsdk:"region"`
	Timeouts         timeouts.Value            `tfsdk:"timeouts"`
}

func getSchemaJupyterHub(ctx context.Context, resp *resource.SchemaResponse) map[string]schema.Attribute {
	jupyterHubAttrs := map[string]schema.Attribute{
		"domain_name": schema.StringAttribute{
			Optional:    true,
			Description: "Domain name. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"admin_name": schema.StringAttribute{
			Optional:    true,
			Description: "JupyterHub admin name. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"admin_password": schema.StringAttribute{
			Optional:    true,
			Sensitive:   true,
			Description: "JupyterHub admin password. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},

		"data_volumes": schema.ListNestedAttribute{
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"size": schema.Int64Attribute{
						Required:    true,
						Description: "Size of the volume",
					},
					"volume_type": schema.StringAttribute{
						Required:    true,
						Description: "Type of the volume",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Name of the volume",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"volume_id": schema.StringAttribute{
						Computed:    true,
						Description: "ID of the volume",
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			Required:    true,
			Description: "Instance's data volumes configuration",
		},

		"s3fs_bucket": schema.StringAttribute{
			Optional:    true,
			Description: "Connect specified s3 bucket to instance as volume. Changing this creates a new resource",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		},
	}

	jupyterHubSchema := getCommonInstanceSchema(ctx, resp)
	maps.Copy(jupyterHubSchema, jupyterHubAttrs)
	return jupyterHubSchema
}

func (r *JupyterHubResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:  getSchemaJupyterHub(ctx, resp),
		Description: "Manages a ML Platform JupyterHub resource.",
	}
}

func (r *JupyterHubResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *JupyterHubResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data JupyterHubResourceModel
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
		imageInfoResp, err := imageinfos.Get(mlPlatformClient, jupyterHubInstanceType).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error creating vkcs_mlplatform_jupyterhub", err.Error())
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
			resp.Diagnostics.AddError("Error creating vkcs_mlplatform_jupyterhub", err.Error())
			return
		}
		domainName = domainNameResp.DNS
	}

	jupyterHubCreateOpts := instances.CreateOpts{
		InstanceName:    data.Name.ValueString(),
		DomainName:      domainName,
		InstanceType:    jupyterHubInstanceType,
		JHAdminName:     data.AdminName.ValueString(),
		JHAdminPassword: data.AdminPassword.ValueString(),
		Flavor:          data.FlavorID.ValueString(),
		Volumes:         volumesOpts,
		Networks:        networksOpts,
		S3FSBucket:      data.S3FSBucket.ValueString(),
	}

	tflog.Debug(ctx, "Creating JupyterHub Instance", map[string]interface{}{"options": fmt.Sprintf("%#v", jupyterHubCreateOpts)})

	jupyterHubResp, err := instances.Create(mlPlatformClient, &jupyterHubCreateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_mlplatform_jupyterhub", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), jupyterHubResp.ID)

	jupyterHubStateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusPrepareDBAAS, instanceStatusCreating, instanceStatusInstallScripts, instanceStatusStarting},
		Target:     []string{instanceStatusRunning},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, jupyterHubResp.ID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the JupyterHub instance to be created", map[string]interface{}{"timeout": timeout})

	_, err = jupyterHubStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}

	jupyterHubResp, err = instances.Get(mlPlatformClient, jupyterHubResp.ID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_jupyterhub", err.Error())
		return
	}

	data.ID = types.StringValue(jupyterHubResp.ID)
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(jupyterHubResp.Name)
	data.FlavorID = types.StringValue(jupyterHubResp.FlavorID)
	data.CreatedAt = types.StringValue(jupyterHubResp.CreatedAt)
	data.PrivateIP = types.StringValue(jupyterHubResp.PrivateIP)
	data.AdminName = types.StringValue(jupyterHubResp.JHAdminName)
	data.DNSName = types.StringValue(jupyterHubResp.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(jupyterHubResp.Volumes)

	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *JupyterHubResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data JupyterHubResourceModel

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

	jupyterHubID := data.ID.ValueString()

	jupyterHubResp, err := instances.Get(mlPlatformClient, jupyterHubID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_jupyterhub", checkDeleted.Error())
		}
		return
	}
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(jupyterHubResp.Name)
	data.FlavorID = types.StringValue(jupyterHubResp.FlavorID)
	data.CreatedAt = types.StringValue(jupyterHubResp.CreatedAt)
	data.PrivateIP = types.StringValue(jupyterHubResp.PrivateIP)
	data.AdminName = types.StringValue(jupyterHubResp.JHAdminName)
	data.DNSName = types.StringValue(jupyterHubResp.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(jupyterHubResp.Volumes)
	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *JupyterHubResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan JupyterHubResourceModel
	var data JupyterHubResourceModel

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

	jupyterHubID := data.ID.ValueString()

	if !plan.FlavorID.Equal(data.FlavorID) {
		err := instanceUpdateFlavor(ctx, mlPlatformClient, jupyterHubID, plan.FlavorID.ValueString(), timeout)

		if err != nil {
			resp.Diagnostics.AddError("Error updating flavor of vkcs_mlplatform_jupyterhub", err.Error())
			return
		}

		data.FlavorID = plan.FlavorID
	}

	var changedVolumes []instances.ResizeVolumeParams
	if !plan.BootVolume.Size.Equal(data.BootVolume.Size) {
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
		err := instanceUpdateVolumes(ctx, mlPlatformClient, jupyterHubID, changedVolumes, timeout)
		if err != nil {
			resp.Diagnostics.AddError("Error updating volumes of vkcs_mlplatform_jupyterhub", err.Error())
			return
		}

		jupyterHubResp, err := instances.Get(mlPlatformClient, jupyterHubID).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_jupyterhub", err.Error())
			return
		}
		bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(jupyterHubResp.Volumes)

		data.BootVolume = bootVolume
		data.DataVolumes = dataVolumes
		data.AvailabilityZone = types.StringValue(availabilityZone)
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *JupyterHubResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data JupyterHubResourceModel

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

	jupyterHubID := data.ID.ValueString()

	err = instances.Delete(mlPlatformClient, jupyterHubID).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_mlplatform_jupyterhub", err.Error())
		return
	}

	jupyterHubStateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusRunning, instanceStatusDeleting},
		Target:     []string{instanceStatusDeleted},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, jupyterHubID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the JupyterHub instance to be deleted")

	_, err = jupyterHubStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}
}

func (r *JupyterHubResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
