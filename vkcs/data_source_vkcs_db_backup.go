package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseBackupRead,

		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dbms_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

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

func dataSourceDatabaseBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS database client: %s", err)
	}

	backupID := d.Get("backup_id").(string)
	backup, err := dbBackupGet(DatabaseV1Client, backupID).extract()
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
	d.SetId(backupID)

	return nil
}
