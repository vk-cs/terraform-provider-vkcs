package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/backups"
)

func dataSourceDatabaseBackup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseBackupRead,

		Schema: map[string]*schema.Schema{
			"backup_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the backup.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the backup.",
			},

			"dbms_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the backed up instance or cluster",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the backup",
			},

			"dbms_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of dbms of the backup, can be \"instance\" or \"cluster\".",
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
							Required:    true,
							Description: "Version of the datastore.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the datastore.",
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
		Description: "Use this data source to get the information on a db backup resource.",
	}
}

func dataSourceDatabaseBackupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS database client: %s", err)
	}

	backupID := d.Get("backup_id").(string)
	backup, err := backups.Get(DatabaseV1Client, backupID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_backup"))
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
	d.SetId(backupID)

	return nil
}
