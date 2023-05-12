package lb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	octaviamonitors "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/monitors"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func ResourceMonitor() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceMonitorCreate,
		ReadContext:   resourceMonitorRead,
		UpdateContext: resourceMonitorUpdate,
		DeleteContext: resourceMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceMonitorImport,
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
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new monitor.",
			},

			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the pool that this monitor will be assigned to.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The Name of the Monitor.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP-CONNECT", "HTTP", "HTTPS", "TLS-HELLO", "PING",
				}, false),
				Description: "The type of probe, which is PING, TCP, HTTP, HTTPS, TLS-HELLO or UDP-CONNECT, that is sent by the load balancer to verify the member state. Changing this creates a new monitor.",
			},

			"delay": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The time, in seconds, between sending probes to members.",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Maximum number of seconds for a monitor to wait for a ping reply before it times out. The value must be less than the delay value.",
			},

			"max_retries": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Number of permissible ping failures before changing the member's status to INACTIVE. Must be a number between 1 and 10.",
			},

			"max_retries_down": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "Number of permissible ping failures befor changing the member's status to ERROR. Must be a number between 1 and 10. Changing this updates the max_retries_down of the existing monitor.",
			},

			"url_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required for HTTP(S) types. URI path that will be accessed if monitor type is HTTP or HTTPS.",
			},

			"http_method": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required for HTTP(S) types. The HTTP method used for requests by the monitor. If this attribute is not specified, it defaults to \"GET\".",
			},

			"expected_codes": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Required for HTTP(S) types. Expected HTTP codes for a passing HTTP(S) monitor. You can either specify a single status like \"200\", or a range like \"200-202\".",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "The administrative state of the monitor. A valid value is true (UP) or false (DOWN).",
			},
		},
		Description: "Manages a monitor resource within VKCS.",
	}
}

func resourceMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)

	createOpts := octaviamonitors.CreateOpts{
		PoolID:         d.Get("pool_id").(string),
		Type:           d.Get("type").(string),
		Delay:          d.Get("delay").(int),
		Timeout:        d.Get("timeout").(int),
		MaxRetries:     d.Get("max_retries").(int),
		MaxRetriesDown: d.Get("max_retries_down").(int),
		URLPath:        d.Get("url_path").(string),
		HTTPMethod:     d.Get("http_method").(string),
		ExpectedCodes:  d.Get("expected_codes").(string),
		Name:           d.Get("name").(string),
		AdminStateUp:   &adminStateUp,
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent vkcs_lb_pool %s: %s", poolID, err)
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] vkcs_lb_monitor create options: %#v", createOpts)
	var monitor *octaviamonitors.Monitor
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		monitor, err = octaviamonitors.Create(lbClient, createOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to create vkcs_lb_monitor: %s", err)
	}

	// Wait for monitor to become active before continuing
	err = waitForLBMonitor(ctx, lbClient, parentPool, monitor, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(monitor.ID)

	return resourceMonitorRead(ctx, d, meta)
}

func resourceMonitorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	monitor, err := octaviamonitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "monitor"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_lb_monitor %s: %#v", d.Id(), monitor)

	d.Set("type", monitor.Type)
	d.Set("delay", monitor.Delay)
	d.Set("timeout", monitor.Timeout)
	d.Set("max_retries", monitor.MaxRetries)
	d.Set("max_retries_down", monitor.MaxRetriesDown)
	d.Set("url_path", monitor.URLPath)
	d.Set("http_method", monitor.HTTPMethod)
	d.Set("expected_codes", monitor.ExpectedCodes)
	d.Set("admin_state_up", monitor.AdminStateUp)
	d.Set("name", monitor.Name)
	d.Set("region", util.GetRegion(d, config))

	// OpenContrail workaround (https://github.com/terraform-provider-openstack/terraform-provider-openstack/issues/762)
	if len(monitor.Pools) > 0 && monitor.Pools[0].ID != "" {
		d.Set("pool_id", monitor.Pools[0].ID)
	}

	return nil

}

func resourceMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	var hasChange bool
	var updateOpts octaviamonitors.UpdateOpts
	if d.HasChange("url_path") {
		hasChange = true
		updateOpts.URLPath = d.Get("url_path").(string)
	}
	if d.HasChange("expected_codes") {
		hasChange = true
		updateOpts.ExpectedCodes = d.Get("expected_codes").(string)
	}
	if d.HasChange("delay") {
		hasChange = true
		updateOpts.Delay = d.Get("delay").(int)
	}
	if d.HasChange("timeout") {
		hasChange = true
		updateOpts.Timeout = d.Get("timeout").(int)
	}
	if d.HasChange("max_retries") {
		hasChange = true
		updateOpts.MaxRetries = d.Get("max_retries").(int)
	}
	if d.HasChange("max_retries_down") {
		hasChange = true
		updateOpts.MaxRetriesDown = d.Get("max_retries_down").(int)
	}
	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}
	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("http_method") {
		hasChange = true
		updateOpts.HTTPMethod = d.Get("http_method").(string)
	}

	if !hasChange {
		log.Printf("[DEBUG] vkcs_lb_monitor %s: nothing to update", d.Id())
		return resourceMonitorRead(ctx, d, meta)
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent vkcs_lb_pool %s: %s", poolID, err)
	}

	// Get a clean copy of the monitor.
	monitor, err := octaviamonitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_lb_monitor %s: %s", d.Id(), err)
	}

	// Wait for parent pool to become active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for monitor to become active before continuing.
	err = waitForLBMonitor(ctx, lbClient, parentPool, monitor, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] vkcs_lb_monitor %s update options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = octaviamonitors.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to update vkcs_lb_monitor %s: %s", d.Id(), err)
	}

	// Wait for monitor to become active before continuing
	err = waitForLBMonitor(ctx, lbClient, parentPool, monitor, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceMonitorRead(ctx, d, meta)
}

func resourceMonitorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Get a clean copy of the parent pool.
	poolID := d.Get("pool_id").(string)
	parentPool, err := pools.Get(lbClient, poolID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent vkcs_lb_pool (%s)"+
			" for the vkcs_lb_monitor: %s", poolID, err)
	}

	// Get a clean copy of the monitor.
	monitor, err := octaviamonitors.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Unable to retrieve vkcs_lb_monitor"))
	}

	// Wait for parent pool to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBPool(ctx, lbClient, parentPool, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting vkcs_lb_monitor %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = octaviamonitors.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_lb_monitor"))
	}

	// Wait for monitor to become DELETED
	err = waitForLBMonitor(ctx, lbClient, parentPool, monitor, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceMonitorImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	monitorID := parts[0]

	if len(monitorID) == 0 {
		return nil, fmt.Errorf("invalid format specified for vkcs_lb_monitor. Format must be <monitorID>[/<poolID>]")
	}

	d.SetId(monitorID)

	if len(parts) == 2 {
		d.Set("pool_id", parts[1])
	}

	return []*schema.ResourceData{d}, nil
}
