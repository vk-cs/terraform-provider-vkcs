package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func resourceNetworkingFloatingIPAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingFloatingIPAssociateCreate,
		ReadContext:   resourceNetworkingFloatingIPAssociateRead,
		UpdateContext: resourceNetworkingFloatingIPAssociateUpdate,
		DeleteContext: resourceNetworkingFloatingIPAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"floating_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"fixed_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func resourceNetworkingFloatingIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	portID := d.Get("port_id").(string)
	fixedIP := d.Get("fixed_ip").(string)

	fipID, err := networkingFloatingIPV2ID(networkingClient, floatingIP)
	if err != nil {
		return diag.Errorf("Unable to get ID of vkcs_networking_floatingip_associate floating_ip %s: %s", floatingIP, err)
	}

	updateOpts := floatingips.UpdateOpts{
		PortID:  &portID,
		FixedIP: fixedIP,
	}

	log.Printf("[DEBUG] vkcs_networking_floatingip_associate create options: %#v", updateOpts)
	_, err = floatingips.Update(networkingClient, fipID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error associating vkcs_networking_floatingip_associate floating_ip %s with port %s: %s", fipID, portID, err)
	}

	d.SetId(fipID)

	return resourceNetworkingFloatingIPAssociateRead(ctx, d, meta)
}

func resourceNetworkingFloatingIPAssociateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	fip, err := floatingips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_floatingip_associate"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_floatingip_associate %s: %#v", d.Id(), fip)

	d.Set("floating_ip", fip.FloatingIP)
	d.Set("port_id", fip.PortID)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}

func resourceNetworkingFloatingIPAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	var updateOpts floatingips.UpdateOpts

	// port_id must always exists
	portID := d.Get("port_id").(string)
	updateOpts.PortID = &portID

	if d.HasChange("fixed_ip") {
		updateOpts.FixedIP = d.Get("fixed_ip").(string)
	}

	log.Printf("[DEBUG] vkcs_networking_floatingip_associate %s update options: %#v", d.Id(), updateOpts)
	_, err = floatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_networking_floatingip_associate %s: %s", d.Id(), err)
	}

	return resourceNetworkingFloatingIPAssociateRead(ctx, d, meta)
}

func resourceNetworkingFloatingIPAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	portID := d.Get("port_id").(string)
	updateOpts := floatingips.UpdateOpts{
		PortID: new(string),
	}

	log.Printf("[DEBUG] vkcs_networking_floatingip_associate disassociating options: %#v", updateOpts)
	_, err = floatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error disassociating vkcs_networking_floatingip_associate floating_ip %s with port %s: %s", d.Id(), portID, err)
	}

	return nil
}
