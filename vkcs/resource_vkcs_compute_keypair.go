package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/keypairs"
)

func resourceComputeKeypair() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeKeypairCreate,
		ReadContext:   resourceComputeKeypairRead,
		DeleteContext: resourceComputeKeypairDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			// computed-only
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceComputeKeypairCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	name := d.Get("name").(string)
	createOpts := ComputeKeyPairV2CreateOpts{
		keypairs.CreateOpts{
			Name:      name,
			PublicKey: d.Get("public_key").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] vkcs_compute_keypair create options: %#v", createOpts)

	kp, err := keypairs.Create(computeClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Unable to create vkcs_compute_keypair %s: %s", name, err)
	}

	d.SetId(kp.Name)

	// Private Key is only available in the response to a create.
	d.Set("private_key", kp.PrivateKey)

	return resourceComputeKeypairRead(ctx, d, meta)
}

func resourceComputeKeypairRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	kp, err := keypairs.Get(computeClient, d.Id(), keypairs.GetOpts{}).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_compute_keypair"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_keypair %s: %#v", d.Id(), kp)

	d.Set("name", kp.Name)
	d.Set("public_key", kp.PublicKey)
	d.Set("fingerprint", kp.Fingerprint)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceComputeKeypairDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	err = keypairs.Delete(computeClient, d.Id(), keypairs.DeleteOpts{}).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_compute_keypair"))
	}

	return nil
}
