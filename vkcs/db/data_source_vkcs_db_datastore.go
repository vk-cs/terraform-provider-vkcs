package db

import (
	"context"
	"log"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/db/v1/datastores"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceDatabaseDatastore() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDatabaseDatastoreRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The `region` to fetch availability zones from, defaults to the provider's `region`",
			},

			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The id of the datastore.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the datastore.",
			},

			"minimum_cpu": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum CPU required for instance of the datastore.",
			},

			"minimum_ram": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Minimum RAM required for instance of the datastore.",
			},

			"volume_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Supported volume types for the datastore.",
			},

			"cluster_volume_types": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Supported volume types for the datastore when used in a cluster.",
			},

			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "ID of a version of the datastore.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of a version of the datastore.",
						},
					},
				},
				Description: "Versions of the datastore.",
			},
		},
		Description: "Use this data source to get information on a VKCS db datastore.",
	}
}

func dataSourceDatabaseDatastoreRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	region := util.GetRegion(d, config)
	dbClient, err := config.DatabaseV1Client(region)
	if err != nil {
		return diag.Errorf("Error creating VKCS database client: %s", err)
	}

	allPages, err := datastores.List(dbClient).AllPages()
	if err != nil {
		return diag.Errorf("Error retrieving datastores: %s", err)
	}

	datastoresInfo, err := datastores.ExtractDatastores(allPages)
	if err != nil {
		return diag.Errorf("Error extracting datastores from response: %s", err)
	}

	id, name := d.Get("id").(string), d.Get("name").(string)
	allDatastores := filterDatabaseDatastores(datastoresInfo, id, name)

	if len(allDatastores) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allDatastores) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allDatastores)
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	dsID := allDatastores[0].ID
	ds, err := datastores.Get(dbClient, dsID).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_db_datastore: %s", err)
	}

	flattenedVersions := flattenDatabaseDatastoreVersions(ds.Versions)
	sort.SliceStable(flattenedVersions, func(i, j int) bool {
		return flattenedVersions[i]["name"].(string) > flattenedVersions[j]["name"].(string)
	})

	d.SetId(ds.ID)
	d.Set("name", ds.Name)
	d.Set("minimum_cpu", ds.MinimumCPU)
	d.Set("minimum_ram", ds.MinimumRAM)
	d.Set("volume_types", ds.VolumeTypes)
	d.Set("cluster_volume_types", ds.ClusterVolumeTypes)
	d.Set("versions", flattenedVersions)

	return nil
}

func filterDatabaseDatastores(dsSlice []datastores.Datastore, id, name string) []datastores.Datastore {
	var res []datastores.Datastore
	for _, ds := range dsSlice {
		if (name == "" || ds.Name == name) && (id == "" || ds.ID == id) {
			res = append(res, ds)
		}
	}
	return res
}

func flattenDatabaseDatastoreVersions(versions []datastores.Version) (r []map[string]interface{}) {
	for _, v := range versions {
		r = append(r, map[string]interface{}{
			"id":   v.ID,
			"name": v.Name,
		})
	}
	return
}
