package vkcs

import (
	"context"
	"sort"

	"github.com/gophercloud/utils/terraform/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDatabaseDatastores() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatastoresRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`",
			},

			"datastores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The id of the datastore.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the datastore.",
						},
					},
				},
			},
		},
		Description: "Use this data source to get a list of datastores from VKCS. **New since v.0.2.0**.",
	}
}

func dataSourceDatabaseDatastoresRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	region := getRegion(d, config)
	dbClient, err := config.DatabaseV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	allPages, err := dataStoreList(dbClient).AllPages()
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_db_datastores: %s", err)
	}

	allDatastores, err := extractDatastores(allPages)
	if err != nil {
		return diag.Errorf("Error extracting vkcs_db_datastores from response: %s", err)
	}

	flattenedDatastores := flattenDatabaseDatastoresDatastores(allDatastores)
	sort.SliceStable(flattenedDatastores, func(i, j int) bool {
		return flattenedDatastores[i]["name"].(string) < flattenedDatastores[j]["name"].(string)
	})

	var names []string
	for _, d := range flattenedDatastores {
		names = append(names, d["name"].(string))
	}

	d.SetId(hashcode.Strings(names))
	d.Set("region", region)
	d.Set("datastores", flattenedDatastores)

	return nil
}

func flattenDatabaseDatastoresDatastores(datastores []dataStore) (r []map[string]interface{}) {
	for _, d := range datastores {
		r = append(r, map[string]interface{}{
			"id":   d.ID,
			"name": d.Name,
		})
	}
	return
}
