package backup

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/objectvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/triggers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"golang.org/x/exp/maps"
)

const (
	planStatusDeleted = "deleted"
)

// Ensure the implementation satisfies the desired interfaces.
var _ resource.Resource = &PlanResource{}

func NewPlanResource() resource.Resource {
	return &PlanResource{}
}

type PlanResource struct {
	config clients.Config
}

func (r *PlanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_backup_plan"
}

type PlanResourceModel struct {
	ID                types.String                    `tfsdk:"id"`
	Name              types.String                    `tfsdk:"name"`
	Schedule          *PlanResourceScheduleModel      `tfsdk:"schedule"`
	FullRetention     *PlanResourceFullRetentionModel `tfsdk:"full_retention"`
	GFSRetention      *PlanResourceGFSRetentionModel  `tfsdk:"gfs_retention"`
	IncrementalBackup types.Bool                      `tfsdk:"incremental_backup"`
	ProviderID        types.String                    `tfsdk:"provider_id"`
	ProviderName      types.String                    `tfsdk:"provider_name"`
	InstanceIDs       types.Set                       `tfsdk:"instance_ids"`
	Region            types.String                    `tfsdk:"region"`
	BackupTargets     types.Set                       `tfsdk:"backup_targets"`
}

type PlanResourceScheduleModel struct {
	Date       []types.String `tfsdk:"date"`
	Time       types.String   `tfsdk:"time"`
	EveryHours types.Int64    `tfsdk:"every_hours"`
}

type PlanResourceFullRetentionModel struct {
	MaxFullBackup types.Int64 `tfsdk:"max_full_backup"`
}

type PlanResourceGFSRetentionModel struct {
	GFSWeekly  types.Int64 `tfsdk:"gfs_weekly"`
	GFSMonthly types.Int64 `tfsdk:"gfs_monthly"`
	GFSYearly  types.Int64 `tfsdk:"gfs_yearly"`
}

type PlanResourceBackupTargetModel struct {
	InstanceID types.String `tfsdk:"instance_id"`
	VolumeIDs  types.Set    `tfsdk:"volume_ids"`
}

func (PlanResourceBackupTargetModel) ObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"instance_id": types.StringType,
			"volume_ids":  types.SetType{ElemType: types.StringType},
		},
	}
}

