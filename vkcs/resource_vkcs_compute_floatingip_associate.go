package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

func resourceComputeFloatingIPAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeFloatingIPAssociateCreate,
		ReadContext:   resourceComputeFloatingIPAssociateRead,
		DeleteContext: resourceComputeFloatingIPAssociateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the V2 Compute client. Keypairs are associated with accounts, but a Compute client is needed to create one. If omitted, the `region` argument of the provider is used. Changing this creates a new floatingip_associate.",
			},

			"floating_ip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The floating IP to associate.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The instance to associate the floating IP with.",
			},

			"fixed_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The specific IP address to direct traffic to.",
			},

			"wait_until_associated": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "In cases where the VKCS environment does not automatically wait until the association has finished, set this option to have Terraform poll the instance until the floating IP has been associated. Defaults to false.",
			},
		},
		Description: "Associate a floating IP to an instance.",
	}
}

func resourceComputeFloatingIPAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	fixedIP := d.Get("fixed_ip").(string)
	instanceID := d.Get("instance_id").(string)

	associateOpts := floatingips.AssociateOpts{
		FloatingIP: floatingIP,
		FixedIP:    fixedIP,
	}
	log.Printf("[DEBUG] vkcs_compute_floatingip_associate create options: %#v", associateOpts)

	err = floatingips.AssociateInstance(computeClient, instanceID, associateOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating vkcs_compute_floatingip_associate: %s", err)
	}

	// This API call should be synchronous, but we've had reports where it isn't.
	// If the user opted in to wait for association, then poll here.
	var waitUntilAssociated bool
	if v, ok := d.GetOk("wait_until_associated"); ok {
		if wua, ok := v.(bool); ok {
			waitUntilAssociated = wua
		}
	}

	if waitUntilAssociated {
		stateConf := &resource.StateChangeConf{
			Pending:    []string{"NOT_ASSOCIATED"},
			Target:     []string{"ASSOCIATED"},
			Refresh:    computeFloatingIPAssociateCheckAssociation(computeClient, instanceID, floatingIP),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// There's an API call to get this information, but it has been
	// deprecated. The Neutron API could be used, but I'm trying not
	// to mix service APIs. Therefore, a faux ID will be used.
	id := fmt.Sprintf("%s/%s/%s", floatingIP, instanceID, fixedIP)
	d.SetId(id)

	return resourceComputeFloatingIPAssociateRead(ctx, d, meta)
}

func resourceComputeFloatingIPAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	// Obtain relevant info from parsing the ID
	floatingIP, instanceID, fixedIP, err := parseComputeFloatingIPAssociateID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	// Now check and see whether the floating IP still exists.
	// First try to do this by querying the Network API.
	networkEnabled := true
	networkClient, err := config.NetworkingV2Client(getRegion(d, config), networking.SearchInAllSDNs)
	if err != nil {
		networkEnabled = false
	}

	var exists bool
	if networkEnabled {
		log.Printf("[DEBUG] Checking for vkcs_compute_floatingip_associate %s existence via Network API", d.Id())
		exists, err = computeFloatingIPAssociateNetworkExists(networkClient, floatingIP)
	} else {
		log.Printf("[DEBUG] Checking for vkcs_compute_floatingip_associate %s existence via Compute API", d.Id())
		exists, err = computeFloatingIPAssociateComputeExists(computeClient, floatingIP)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	if !exists {
		d.SetId("")
	}

	// Next, see if the instance still exists
	instance, err := servers.Get(computeClient, instanceID).Extract()
	if err != nil {
		if checkDeleted(d, err, "instance") == nil {
			return nil
		}
	}

	// Finally, check and see if the floating ip is still associated with the instance.
	var associated bool
	for _, networkAddresses := range instance.Addresses {
		for _, element := range networkAddresses.([]interface{}) {
			address := element.(map[string]interface{})
			if address["OS-EXT-IPS:type"] == "floating" && address["addr"] == floatingIP {
				associated = true
			}
		}
	}

	if !associated {
		d.SetId("")
	}

	// Set the attributes pulled from the composed resource ID
	d.Set("floating_ip", floatingIP)
	d.Set("instance_id", instanceID)
	d.Set("fixed_ip", fixedIP)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceComputeFloatingIPAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	floatingIP := d.Get("floating_ip").(string)
	instanceID := d.Get("instance_id").(string)

	disassociateOpts := floatingips.DisassociateOpts{
		FloatingIP: floatingIP,
	}
	log.Printf("[DEBUG] vkcs_compute_floatingip_associate %s delete options: %#v", d.Id(), disassociateOpts)

	err = floatingips.DisassociateInstance(computeClient, instanceID, disassociateOpts).ExtractErr()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault409); ok {
			// 409 is returned when floating ip address is not associated with an instance.
			log.Printf("[DEBUG] vkcs_compute_floatingip_associate %s is not associated with instance %s", d.Id(), instanceID)
		} else {
			return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_compute_floatingip_associate"))
		}
	}

	return nil
}
