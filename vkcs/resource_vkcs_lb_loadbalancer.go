package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	neutronloadbalancers "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/loadbalancers"
)

func resourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerCreate,
		ReadContext:   resourceLoadBalancerRead,
		UpdateContext: resourceLoadBalancerUpdate,
		DeleteContext: resourceLoadBalancerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"vip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"vip_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vip_subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"vip_port_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		lbID      string
		vipPortID string
	)

	adminStateUp := d.Get("admin_state_up").(bool)

	if lbClient.Type == octaviaLBClientType {
		createOpts := octavialoadbalancers.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			VipNetworkID: d.Get("vip_network_id").(string),
			VipSubnetID:  d.Get("vip_subnet_id").(string),
			VipPortID:    d.Get("vip_port_id").(string),
			VipAddress:   d.Get("vip_address").(string),
			AdminStateUp: &adminStateUp,
		}

		// availability_zone requires octavia minor version 2.14. Only set when specified.
		if v, ok := d.GetOk("availability_zone"); ok {
			aZ := v.(string)
			createOpts.AvailabilityZone = aZ
		}

		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			createOpts.Tags = expandToStringSlice(tags)
		}

		log.Printf("[DEBUG][Octavia] vkcs_lb_loadbalancer create options: %#v", createOpts)
		lb, err := octavialoadbalancers.Create(lbClient, createOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating vkcs_lb_loadbalancer: %s", err)
		}
		lbID = lb.ID
		vipPortID = lb.VipPortID
	} else {
		createOpts := neutronloadbalancers.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			VipSubnetID:  d.Get("vip_subnet_id").(string),
			VipAddress:   d.Get("vip_address").(string),
			AdminStateUp: &adminStateUp,
		}

		log.Printf("[DEBUG][Neutron] vkcs_lb_loadbalancer create options: %#v", createOpts)
		lb, err := neutronloadbalancers.Create(lbClient, createOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating vkcs_lb_loadbalancer: %s", err)
		}
		lbID = lb.ID
		vipPortID = lb.VipPortID
	}

	// Wait for load-balancer to become active before continuing.
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBLoadBalancer(ctx, lbClient, lbID, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Once the load-balancer has been created, apply any requested security groups
	// to the port that was created behind the scenes.
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}
	if err := resourceLoadBalancerSetSecurityGroups(networkingClient, vipPortID, d); err != nil {
		return diag.Errorf("Error setting vkcs_lb_loadbalancer security groups: %s", err)
	}

	d.SetId(lbID)

	return resourceLoadBalancerRead(ctx, d, meta)
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var vipPortID string

	if lbClient.Type == octaviaLBClientType {
		lb, err := octavialoadbalancers.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(checkDeleted(d, err, "Unable to retrieve vkcs_lb_loadbalancer"))
		}

		log.Printf("[DEBUG][Octavia] Retrieved vkcs_lb_loadbalancer %s: %#v", d.Id(), lb)

		d.Set("name", lb.Name)
		d.Set("description", lb.Description)
		d.Set("vip_subnet_id", lb.VipSubnetID)
		d.Set("vip_network_id", lb.VipNetworkID)
		d.Set("vip_port_id", lb.VipPortID)
		d.Set("vip_address", lb.VipAddress)
		d.Set("admin_state_up", lb.AdminStateUp)
		d.Set("availability_zone", lb.AvailabilityZone)
		d.Set("region", getRegion(d, config))
		d.Set("tags", lb.Tags)
		vipPortID = lb.VipPortID
	} else {
		lb, err := neutronloadbalancers.Get(lbClient, d.Id()).Extract()
		if err != nil {
			return diag.FromErr(checkDeleted(d, err, "Unable to retrieve vkcs_lb_loadbalancer"))
		}

		log.Printf("[DEBUG][Neutron] Retrieved vkcs_lb_loadbalancer %s: %#v", d.Id(), lb)

		d.Set("name", lb.Name)
		d.Set("description", lb.Description)
		d.Set("vip_subnet_id", lb.VipSubnetID)
		d.Set("vip_port_id", lb.VipPortID)
		d.Set("vip_address", lb.VipAddress)
		d.Set("admin_state_up", lb.AdminStateUp)
		d.Set("region", getRegion(d, config))
		vipPortID = lb.VipPortID
	}

	// Get any security groups on the VIP Port.
	if vipPortID != "" {
		networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
		if err != nil {
			return diag.Errorf("Error creating OpenStack networking client: %s", err)
		}
		if err := resourceLoadBalancerGetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return diag.Errorf("Error getting port security groups for vkcs_lb_loadbalancer: %s", err)
		}
	}

	return nil
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	updateOpts, err := chooseLBLoadbalancerUpdateOpts(d, config)
	if err != nil {
		return diag.Errorf("Error building vkcs_lb_loadbalancer update options: %s", err)
	}

	if updateOpts != nil {
		// Wait for load-balancer to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBLoadBalancer(ctx, lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating vkcs_lb_loadbalancer %s with options: %#v", d.Id(), updateOpts)
		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err = neutronloadbalancers.Update(lbClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return diag.Errorf("Error updating vkcs_lb_loadbalancer %s: %s", d.Id(), err)
		}

		// Wait for load-balancer to become active before continuing.
		err = waitForLBLoadBalancer(ctx, lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Security Groups get updated separately.
	if d.HasChange("security_group_ids") {
		networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
		if err != nil {
			return diag.Errorf("Error creating OpenStack networking client: %s", err)
		}
		vipPortID := d.Get("vip_port_id").(string)
		if err := resourceLoadBalancerSetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return diag.Errorf("Error setting vkcs_lb_loadbalancer security groups: %s", err)
		}
	}

	return resourceLoadBalancerRead(ctx, d, meta)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := chooseLBClient(d, config)
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting vkcs_lb_loadbalancer %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = neutronloadbalancers.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_lb_loadbalancer"))
	}

	// Wait for load-balancer to become deleted.
	err = waitForLBLoadBalancer(ctx, lbClient, d.Id(), "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
