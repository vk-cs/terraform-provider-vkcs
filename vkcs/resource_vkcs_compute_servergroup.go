package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
)

func resourceComputeServerGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeServerGroupCreate,
		ReadContext:   resourceComputeServerGroupRead,
		Update:        nil,
		DeleteContext: resourceComputeServerGroupDelete,
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
				ForceNew: true,
				Required: true,
			},

			"policies": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeServerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	name := d.Get("name").(string)

	rawPolicies := d.Get("policies").([]interface{})
	policies := expandComputeServerGroupPolicies(computeClient, rawPolicies)

	createOpts := ComputeServerGroupCreateOpts{
		servergroups.CreateOpts{
			Name:     name,
			Policies: policies,
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] vkcs_compute_servergroup create options: %#v", createOpts)
	newSG, err := servergroups.Create(computeClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_compute_servergroup %s: %s", name, err)
	}

	d.SetId(newSG.ID)

	return resourceComputeServerGroupRead(ctx, d, meta)
}

func resourceComputeServerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	sg, err := servergroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_compute_servergroup"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_servergroup %s: %#v", d.Id(), sg)

	d.Set("name", sg.Name)
	d.Set("policies", sg.Policies)
	d.Set("members", sg.Members)

	d.Set("region", getRegion(d, config))

	return nil
}

func resourceComputeServerGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	if err := servergroups.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_compute_servergroup"))
	}

	return nil
}
