package regions

import (
	"context"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

func DataSourceVkcsRegions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVkcsRegionsRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names of regions that meets the criteria.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Random identifier of the data source.",
			},
		},
		Description: "`vkcs_regions` provides information about VKCS regions. To get details of each region the data source can be combined with the `vkcs_region` data source.",
	}
}

func dataSourceVkcsRegionsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.IdentityV3Client(config.GetRegion())
	if err != nil {
		return diag.Errorf("failed to init identity v3 client: %s", err)
	}

	opts := regions.ListOpts{}

	allPages, err := regions.List(client, opts).AllPages()
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
