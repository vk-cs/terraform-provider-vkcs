package mlplatform

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/mlplatform/v1/backups"
)

// Ensure the implementation satisfies the desired interfaces.
var _ datasource.DataSource = &BackupDataSource{}

func NewBackupDataSource() datasource.DataSource {
	return &BackupDataSource{}
}

type BackupDataSource struct {
	config clients.Config
}

func (d *BackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_mlplatform_backup"
}

type BackupDataSourceModel struct {
	ID         types.String             `tfsdk:"id"`
	InstanceID types.String             `tfsdk:"instance_id"`
	Backups    []*MLPlatformBackupModel `tfsdk:"backups"`
	Region     types.String             `tfsdk:"region"`
}

type MLPlatformBackupModel struct {
	VolumeID  types.String `tfsdk:"volume_id"`
	CreatedAt types.String `tfsdk:"created_at"`
	BackupID  types.String `tfsdk:"backup_id"`
	Comment   types.String `tfsdk:"comment"`
}

func (d *BackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the data source",
			},

			"instance_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the instance to get backups for",
			},

			"backups": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"backup_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the backup",
						},
						"volume_id": schema.StringAttribute{
							Computed:    true,
							Description: "ID of the volume",
						},
						"created_at": schema.StringAttribute{
							Computed:    true,
							Description: "Creation timestamp",
						},
						"comment": schema.StringAttribute{
							Computed:    true,
							Description: "Backup comment",
						},
					},
				},
				Computed:    true,
				Description: "Backups info",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},
		},
		Description: "Use this data source to get backups for instance.",
	}
}

func (d *BackupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	d.config = req.ProviderData.(clients.Config)
}

func (d *BackupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data BackupDataSourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	mlPlatformClient, err := d.config.MLPlatformV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS ML Platform client", err.Error())
		return
	}

	backupResp, err := backups.Get(mlPlatformClient, data.InstanceID.ValueString()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving vkcs_mlplatform_backup", err.Error())
		return
	}

	data.ID = types.StringValue(uuid.New().String())

	data.Backups = flattenBackupOpts(backupResp)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
