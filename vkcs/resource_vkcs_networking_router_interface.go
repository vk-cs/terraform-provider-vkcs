package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func resourceNetworkingRouterInterface() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterInterfaceCreate,
		ReadContext:   resourceNetworkingRouterInterfaceRead,
		DeleteContext: resourceNetworkingRouterInterfaceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"sdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateSDN(),
			},
		},
	}
}

func resourceNetworkingRouterInterfaceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := routers.AddInterfaceOpts{
		SubnetID: d.Get("subnet_id").(string),
		PortID:   d.Get("port_id").(string),
	}

	log.Printf("[DEBUG] vkcs_networking_router_interface create options: %#v", createOpts)
	r, err := routers.AddInterface(networkingClient, d.Get("router_id").(string), createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_router_interface: %s", err)
	}

	log.Printf("[DEBUG] Waiting for vkcs_networking_router_interface %s to become available", r.PortID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    resourceNetworkingRouterInterfaceStateRefreshFunc(networkingClient, r.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_router_interface %s to become available: %s", r.ID, err)
	}

	d.SetId(r.PortID)

	log.Printf("[DEBUG] Created vkcs_networking_router_interface %s: %#v", r.ID, r)
	return resourceNetworkingRouterInterfaceRead(ctx, d, meta)
}

func resourceNetworkingRouterInterfaceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	r, err := ports.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving vkcs_networking_router_interface: %s", err)
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_router_interface %s: %#v", d.Id(), r)

	d.Set("router_id", r.DeviceID)
	d.Set("port_id", r.ID)
	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	// Set the subnet ID by looking at the port's FixedIPs.
	// If there's more than one FixedIP, do not set the subnet
	// as it's not possible to confidently determine which subnet
	// belongs to this interface. However, that situation should
	// not happen.
	if len(r.FixedIPs) != 1 {
		log.Printf("[DEBUG] Unable to set vkcs_networking_router_interface %s subnet_id", d.Id())
	} else {
		d.Set("subnet_id", r.FixedIPs[0].SubnetID)
	}

	return nil
}

func resourceNetworkingRouterInterfaceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingRouterInterfaceDeleteRefreshFunc(networkingClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_router_interface %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
