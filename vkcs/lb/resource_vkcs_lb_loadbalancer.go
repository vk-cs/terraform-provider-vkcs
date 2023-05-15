package lb

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
)

func ResourceLoadBalancer() *schema.Resource {
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new LB loadbalancer.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable name for the Loadbalancer. Does not have to be unique.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description for the Loadbalancer.",
			},

			"vip_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The ip address of the load balancer. Changing this creates a new loadbalancer.",
			},

			"vip_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The network on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.",
			},

			"vip_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The subnet on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.",
			},

			"vip_port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The port UUID that the loadbalancer will use. Changing this creates a new loadbalancer.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "The administrative state of the Loadbalancer. A valid value is true (UP) or false (DOWN).",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The availability zone of the Loadbalancer. Changing this creates a new loadbalancer.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "A list of security group IDs to apply to the loadbalancer. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance).",
				Deprecated:  "This argument is deprecated, please do not use it.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "A list of simple strings assigned to the loadbalancer.",
			},
		},
		Description: "Manages a loadbalancer resource within VKCS.",
	}
}

func resourceLoadBalancerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	var lbID string

	adminStateUp := d.Get("admin_state_up").(bool)

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
		createOpts.Tags = util.ExpandToStringSlice(tags)
	}

	log.Printf("[DEBUG][Octavia] vkcs_lb_loadbalancer create options: %#v", createOpts)
	lb, err := octavialoadbalancers.Create(lbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_lb_loadbalancer: %s", err)
	}
	lbID = lb.ID

	// Wait for load-balancer to become active before continuing.
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBLoadBalancer(ctx, lbClient, lbID, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(lbID)

	return resourceLoadBalancerRead(ctx, d, meta)
}

func resourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	var vipPortID string

	lb, err := octavialoadbalancers.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Unable to retrieve vkcs_lb_loadbalancer"))
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
	d.Set("region", util.GetRegion(d, config))
	d.Set("tags", lb.Tags)
	vipPortID = lb.VipPortID

	// Get any security groups on the VIP Port.
	if vipPortID != "" {
		networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.SearchInAllSDNs)
		if err != nil {
			return diag.Errorf("Error creating VKCS networking client: %s", err)
		}
		if err := resourceLoadBalancerGetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return diag.Errorf("Error getting port security groups for vkcs_lb_loadbalancer: %s", err)
		}
	}

	return nil
}

func resourceLoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}
	var updateOpts octavialoadbalancers.UpdateOpts
	var hasChange bool
	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if d.HasChange("tags") {
		hasChange = true
		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			tagsToUpdate := util.ExpandToStringSlice(tags)
			updateOpts.Tags = &tagsToUpdate
		} else {
			updateOpts.Tags = &[]string{}
		}
	}

	if hasChange {
		// Wait for load-balancer to become active before continuing.
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBLoadBalancer(ctx, lbClient, d.Id(), "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updating vkcs_lb_loadbalancer %s with options: %#v", d.Id(), updateOpts)
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err = octavialoadbalancers.Update(lbClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return util.CheckForRetryableError(err)
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

	return resourceLoadBalancerRead(ctx, d, meta)
}

func resourceLoadBalancerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting vkcs_lb_loadbalancer %s", d.Id())
	timeout := d.Timeout(schema.TimeoutDelete)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = octavialoadbalancers.Delete(lbClient, d.Id(), nil).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_lb_loadbalancer"))
	}

	// Wait for load-balancer to become deleted.
	err = waitForLBLoadBalancer(ctx, lbClient, d.Id(), "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
