package backup

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/backup/v1/restores"
)

var _ datasource.DataSource = &RestoreDataSource{}

func NewRestoreDataSource() datasource.DataSource {
	return &RestoreDataSource{}
}

type RestoreDataSource struct {
	config clients.Config
}

type RestoreDataSourceModel struct {
	ID              types.String `tfsdk:"id"`
	RestoreID       types.String `tfsdk:"restore_id"`
	CheckpointID    types.String `tfsdk:"checkpoint_id"`
	ProviderID      types.String `tfsdk:"provider_id"`
	ResourcesStatus types.Map    `tfsdk:"resources_status"`
	ResourcesReason types.Map    `tfsdk:"resources_reason"`
	Region          types.String `tfsdk:"region"`
}

func (d *RestoreDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_backup_restore"
}

func (d *RestoreDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the restore operation.",
			},
			"restore_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the restore operation to retrieve.",
			},
			"checkpoint_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the restored backup checkpoint.",
			},
			"provider_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the backup provider used for restoration.",
			},
			"resources_status": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Restore status for each affected resource.",
			},
			"resources_reason": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
				Description: "Reason for each resource that could not be restored.",
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region from which to obtain the backup service endpoint. Defaults to the provider region.",
			},
		},
		Description: "Use this data source to read a restore operation and its per-resource statuses.",
	}
}

func (d *RestoreDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *RestoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data RestoreDataSourceModel
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

	restore, err := restores.Get(backupClient, data.RestoreID.ValueString()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_backup_restore", err.Error())
		return
	}

	data.ID = types.StringValue(restore.ID)
	data.Region = types.StringValue(region)
	data.CheckpointID = types.StringValue(restore.CheckpointID)
	data.ProviderID = types.StringValue(restore.ProviderID)
	state := RestoreResourceModel{}
	setRestoreStatuses(ctx, &state, restore, &resp.Diagnostics)
	data.ResourcesStatus = state.ResourcesStatus
	data.ResourcesReason = state.ResourcesReason
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
