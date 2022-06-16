package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseConfigGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseConfigGroupRead,

		Schema: map[string]*schema.Schema{
			"config_group_id": {
				Type:     schema.TypeString,
				Required: true,
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
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"values": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDatabaseConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}
	configGroupID := d.Get("config_group_id").(string)
	configGroup, err := dbConfigGroupGet(DatabaseV1Client, configGroupID).extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_config_group"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_config_group %s: %#v", configGroupID, configGroup)

	d.Set("name", configGroup.Name)
	ds := dataStore{
		Type:    configGroup.DatastoreName,
		Version: configGroup.DatastoreVersionName,
	}
	d.Set("datastore", flattenDatabaseInstanceDatastore(ds))
	d.Set("values", flattenDatabaseConfigGroupValues(configGroup.Values))

	d.Set("updated", configGroup.Updated)
	d.Set("created", configGroup.Created)
	d.Set("description", configGroup.Description)
	d.SetId(configGroupID)

	return nil
}
