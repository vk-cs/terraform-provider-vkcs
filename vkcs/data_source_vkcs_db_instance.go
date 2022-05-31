package vkcs

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ip": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"datastore": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
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

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"volume": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Set:      volInstHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},

						"used": {
							Type:     schema.TypeFloat,
							Required: true,
						},

						"volume_id": {
							Type:     schema.TypeString,
							Required: true,
						},

						"volume_type": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"backup_schedule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"start_hours": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"start_minutes": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"interval_hours": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"keep_count": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDatabaseInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS database client: %s", err)
	}

	instance, err := instanceGet(DatabaseV1Client, d.Get("id").(string)).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_instance"))
	}

	d.SetId(instance.ID)

	d.Set("name", instance.Name)
	d.Set("flavor_id", instance.Flavor)
	d.Set("datastore", instance.DataStore)
	d.Set("region", getRegion(d, config))
	d.Set("ip", instance.IP)
	d.Set("status", instance.Status)

	m := map[string]interface{}{
		"size":        *instance.Volume.Size,
		"used":        *instance.Volume.Used,
		"volume_id":   instance.Volume.VolumeID,
		"volume_type": instance.Volume.VolumeType,
	}

	d.Set("volume", schema.NewSet(volInstHash, []interface{}{m}))

	backupSchedule, err := instanceGetBackupSchedule(DatabaseV1Client, d.Id()).extract()
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
