package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func resourceNetworkingPort() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingPortCreate,
		ReadContext:   resourceNetworkingPortRead,
		UpdateContext: resourceNetworkingPortUpdate,
		DeleteContext: resourceNetworkingPortDelete,
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
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a port. If omitted, the `region` argument of the provider is used. Changing this creates a new port.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "A unique name for the port. Changing this updates the `name` of an existing port.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description of the port. Changing this updates the `description` of an existing port.",
			},

			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the network to attach the port to. Changing this creates a new port.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     true,
				Description: "Administrative up/down status for the port (must be `true` or `false` if provided). Changing this updates the `admin_state_up` of an existing port.",
			},

			"mac_address": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "Specify a specific MAC address for the port. Changing this creates a new port.",
			},

			"device_owner": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The device owner of the port. Changing this creates a new port.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "(Conflicts with `no_security_groups`) A list of security group IDs to apply to the port. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance).",
			},

			"no_security_groups": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "(Conflicts with `security_group_ids`) If set to `true`, then no security groups are applied to the port. If set to `false` and no `security_group_ids` are specified, then the port will yield to the default behavior of the Networking service, which is to usually apply the \"default\" security group.",
			},

			"device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The ID of the device attached to the port. Changing this creates a new port.",
			},

			"fixed_ip": {
				Type:          schema.TypeList,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"no_fixed_ip"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subnet_id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Subnet in which to allocate IP address for this port.",
						},
						"ip_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "IP address desired in the subnet for this port. If you don't specify `ip_address`, an available IP address from the specified subnet will be allocated to this port. This field will not be populated if it is left blank or omitted. To retrieve the assigned IP address, use the `all_fixed_ips` attribute.",
						},
					},
				},
				Description: "(Conflicts with `no_fixed_ip`) An array of desired IPs for this port. The structure is described below.",
			},

			"no_fixed_ip": {
				Type:          schema.TypeBool,
				Optional:      true,
				ForceNew:      false,
				ConflictsWith: []string{"fixed_ip"},
				Description:   "(Conflicts with `fixed_ip`) Create a port with no fixed IP address. This will also remove any fixed IPs previously set on a port. `true` is the only valid value for this argument.",
			},

			"allowed_address_pairs": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Set:      resourceNetworkingPortAllowedAddressPairsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The additional IP address.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The additional MAC address.",
						},
					},
				},
				Description: "An IP/MAC Address pair of additional IP addresses that can be active on this port. The structure is described below.",
			},

			"extra_dhcp_option": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Name of the DHCP option.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Value of the DHCP option.",
						},
					},
				},
				Description: "An extra DHCP option that needs to be configured on the port. The structure is described below. Can be specified multiple times.",
			},

			"value_specs": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of additional options.",
			},

			"all_fixed_ips": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of Fixed IP addresses on the port in the order returned by the Network v2 API.",
			},

			"all_security_group_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The collection of Security Group IDs on the port which have been explicitly and implicitly added.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags for the port.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of tags assigned on the port, which have been explicitly and implicitly added.",
			},

			"port_security_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to explicitly enable or disable port security on the port. Port Security is usually enabled by default, so omitting argument will usually result in a value of `true`. Setting this explicitly to `false` will disable port security. In order to disable port security, the port must not have any security groups. Valid values are `true` and `false`.",
			},

			"dns_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The port DNS name.",
			},

			"dns_assignment": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeMap},
				Description: "The list of maps representing port DNS assignments.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: validateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},
		},
		Description: "Manages a port resource within VKCS.",
	}
}

func resourceNetworkingPortCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	noSecurityGroups := d.Get("no_security_groups").(bool)

	// Check and make sure an invalid security group configuration wasn't given.
	if noSecurityGroups && len(securityGroups) > 0 {
		return diag.Errorf("Cannot have both no_security_groups and security_group_ids set for vkcs_networking_port")
	}

	allowedAddressPairs := d.Get("allowed_address_pairs").(*schema.Set)
	createOpts := PortCreateOpts{
		ports.CreateOpts{
			Name:                d.Get("name").(string),
			Description:         d.Get("description").(string),
			NetworkID:           d.Get("network_id").(string),
			MACAddress:          d.Get("mac_address").(string),
			DeviceOwner:         d.Get("device_owner").(string),
			DeviceID:            d.Get("device_id").(string),
			FixedIPs:            expandNetworkingPortFixedIP(d),
			AllowedAddressPairs: expandNetworkingPortAllowedAddressPairs(allowedAddressPairs),
		},
		MapValueSpecs(d),
	}

	asu := d.Get("admin_state_up").(bool)
	createOpts.AdminStateUp = &asu

	if noSecurityGroups {
		securityGroups = []string{}
		createOpts.SecurityGroups = &securityGroups
	}

	// Only set SecurityGroups if one was specified.
	// Otherwise this would mimic the no_security_groups action.
	if len(securityGroups) > 0 {
		createOpts.SecurityGroups = &securityGroups
	}

	// Declare a finalCreateOpts interface to hold either the
	// base create options or the extended DHCP options.
	var finalCreateOpts ports.CreateOptsBuilder
	finalCreateOpts = createOpts

	dhcpOpts := d.Get("extra_dhcp_option").(*schema.Set)
	if dhcpOpts.Len() > 0 {
		finalCreateOpts = extradhcpopts.CreateOptsExt{
			CreateOptsBuilder: createOpts,
			ExtraDHCPOpts:     expandNetworkingPortDHCPOptsCreate(dhcpOpts),
		}
	}

	// Add the port security attribute
	pse := d.Get("port_security_enabled").(bool)
	finalCreateOpts = portsecurity.PortCreateOptsExt{
		CreateOptsBuilder:   finalCreateOpts,
		PortSecurityEnabled: &pse,
	}

	if dnsName := d.Get("dns_name").(string); dnsName != "" {
		finalCreateOpts = dns.PortCreateOptsExt{
			CreateOptsBuilder: finalCreateOpts,
			DNSName:           dnsName,
		}
	}

	log.Printf("[DEBUG] vkcs_networking_port create options: %#v", finalCreateOpts)

	// Create a Neutron port and set extra options if they're specified.
	var port portExtended

	err = ports.Create(networkingClient, finalCreateOpts).ExtractInto(&port)
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_port: %s", err)
	}

	log.Printf("[DEBUG] Waiting for vkcs_networking_port %s to become available.", port.ID)

	stateConf := &resource.StateChangeConf{
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    resourceNetworkingPortStateRefreshFunc(networkingClient, port.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_port %s to become available: %s", port.ID, err)
	}

	d.SetId(port.ID)

	tags := networkingAttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "ports", port.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_port %s: %s", port.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_port %s", tags, port.ID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_port %s: %#v", port.ID, port)
	return resourceNetworkingPortRead(ctx, d, meta)
}

func resourceNetworkingPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var port portExtended
	err = ports.Get(networkingClient, d.Id()).ExtractInto(&port)
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_port"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_port %s: %#v", d.Id(), port)

	d.Set("name", port.Name)
	d.Set("description", port.Description)
	d.Set("admin_state_up", port.AdminStateUp)
	d.Set("network_id", port.NetworkID)
	d.Set("mac_address", port.MACAddress)
	d.Set("device_owner", port.DeviceOwner)
	d.Set("device_id", port.DeviceID)

	networkingReadAttributesTags(d, port.Tags)

	// Set a slice of all returned Fixed IPs.
	// This will be in the order returned by the API,
	// which is usually alpha-numeric.
	d.Set("all_fixed_ips", expandNetworkingPortFixedIPToStringSlice(port.FixedIPs))

	// Set all security groups.
	// This can be different from what the user specified since
	// the port can have the "default" group automatically applied.
	d.Set("all_security_group_ids", port.SecurityGroups)

	d.Set("allowed_address_pairs", flattenNetworkingPortAllowedAddressPairs(port.MACAddress, port.AllowedAddressPairs))
	d.Set("extra_dhcp_option", flattenNetworkingPortDHCPOpts(port.ExtraDHCPOptsExt))
	d.Set("port_security_enabled", port.PortSecurityEnabled)
	d.Set("dns_name", port.DNSName)
	d.Set("dns_assignment", port.DNSAssignment)

	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}

func resourceNetworkingPortUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	noSecurityGroups := d.Get("no_security_groups").(bool)

	// Check and make sure an invalid security group configuration wasn't given.
	if noSecurityGroups && len(securityGroups) > 0 {
		return diag.Errorf("Cannot have both no_security_groups and security_group_ids set for vkcs_networking_port")
	}

	var hasChange bool
	var updateOpts ports.UpdateOpts

	if d.HasChange("allowed_address_pairs") {
		hasChange = true
		allowedAddressPairs := d.Get("allowed_address_pairs").(*schema.Set)
		aap := expandNetworkingPortAllowedAddressPairs(allowedAddressPairs)
		updateOpts.AllowedAddressPairs = &aap
	}

	if d.HasChange("no_security_groups") {
		if noSecurityGroups {
			hasChange = true
			v := []string{}
			updateOpts.SecurityGroups = &v
		}
	}

	if d.HasChange("security_group_ids") {
		hasChange = true
		updateOpts.SecurityGroups = &securityGroups
	}

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

	if d.HasChange("device_owner") {
		hasChange = true
		deviceOwner := d.Get("device_owner").(string)
		updateOpts.DeviceOwner = &deviceOwner
	}

	if d.HasChange("device_id") {
		hasChange = true
		deviceID := d.Get("device_id").(string)
		updateOpts.DeviceID = &deviceID
	}

	if d.HasChange("fixed_ip") || d.HasChange("no_fixed_ip") {
		fixedIPs := expandNetworkingPortFixedIP(d)
		if fixedIPs != nil {
			hasChange = true
			updateOpts.FixedIPs = fixedIPs
		}
	}

	var finalUpdateOpts ports.UpdateOptsBuilder
	finalUpdateOpts = updateOpts

	if d.HasChange("port_security_enabled") {
		hasChange = true
		portSecurityEnabled := d.Get("port_security_enabled").(bool)
		finalUpdateOpts = portsecurity.PortUpdateOptsExt{
			UpdateOptsBuilder:   finalUpdateOpts,
			PortSecurityEnabled: &portSecurityEnabled,
		}
	}

	// Next, perform any dhcp option changes.
	if d.HasChange("extra_dhcp_option") {
		hasChange = true

		o, n := d.GetChange("extra_dhcp_option")
		oldDHCPOpts := o.(*schema.Set)
		newDHCPOpts := n.(*schema.Set)

		deleteDHCPOpts := oldDHCPOpts.Difference(newDHCPOpts)
		addDHCPOpts := newDHCPOpts.Difference(oldDHCPOpts)

		updateExtraDHCPOpts := expandNetworkingPortDHCPOptsUpdate(deleteDHCPOpts, addDHCPOpts)
		finalUpdateOpts = extradhcpopts.UpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			ExtraDHCPOpts:     updateExtraDHCPOpts,
		}
	}

	if d.HasChange("dns_name") {
		hasChange = true

		dnsName := d.Get("dns_name").(string)
		finalUpdateOpts = dns.PortUpdateOptsExt{
			UpdateOptsBuilder: finalUpdateOpts,
			DNSName:           &dnsName,
		}
	}

	// At this point, perform the update for all "standard" port changes.
	if hasChange {
		log.Printf("[DEBUG] vkcs_networking_port %s update options: %#v", d.Id(), finalUpdateOpts)
		_, err = ports.Update(networkingClient, d.Id(), finalUpdateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VKCS networking Port: %s", err)
		}
	}

	// Next, perform any required updates to the tags.
	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "ports", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_port %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_port %s", tags, d.Id())
	}

	return resourceNetworkingPortRead(ctx, d, meta)
}

func resourceNetworkingPortDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	if err := ports.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_networking_port"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingPortStateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_port %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
