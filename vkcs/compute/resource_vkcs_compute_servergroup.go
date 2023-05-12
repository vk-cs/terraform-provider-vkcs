package compute

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
)

func ResourceComputeServerGroup() *schema.Resource {
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used. Changing this creates a new server group.",
			},

			"name": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "A unique name for the server group. Changing this creates a new server group.",
			},

			"policies": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The set of policies for the server group. All policies are mutually exclusive. See the Policies section for more information. Changing this creates a new server group.",
			},

			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The instances that are part of this server group.",
			},

			"value_specs": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of additional options.",
			},
		},
		Description: "Manages a Server Group resource within VKCS.",
	}
}

func resourceComputeServerGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	name := d.Get("name").(string)

	rawPolicies := d.Get("policies").([]interface{})
	policies := ExpandComputeServerGroupPolicies(computeClient, rawPolicies)

	createOpts := ComputeServerGroupCreateOpts{
		servergroups.CreateOpts{
			Name:     name,
			Policies: policies,
		},
		util.MapValueSpecs(d),
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
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	sg, err := servergroups.Get(computeClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_compute_servergroup"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_servergroup %s: %#v", d.Id(), sg)

	d.Set("name", sg.Name)
	d.Set("policies", sg.Policies)
	d.Set("members", sg.Members)

	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourceComputeServerGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	if err := servergroups.Delete(computeClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_compute_servergroup"))
	}

	return nil
}
