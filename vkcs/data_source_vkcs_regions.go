package vkcs

import (
	"context"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceVkcsRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVkcsRegionsRead,
		Schema: map[string]*schema.Schema{
			"parent_region_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"names": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceVkcsRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	client, err := config.IdentityV3Client(config.GetRegion())
	if err != nil {
		return diag.Errorf("failed to init identity v3 client: %s", err)
	}

	opts := regions.ListOpts{}
	if parentRegion, ok := d.GetOk("parent_region_id"); ok {
		opts.ParentRegionID = parentRegion.(string)
	}

	allPages, err := regions.List(client.(*gophercloud.ServiceClient), opts).AllPages()
	if err != nil {
		return diag.Errorf("failed to list regions: %s", err)
	}

	allRegions, err := regions.ExtractRegions(allPages)
	if err != nil {
		return diag.Errorf("failed to extract regions: %s", err)
	}

	names := make([]string, 0, len(allRegions))
	for _, r := range allRegions {
		names = append(names, r.ID)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	if err := d.Set("names", names); err != nil {
		return diag.Errorf("failed to set names: %s", err)
	}
	return nil
}
