package db

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/framework/planmodifiers"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	backupCreateTimeout = 30 * time.Minute
	backupDeleteTimeout = 30 * time.Minute
)

const (
	dbBackupStatusBuild    = "BUILDING"
	dbBackupStatusNew      = "NEW"
	dbBackupStatusActive   = "COMPLETED"
	dbBackupStatusError    = "ERROR"
	dbBackupStatusToDelete = "TO_DELETE"
	dbBackupStatusDeleted  = "DELETED"
)

var (
	_ resource.Resource                = &BackupResource{}
	_ resource.ResourceWithConfigure   = &BackupResource{}
	_ resource.ResourceWithImportState = &BackupResource{}
)

func NewBackupResource() resource.Resource {
	return &BackupResource{}
}

type BackupResource struct {
	config clients.Config
}

type BackupResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Region types.String `tfsdk:"region"`

	ContainerPrefix types.String  `tfsdk:"container_prefix"`
	Created         types.String  `tfsdk:"created"`
	Datastore       types.List    `tfsdk:"datastore"`
	DbmsID          types.String  `tfsdk:"dbms_id"`
	DbmsType        types.String  `tfsdk:"dbms_type"`
	Description     types.String  `tfsdk:"description"`
	LocationRef     types.String  `tfsdk:"location_ref"`
	Meta            types.String  `tfsdk:"meta"`
	Name            types.String  `tfsdk:"name"`
	Size            types.Float64 `tfsdk:"size"`
	Updated         types.String  `tfsdk:"updated"`
	WalSize         types.Float64 `tfsdk:"wal_size"`

	Timeouts timeouts.Value `tfsdk:"timeouts"`
}

func (r *BackupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "vkcs_db_backup"
}

func (r *BackupResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	s := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the resource.",
			},

			"region": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(planmodifiers.GetRegionPlanModifier(resp),
						"require replacement if configuration value changes", "require replacement if configuration value changes"),
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "The region in which to obtain the service client. If omitted, the `region` argument of the provider is used. Changing this creates a new backup.",
			},

			"container_prefix": schema.StringAttribute{
				Optional:    true,
				Description: "Prefix of S3 bucket ([prefix] - [project_id]) to store backup data. Default: databasebackups",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Backup creation timestamp",
			},

			"datastore": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Version of the datastore. Changing this creates a new instance.",
						},

						"version": schema.StringAttribute{
							Computed:    true,
							Description: "Type of the datastore. Changing this creates a new instance.",
						},
					},
				},
				Computed:    true,
				Description: "Object that represents datastore of backup",
			},

			"dbms_id": schema.StringAttribute{
				Required:    true,
				Description: "ID of the instance or cluster, to create backup of.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"dbms_type": schema.StringAttribute{
				Computed:    true,
				Description: "Type of dbms for the backup, can be \"instance\" or \"cluster\".",
			},

			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The description of the backup",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				Required:    true,
				Description: "The name of the backup. Changing this creates a new backup",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
		Description: "Provides a db backup resource. This can be used to create and delete db backup.",
	}

	if s.Blocks == nil {
		s.Blocks = make(map[string]schema.Block)
	}
	s.Blocks["timeouts"] = timeouts.Block(ctx, timeouts.Opts{
		Create: true,
		Delete: true,
	})

	resp.Schema = s
}

func (r *BackupResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.config = req.ProviderData.(clients.Config)
}

