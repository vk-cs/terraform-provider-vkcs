package mlplatform

import (
	"context"
	"fmt"

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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/imageinfos"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	"golang.org/x/exp/maps"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &MLFlowDeployResource{}

func NewMLFlowDeployResource() resource.Resource {
	return &MLFlowDeployResource{}
}

type MLFlowDeployResource struct {
	config clients.Config
}

func (r *MLFlowDeployResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_mlplatform_mlflow_deploy"
}

type MLFlowDeployResourceModel struct {
	ID               types.String              `tfsdk:"id"`
	Name             types.String              `tfsdk:"name"`
	FlavorID         types.String              `tfsdk:"flavor_id"`
	MLFlowInstanceID types.String              `tfsdk:"mlflow_instance_id"`
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

func getSchemaMLFlowDeploy(ctx context.Context, resp *resource.SchemaResponse) map[string]schema.Attribute {
	mlFlowDeployAttrs := map[string]schema.Attribute{
		"mlflow_instance_id": schema.StringAttribute{
			Required:    true,
			Description: "MLFlow instance ID",
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
			Optional:    true,
			Description: "Instance's data volumes configuration",
		},
	}

	maps.Copy(mlFlowDeployAttrs, getCommonInstanceSchema(ctx, resp))
	return mlFlowDeployAttrs
}

func (r *MLFlowDeployResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes:  getSchemaMLFlowDeploy(ctx, resp),
		Description: "Manages a ML Platform Deploy resource.",
	}
}

func (r *MLFlowDeployResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *MLFlowDeployResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data MLFlowDeployResourceModel
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
		imageInfoResp, err := imageinfos.Get(mlPlatformClient, mlFlowDeployInstanceType).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error creating vkcs_mlplatform_jupyterhub", err.Error())
			return
		}
		bootVolumeSize = imageInfoResp.VolumeSize
	}

	volumesOpts := expandVolumeOpts(data.BootVolume, bootVolumeSize, data.DataVolumes, data.AvailabilityZone.ValueString())
	networksOpts := expandNetworkOpts(data.Networks)
	deployCreateOpts := instances.CreateOpts{
		InstanceName:           data.Name.ValueString(),
		DomainName:             "",
		InstanceType:           mlFlowDeployInstanceType,
		DeployMLFlowInstanceID: data.MLFlowInstanceID.ValueString(),
		Flavor:                 data.FlavorID.ValueString(),
		Volumes:                volumesOpts,
		Networks:               networksOpts,
	}

	tflog.Debug(ctx, "Creating Deploy Instance", map[string]interface{}{"options": fmt.Sprintf("%#v", deployCreateOpts)})

	deployResp, err := instances.Create(mlPlatformClient, &deployCreateOpts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_mlplatform_mlflow_deploy", err.Error())
		return
	}
	resp.State.SetAttribute(ctx, path.Root("id"), deployResp.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusPrepareDBAAS, instanceStatusCreating, instanceStatusInstallScripts, instanceStatusStarting},
		Target:     []string{instanceStatusRunning},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, deployResp.ID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the Deploy instance to be created", map[string]interface{}{"timeout": timeout})

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}

	deployResp, err = instances.Get(mlPlatformClient, deployResp.ID).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_deploy", err.Error())
		return
	}

	data.ID = types.StringValue(deployResp.ID)
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(deployResp.Name)
	data.FlavorID = types.StringValue(deployResp.FlavorID)
	data.CreatedAt = types.StringValue(deployResp.CreatedAt)
	data.PrivateIP = types.StringValue(deployResp.PrivateIP)
	data.DNSName = types.StringValue(deployResp.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(deployResp.Volumes)

	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *MLFlowDeployResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data MLFlowDeployResourceModel

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

	deployID := data.ID.ValueString()

	deployResp, err := instances.Get(mlPlatformClient, deployID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_deploy", checkDeleted.Error())
		}
		return
	}
	data.Region = types.StringValue(region)

	data.Name = types.StringValue(deployResp.Name)
	data.FlavorID = types.StringValue(deployResp.FlavorID)
	data.CreatedAt = types.StringValue(deployResp.CreatedAt)
	data.PrivateIP = types.StringValue(deployResp.PrivateIP)
	data.MLFlowInstanceID = types.StringValue(deployResp.DeployMLFlowInstanceID)
	data.DNSName = types.StringValue(deployResp.DomainName)

	bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(deployResp.Volumes)
	data.BootVolume = bootVolume
	data.DataVolumes = dataVolumes
	data.AvailabilityZone = types.StringValue(availabilityZone)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *MLFlowDeployResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MLFlowDeployResourceModel
	var data MLFlowDeployResourceModel

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

	deployID := data.ID.ValueString()

	if !plan.FlavorID.Equal(data.FlavorID) {
		err := instanceUpdateFlavor(ctx, mlPlatformClient, deployID, plan.FlavorID.ValueString(), timeout)

		if err != nil {
			resp.Diagnostics.AddError("Error updating flavor of vkcs_mlplatform_deploy", err.Error())
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
		err := instanceUpdateVolumes(ctx, mlPlatformClient, deployID, changedVolumes, timeout)
		if err != nil {
			resp.Diagnostics.AddError("Error updating volumes of vkcs_mlplatform_deploy", err.Error())
			return
		}

		deployResp, err := instances.Get(mlPlatformClient, deployID).Extract()
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving deploy", err.Error())
			return
		}
		bootVolume, dataVolumes, availabilityZone := flattenVolumeOpts(deployResp.Volumes)

		data.BootVolume = bootVolume
		data.DataVolumes = dataVolumes
		data.AvailabilityZone = types.StringValue(availabilityZone)
	}

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func (r *MLFlowDeployResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data MLFlowDeployResourceModel

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

	deployID := data.ID.ValueString()

	err = instances.Delete(mlPlatformClient, deployID).ExtractErr()
	if errutil.IsNotFound(err) {
		return
	} else if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_mlplatform_deploy", err.Error())
		return
	}

	deployStateConf := &retry.StateChangeConf{
		Pending:    []string{instanceStatusRunning, instanceStatusDeleting},
		Target:     []string{instanceStatusDeleted},
		Refresh:    instanceStateRefreshFunc(mlPlatformClient, deployID),
		Timeout:    timeout,
		Delay:      instanceDelay,
		MinTimeout: instanceMinTimeout,
	}

	tflog.Debug(ctx, "Waiting for the Deploy instance to be deleted")

	_, err = deployStateConf.WaitForStateContext(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for the instance to become ready", err.Error())
		return
	}
}

func (r *MLFlowDeployResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
