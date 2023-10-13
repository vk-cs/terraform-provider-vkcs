package networking

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func ResourceNetworkingFloating() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkFloatingIPCreate,
		ReadContext:   resourceNetworkFloatingIPRead,
		UpdateContext: resourceNetworkFloatingIPUpdate,
		DeleteContext: resourceNetworkFloatingIPDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a floating IP that can be used with another networking resource, such as a load balancer. If omitted, the `region` argument of the provider is used. Changing this creates a new floating IP (which may or may not have a different address).",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description for the floating IP.",
			},

			"address": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The actual floating IP address itself.",
			},

			"pool": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("OS_POOL_NAME", nil),
				Description: "The name of the pool from which to obtain the floating IP. Changing this creates a new floating IP.",
			},

			"port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "ID of an existing port with at least one IP address to associate with this floating IP.",
			},

			"fixed_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Fixed IP of the port to associate with this floating IP. Required if the port has multiple fixed IPs.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The subnet ID of the floating IP pool. Specify this if the floating IP network has multiple subnets.",
			},

			"subnet_ids": {
				Type:          schema.TypeList,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"subnet_id"},
				Description:   "A list of external subnet IDs to try over each to allocate a floating IP address. If a subnet ID in a list has exhausted floating IP pool, the next subnet ID will be tried. This argument is used only during the resource creation. Conflicts with a `subnet_id` argument.",
			},

			"value_specs": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of additional options.",
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
		Description: "Manages a floating IP resource within VKCS that can be used for load balancers.",
	}
}

func resourceNetworkFloatingIPCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	poolName := d.Get("pool").(string)
	poolID, err := networkingNetworkID(d, meta, poolName)
	if err != nil {
		return diag.Errorf("Error retrieving ID for vkcs_networking_floatingip pool name %s: %s", poolName, err)
	}
	if len(poolID) == 0 {
		return diag.Errorf("No network found with name: %s", poolName)
	}

	subnetID := d.Get("subnet_id").(string)
	var subnetIDs []string
	if v, ok := d.Get("subnet_ids").([]interface{}); ok {
		subnetIDs = make([]string, len(v))
		for i, v := range v {
			subnetIDs[i] = v.(string)
		}
	}

	if subnetID == "" && len(subnetIDs) > 0 {
		subnetID = subnetIDs[0]
	}

	createOpts := &floatingips.CreateOpts{
		FloatingNetworkID: poolID,
		Description:       d.Get("description").(string),
		PortID:            d.Get("port_id").(string),
		FixedIP:           d.Get("fixed_ip").(string),
		SubnetID:          subnetID,
	}

	finalCreateOpts := FloatingIPCreateOpts{
		CreateOpts: createOpts,
		ValueSpecs: util.MapValueSpecs(d),
	}

	var fip floatingIPExtended

	log.Printf("[DEBUG] vkcs_networking_floatingip create options: %#v", finalCreateOpts)

	if len(subnetIDs) == 0 {
		// floating IP allocation without a retry
		err = floatingips.Create(networkingClient, finalCreateOpts).ExtractInto(&fip)
		if err != nil {
			return diag.Errorf("Error creating vkcs_networking_floatingip: %s", err)
		}
	} else {
		// create a floatingip in a loop with the first available external subnet
		for i, subnetID := range subnetIDs {
			createOpts.SubnetID = subnetID

			log.Printf("[DEBUG] vkcs_networking_floatingip create options (try %d): %#v", i+1, finalCreateOpts)

			err = floatingips.Create(networkingClient, finalCreateOpts).ExtractInto(&fip)
			if err != nil {
				if retryOn409(err) {
					continue
				}
				return diag.Errorf("Error creating vkcs_networking_floatingip: %s", err)
			}
			break
		}
		// handle the last error
		if err != nil {
			return diag.Errorf("Error creating vkcs_networking_floatingip: %d subnets exhausted: %s", len(subnetIDs), err)
		}
	}

	d.SetId(fip.ID)

	log.Printf("[DEBUG] Waiting for vkcs_networking_floatingip %s to become available.", fip.ID)

	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    networkingFloatingIPV2StateRefreshFunc(networkingClient, fip.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_floatingip %s to become available: %s", fip.ID, err)
	}

	if createOpts.SubnetID != "" {
		// resourceNetworkFloatingIPRead doesn't handle this, since FIP GET request doesn't provide this info.
		d.Set("subnet_id", createOpts.SubnetID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_floatingip %s: %#v", fip.ID, fip)
	return resourceNetworkFloatingIPRead(ctx, d, meta)
}

func resourceNetworkFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	var fip floatingIPExtended

	err = floatingips.Get(networkingClient, d.Id()).ExtractInto(&fip)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_networking_floatingip"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_floatingip %s: %#v", d.Id(), fip)

	d.Set("description", fip.Description)
	d.Set("address", fip.FloatingIP.FloatingIP)
	d.Set("port_id", fip.PortID)
	d.Set("fixed_ip", fip.FixedIP)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", fip.SDN)

	poolName, err := networkingNetworkName(d, meta, fip.FloatingNetworkID)
	if err != nil {
		return diag.Errorf("Error retrieving pool name for vkcs_networking_floatingip %s: %s", d.Id(), err)
	}
	d.Set("pool", poolName)

	return nil
}

func resourceNetworkFloatingIPUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	var hasChange bool
	var updateOpts floatingips.UpdateOpts

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	// fixed_ip_address cannot be specified without a port_id
	if d.HasChange("port_id") || d.HasChange("fixed_ip") {
		hasChange = true
		portID := d.Get("port_id").(string)
		updateOpts.PortID = &portID
	}

	if d.HasChange("fixed_ip") {
		hasChange = true
		fixedIP := d.Get("fixed_ip").(string)
		updateOpts.FixedIP = fixedIP
	}

	if hasChange {
		log.Printf("[DEBUG] vkcs_networking_floatingip %s update options: %#v", d.Id(), updateOpts)
		_, err = floatingips.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating vkcs_networking_floatingip %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkFloatingIPRead(ctx, d, meta)
}

func resourceNetworkFloatingIPDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS network client: %s", err)
	}

	if err := floatingips.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_networking_floatingip"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE", "DOWN"},
		Target:     []string{"DELETED"},
		Refresh:    networkingFloatingIPV2StateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_floatingip %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
