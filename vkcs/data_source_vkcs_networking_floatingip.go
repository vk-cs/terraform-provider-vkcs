package vkcs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func dataSourceNetworkingFloatingIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingFloatingIPRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"pool": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateSDN(),
			},
		},
	}
}

func dataSourceNetworkingFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
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

	tags := networkingAttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
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
	d.Set("all_tags", fip.Tags)
	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}
