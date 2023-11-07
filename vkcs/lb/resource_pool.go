package lb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/lb/v2/listeners"
	ipools "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/lb/v2/pools"
)

func ResourcePool() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePoolCreate,
		ReadContext:   resourcePoolRead,
		UpdateContext: resourcePoolUpdate,
		DeleteContext: resourcePoolDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePoolImport,
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
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new pool.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable name for the pool.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description for the pool.",
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS", "PROXY",
				}, false),
				Description: "The protocol - can either be TCP, HTTP, HTTPS, PROXY, or UDP. Changing this creates a new pool.",
			},

			// One of loadbalancer_id or listener_id must be provided
			"loadbalancer_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The load balancer on which to provision this pool. Changing this creates a new pool. _note_ One of LoadbalancerID or ListenerID must be provided.",
			},

			// One of loadbalancer_id or listener_id must be provided
			"listener_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The Listener on which the members of the pool will be associated with. Changing this creates a new pool. _note_ One of LoadbalancerID or ListenerID must be provided.",
			},

			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP",
				}, false),
				Description: "The load balancing algorithm to distribute traffic to the pool's members. Must be one of ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP, or SOURCE_IP_PORT.",
			},

			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE",
							}, false),
							Description: "The type of persistence mode. The current specification supports SOURCE_IP, HTTP_COOKIE, and APP_COOKIE.",
						},

						"cookie_name": {
							Type:        schema.TypeString,
							Optional:    true,
							ForceNew:    true,
							Description: "The name of the cookie if persistence mode is set appropriately. Required if `type = APP_COOKIE`.",
						},
					},
				},
				Description: "Omit this field to prevent session persistence. Indicates whether connections in the same session will be processed by the same Pool member or not. Changing this creates a new pool.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "The administrative state of the pool. A valid value is true (UP) or false (DOWN).",
			},
		},
		Description: "Manages a pool resource within VKCS.",
	}
}

func resourcePoolCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	lbID := d.Get("loadbalancer_id").(string)
	listenerID := d.Get("listener_id").(string)
	var persistence pools.SessionPersistence
	if p, ok := d.GetOk("persistence"); ok {
		pV := (p.([]interface{}))[0].(map[string]interface{})

		persistence = pools.SessionPersistence{
			Type: pV["type"].(string),
		}

		if persistence.Type == "APP_COOKIE" {
			if pV["cookie_name"].(string) == "" {
				return diag.Errorf(
					"Persistence cookie_name needs to be set if using 'APP_COOKIE' persistence type")
			}
			persistence.CookieName = pV["cookie_name"].(string)
		} else if pV["cookie_name"].(string) != "" {
			return diag.Errorf(
				"Persistence cookie_name can only be set if using 'APP_COOKIE' persistence type")
		}
	}

	createOpts := pools.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Protocol:       pools.Protocol(d.Get("protocol").(string)),
		LoadbalancerID: lbID,
		ListenerID:     listenerID,
		LBMethod:       pools.LBMethod(d.Get("lb_method").(string)),
		AdminStateUp:   &adminStateUp,
	}

	// Must omit if not set
	if persistence != (pools.SessionPersistence{}) {
		createOpts.Persistence = &persistence
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for Listener or LoadBalancer to become active before continuing
	if listenerID != "" {
		listener, err := listeners.Get(lbClient, listenerID).Extract()
		if err != nil {
			return diag.Errorf("Unable to get vkcs_lb_listener %s: %s", listenerID, err)
		}

		waitErr := waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
		if waitErr != nil {
			return diag.Errorf(
				"Error waiting for vkcs_lb_listener %s to become active: %s", listenerID, err)
		}
	} else {
		waitErr := waitForLBLoadBalancer(ctx, lbClient, lbID, "ACTIVE", getLbPendingStatuses(), timeout)
		if waitErr != nil {
			return diag.Errorf(
				"Error waiting for vkcs_lb_loadbalancer %s to become active: %s", lbID, err)
		}
	}

	log.Printf("[DEBUG] Attempting to create pool")
	var pool *pools.Pool
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		pool, err = ipools.Create(lbClient, createOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating pool: %s", err)
	}

	d.SetId(pool.ID)

	// Pool was successfully created
	// Wait for pool to become active before continuing
	err = waitForLBPool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolRead(ctx, d, meta)
}

func resourcePoolRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	pool, err := ipools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "pool"))
	}

	log.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	d.Set("lb_method", pool.LBMethod)
	d.Set("protocol", pool.Protocol)
	d.Set("description", pool.Description)
	d.Set("admin_state_up", pool.AdminStateUp)
	d.Set("name", pool.Name)
	d.Set("persistence", flattenLBPoolPersistence(pool.Persistence))
	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourcePoolUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = pools.LBMethod(d.Get("lb_method").(string))
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	// Get a clean copy of the pool.
	pool, err := ipools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve pool %s: %s", d.Id(), err)
	}

	// Wait for pool to become active before continuing
	err = waitForLBPool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = ipools.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to update pool %s: %s", d.Id(), err)
	}

	// Wait for pool to become active before continuing
	err = waitForLBPool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourcePoolRead(ctx, d, meta)
}

func resourcePoolDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	// Get a clean copy of the pool.
	pool, err := ipools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Unable to retrieve pool"))
	}

	log.Printf("[DEBUG] Attempting to delete pool %s", d.Id())
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = ipools.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting pool"))
	}

	// Wait for Pool to delete
	err = waitForLBPool(ctx, lbClient, pool, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourcePoolImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS networking client: %s", err)
	}

	pool, err := ipools.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return nil, util.CheckDeleted(d, err, "pool")
	}

	log.Printf("[DEBUG] Retrieved pool %s during the import: %#v", d.Id(), pool)

	switch {
	case len(pool.Listeners) > 0 && pool.Listeners[0].ID != "":
		d.Set("listener_id", pool.Listeners[0].ID)
	case len(pool.Loadbalancers) > 0 && pool.Loadbalancers[0].ID != "":
		d.Set("loadbalancer_id", pool.Loadbalancers[0].ID)
	default:
		return nil, fmt.Errorf("unable to detect pool's Listener ID or Load Balancer ID")
	}

	return []*schema.ResourceData{d}, nil
}
