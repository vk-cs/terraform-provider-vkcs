package vkcs

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/utils/terraform/hashcode"
)

func dataSourceComputeAvailabilityZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeAvailabilityZonesRead,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"state": {
				Type:         schema.TypeString,
				Default:      "available",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"available", "unavailable"}, true),
			},
		},
	}
}

func dataSourceComputeAvailabilityZonesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	region := getRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	allPages, err := availabilityzones.List(computeClient).AllPages()
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_compute_availability_zones: %s", err)
	}
	zoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return diag.Errorf("Error extracting vkcs_compute_availability_zones from response: %s", err)
	}

	stateBool := d.Get("state").(string) == "available"
	zones := make([]string, 0, len(zoneInfo))
	for _, z := range zoneInfo {
		if z.ZoneState.Available == stateBool {
			zones = append(zones, z.ZoneName)
		}
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(zones)

	d.SetId(hashcode.Strings(zones))
	d.Set("names", zones)
	d.Set("region", region)

	return nil
}
