package lb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	ilisteners "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/lb/v2/listeners"
)

func ResourceListener() *schema.Resource {
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new Listener.",
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS", "TERMINATED_HTTPS",
				}, false),
				Description: "The protocol - can either be TCP, HTTP, HTTPS, TERMINATED_HTTPS, UDP. Changing this creates a new Listener.",
			},

			"protocol_port": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "The port on which to listen for client traffic. Changing this creates a new Listener.",
			},

			"loadbalancer_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The load balancer on which to provision this Listener. Changing this creates a new Listener.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Human-readable name for the Listener. Does not have to be unique.",
			},

			"default_pool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the default pool with which the Listener is associated.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description for the Listener.",
			},

			"connection_limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The maximum number of connections allowed for the Listener.",
			},

			"default_tls_container_ref": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A reference to a Keymanager Secrets container which stores TLS information. This is required if the protocol is `TERMINATED_HTTPS`.",
			},

			"sni_container_refs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of references to Keymanager Secrets containers which store SNI information.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "The administrative state of the Listener. A valid value is true (UP) or false (DOWN).",
			},

			"timeout_client_data": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The client inactivity timeout in milliseconds.",
			},

			"timeout_member_connect": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The member connection timeout in milliseconds.",
			},

			"timeout_member_data": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The member inactivity timeout in milliseconds.",
			},

			"timeout_tcp_inspect": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The time in milliseconds, to wait for additional TCP packets for content inspection.",
			},

			"insert_headers": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    false,
				Description: "The list of key value pairs representing headers to insert into the request before it is sent to the backend members. Changing this updates the headers of the existing listener.",
			},

			"allowed_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of CIDR blocks that are permitted to connect to this listener, denying all other source addresses. If not present, defaults to allow all.",
			},
		},
		Description: "Manages a listener resource within VKCS.",
	}
}

func resourceListenerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)

	// Wait for LoadBalancer to become active before continuing.
	err = waitForLBLoadBalancer(ctx, lbClient, d.Get("loadbalancer_id").(string), "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)

	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}

	var createOpts listeners.CreateOptsBuilder
	opts := listeners.CreateOpts{
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		SniContainerRefs:       sniContainerRefs,
		AdminStateUp:           &adminStateUp,
	}

	if v, ok := d.GetOk("connection_limit"); ok {
		connectionLimit := v.(int)
		opts.ConnLimit = &connectionLimit
	}

	if v, ok := d.GetOk("timeout_client_data"); ok {
		timeoutClientData := v.(int)
		opts.TimeoutClientData = &timeoutClientData
	}

	if v, ok := d.GetOk("timeout_member_connect"); ok {
		timeoutMemberConnect := v.(int)
		opts.TimeoutMemberConnect = &timeoutMemberConnect
	}

	if v, ok := d.GetOk("timeout_member_data"); ok {
		timeoutMemberData := v.(int)
		opts.TimeoutMemberData = &timeoutMemberData
	}

	if v, ok := d.GetOk("timeout_tcp_inspect"); ok {
		timeoutTCPInspect := v.(int)
		opts.TimeoutTCPInspect = &timeoutTCPInspect
	}

	// Get and check insert  headers map.
	rawHeaders := d.Get("insert_headers").(map[string]interface{})
	headers, err := expandLBListenerHeadersMap(rawHeaders)
	if err != nil {
		return diag.Errorf("unable to parse insert_headers argument: %s", err)
	}

	opts.InsertHeaders = headers

	if raw, ok := d.GetOk("allowed_cidrs"); ok {
		allowedCidrs := make([]string, len(raw.([]interface{})))
		for i, v := range raw.([]interface{}) {
			allowedCidrs[i] = v.(string)
		}
		opts.AllowedCIDRs = allowedCidrs
	}

	createOpts = opts

	log.Printf("[DEBUG] vkcs_lb_listener create options: %#v", createOpts)
	var listener *listeners.Listener
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		listener, err = ilisteners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating vkcs_lb_listener: %s", err)
	}

	d.SetId(listener.ID)

	// Wait for the listener to become ACTIVE.
	err = waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceListenerRead(ctx, d, meta)
}

func resourceListenerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	listener, err := ilisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "vkcs_lb_listener"))
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
	d.Set("region", util.GetRegion(d, config))

	// Required by import.
	if len(listener.Loadbalancers) > 0 {
		d.Set("loadbalancer_id", listener.Loadbalancers[0].ID)
	}

	if err := d.Set("insert_headers", listener.InsertHeaders); err != nil {
		return diag.Errorf("Unable to set vkcs_lb_listener insert_headers: %s", err)
	}

	return nil

}

func resourceListenerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := ilisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_lb_listener %s: %s", d.Id(), err)
	}

	// Wait for the listener to become ACTIVE.
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}
	var hasChange bool
	var opts listeners.UpdateOpts
	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		opts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		opts.Description = &description
	}

	if d.HasChange("connection_limit") {
		hasChange = true
		connLimit := d.Get("connection_limit").(int)
		opts.ConnLimit = &connLimit
	}

	if d.HasChange("timeout_client_data") {
		hasChange = true
		timeoutClientData := d.Get("timeout_client_data").(int)
		opts.TimeoutClientData = &timeoutClientData
	}

	if d.HasChange("timeout_member_connect") {
		hasChange = true
		timeoutMemberConnect := d.Get("timeout_member_connect").(int)
		opts.TimeoutMemberConnect = &timeoutMemberConnect
	}

	if d.HasChange("timeout_member_data") {
		hasChange = true
		timeoutMemberData := d.Get("timeout_member_data").(int)
		opts.TimeoutMemberData = &timeoutMemberData
	}

	if d.HasChange("timeout_tcp_inspect") {
		hasChange = true
		timeoutTCPInspect := d.Get("timeout_tcp_inspect").(int)
		opts.TimeoutTCPInspect = &timeoutTCPInspect
	}

	if d.HasChange("default_pool_id") {
		hasChange = true
		defaultPoolID := d.Get("default_pool_id").(string)
		opts.DefaultPoolID = &defaultPoolID
	}

	if d.HasChange("default_tls_container_ref") {
		hasChange = true
		defaultTLSContainerRef := d.Get("default_tls_container_ref").(string)
		opts.DefaultTlsContainerRef = &defaultTLSContainerRef
	}

	if d.HasChange("sni_container_refs") {
		hasChange = true
		var sniContainerRefs []string
		if raw, ok := d.GetOk("sni_container_refs"); ok {
			for _, v := range raw.([]interface{}) {
				sniContainerRefs = append(sniContainerRefs, v.(string))
			}
		}
		opts.SniContainerRefs = &sniContainerRefs
	}

	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &asu
	}

	if d.HasChange("insert_headers") {
		hasChange = true

		// Get and check insert headers map.
		rawHeaders := d.Get("insert_headers").(map[string]interface{})
		headers, err := expandLBListenerHeadersMap(rawHeaders)
		if err != nil {
			return diag.Errorf("unable to parse insert_headers argument: %s", err)
		}

		opts.InsertHeaders = &headers
	}

	if d.HasChange("allowed_cidrs") {
		hasChange = true
		var allowedCidrs []string
		if raw, ok := d.GetOk("allowed_cidrs"); ok {
			for _, v := range raw.([]interface{}) {
				allowedCidrs = append(allowedCidrs, v.(string))
			}
		}
		opts.AllowedCIDRs = &allowedCidrs
	}
	updateOpts := opts
	if !hasChange {
		log.Printf("[DEBUG] vkcs_lb_listener %s: nothing to update", d.Id())
		return resourceListenerRead(ctx, d, meta)
	}

	log.Printf("[DEBUG] vkcs_lb_listener %s update options: %#v", d.Id(), updateOpts)
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err = ilisteners.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
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
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Get a clean copy of the listener.
	listener, err := ilisteners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Unable to retrieve vkcs_lb_listener"))
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Deleting vkcs_lb_listener %s", d.Id())
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = ilisteners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_lb_listener"))
	}

	// Wait for the listener to become DELETED.
	err = waitForLBListener(ctx, lbClient, listener, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
