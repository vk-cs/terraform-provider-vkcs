package db

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

var (
	_ datasource.DataSource              = &BackupDataSource{}
	_ datasource.DataSourceWithConfigure = &BackupDataSource{}
)

func NewBackupDataSource() datasource.DataSource {
	return &BackupDataSource{}
}

type BackupDataSource struct {
	config clients.Config
}

type BackupDataSourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	BackupID    types.String           `tfsdk:"backup_id"`
	Created     types.String           `tfsdk:"created"`
	Datastore   []BackupDataStoreModel `tfsdk:"datastore"`
	DbmsID      types.String           `tfsdk:"dbms_id"`
	DbmsType    types.String           `tfsdk:"dbms_type"`
	Description types.String           `tfsdk:"description"`
	LocationRef types.String           `tfsdk:"location_ref"`
	Meta        types.String           `tfsdk:"meta"`
	Name        types.String           `tfsdk:"name"`
	Size        types.Float64          `tfsdk:"size"`
	Updated     types.String           `tfsdk:"updated"`
	WalSize     types.Float64          `tfsdk:"wal_size"`
}

type BackupDataStoreModel struct {
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
}

func (d *BackupDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "vkcs_db_backup"
}

func (d *BackupDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used.",
			},

			"backup_id": schema.StringAttribute{
				Required:    true,
				Description: "The UUID of the backup.",
			},

			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Backup creation timestamp",
			},

			// TODO: change type to SingleNestedAttribute #BREAKING_CHANGE
			"datastore": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Required:    true,
							Description: "Type of the datastore.",
						},

						"version": schema.StringAttribute{
							Required:    true,
							Description: "Version of the datastore.",
						},
					},
				},
				Computed:    true,
				Description: "Object that represents datastore of backup",
			},

			"dbms_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the backed up instance or cluster",
			},

			"dbms_type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of dbms of the backup, can be \"instance\" or \"cluster\".",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Description: "The description of the backup",
			},

			"location_ref": schema.StringAttribute{
				Computed:    true,
				Description: "Location of backup data on backup storage",
			},

			"meta": schema.StringAttribute{
				Computed:    true,
				Description: "Metadata of the backup",
			},

			"name": schema.StringAttribute{
				Computed:    true,
				Description: "The name of the backup.",
			},

			"size": schema.Float64Attribute{
				Computed:    true,
				Description: "Backup's volume size",
			},

			"updated": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp of backup's last update",
			},

			"wal_size": schema.Float64Attribute{
				Computed:    true,
				Description: "Backup's WAL volume size",
			},
		},
		Description: "Use this data source to get the information on a db backup resource.",
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

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = d.config.GetRegion()
	}

	client, err := d.config.DatabaseV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Databases API client", err.Error())
		return
	}

	tflog.Debug(ctx, "Calling Databases API to get the backup")

	backup, err := backups.Get(client, data.BackupID.ValueString()).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to get the backup", map[string]interface{}{"backup": fmt.Sprintf("%#v", backup)})

	data.ID = types.StringValue(backup.ID)
	data.Region = types.StringValue(region)
	data.Created = types.StringValue(backup.Created)
	data.Datastore = flattenBackupDatastore(*backup.Datastore)
	if backup.InstanceID != "" {
		data.DbmsID = types.StringValue(backup.InstanceID)
		data.DbmsType = types.StringValue(db.DBMSTypeInstance)
	} else {
		data.DbmsID = types.StringValue(backup.ClusterID)
		data.DbmsType = types.StringValue(db.DBMSTypeInstance)
	}
	data.Description = types.StringValue(backup.Description)
	data.LocationRef = types.StringValue(backup.LocationRef)
	data.Meta = types.StringValue(backup.Meta)
	data.Name = types.StringValue(backup.Name)
	data.Size = types.Float64Value(backup.Size)
	data.Updated = types.StringValue(backup.Updated)
	data.WalSize = types.Float64Value(backup.WalSize)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func flattenBackupDatastore(d datastores.DatastoreShort) []BackupDataStoreModel {
	return []BackupDataStoreModel{
		{
			Type:    types.StringValue(d.Type),
			Version: types.StringValue(d.Version),
		},
	}
}
