package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	octaviapools "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func resourceMembers() *schema.Resource {
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"pool_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"member": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"address": {
							Type:     schema.TypeString,
							Required: true,
						},

						"protocol_port": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntBetween(1, 65535),
						},

						"weight": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      1,
							ValidateFunc: validation.IntBetween(0, 256),
						},

						"subnet_id": {
							Type:     schema.TypeString,
							Optional: true,
						},

						"backup": {
							Type:     schema.TypeBool,
							Optional: true,
						},

						"admin_state_up": {
							Type:     schema.TypeBool,
							Default:  true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceMembersCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
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
			return checkForRetryableError(err)
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
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	allPages, err := octaviapools.ListMembers(lbClient, d.Id(), octaviapools.ListMembersOpts{}).AllPages()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_lb_members"))
	}

	members, err := octaviapools.ExtractMembers(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_lb_members: %s", err)
	}

	log.Printf("[DEBUG] Retrieved members for the %s pool: %#v", d.Id(), members)

	d.Set("pool_id", d.Id())
	d.Set("member", flattenLBMembers(members))
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceMembersUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
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
				return checkForRetryableError(err)
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
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Get a clean copy of the parent pool.
	parentPool, err := octaviapools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, fmt.Sprintf("Unable to retrieve parent pool (%s) for the member", d.Id())))
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error waiting for the members' pool status"))
	}

	log.Printf("[DEBUG] Attempting to delete %s pool members", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = octaviapools.BatchUpdateMembers(lbClient, d.Id(), []octaviapools.BatchUpdateMemberOpts{}).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting members"))
	}

	// Wait for parent pool to become active before continuing.
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error waiting for the members' pool status"))
	}

	return nil
}
