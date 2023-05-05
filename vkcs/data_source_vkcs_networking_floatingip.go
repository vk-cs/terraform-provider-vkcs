package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func dataSourceNetworkingFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingFloatingIPRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve floating IP ids. If omitted, the `region` argument of the provider is used.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description of the floating IP.",
			},

			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IP address of the floating IP.",
			},

			"pool": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the pool from which the floating IP belongs to.",
			},

			"port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the port the floating IP is attached.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The owner of the floating IP.",
			},

			"fixed_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The specific IP address of the internal port which should be associated with the floating IP.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Status of the floating IP (ACTIVE/DOWN).",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the found floating IP.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS floating IP.",
	}
}

func dataSourceNetworkingFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	listOpts := floatingips.ListOpts{}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("address"); ok {
		listOpts.FloatingIP = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("pool"); ok {
		listOpts.FloatingNetworkID = v.(string)
	}

	if v, ok := d.GetOk("port_id"); ok {
		listOpts.PortID = v.(string)
	}

	if v, ok := d.GetOk("fixed_ip"); ok {
		listOpts.FixedIP = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	pages, err := floatingips.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list vkcs_networking_floatingips: %s", err)
	}

	var allFloatingIPs []floatingIPExtended

	err = floatingips.ExtractFloatingIPsInto(pages, &allFloatingIPs)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_networking_floatingips: %s", err)
	}

	if len(allFloatingIPs) < 1 {
		return diag.Errorf("No vkcs_networking_floatingip found")
	}

	if len(allFloatingIPs) > 1 {
		return diag.Errorf("More than one vkcs_networking_floatingip found")
	}

	fip := allFloatingIPs[0]

	log.Printf("[DEBUG] Retrieved vkcs_networking_floatingip %s: %+v", fip.ID, fip)
	d.SetId(fip.ID)

	d.Set("description", fip.Description)
	d.Set("address", fip.FloatingIP.FloatingIP)
	d.Set("pool", fip.FloatingNetworkID)
	d.Set("port_id", fip.PortID)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("tenant_id", fip.TenantID)
	d.Set("status", fip.Status)
	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}
