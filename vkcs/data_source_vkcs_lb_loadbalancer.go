package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	octavialoadbalancers "github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
)

func dataSourceLoadBalancer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLoadBalancerRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used.",
			},

			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the Loadbalancer",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the Loadbalancer.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Human-readable description of the Loadbalancer.",
			},

			"vip_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ip address of the Loadbalancer.",
			},

			"vip_network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network on which to allocate the Loadbalancer's address. A tenant can only create Loadbalancers on networks authorized by policy (e.g. networks that belong to them or networks that are shared).  Changing this creates a new loadbalancer.",
			},

			"vip_subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The subnet on which the Loadbalancer's address is allocated.",
			},

			"vip_port_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The port UUID of the Loadbalancer.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "The administrative state of the Loadbalancer.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability zone of the Loadbalancer.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "A list of security group IDs applied to the Loadbalancer.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "A list of simple strings assigned to the loadbalancer.",
			},
		},
		Description: "Use this data source to get the details of a loadbalancer",
	}
}

func dataSourceLoadBalancerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	var vipPortID string
	lbID := d.Get("id").(string)
	lb, err := octavialoadbalancers.Get(lbClient, lbID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Unable to retrieve vkcs_lb_loadbalancer"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_lb_loadbalancer %s: %#v", lbID, lb)

	d.SetId(lb.ID)
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

	// Get any security groups on the VIP Port.
	if vipPortID != "" {
		networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
		if err != nil {
			return diag.Errorf("Error creating VKCS networking client: %s", err)
		}
		if err := resourceLoadBalancerGetSecurityGroups(networkingClient, vipPortID, d); err != nil {
			return diag.Errorf("Error getting port security groups for vkcs_lb_loadbalancer: %s", err)
		}
	}

	return nil
}
