package vkcs

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseDatastores() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatastoresRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"datastores": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"minimum_cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minimum_ram": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"versions": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"volume_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"cluster_volume_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDatabaseDatastoresRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	DatabaseV1Client, err := config.DatabaseV1Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating vkcs database client: %s", err)
	}

	filter := d.Get("filter").(string)

	var datastoreMap []map[string]interface{}
	if filter != "" {
		datastore, err := datastoreGet(DatabaseV1Client, filter).extract()
		if err != nil {
			return diag.Errorf("Error retrieving vkcs_db_datastores: %s", err)
		}
		datastoreMap = make([]map[string]interface{}, 1)
		datastoreMap[0] = flattenDatabaseDatastore(datastore)
		d.SetId(datastore.ID)
	} else {
		datastores, err := datastoresGet(DatabaseV1Client).extract()
		if err != nil {
			return diag.Errorf("Error retrieving vkcs_db_datastores: %s", err)
		}
		datastoreMap = make([]map[string]interface{}, len(*datastores))
		for i, ds := range *datastores {
			datastoreMap[i] = flattenDatabaseDatastore(&ds)
		}
		d.SetId((*datastores)[0].ID)
	}

	d.Set("datastores", datastoreMap)

	return nil
}
