package regions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/regions/regions"
)

func DataSourceVkcsRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVkcsRegionRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of the region to learn or use. Use empty value to learn current region on the provider.",
			},
			"parent_region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Parent of the region.",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description of the region.",
			},
		},
		Description: "`vkcs_region` provides details about a specific VKCS region. As well as validating a given region name this resource can be used to discover the name of the region configured within the provider.",
	}
}

func dataSourceVkcsRegionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
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

	region, err := regions.Get(client, regionName).Extract()
	if err != nil {
		return diag.Errorf("failed to get region for %s: %s", regionName, err)
	}

	d.SetId(region.ID)
	d.Set("parent_region", region.ParentRegionID)
	d.Set("description", region.Description)
	return nil
}
