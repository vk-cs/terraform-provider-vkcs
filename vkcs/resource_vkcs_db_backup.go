package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func resourceDatabaseBackup() *schema.Resource {
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dbms_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"container_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			// Computed fields
			"dbms_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"location_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"wal_size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},

			"datastore": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"meta": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDatabaseBackupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	dbmsID := d.Get("dbms_id").(string)

	dbmsResp, err := getDBMSResource(DatabaseV1Client, dbmsID)
	if err != nil {
		return diag.Errorf("error while getting resource: %s", err)
	}

	var dbmsType string
	if instanceResource, ok := dbmsResp.(*instanceResp); ok {
		if isOperationNotSupported(instanceResource.DataStore.Type, Redis, Tarantool) {
			return diag.Errorf("operation not supported for this datastore")
		}
		if instanceResource.ReplicaOf != nil {
			return diag.Errorf("operation not supported for replica")
		}
		dbmsType = dbmsTypeInstance
	}
	if clusterResource, ok := dbmsResp.(*dbClusterResp); ok {
		if isOperationNotSupported(clusterResource.DataStore.Type, Redis, Tarantool) {
			return diag.Errorf("operation not supported for this datastore")
		}
		dbmsType = dbmsTypeCluster
	}

	b := dbBackupCreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		ContainerPrefix: d.Get("container_prefix").(string),
	}

	if dbmsType == dbmsTypeInstance {
		b.Instance = d.Get("dbms_id").(string)
	} else {
		b.Cluster = d.Get("dbms_id").(string)
	}

	log.Printf("[DEBUG] vkcs_db_backup create options: %#v", b)

	back := dbBackup{
		Backup: &b,
	}
	backup, err := dbBackupCreate(DatabaseV1Client, &back).extract()
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
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	backup, err := dbBackupGet(DatabaseV1Client, d.Id()).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_backup"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_backup %s: %#v", d.Id(), backup)

	d.Set("name", backup.Name)
	if backup.InstanceID != "" {
		d.Set("dbms_id", backup.InstanceID)
		d.Set("dbms_type", dbmsTypeInstance)
	} else {
		d.Set("dbms_id", backup.ClusterID)
		d.Set("dbms_type", dbmsTypeCluster)
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
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	err = dbBackupDelete(DatabaseV1Client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_db_backup"))
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
