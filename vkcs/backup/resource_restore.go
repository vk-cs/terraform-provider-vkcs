package backup

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

var (
	_ resource.Resource = &RestoreResource{}
)

func NewRestoreResource() resource.Resource {
	return &RestoreResource{}
}

type RestoreResource struct {
	config clients.Config
}

type RestoreResourceModel struct {
	ID                      types.String  `tfsdk:"id"`
	Region                  types.String  `tfsdk:"region"`
	CheckpointID            types.String  `tfsdk:"checkpoint_id"`
	SourceResourceID        types.String  `tfsdk:"source_resource_id"`
	ProviderID              types.String  `tfsdk:"provider_id"`
	RestoreName             types.String  `tfsdk:"restore_name"`
	RestoreKeypair          types.String  `tfsdk:"restore_keypair"`
	AvailabilityZone        types.String  `tfsdk:"availability_zone"`
	RestoreNetID            types.String  `tfsdk:"restore_net_id"`
	Flavor                  types.String  `tfsdk:"flavor"`
	VolumeType              types.String  `tfsdk:"volume_type"`
	VolumeSize              types.Int64   `tfsdk:"volume_size"`
	EnableFIP               types.Bool    `tfsdk:"enable_fip"`
	SecurityGroups          types.Set     `tfsdk:"security_groups"`
	VolumeTargetVMID        types.String  `tfsdk:"volume_target_vm_id"`
	VolumeAutoresizeEnabled types.Int64   `tfsdk:"volume_autoresize_enabled"`
	VolumeAutoresizeMaxSize types.Int64   `tfsdk:"volume_autoresize_max_size"`
	CloudMonitoringEnabled  types.Bool    `tfsdk:"cloud_monitoring_enabled"`
	ResourcesStatus         types.Map     `tfsdk:"resources_status"`
	ResourcesReason         types.Map     `tfsdk:"resources_reason"`
	Timeouts                timeouts.Value `tfsdk:"timeouts"`
}

func (r *RestoreResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_backup_restore"
}

func (r *RestoreResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the restore operation.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region to fetch the backup service endpoint from. Defaults to the provider region.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"checkpoint_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the backup checkpoint to restore.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_resource_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the VM resource in the checkpoint. Required only when the checkpoint contains multiple VMs.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"provider_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the Nova backup provider used for restoration.",
			},
			"restore_name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the restored VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"restore_keypair": schema.StringAttribute{
				Optional:    true,
				Description: "Keypair to associate with the restored VM. The backup service validates this value and may require it for a new VM restore.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"availability_zone": schema.StringAttribute{
				Optional:    true,
				Description: "Availability zone for the restored VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"restore_net_id": schema.StringAttribute{
				Optional:    true,
				Description: "Network or subnet ID for the restored VM. The backup service validates this value and may require it for a new VM restore.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"flavor": schema.StringAttribute{
				Optional:    true,
				Description: "Flavor ID for the restored VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_type": schema.StringAttribute{
				Optional:    true,
				Description: "Volume type for restored disks.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_size": schema.Int64Attribute{
				Optional:    true,
				Description: "Size in GiB for restored disks.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"enable_fip": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to allocate a floating IP for the restored VM.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"security_groups": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Security group IDs to apply to the restored VM.",
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
			},
			"volume_target_vm_id": schema.StringAttribute{
				Optional:    true,
				Description: "Target VM ID for restored volumes, when supported by the checkpoint.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"volume_autoresize_enabled": schema.Int64Attribute{
				Optional:    true,
				Description: "Volume auto-resize setting passed to the backup service.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"volume_autoresize_max_size": schema.Int64Attribute{
				Optional:    true,
				Description: "Maximum volume size for auto-resize, in GiB.",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"cloud_monitoring_enabled": schema.BoolAttribute{
				Optional:    true,
				Description: "Whether to enable cloud monitoring on the restored VM.",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"resources_status": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Read-only restore status for each affected resource.",
			},
			"resources_reason": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Read-only restore reason for each resource that could not be restored.",
			},
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{Create: true}),
		},
		Description: "Restores a VM from a backup checkpoint. This is a one-time operation: removing the resource from configuration only removes it from Terraform state and does not delete the restored VM.",
	}
}