func (r *BackupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Create(ctx, backupCreateTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Consider adding "region" attribute if it is not present.
	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.DatabaseV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Databases API client", err.Error())
		return
	}

	dbmsID := data.DbmsID.ValueString()
	ctx = tflog.SetField(ctx, "dbms_id", dbmsID)

	tflog.Debug(ctx, "Calling Databases API to get DBMS resource")

	dbmsResource, err := getDBMSResource(client, dbmsID)
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to get DBMS resource", map[string]interface{}{"dbms_resource": fmt.Sprintf("%#v", dbmsResource)})

	var dbmsType string
	if instanceResource, ok := dbmsResource.(*instances.InstanceResp); ok {
		if util.IsOperationNotSupported(instanceResource.DataStore.Type, Redis, Tarantool) {
			resp.Diagnostics.AddError("Unable to create a backup of the resource",
				"Operation is not supported for this datastore")
			return
		}
		if instanceResource.ReplicaOf != nil {
			resp.Diagnostics.AddError("Unable to create a backup of the resource",
				"Operation is not supported for a replica")
			return
		}
		dbmsType = db.DBMSTypeInstance
	} else if clusterResource, ok := dbmsResource.(*clusters.ClusterResp); ok {
		if util.IsOperationNotSupported(clusterResource.DataStore.Type, Redis, Tarantool) {
			resp.Diagnostics.AddError("Unable to create a backup of the resource",
				"Operation is not supported for this datastore")
			return
		}
		dbmsType = db.DBMSTypeCluster
	}

	ctx = tflog.SetField(ctx, "dbms_type", dbmsType)

	opts := backups.Backup{
		Backup: &backups.BackupCreateOpts{
			Name:            data.Name.ValueString(),
			Description:     data.Description.ValueString(),
			ContainerPrefix: data.ContainerPrefix.ValueString(),
		},
	}
	if dbmsType == db.DBMSTypeInstance {
		opts.Backup.Instance = dbmsID
	} else {
		opts.Backup.Cluster = dbmsID
	}

	tflog.Debug(ctx, "Calling Databases API to create the backup", map[string]interface{}{"opts": fmt.Sprintf("%#v", opts)})

	backup, err := backups.Create(client, &opts).Extract()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API to create a backup", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to create the backup", map[string]interface{}{"backup": fmt.Sprintf("%#v", backup)})

	id := backup.ID
	resp.State.SetAttribute(ctx, path.Root("id"), id)
	ctx = tflog.SetField(ctx, "id", backup.ID)

	tflog.Debug(ctx, "Waiting for backup to become active")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{dbBackupStatusNew, dbBackupStatusBuild},
		Target:     []string{dbBackupStatusActive},
		Refresh:    backupStateRefreshFunc(client, backup.ID),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError("Error waiting for the backup to become ready", err.Error())
		return
	}

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.Created = types.StringValue(backup.Created)
	data.Datastore = flattenBackupDatastore(ctx, *backup.Datastore, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if backup.ClusterID != "" {
		data.DbmsID = types.StringValue(backup.ClusterID)
		data.DbmsType = types.StringValue(db.DBMSTypeCluster)
	} else {
		data.DbmsID = types.StringValue(backup.InstanceID)
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

func (r *BackupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Consider adding "region" attribute if it is not present.
	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.DatabaseV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Databases API client", err.Error())
		return
	}

	id := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "id", id)

	tflog.Debug(ctx, "Calling Databases API to read the backup")

	backup, err := backups.Get(client, id).Extract()
	if errutil.IsNotFound(err) {
		resp.State.RemoveResource(ctx)
		return
	}
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to read the backup", map[string]interface{}{"backup": fmt.Sprintf("%#v", backup)})

	data.ID = types.StringValue(id)
	data.Region = types.StringValue(region)
	data.Created = types.StringValue(backup.Created)
	data.Datastore = flattenBackupDatastore(ctx, *backup.Datastore, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	if backup.ClusterID != "" {
		data.DbmsID = types.StringValue(backup.ClusterID)
		data.DbmsType = types.StringValue(db.DBMSTypeCluster)
	} else {
		data.DbmsID = types.StringValue(backup.InstanceID)
		data.DbmsType = types.StringValue(db.DBMSTypeInstance)
	}
	data.LocationRef = types.StringValue(backup.LocationRef)
	data.Meta = types.StringValue(backup.Meta)
	data.Name = types.StringValue(backup.Name)
	data.Size = types.Float64Value(backup.Size)
	data.Updated = types.StringValue(backup.Updated)
	data.WalSize = types.Float64Value(backup.WalSize)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *BackupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.AddError("Unable to update the resource",
		"Not implemented. Please report this issue to the provider developers.")
}

func (r *BackupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data BackupResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	timeout, diags := data.Timeouts.Delete(ctx, backupDeleteTimeout)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	region := data.Region.ValueString()
	if region == "" {
		region = r.config.GetRegion()
	}

	client, err := r.config.DatabaseV1Client(region)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VKCS Databases API client", err.Error())
		return
	}

	id := data.ID.ValueString()
	ctx = tflog.SetField(ctx, "id", id)

	tflog.Debug(ctx, "Calling Databases API to delete the backup")

	err = backups.Delete(client, id).ExtractErr()
	if err != nil {
		resp.Diagnostics.AddError("Error calling VKCS Databases API", err.Error())
		return
	}

	tflog.Debug(ctx, "Called Databases API to delete the backup")

	stateConf := &retry.StateChangeConf{
		Pending:    []string{dbBackupStatusActive, dbBackupStatusToDelete},
		Target:     []string{dbBackupStatusDeleted},
		Refresh:    backupStateRefreshFunc(client, id),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		resp.Diagnostics.AddError("Error waiting for the backup deletion", err.Error())
		return
	}
}

func (r *BackupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
