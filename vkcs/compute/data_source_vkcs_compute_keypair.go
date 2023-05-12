package compute

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
)

func DataSourceComputeKeypair() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeKeypairRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The unique name of the keypair.",
			},

			// computed-only
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fingerprint of the OpenSSH key.",
			},

			"public_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The OpenSSH-formatted public key of the keypair.",
			},
		},
		Description: "Use this data source to get the ID and public key of an VKCS keypair.",
	}
}

func dataSourceComputeKeypairRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
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
	d.Set("region", util.GetRegion(d, config))

	return nil
}
