package mlplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/imageinfos"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	"golang.org/x/exp/maps"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &MLFlowResource{}

func NewMLFlowResource() resource.Resource {
	return &MLFlowResource{}
}

type MLFlowResource struct {
	config clients.Config
}

func (r *MLFlowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_mlplatform_mlflow"
}

type MLFlowResourceModel struct {
	ID               types.String              `tfsdk:"id"`
	Name             types.String              `tfsdk:"name"`
	FlavorID         types.String              `tfsdk:"flavor_id"`
	JHInstanceID     types.String              `tfsdk:"jh_instance_id"`
	AvailabilityZone types.String              `tfsdk:"availability_zone"`
	BootVolume       *MLPlatformVolumeModel    `tfsdk:"boot_volume"`
	DataVolumes      []*MLPlatformVolumeModel  `tfsdk:"data_volumes"`
	Networks         []*MLPlatformNetworkModel `tfsdk:"networks"`
	DemoMode         types.Bool                `tfsdk:"demo_mode"`
	CreatedAt        types.String              `tfsdk:"created_at"`
	PrivateIP        types.String              `tfsdk:"private_ip"`
	DNSName          types.String              `tfsdk:"dns_name"`
	Region           types.String              `tfsdk:"region"`
	Timeouts         timeouts.Value            `tfsdk:"timeouts"`
}

func getSchemaMLFlow(ctx context.Context, resp *resource.SchemaResponse) map[string]schema.Attribute {
	mlFlowAttrs := map[string]schema.Attribute{
		"jh_instance_id": schema.StringAttribute{
			Required:    true,
			Description: "JupyterHub instance ID",
		},

		"demo_mode": schema.BoolAttribute{
			Optional:    true,
			Description: "Controls whether demo mode is enabled. If true, data will be stored on mlflow virtual machine. If false, s3 bucket will be used alongside dbaas postgres database.",
		},
	}

	mlFlowSchema := getCommonInstanceSchema(ctx, resp)
	maps.Copy(mlFlowSchema, mlFlowAttrs)
	return mlFlowSchema
}

func (r *MLFlowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:  getSchemaMLFlow(ctx, resp),
		Description: "Manages a ML Platform MLFlow resource.",
	}
}

func (r *MLFlowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *MLFlowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MLFlowResourceModel
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
		imageInfoResp, err := imageinfos.Get(mlPlatformClient, mlFlowInstanceType).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error creating vkcs_mlplatform_jupyterhub", err.Error())
			return
		}
		bootVolumeSize = imageInfoResp.VolumeSize
	}

	volumesOpts := expandVolumeOpts(data.BootVolume, bootVolumeSize, data.DataVolumes, data.AvailabilityZone.ValueString())
	networksOpts := expandNetworkOpts(data.Networks)
	mlFlowCreateOpts := instances.CreateOpts{
		InstanceName:       data.Name.ValueString(),
		DomainName:         "",
		InstanceType:       mlFlowInstanceType,
		MLFlowJHInstanceID: data.JHInstanceID.ValueString(),
		ISMLFlowDemoMode:   data.DemoMode.ValueBool(),
		Flavor:             data.FlavorID.ValueString(),
		Volumes:            volumesOpts,
		Networks:           networksOpts,
	}

	tflog.Debug(ctx, "Creating MLFlow Instance", map[string]interface{}{"options": fmt.Sprintf("%#v", mlFlowCreateOpts)})

	mlFlowResp, err := instances.Create(mlPlatformClient, &mlFlowCreateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_mlplatform_mlflow", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), mlFlowResp.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusPrepareDBAAS, instanceStatusCreating, instanceStatusInstallScripts, instanceStatusStarting},
		Target:     []string{instanceStatusRunning},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, mlFlowResp.ID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the MLFlow instance to be created", map[string]interface{}{"timeout": timeout})

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}

	mlFlowResp, err = instances.Get(mlPlatformClient, mlFlowResp.ID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_mlflow", err.Error())
		return
	}

	data.ID = types.StringValue(mlFlowResp.ID)
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(mlFlowResp.Name)
	data.FlavorID = types.StringValue(mlFlowResp.FlavorID)
	data.CreatedAt = types.StringValue(mlFlowResp.CreatedAt)
	data.PrivateIP = types.StringValue(mlFlowResp.PrivateIP)
	data.DNSName = types.StringValue(mlFlowResp.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(mlFlowResp.Volumes)

	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *MLFlowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MLFlowResourceModel

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
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	mlFlowID := data.ID.ValueString()

	mlFlowResp, err := instances.Get(mlPlatformClient, mlFlowID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_mlflow", checkDeleted.Error())
		}
		return
	}
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(mlFlowResp.Name)
	data.FlavorID = types.StringValue(mlFlowResp.FlavorID)
	data.CreatedAt = types.StringValue(mlFlowResp.CreatedAt)
	data.PrivateIP = types.StringValue(mlFlowResp.PrivateIP)
	data.JHInstanceID = types.StringValue(mlFlowResp.MLFlowJHInstanceID)
	data.DNSName = types.StringValue(mlFlowResp.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(mlFlowResp.Volumes)
	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MLFlowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MLFlowResourceModel
	var data MLFlowResourceModel

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

	mlFlowID := data.ID.ValueString()

	if !plan.FlavorID.Equal(data.FlavorID) {
		err := instanceUpdateFlavor(ctx, mlPlatformClient, mlFlowID, plan.FlavorID.ValueString(), timeout)

		if err != nil {
			resp.Diagnostics.AddError("Error updating flavor of vkcs_mlplatform_mlflow", err.Error())
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
		err := instanceUpdateVolumes(ctx, mlPlatformClient, mlFlowID, changedVolumes, timeout)
		if err != nil {
			resp.Diagnostics.AddError("Error updating volumes of vkcs_mlplatform_mlflow", err.Error())
			return
		}

		mlFlowResp, err := instances.Get(mlPlatformClient, mlFlowID).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving mlflow", err.Error())
			return
		}
		bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(mlFlowResp.Volumes)

		data.BootVolume = bootVolume
		data.DataVolumes = dataVolumes
		data.AvailabilityZone = types.StringValue(availabilityZone)
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *MLFlowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MLFlowResourceModel

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

	mlFlowID := data.ID.ValueString()

	err = instances.Delete(mlPlatformClient, mlFlowID).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_mlplatform_mlflow", err.Error())
		return
	}

	mlFlowStateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusRunning, instanceStatusDeleting},
		Target:     []string{instanceStatusDeleted},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, mlFlowID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the MLFlow instance to be deleted")

	_, err = mlFlowStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}
}

func (r *MLFlowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
