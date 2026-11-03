package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/checkpoints"
)

var _ datasource.DataSource = &CheckpointsDataSource{}

func NewCheckpointsDataSource() datasource.DataSource {
	return &CheckpointsDataSource{}
}

type CheckpointsDataSource struct {
	config clients.Config
}

type CheckpointsDataSourceModel struct {
	ID        types.String                              `tfsdk:"id"`
	PlanID    types.String                              `tfsdk:"plan_id"`
	ProjectID types.String                              `tfsdk:"project_id"`
	Resources []CheckpointsDataSourceResourceModel      `tfsdk:"resources"`
	Region    types.String                              `tfsdk:"region"`
}

type CheckpointsDataSourceResourceModel struct {
	ID           types.String                         `tfsdk:"id"`
	Name         types.String                         `tfsdk:"name"`
	ResourceType types.String                         `tfsdk:"resource_type"`
	Backups      []CheckpointsDataSourceBackupModel   `tfsdk:"backups"`
}

type CheckpointsDataSourceBackupModel struct {
	CheckpointID types.String                 `tfsdk:"checkpoint_id"`
	Status       types.String                 `tfsdk:"status"`
	CreatedAt    types.String                 `tfsdk:"created_at"`
	ObjectLock   *PlanResourceObjectLockModel `tfsdk:"object_lock"`
}

func (d *CheckpointsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_backup_checkpoints"
}

func (d *CheckpointsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the backup plan used to query checkpoints.",
			},
			"plan_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the backup plan whose checkpoints should be retrieved.",
			},
			"project_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the project containing the plan and its checkpoints.",
			},
			"resources": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Resources protected by the plan, with their backup checkpoints.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the protected resource.",
						},
						"name": schema.StringAttribute{
							Computed:    true,
							Description: "Name of the protected resource.",
						},
						"resource_type": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the protected resource.",
						},
						"backups": schema.ListNestedAttribute{
							Computed:    true,
							Description: "Backup checkpoints created for this resource.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"checkpoint_id": schema.StringAttribute{
										Computed:    true,
										Description: "ID of the backup checkpoint.",
									},
									"status": schema.StringAttribute{
										Computed:    true,
										Description: "Status of the backup checkpoint.",
									},
									"created_at": schema.StringAttribute{
										Computed:    true,
										Description: "Checkpoint creation timestamp.",
									},
									"object_lock": schema.SingleNestedAttribute{
										Computed:    true,
										Description: "Object Lock protection for the checkpoint; read-only.",
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												Computed:    true,
												Description: "Object Lock mode: compliance, governance, or legalhold.",
											},
											"days": schema.Int64Attribute{
												Computed:    true,
												Description: "Retention period in days, if configured.",
											},
											"retain_until": schema.Int64Attribute{
												Computed:    true,
												Description: "Retention end date as a Unix timestamp, if configured.",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region in which the backup plan exists. Defaults to the provider region.",
			},
		},
		Description: "Use this data source to read backup checkpoints grouped by plan. The object_lock values are read-only.",
	}
}

func (d *CheckpointsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *CheckpointsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data CheckpointsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	backupClient, err := d.config.BackupV1Client(region, d.config.GetProjectID())
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS backup client", err.Error())
		return
	}

	response, err := checkpoints.ListGroupedByPlanID(backupClient, checkpoints.ListOpts{PlanID: data.PlanID.ValueString()}).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_checkpoints", err.Error())
		return
	}

	data.ID = types.StringValue(response.PlanID)
	data.ProjectID = types.StringValue(response.ProjectID)
	data.Region = types.StringValue(region)
	data.Resources = flattenCheckpointResources(response.Resources)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenCheckpointResources(resources []checkpoints.CheckpointResource) []CheckpointsDataSourceResourceModel {
	flattened := make([]CheckpointsDataSourceResourceModel, len(resources))
	for resourceIndex, resource := range resources {
		backups := make([]CheckpointsDataSourceBackupModel, len(resource.Backups))
		for backupIndex, backup := range resource.Backups {
			backups[backupIndex] = CheckpointsDataSourceBackupModel{
				CheckpointID: types.StringValue(backup.CheckpointID),
				Status:       types.StringValue(backup.Status),
				CreatedAt:    types.StringValue(backup.CreatedAt),
				ObjectLock:   flattenObjectLock(backup.ObjectLock),
			}
		}
		flattened[resourceIndex] = CheckpointsDataSourceResourceModel{
			ID:           types.StringValue(resource.ID),
			Name:         types.StringValue(resource.Name),
			ResourceType: types.StringValue(resource.ResourceType),
			Backups:      backups,
		}
	}
	return flattened
}
