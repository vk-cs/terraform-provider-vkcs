package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	configgroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/config_groups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

func dataSourceDatabaseConfigGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseConfigGroupRead,

		Schema: map[string]*schema.Schema{
			"config_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the config_group.",
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
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the config group.",
			},
			"values": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Map of configuration parameters in format \"key\": \"value\".",
			},
			"updated": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of config group's last update.",
			},
			"created": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of config group's creation.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The description of the config group.",
			},
		},
		Description: "Use this data source to get the information on a db config group resource.\n" +
			"**New since v.0.1.7**.",
	}
}

func dataSourceDatabaseConfigGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}
	configGroupID := d.Get("config_group_id").(string)
	configGroup, err := configgroups.Get(DatabaseV1Client, configGroupID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_config_group"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_db_config_group %s: %#v", configGroupID, configGroup)

	d.Set("name", configGroup.Name)
	ds := datastores.DatastoreShort{
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