func (r *RestoreResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *RestoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan RestoreResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout, diags := plan.Timeouts.Create(ctx, restoreCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	createDeadline := time.Now().Add(createTimeout)

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	backupClient, err := r.config.BackupV1Client(region, r.config.GetProjectID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	providerInfo, err := findProvider(backupClient, "", ProviderNameNova)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_backup_restore", err.Error())
		return
	}

	checkpoint, err := checkpoints.Get(backupClient, providerInfo.ID, plan.CheckpointID.ValueString()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving checkpoint for vkcs_backup_restore", err.Error())
		return
	}
	sourceResourceID, err := selectRestoreResourceID(checkpoint.Resources, plan.SourceResourceID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Invalid checkpoint for vkcs_backup_restore", err.Error())
		return
	}

	parameters, diags := expandRestoreParameters(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resourceKey := fmt.Sprintf("%s#%s", NovaInstance, sourceResourceID)
	mutexKey := fmt.Sprintf("backup-restore:%s:%s:%s", plan.CheckpointID.ValueString(), sourceResourceID, parameters.RestoreName)
	mutex := r.config.GetMutex()
	mutex.Lock(mutexKey)
	defer mutex.Unlock(mutexKey)

	restoreOpts := restores.RestoreCreateOpts{
		ProviderID:   providerInfo.ID,
		CheckpointID: plan.CheckpointID.ValueString(),
		Parameters:   map[string]restores.RestoreParameters{resourceKey: parameters},
	}

	// Failed create responses taint Terraform resources. Reuse an existing matching
	// one-shot restore on retry instead of issuing a duplicate remote restore.
	existingRestores, err := restores.ListAll(backupClient)
	if err != nil {
		resp.Diagnostics.AddError("Error checking existing vkcs_backup_restore requests", err.Error())
		return
	}

	restoreResponse, err := getOrCreateRestore(existingRestores, restoreOpts, func() (*restores.Restore, error) {
		return restores.Create(backupClient, &restores.RestoreCreate{Restore: &restoreOpts}).Extract()
	})
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_backup_restore", err.Error())
		return
	}
	if restoreResponse.ID == "" {
		resp.Diagnostics.AddError("Error creating vkcs_backup_restore", "backup API returned an empty restore ID")
		return
	}

	plan.ID = types.StringValue(restoreResponse.ID)
	plan.Region = types.StringValue(region)
	plan.ProviderID = types.StringValue(providerInfo.ID)
	setRestoreStatuses(ctx, &plan, restoreResponse, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	waitCtx, cancel := context.WithDeadline(ctx, createDeadline)
	defer cancel()
	completedRestore, err := waitForRestore(waitCtx, backupClient, restoreResponse.ID)
	if completedRestore != nil {
		setRestoreStatuses(ctx, &plan, completedRestore, &resp.Diagnostics)
		resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for vkcs_backup_restore", err.Error())
		return
	}

	if err := restoreOutcomeError(completedRestore.ID, completedRestore.ResourcesStatus, completedRestore.ResourcesReason); err != nil {
		resp.Diagnostics.AddError("VM restore failed", err.Error())
	}
}

func (r *RestoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RestoreResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := state.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}
	backupClient, err := r.config.BackupV1Client(region, r.config.GetProjectID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	restoreResponse, err := restores.Get(backupClient, state.ID.ValueString()).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_backup_restore", checkDeleted.Error())
		}
		return
	}

	state.ID = types.StringValue(restoreResponse.ID)
	state.Region = types.StringValue(region)
	state.ProviderID = types.StringValue(restoreResponse.ProviderID)
	state.CheckpointID = types.StringValue(restoreResponse.CheckpointID)
	setRestoreStatuses(ctx, &state, restoreResponse, &resp.Diagnostics)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RestoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unsupported vkcs_backup_restore update", "Restore operations are one-time actions and must be replaced when configuration changes.")
}

