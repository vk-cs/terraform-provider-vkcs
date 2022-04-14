package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
)

func dataSourceComputeKeypair() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeKeypairRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// computed-only
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeKeypairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	name := d.Get("name").(string)
	kp, err := keypairs.Get(computeClient, name, keypairs.GetOpts{}).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_compute_keypair %s: %s", name, err)
	}

	d.SetId(name)

	log.Printf("[DEBUG] Retrieved vkcs_compute_keypair %s: %#v", d.Id(), kp)

	d.Set("fingerprint", kp.Fingerprint)
	d.Set("public_key", kp.PublicKey)
	d.Set("region", getRegion(d, config))

	return nil
}
