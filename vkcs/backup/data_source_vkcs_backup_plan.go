package backup

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/plans"
)

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &PlanDataSource{}

func NewPlanDataSource() datasource.DataSource {
	return &PlanDataSource{}
}

type PlanDataSource struct {
	config clients.Config
}

func (d *PlanDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_backup_plan"
}

type PlanDataSourceModel struct {
	ID                types.String                    `tfsdk:"id"`
	Name              types.String                    `tfsdk:"name"`
	InstanceID        types.String                    `tfsdk:"instance_id"`
	Schedule          *PlanResourceScheduleModel      `tfsdk:"schedule"`
	FullRetention     *PlanResourceFullRetentionModel `tfsdk:"full_retention"`
	GFSRetention      *PlanResourceGFSRetentionModel  `tfsdk:"gfs_retention"`
	IncrementalBackup types.Bool                      `tfsdk:"incremental_backup"`
	ProviderID        types.String                    `tfsdk:"provider_id"`
	InstanceIDs       []types.String                  `tfsdk:"instance_ids"`
	Region            types.String                    `tfsdk:"region"`
}

func (d *PlanDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource",
			},

			"name": schema.StringAttribute{
				Optional:    true,
				Description: "Name of the backup plan",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("instance_id"),
					}...),
				},
			},

			"instance_id": schema.StringAttribute{
				Optional:    true,
				Description: "ID of the instance that should be included in backup plan",
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.Expressions{
						path.MatchRoot("name"),
					}...),
				},
			},

			"schedule": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"date": schema.ListAttribute{
						ElementType: types.StringType,
						Computed:    true,
						Description: "List of days when to perform backups. If incremental_backups is enabled, this field contains day of full backup",
					},
					"time": schema.StringAttribute{
						Computed:    true,
						Description: "Time of backup in format hh:mm, using UTC timezone",
					},
					"every_hours": schema.Int64Attribute{
						Computed:    true,
						Description: "Hour period of backups",
					},
				},
			},

			"full_retention": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"max_full_backup": schema.Int64Attribute{
						Computed:    true,
						Description: "Maximum number of backups",
					},
				},
				Description: "Parameters for full retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full",
			},

			"gfs_retention": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"gfs_weekly": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of weeks to store backups",
					},
					"gfs_monthly": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of months to store backups",
					},
					"gfs_yearly": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of years to store backups",
					},
				},
				Description: "Parameters for gfs retention policy. Specifies number of full backups stored. Incremental backups (if enabled) are not counted as full",
			},

			"incremental_backup": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether incremental backups should be stored",
			},

			"provider_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of backup provider",
			},

			"instance_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Computed:    true,
				Description: "List of ids of backed up instances",
			},

			"region": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},
		},
		Description: "Use this data source to get backup plan info",
	}
}

func (d *PlanDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *PlanDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PlanDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	backupClient, err := d.config.BackupV1Client(region, d.config.GetTenantID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	allPages, err := plans.List(backupClient).AllPages()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
		return
	}
	allPlans, err := plans.ExtractPlans(allPages)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
		return
	}

	name := data.Name.ValueString()
	instanceID := data.InstanceID.ValueString()

	foundPlanCount := 0
	var planResp *plans.PlanResponse
	for _, pl := range allPlans {
		if name != "" && pl.Name != name {
			continue
		}
		if instanceID != "" {
			foundRes := false
			for _, res := range pl.Resources {
				if res.ID == instanceID {
					foundRes = true
					break
				}
			}
			if !foundRes {
				continue
			}
		}
		planResp = &pl
		foundPlanCount++
	}
	if foundPlanCount == 0 {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", "No suitable plans found")
		return
	} else if foundPlanCount > 1 {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", "More than one suitable plan found")
	}
	planID := planResp.ID
	data.ID = types.StringValue(planID)
	triggerResp, err := findTrigger(backupClient, planID)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_plan", err.Error())
		return
	}

	data.Name = types.StringValue(planResp.Name)
	data.ProviderID = types.StringValue(planResp.ProviderID)

	resources := make([]types.String, len(planResp.Resources))
	for i, respResource := range planResp.Resources {
		resources[i] = types.StringValue(respResource.ID)
	}
	data.InstanceIDs = resources
	data.Region = types.StringValue(region)

	if planResp.FullDay != nil {
		data.IncrementalBackup = types.BoolValue(true)
	} else {
		data.IncrementalBackup = types.BoolValue(false)
	}

	if planResp.RetentionType == RetentionFull {
		fullRetention := PlanResourceFullRetentionModel{
			MaxFullBackup: types.Int64Value(int64(triggerResp.Properties.MaxBackups)),
		}
		data.FullRetention = &fullRetention
		data.GFSRetention = nil
	} else if planResp.RetentionType == RetentionGFS {
		gfsRetention := flattenGFS(planResp)
		data.GFSRetention = gfsRetention
		data.FullRetention = nil
	}

	schedule := flattenSchedule(planResp, triggerResp, time.UTC)
	data.Schedule = schedule

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
