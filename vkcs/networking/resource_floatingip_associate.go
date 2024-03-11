package networking

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	ifloatingips "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/floatingips"
)

func ResourceNetworkingFloatingIPAssociate() *schema.Resource {
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a floating IP that can be used with another networking resource, such as a load balancer. If omitted, the `region` argument of the provider is used. Changing this creates a new floating IP (which may or may not have a different address).",
			},

			"floating_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "IP Address of an existing floating IP.",
			},

			"port_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ID of an existing port with at least one IP address to associate with this floating IP.",
			},

			"fixed_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "One of the port's IP addresses.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
			},
		},
		Description: "Associates a floating IP to a port. This can be done only if port is assigned to router connected to external network. This is useful for situations where you have a pre-allocated floating IP or are unable to use the `vkcs_networking_floatingip` resource to create a floating IP.",
	}
}

func resourceNetworkingFloatingIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
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
	_, err = ifloatingips.Update(networkingClient, fipID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error associating vkcs_networking_floatingip_associate floating_ip %s with port %s: %s", fipID, portID, err)
	}

	d.SetId(fipID)

	return resourceNetworkingFloatingIPAssociateRead(ctx, d, meta)
}

func resourceNetworkingFloatingIPAssociateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	var fip floatingIPExtended

	err = ifloatingips.Get(networkingClient, d.Id()).ExtractInto(&fip)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_networking_floatingip_associate"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_floatingip_associate %s: %#v", d.Id(), fip)

	d.Set("floating_ip", fip.FloatingIP.FloatingIP)
	d.Set("port_id", fip.PortID)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", fip.SDN)

	return nil
}

func resourceNetworkingFloatingIPAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	var updateOpts floatingips.UpdateOpts

	// port_id must always exists
	portID := d.Get("port_id").(string)
	updateOpts.PortID = &portID

	if d.HasChange("fixed_ip") {
		updateOpts.FixedIP = d.Get("fixed_ip").(string)
	}

	log.Printf("[DEBUG] vkcs_networking_floatingip_associate %s update options: %#v", d.Id(), updateOpts)
	_, err = ifloatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_networking_floatingip_associate %s: %s", d.Id(), err)
	}

	return resourceNetworkingFloatingIPAssociateRead(ctx, d, meta)
}

func resourceNetworkingFloatingIPAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	portID := d.Get("port_id").(string)
	updateOpts := floatingips.UpdateOpts{
		PortID: new(string),
	}

	log.Printf("[DEBUG] vkcs_networking_floatingip_associate disassociating options: %#v", updateOpts)
	_, err = ifloatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error disassociating vkcs_networking_floatingip_associate floating_ip %s with port %s: %s", d.Id(), portID, err)
	}

	return nil
}
