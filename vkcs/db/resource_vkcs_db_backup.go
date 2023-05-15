package db

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/clusters"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	dbBackupDelay         = 10 * time.Second
	dbBackupMinTimeout    = 3 * time.Second
	dbBackupCreateTimeout = 30 * time.Minute
	dbBackupDeleteTimeout = 30 * time.Minute
)

type dbBackupStatus string

var (
	dbBackupStatusBuild   dbBackupStatus = "BUILDING"
	dbBackupStatusActive  dbBackupStatus = "COMPLETED"
	dbBackupStatusError   dbBackupStatus = "ERROR"
	dbBackupStatusDeleted dbBackupStatus = "DELETED"
)

func ResourceDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseBackupCreate,
		ReadContext:   resourceDatabaseBackupRead,
		DeleteContext: resourceDatabaseBackupDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(dbBackupCreateTimeout),
			Delete: schema.DefaultTimeout(dbBackupDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the backup. Changing this creates a new backup",
			},

			"dbms_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the instance or cluster, to create backup of.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The description of the backup",
			},

			"container_prefix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Prefix of S3 bucket ([prefix] - [project_id]) to store backup data. Default: databasebackups",
			},
			// Computed fields
			"dbms_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of dbms for the backup, can be \"instance\" or \"cluster\".",
			},

			"location_ref": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Location of backup data on backup storage",
			},

			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Backup creation timestamp",
			},

			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of backup's last update",
			},

			"size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Backup's volume size",
			},

			"wal_size": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "Backup's WAL volume size",
			},

			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of the datastore. Changing this creates a new instance.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Version of the datastore. Changing this creates a new instance.",
						},
					},
				},
				Description: "Object that represents datastore of backup",
			},

			"meta": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Metadata of the backup",
			},
		},
		Description: "Provides a db backup resource. This can be used to create and delete db backup.",
	}
}

func resourceDatabaseBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	dbmsID := d.Get("dbms_id").(string)

	dbmsResp, err := getDBMSResource(DatabaseV1Client, dbmsID)
	if err != nil {
		return diag.Errorf("error while getting resource: %s", err)
	}

	var dbmsType string
	if instanceResource, ok := dbmsResp.(*instances.InstanceResp); ok {
		if util.IsOperationNotSupported(instanceResource.DataStore.Type, Redis, Tarantool) {
			return diag.Errorf("operation not supported for this datastore")
		}
		if instanceResource.ReplicaOf != nil {
			return diag.Errorf("operation not supported for replica")
		}
		dbmsType = db.DBMSTypeInstance
	}
	if clusterResource, ok := dbmsResp.(*clusters.ClusterResp); ok {
		if util.IsOperationNotSupported(clusterResource.DataStore.Type, Redis, Tarantool) {
			return diag.Errorf("operation not supported for this datastore")
		}
		dbmsType = db.DBMSTypeCluster
	}

	b := backups.BackupCreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		ContainerPrefix: d.Get("container_prefix").(string),
	}

	if dbmsType == db.DBMSTypeInstance {
		b.Instance = d.Get("dbms_id").(string)
	} else {
		b.Cluster = d.Get("dbms_id").(string)
	}

	log.Printf("[DEBUG] vkcs_db_backup create options: %#v", b)

	back := backups.Backup{
		Backup: &b,
	}
	backup, err := backups.Create(DatabaseV1Client, &back).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_db_backup: %s", err)
	}

	// Wait for the backup to become available.
	log.Printf("[DEBUG] Waiting for vkcs_db_backup %s to become available", backup.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbBackupStatusBuild)},
		Target:     []string{string(dbBackupStatusActive)},
		Refresh:    databaseBackupStateRefreshFunc(DatabaseV1Client, backup.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      dbBackupDelay,
		MinTimeout: dbBackupMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_backup %s to become ready: %s", backup.ID, err)
	}

	// Store the ID now
	d.SetId(backup.ID)

	return resourceDatabaseBackupRead(ctx, d, meta)
}

func resourceDatabaseBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	backup, err := backups.Get(DatabaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_db_backup"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_backup %s: %#v", d.Id(), backup)

	d.Set("name", backup.Name)
	if backup.InstanceID != "" {
		d.Set("dbms_id", backup.InstanceID)
		d.Set("dbms_type", db.DBMSTypeInstance)
	} else {
		d.Set("dbms_id", backup.ClusterID)
		d.Set("dbms_type", db.DBMSTypeCluster)
	}
	d.Set("description", backup.Description)
	d.Set("location_ref", backup.LocationRef)
	d.Set("created", backup.Created)
	d.Set("updated", backup.Updated)
	d.Set("size", backup.Size)
	d.Set("wal_size", backup.WalSize)
	d.Set("datastore", flattenDatabaseInstanceDatastore(*backup.Datastore))
	d.Set("meta", backup.Meta)

	return nil
}

func resourceDatabaseBackupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = backups.Delete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_db_backup"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(dbBackupStatusActive)},
		Target:     []string{string(dbBackupStatusDeleted)},
		Refresh:    databaseBackupStateRefreshFunc(DatabaseV1Client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      dbInstanceDelay,
		MinTimeout: dbInstanceMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_db_backup %s to delete: %s", d.Id(), err)
	}

	return nil
}
