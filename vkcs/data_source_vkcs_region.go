package vkcs

import (
	"context"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVkcsRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVkcsRegionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent_region": {
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

func dataSourceVkcsRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	client, err := config.IdentityV3Client(config.GetRegion())
	if err != nil {
		return diag.Errorf("failed to init identity v3 client: %s", err)
	}

	// default region
	regionName := config.GetRegion()
	// or passed from config
	if v, ok := d.GetOk("id"); ok {
		regionName = v.(string)
	}

	region, err := regions.Get(client.(*gophercloud.ServiceClient), regionName).Extract()
	if err != nil {
		return diag.Errorf("failed to get region for %s: %s", regionName, err)
	}

	d.SetId(region.ID)
	d.Set("parent_region", region.ParentRegionID)
	d.Set("description", region.Description)
	return nil
}
