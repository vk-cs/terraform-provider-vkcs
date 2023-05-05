package vkcs

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/instances"
)

// volInstHash calculates hash of the volume of instance
func volInstHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%d-", m["size"].(int)))
	switch used := m["used"].(type) {
	case float32:
		buf.WriteString(fmt.Sprintf("%.2f-", used))
	default:
		buf.WriteString(fmt.Sprintf("%.2f-", used))
	}
	buf.WriteString(fmt.Sprintf("%s-", m["volume_id"].(string)))
	// TODO(irlndts): the function is deprecated, replace it.
	// nolint:staticcheck
	return hashcode.String(buf.String())
}

func dataSourceDatabaseInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseInstanceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The id of the instance.",
			},

			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Region of the resource.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the instance.",
			},

			"flavor_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of flavor for the instance.",
			},

			"hostname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The hostname of the instance.",
			},

			"ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "IP address of the instance.",
			},

			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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
				Description: "Object that represents datastore of the instance.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Instance status.",
			},

			"volume": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set:      volInstHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Size of the instance volume.",
						},

						"used": {
							Type:        schema.TypeFloat,
							Required:    true,
							Description: "Size of the used volume space.",
						},

						"volume_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "ID of the instance volume.",
						},

						"volume_type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Type of the instance volume.",
						},
					},
				},
				Description: "Object that describes volume of the instance.",
			},
			"backup_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the schedule.",
						},
						"start_hours": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Hours part of timestamp of initial backup.",
						},
						"start_minutes": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Minutes part of timestamp of initial backup.",
						},
						"interval_hours": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Time interval between backups, specified in hours. Available values: 3, 6, 8, 12, 24.",
						},
						"keep_count": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Number of backups to be stored.",
						},
					},
				},
				Description: "Object that represents configuration of PITR backup. This functionality is available only for postgres datastore. **New since v.0.1.4**.",
			},
		},
		Description: "Use this data source to get the information on a db instance resource.",
	}
}

func dataSourceDatabaseInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS database client: %s", err)
	}

	instance, err := instances.Get(DatabaseV1Client, d.Get("id").(string)).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_instance"))
	}

	d.SetId(instance.ID)

	d.Set("name", instance.Name)
	d.Set("flavor_id", instance.Flavor.ID)
	d.Set("datastore", flattenDatabaseInstanceDatastore(*instance.DataStore))
	d.Set("region", getRegion(d, config))
	d.Set("ip", instance.IP)
	d.Set("status", instance.Status)

	m := map[string]interface{}{
		"size":        *instance.Volume.Size,
		"used":        0,
		"volume_id":   instance.Volume.VolumeID,
		"volume_type": instance.Volume.VolumeType,
	}
	if instance.Volume.Used != nil {
		m["used"] = *instance.Volume.Used
	}

	d.Set("volume", schema.NewSet(volInstHash, []interface{}{m}))

	backupSchedule, err := instances.GetBackupSchedule(DatabaseV1Client, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("error getting backup schedule for instance: %s: %s", d.Id(), err)
	}
	if backupSchedule != nil {
		flattened := flattenDatabaseBackupSchedule(*backupSchedule)
		d.Set("backup_schedule", flattened)
	} else {
		d.Set("backup_schedule", nil)
	}

	return nil
}
