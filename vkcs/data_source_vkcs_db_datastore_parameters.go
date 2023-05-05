package vkcs

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
)

func dataSourceDatabaseDatastoreParameters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatastoreParametersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`.",
			},

			"datastore_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the data store.",
			},

			"datastore_version_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of the version of the data store.",
			},

			"parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a configuration parameter.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Type of a configuration parameter.",
						},
						"min": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Minimum value of a configuration parameter.",
						},
						"max": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Maximum value of a configuration parameter.",
						},
						"restart_required": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This attribute indicates whether a restart required when a parameter is set.",
						},
					},
				},
				Description: "Versions of the datastore.",
			},
		},
		Description: "Use this data source to get configuration parameters supported for a VKCS datastore. **New since v.0.2.0**.",
	}
}

func dataSourceDatabaseDatastoreParametersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	region := getRegion(d, config)
	dbClient, err := config.DatabaseV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	dsName := d.Get("datastore_name").(string)
	dsVersionID := d.Get("datastore_version_id").(string)

	params, err := datastores.ListParameters(dbClient, dsName, dsVersionID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_db_backup"))
	}

	flattenedParams := flattenDatabaseDatastoreParameters(params)

	d.SetId(fmt.Sprintf("%s/%s/params", dsName, dsVersionID))
	d.Set("parameters", flattenedParams)

	return nil
}

func flattenDatabaseDatastoreParameters(params []datastores.Param) (r []map[string]interface{}) {
	for _, p := range params {
		r = append(r, map[string]interface{}{
			"name":             p.Name,
			"type":             p.Type,
			"min":              p.MinValue,
			"max":              p.MaxValue,
			"restart_required": p.RestartRequried,
		})
	}
	return
}