func (r *PlanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of the backup plan",
			},

			"schedule": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"date": schema.ListAttribute{
						ElementType: types.StringType,
						Optional:    true,
						Validators: []validator.List{
							listvalidator.ValueStringsAre(
								stringvalidator.OneOf([]string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}...),
							),
							listvalidator.ConflictsWith(path.Expressions{
								path.MatchRoot("schedule").AtName("every_hours"),
							}...),
							listvalidator.AlsoRequires(path.Expressions{
								path.MatchRoot("schedule").AtName("time"),
							}...),
							listvalidator.SizeAtLeast(1),
						},
						Description: "List of days when to perform backups. If incremental_backups is enabled, only one day should be specified",
					},
					"time": schema.StringAttribute{
						Optional:    true,
						Description: "Time of backup in format hh:mm (for UTC timezone) or hh:mm+tz (for other timezones, e.g. 10:00+03 for MSK, 10:00-04 for ET)",
						Validators: []validator.String{
							stringvalidator.ConflictsWith(path.Expressions{
								path.MatchRoot("schedule").AtName("every_hours"),
							}...),
							stringvalidator.AlsoRequires(path.Expressions{
								path.MatchRoot("schedule").AtName("date"),
							}...),
							timeValidator{},
						},
					},
					"every_hours": schema.Int64Attribute{
						Optional: true,
						Validators: []validator.Int64{
							int64validator.OneOf([]int64{3, 12, 24}...),
							int64validator.ConflictsWith(path.Expressions{
								path.MatchRoot("schedule").AtName("date"),
							}...),
							int64validator.ConflictsWith(path.Expressions{
								path.MatchRoot("schedule").AtName("time"),
							}...),
						},
						Description: "Hour interval of backups, must be one of: 3, 12, 24. This field is incompatible with date/time fields",
					},
				},
			},

			"full_retention": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"max_full_backup": schema.Int64Attribute{
						Required:    true,
						Description: "Maximum number of backups",
					},
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("gfs_retention"),
					}...),
				},
				Description: "Parameters for full retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full. Incompatible with gfs_retention",
			},

			"gfs_retention": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"gfs_weekly": schema.Int64Attribute{
						Required:    true,
						Description: "Number of weeks to store backups",
					},
					"gfs_monthly": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of months to store backups",
					},
					"gfs_yearly": schema.Int64Attribute{
						Optional:    true,
						Description: "Number of years to store backups",
					},
				},
				Validators: []validator.Object{
					objectvalidator.ExactlyOneOf(path.Expressions{
						path.MatchRoot("full_retention"),
					}...),
				},
				Description: "Parameters for gfs retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full. Incompatible with full_retention",
			},

			"incremental_backup": schema.BoolAttribute{
				Required:    true,
				Description: "Whether incremental backups strategy should be used. If enabled, the schedule.date field must specify one day, on which full backup will be created. On other days, incremental backups will be created. _note_ This option may be enabled for only for 'cloud_servers' provider.",
			},

			"provider_id": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "ID of backup provider",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("provider_name"),
					}...),
				},
			},

			"provider_name": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: fmt.Sprintf("Name of backup provider, must be one of: %s", strings.Join(getProviderNames(), ", ")),
				Validators: []validator.String{
					stringvalidator.OneOf(maps.Values(providerNameMapping)...),
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("provider_id"),
					}...),
				},
			},

			"backup_targets": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"instance_id": schema.StringAttribute{
							Required:    true,
							Description: "ID of the instance for which specific volumes are backed up.",
						},
						"volume_ids": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Description: "Set of volume IDs to back up for the instance. If no list is specified, backups will be created for all disks.",
						},
					},
				},
				Description: "Set of backup targets specifying instance_id and volume_ids for each instance. Either backup_targets or instance_ids must be specified, but not both.",
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("instance_ids"),
					}...),
				},
			},

			"instance_ids": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
				Description: "Set of ids of instances to make backup for. Either backup_targets or instance_ids must be specified, but not both.",
				Validators: []validator.Set{
					setvalidator.ConflictsWith(path.Expressions{
						path.MatchRoot("backup_targets"),
					}...),
				},
			},

			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
						"require replacement if configuration value changes", "require replacement if configuration value changes"),
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`. Changing this creates a new plan.",
			},
		},
		Description: "Manages a backup plan resource.",
	}
}

func (r *PlanResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *PlanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan PlanResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	backupClient, err := r.config.BackupV1Client(region, r.config.GetTenantID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	providerInfo, err := findProvider(backupClient, plan.ProviderID.ValueString(), plan.ProviderName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_backup_plan", err.Error())
		return
	}

	var backupTargets []PlanResourceBackupTargetModel
	if !plan.BackupTargets.IsNull() && !plan.BackupTargets.IsUnknown() {
		diags := plan.BackupTargets.ElementsAs(ctx, &backupTargets, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	var instanceIds []types.String
	if len(backupTargets) > 0 {
		for _, target := range backupTargets {
			instanceIds = append(instanceIds, target.InstanceID)
		}
	} else {
		diags := plan.InstanceIDs.ElementsAs(ctx, &instanceIds, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resourcesInfo, err := getResourcesInfo(r.config, region, instanceIds, providerInfo.Name)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_backup_plan", err.Error())
		return
	}

	if providerInfo.Name == ProviderNameNova && len(backupTargets) > 0 {
		resourcesInfo, diags = enrichWithVolumes(ctx, resourcesInfo, backupTargets)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	planCreateOpts := plans.CreateOpts{
		Name:       plan.Name.ValueString(),
		Resources:  resourcesInfo,
		ProviderID: providerInfo.ID,
	}

	incrementalBackups := plan.IncrementalBackup.ValueBool()
	if incrementalBackups {
		fullDay, err := expandIncrementalFullDay(plan)
		if err != nil {
			resp.Diagnostics.AddError("Error creating vkcs_backup_plan", err.Error())
			return
		}
		planCreateOpts.FullDay = &fullDay
	}

	if plan.FullRetention != nil {
		planCreateOpts.RetentionType = RetentionFull
	}
	if plan.GFSRetention != nil {
		planCreateOpts.RetentionType = RetentionGFS
		planCreateOpts.GFS = expandGFS(plan)
	}

	triggerCreateOpts := triggers.CreateOpts{
		Name: fmt.Sprintf("%s_trigger", plan.Name.ValueString()),
	}
	triggerProperties := triggers.PropertiesOpts{}

	triggerSchedule, err := expandTriggerSchedule(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error getting schedule", err.Error())
		return
	}
	triggerProperties.Pattern = triggerSchedule

	if plan.FullRetention != nil {
		triggerProperties.MaxBackups = int(plan.FullRetention.MaxFullBackup.ValueInt64())
	} else {
		// Send default value for MaxBackups because backend requires this field
		triggerProperties.MaxBackups = 30
	}

	planResp, err := plans.Create(backupClient, &plans.PlanCreate{Plan: &planCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_backup_plan", err.Error())
		return
	}
	planID := planResp.ID
	resp.State.SetAttribute(ctx, path.Root("id"), planID)

	triggerCreateOpts.PlanID = planID
	triggerCreateOpts.Properties = &triggerProperties
	_, err = triggers.Create(backupClient, &triggers.TriggerCreate{TriggerInfo: &triggerCreateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error creating vkcs_backup_plan", err.Error())
		return
	}

	plan.ID = types.StringValue(planID)
	plan.Region = types.StringValue(region)
	plan.ProviderID = types.StringValue(providerInfo.ID)
	plan.ProviderName = types.StringValue(providerInfo.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PlanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan PlanResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	backupClient, err := r.config.BackupV1Client(region, r.config.GetTenantID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	planID := plan.ID.ValueString()

	planResp, err := plans.Get(backupClient, planID).Extract()
	if err != nil {
		checkDeleted := util.CheckDeletedResource(ctx, resp, err)
		if checkDeleted != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", checkDeleted.Error())
		}
		return
	}
	if planResp.Status == planStatusDeleted {
		resp.State.RemoveResource(ctx)
		return
	}

	triggerResp, err := findTrigger(backupClient, planID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
		return
	}

	providerInfo, err := findProvider(backupClient, plan.ProviderID.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
		return
	}
	plan.ProviderID = types.StringValue(providerInfo.ID)
	plan.ProviderName = types.StringValue(providerInfo.Name)

	plan.Name = types.StringValue(planResp.Name)

	if !plan.BackupTargets.IsNull() && !plan.BackupTargets.IsUnknown() {
		resources := make([]PlanResourceBackupTargetModel, len(planResp.Resources))
		for i, respResource := range planResp.Resources {
			var volumeIDsSet types.Set
			if len(respResource.Resources) > 0 {
				volumeIDs := make([]string, len(respResource.Resources))
				for j, resource := range respResource.Resources {
					volumeIDs[j] = resource.ID
				}
				var diags diag.Diagnostics
				volumeIDsSet, diags = types.SetValueFrom(ctx, types.StringType, volumeIDs)
				if diags.HasError() {
					resp.Diagnostics.Append(diags...)
					return
				}
			} else {
				volumeIDsSet = types.SetNull(types.StringType)
			}

			resources[i] = PlanResourceBackupTargetModel{
				InstanceID: types.StringValue(respResource.ID),
				VolumeIDs:  volumeIDsSet,
			}
		}

		backupTargetsSet, diags := types.SetValueFrom(ctx, PlanResourceBackupTargetModel{}.ObjectType(), resources)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.BackupTargets = backupTargetsSet
	} else {
		resources := make([]types.String, len(planResp.Resources))
		for i, respResource := range planResp.Resources {
			resources[i] = types.StringValue(respResource.ID)
		}

		instanceSet, diags := types.SetValueFrom(ctx, types.StringType, resources)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
		plan.InstanceIDs = instanceSet
	}

	plan.Region = types.StringValue(region)

	if planResp.FullDay != nil {
		plan.IncrementalBackup = types.BoolValue(true)
	} else {
		plan.IncrementalBackup = types.BoolValue(false)
	}

	switch planResp.RetentionType {
	case RetentionFull:
		fullRetention := PlanResourceFullRetentionModel{
			MaxFullBackup: types.Int64Value(int64(triggerResp.Properties.MaxBackups)),
		}
		plan.FullRetention = &fullRetention
		plan.GFSRetention = nil
	case RetentionGFS:
		gfsRetention := flattenGFS(planResp)
		plan.GFSRetention = gfsRetention
		plan.FullRetention = nil
	}

	var location *time.Location
	if plan.Schedule != nil && !plan.Schedule.Time.IsNull() {
		oldTimeParsed, err := parseTime(plan.Schedule.Time.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
			return
		}
		location = oldTimeParsed.Location()
	} else {
		location = time.UTC
	}
	schedule := flattenSchedule(planResp, triggerResp, location)
	plan.Schedule = schedule

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *PlanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PlanResourceModel
	var state PlanResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	backupClient, err := r.config.BackupV1Client(region, r.config.GetTenantID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	planID := state.ID.ValueString()

	triggerResp, err := findTrigger(backupClient, planID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
		return
	}
	triggerID := triggerResp.ID

	providerInfo, err := findProvider(backupClient, plan.ProviderID.ValueString(), plan.ProviderName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_backup_plan", err.Error())
		return
	}

	var backupTargets []PlanResourceBackupTargetModel
	if !plan.BackupTargets.IsNull() && !plan.BackupTargets.IsUnknown() {
		diags := plan.BackupTargets.ElementsAs(ctx, &backupTargets, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	var instanceIds []types.String
	if len(backupTargets) > 0 {
		for _, target := range backupTargets {
			instanceIds = append(instanceIds, target.InstanceID)
		}
	} else {
		diags := plan.InstanceIDs.ElementsAs(ctx, &instanceIds, false)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	resourcesInfo, err := getResourcesInfo(r.config, region, instanceIds, providerInfo.Name)
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_backup_plan", err.Error())
		return
	}

	if providerInfo.Name == ProviderNameNova && len(backupTargets) > 0 {
		resourcesInfo, diags = enrichWithVolumes(ctx, resourcesInfo, backupTargets)
		if diags.HasError() {
			resp.Diagnostics.Append(diags...)
			return
		}
	}

	planUpdateOpts := plans.UpdateOpts{
		Name:      plan.Name.ValueString(),
		Status:    "running",
		Resources: resourcesInfo,
	}

	incrementalBackups := plan.IncrementalBackup.ValueBool()
	if incrementalBackups {
		fullDay, err := expandIncrementalFullDay(plan)
		if err != nil {
			resp.Diagnostics.AddError("Error updating vkcs_backup_plan", err.Error())
			return
		}
		planUpdateOpts.FullDay = fullDay
	}

	if plan.FullRetention != nil {
		planUpdateOpts.RetentionType = RetentionFull
	}
	if plan.GFSRetention != nil {
		planUpdateOpts.RetentionType = RetentionGFS
		planUpdateOpts.GFS = expandGFS(plan)
	}

	triggerUpdateOpts := triggers.UpdateOpts{
		Name: fmt.Sprintf("%s_trigger", plan.Name.ValueString()),
	}
	if plan.FullRetention != nil {
		triggerUpdateOpts.MaxBackups = int(plan.FullRetention.MaxFullBackup.ValueInt64())
	} else {
		// Send default value for MaxBackups because backend requires this field
		triggerUpdateOpts.MaxBackups = 30
	}
	triggerSchedule, err := expandTriggerSchedule(plan)
	if err != nil {
		resp.Diagnostics.AddError("Invalid resource schema", "One of full_retention, gfs_retention must be specified")
		return
	}
	triggerUpdateOpts.Pattern = triggerSchedule

	planResp, err := plans.Update(backupClient, planID, &plans.PlanUpdate{Plan: &planUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_backup_plan", err.Error())
		return
	}

	_, err = triggers.Update(backupClient, triggerID, &triggers.TriggerUpdate{TriggerInfo: &triggerUpdateOpts}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error updating vkcs_backup_plan", err.Error())
		return
	}
	plan.ID = types.StringValue(planResp.ID)
	plan.Region = types.StringValue(region)
	plan.ProviderID = types.StringValue(providerInfo.ID)
	plan.ProviderName = types.StringValue(providerInfo.Name)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *PlanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan PlanResourceModel

	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := plan.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	backupClient, err := r.config.BackupV1Client(region, r.config.GetTenantID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	planID := plan.ID.ValueString()

	err = plans.Delete(backupClient, planID).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Unable to delete resource vkcs_backup_plan", err.Error())
		return
	}
}

func (r *PlanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *PlanResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("schedule").AtName("date"),
			path.MatchRoot("schedule").AtName("time"),
			path.MatchRoot("schedule").AtName("every_hours"),
		),
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("instance_ids"),
			path.MatchRoot("backup_targets"),
		),
	}
}