func (r *RestoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// The backup API has no restore-delete operation. Terraform removes state only.
}

func expandRestoreParameters(ctx context.Context, model RestoreResourceModel) (restores.RestoreParameters, diag.Diagnostics) {
	parameters := restores.RestoreParameters{
		RestoreName: model.RestoreName.ValueString(),
		Destination: "new",
	}
	if !model.RestoreKeypair.IsNull() && !model.RestoreKeypair.IsUnknown() {
		parameters.RestoreKeypair = model.RestoreKeypair.ValueString()
	}
	if !model.AvailabilityZone.IsNull() && !model.AvailabilityZone.IsUnknown() {
		parameters.AvailabilityZone = model.AvailabilityZone.ValueString()
	}
	if !model.RestoreNetID.IsNull() && !model.RestoreNetID.IsUnknown() {
		parameters.RestoreNetID = model.RestoreNetID.ValueString()
	}
	if !model.Flavor.IsNull() && !model.Flavor.IsUnknown() {
		parameters.Flavor = model.Flavor.ValueString()
	}
	if !model.VolumeType.IsNull() && !model.VolumeType.IsUnknown() {
		parameters.VolumeType = model.VolumeType.ValueString()
	}
	if !model.VolumeSize.IsNull() && !model.VolumeSize.IsUnknown() {
		volumeSize := model.VolumeSize.ValueInt64()
		parameters.VolumeSize = &volumeSize
	}
	if !model.EnableFIP.IsNull() && !model.EnableFIP.IsUnknown() {
		enableFIP := model.EnableFIP.ValueBool()
		parameters.EnableFIP = &enableFIP
	}
	if !model.VolumeTargetVMID.IsNull() && !model.VolumeTargetVMID.IsUnknown() {
		parameters.VolumeTargetVMID = model.VolumeTargetVMID.ValueString()
	}
	if !model.VolumeAutoresizeEnabled.IsNull() && !model.VolumeAutoresizeEnabled.IsUnknown() {
		volumeAutoresizeEnabled := model.VolumeAutoresizeEnabled.ValueInt64()
		parameters.VolumeAutoresizeEnabled = &volumeAutoresizeEnabled
	}
	if !model.VolumeAutoresizeMaxSize.IsNull() && !model.VolumeAutoresizeMaxSize.IsUnknown() {
		volumeAutoresizeMaxSize := model.VolumeAutoresizeMaxSize.ValueInt64()
		parameters.VolumeAutoresizeMaxSize = &volumeAutoresizeMaxSize
	}
	if !model.CloudMonitoringEnabled.IsNull() && !model.CloudMonitoringEnabled.IsUnknown() {
		cloudMonitoringEnabled := model.CloudMonitoringEnabled.ValueBool()
		parameters.CloudMonitoringEnabled = &cloudMonitoringEnabled
	}

	var securityGroups []string
	var diags diag.Diagnostics
	if !model.SecurityGroups.IsNull() && !model.SecurityGroups.IsUnknown() {
		diags.Append(model.SecurityGroups.ElementsAs(ctx, &securityGroups, false)...)
		if diags.HasError() {
			return parameters, diags
		}
	}
	parameters.SecurityGroups = securityGroups
	return parameters, diags
}

func setRestoreStatuses(ctx context.Context, model *RestoreResourceModel, restore *restores.Restore, diags *diag.Diagnostics) {
	resourcesStatus := restore.ResourcesStatus
	if resourcesStatus == nil {
		resourcesStatus = map[string]string{}
	}
	resourcesReason := restore.ResourcesReason
	if resourcesReason == nil {
		resourcesReason = map[string]string{}
	}

	statusMap, statusDiags := types.MapValueFrom(ctx, types.StringType, resourcesStatus)
	diags.Append(statusDiags...)
	model.ResourcesStatus = statusMap

	reasonMap, reasonDiags := types.MapValueFrom(ctx, types.StringType, resourcesReason)
	diags.Append(reasonDiags...)
	model.ResourcesReason = reasonMap
}
