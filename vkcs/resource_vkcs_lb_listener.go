package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	octavialisteners "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	neutronlisteners "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func resourceListener() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceListenerCreate,
		ReadContext:   resourceListenerRead,
		UpdateContext: resourceListenerUpdate,
		DeleteContext: resourceListenerDelete,
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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "SCTP", "HTTP", "HTTPS", "TERMINATED_HTTPS",
				}, false),
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"connection_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sni_container_refs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"timeout_client_data": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_member_connect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_member_data": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"timeout_tcp_inspect": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"insert_headers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},

			"allowed_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceListenerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for LoadBalancer to become active before continuing.
	err = waitForLBLoadBalancer(ctx, lbClient, d.Get("loadbalancer_id").(string), "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Choose either the Octavia or Neutron create options.
	createOpts, err := chooseLBListenerCreateOpts(d, config)
	if err != nil {
		return diag.Errorf("Error building vkcs_lb_listener create options: %s", err)
	}

	log.Printf("[DEBUG] vkcs_lb_listener create options: %#v", createOpts)
	var listener *neutronlisteners.Listener
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = neutronlisteners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating vkcs_lb_listener: %s", err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(listener.ID)

	return resourceListenerRead(ctx, d, meta)
}

func resourceListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Use Octavia listener body if Octavia/LBaaS is enabled.
	if config.UseOctavia {
		listener, err := octavialisteners.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(checkDeleted(d, err, "vkcs_lb_listener"))
		}

		log.Printf("[DEBUG] Retrieved vkcs_lb_listener %s: %#v", d.Id(), listener)

		d.Set("name", listener.Name)
		d.Set("protocol", listener.Protocol)
		d.Set("description", listener.Description)
		d.Set("protocol_port", listener.ProtocolPort)
		d.Set("admin_state_up", listener.AdminStateUp)
		d.Set("default_pool_id", listener.DefaultPoolID)
		d.Set("connection_limit", listener.ConnLimit)
		d.Set("timeout_client_data", listener.TimeoutClientData)
		d.Set("timeout_member_connect", listener.TimeoutMemberConnect)
		d.Set("timeout_member_data", listener.TimeoutMemberData)
		d.Set("timeout_tcp_inspect", listener.TimeoutTCPInspect)
		d.Set("sni_container_refs", listener.SniContainerRefs)
		d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
		d.Set("allowed_cidrs", listener.AllowedCIDRs)
		d.Set("region", getRegion(d, config))

		// Required by import.
		if len(listener.Loadbalancers) > 0 {
			d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
		}

		if err := d.Set("insert_headers", listener.InsertHeaders); err != nil {
			return diag.Errorf("Unable to set vkcs_lb_listener insert_headers: %s", err)
		}

		return nil
	}

	// Use Neutron/Networking in other case.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "vkcs_lb_listener"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_lb_listener %s: %#v", d.Id(), listener)

	// Required by import.
	if len(listener.Loadbalancers) > 0 {
		d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
	}

	d.Set("name", listener.Name)
	d.Set("protocol", listener.Protocol)
	d.Set("description", listener.Description)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("admin_state_up", listener.AdminStateUp)
	d.Set("default_pool_id", listener.DefaultPoolID)
	d.Set("connection_limit", listener.ConnLimit)
	d.Set("sni_container_refs", listener.SniContainerRefs)
	d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_lb_listener %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts, err := chooseLBListenerUpdateOpts(d, config)
	if err != nil {
		return diag.Errorf("Error building vkcs_lb_listener update options: %s", err)
	}
	if updateOpts == nil {
		log.Printf("[DEBUG] vkcs_lb_listener %s: nothing to update", d.Id())
		return resourceListenerRead(ctx, d, meta)
	}

	log.Printf("[DEBUG] vkcs_lb_listener %s update options: %#v", d.Id(), updateOpts)
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err = neutronlisteners.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error updating vkcs_lb_listener %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	err = waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceListenerRead(ctx, d, meta)
}

func resourceListenerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := neutronlisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Unable to retrieve vkcs_lb_listener"))
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting vkcs_lb_listener %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronlisteners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_lb_listener"))
	}

	// Wait for the listener to become DELETED.
	err = waitForLBListener(ctx, lbClient, listener, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
