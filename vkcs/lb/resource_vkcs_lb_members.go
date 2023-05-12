package lb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func ResourceMembers() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMembersCreate,
		ReadContext:   resourceMembersRead,
		UpdateContext: resourceMembersUpdate,
		DeleteContext: resourceMembersDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new members resource.",
			},

			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the pool that members will be assigned to. Changing this creates a new members resource.",
			},

			"member": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID for the member.",
						},

						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Human-readable name for the member.",
						},

						"address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address of the members to receive traffic from the load balancer.",
						},

						"protocol_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
							Description:  "The port on which to listen for client traffic.",
						},

						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 256),
							Description:  "A positive integer value that indicates the relative portion of traffic that this members should receive from the pool. For example, a member with a weight of 10 receives five times as much traffic as a member with a weight of 2. Defaults to 1.",
						},

						"subnet_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The subnet in which to access the member.",
						},

						"backup": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "A bool that indicates whether the member is backup.",
						},

						"admin_state_up": {
							Type:        schema.TypeBool,
							Default:     true,
							Optional:    true,
							Description: "The administrative state of the member. A valid value is true (UP) or false (DOWN). Defaults to true.",
						},
					},
				},
				Description: "A set of dictionaries containing member parameters. The structure is described below.",
			},
		},
	}
}

func resourceMembersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	createOpts := expandLBMembers(d.Get("member").(*schema.Set), lbClient)
	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := octaviapools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent pool %s: %s", poolID, err)
	}

	// Wait for parent pool to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create members")
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = octaviapools.BatchUpdateMembers(lbClient, poolID, createOpts).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating members: %s", err)
	}

	// Wait for parent pool to become active before continuing
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(poolID)

	return resourceMembersRead(ctx, d, meta)
}

func resourceMembersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	allPages, err := octaviapools.ListMembers(lbClient, d.Id(), octaviapools.ListMembersOpts{}).AllPages()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_lb_members"))
	}

	members, err := octaviapools.ExtractMembers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_lb_members: %s", err)
	}

	log.Printf("[DEBUG] Retrieved members for the %s pool: %#v", d.Id(), members)

	d.Set("pool_id", d.Id())
	d.Set("member", FlattenLBMembers(members))
	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourceMembersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	if d.HasChange("member") {
		updateOpts := expandLBMembers(d.Get("member").(*schema.Set), lbClient)

		// Get a clean copy of the parent pool.
		parentPool, err := octaviapools.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve parent pool %s: %s", d.Id(), err)
		}

		// Wait for parent pool to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating %s pool members with options: %#v", d.Id(), updateOpts)
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			err = octaviapools.BatchUpdateMembers(lbClient, d.Id(), updateOpts).ExtractErr()
			if err != nil {
				return util.CheckForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return diag.Errorf("Unable to update member %s: %s", d.Id(), err)
		}

		// Wait for parent pool to become active before continuing
		err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceMembersRead(ctx, d, meta)
}

func resourceMembersDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Get a clean copy of the parent pool.
	parentPool, err := octaviapools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, fmt.Sprintf("Unable to retrieve parent pool (%s) for the member", d.Id())))
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error waiting for the members' pool status"))
	}

	log.Printf("[DEBUG] Attempting to delete %s pool members", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = octaviapools.BatchUpdateMembers(lbClient, d.Id(), []octaviapools.BatchUpdateMemberOpts{}).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting members"))
	}

	// Wait for parent pool to become active before continuing.
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error waiting for the members' pool status"))
	}

	return nil
}
