package networking

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iattributestags "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/attributestags"
	isubnets "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/subnets"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
)

func ResourceNetworkingSubnet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSubnetCreate,
		ReadContext:   resourceNetworkingSubnetRead,
		UpdateContext: resourceNetworkingSubnetUpdate,
		DeleteContext: resourceNetworkingSubnetDelete,
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
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a subnet. If omitted, the `region` argument of the provider is used. Changing this creates a new subnet.",
			},

			"network_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the parent network. Changing this creates a new subnet.",
			},

			"cidr": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"prefix_length"},
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.IsCIDR,
				},
				Description: "CIDR representing IP range for this subnet, based on IP version. You can omit this option if you are creating a subnet from a subnet pool.",
			},

			"prefix_length": {
				Type:          schema.TypeInt,
				ConflictsWith: []string{"cidr"},
				Optional:      true,
				ForceNew:      true,
				Description:   "The prefix length to use when creating a subnet from a subnet pool. The default subnet pool prefix length that was defined when creating the subnet pool will be used if not provided. Changing this creates a new subnet.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The name of the subnet. Changing this updates the name of the existing subnet.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Human-readable description of the subnet. Changing this updates the name of the existing subnet.",
			},

			"allocation_pool": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The starting address.",
						},
						"end": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ending address.",
						},
					},
				},
				Description: "A block declaring the start and end range of the IP addresses available for use with DHCP in this subnet. Multiple `allocation_pool` blocks can be declared, providing the subnet with more than one range of IP addresses to use with DHCP. However, each IP range must be from the same CIDR that the subnet is part of. The `allocation_pool` block is documented below.",
			},

			"gateway_ip": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"no_gateway"},
				Optional:      true,
				ForceNew:      false,
				Computed:      true,
				Description:   "Default gateway used by devices in this subnet. Leaving this blank and not setting `no_gateway` will cause a default gateway of `.1` to be used. Changing this updates the gateway IP of the existing subnet.",
			},

			"no_gateway": {
				Type:          schema.TypeBool,
				ConflictsWith: []string{"gateway_ip"},
				Optional:      true,
				Default:       false,
				ForceNew:      false,
				Description:   "Do not set a gateway IP on this subnet. Changing this removes or adds a default gateway IP of the existing subnet.",
			},

			"enable_dhcp": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     true,
				Description: "The administrative state of the network. Acceptable values are \"true\" and \"false\". Changing this value enables or disables the DHCP capabilities of the existing subnet. Defaults to true.",
			},

			"dns_nameservers": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    false,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "An array of DNS name server names used by hosts in this subnet. Changing this updates the DNS name servers for the existing subnet.",
			},

			"subnetpool_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the subnetpool associated with the subnet.",
			},

			"value_specs": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of additional options.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags for the subnet.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of ags assigned on the subnet, which have been explicitly and implicitly added.",
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
		Description: "Manages a subnet resource within VKCS.",
	}
}

func resourceNetworkingSubnetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	// Check nameservers.
	if err := networkingSubnetDNSNameserverAreUnique(d.Get("dns_nameservers").([]interface{})); err != nil {
		return diag.Errorf("vkcs_networking_subnet dns_nameservers argument is invalid: %s", err)
	}

	// Get raw allocation pool value.
	allocationPool := networkingSubnetGetRawAllocationPoolsValueToExpand(d)

	// Set basic options.
	createOpts := isubnets.SubnetCreateOpts{
		CreateOpts: subnets.CreateOpts{
			NetworkID:       d.Get("network_id").(string),
			Name:            d.Get("name").(string),
			Description:     d.Get("description").(string),
			AllocationPools: expandNetworkingSubnetAllocationPools(allocationPool),
			DNSNameservers:  util.ExpandToStringSlice(d.Get("dns_nameservers").([]interface{})),
			SubnetPoolID:    d.Get("subnetpool_id").(string),
			IPVersion:       gophercloud.IPVersion(4),
		},
		ValueSpecs: util.MapValueSpecs(d),
	}

	// Set CIDR if provided. Check if inferred subnet would match the provided cidr.
	if v, ok := d.GetOk("cidr"); ok {
		cidr := v.(string)
		_, netAddr, _ := net.ParseCIDR(cidr)
		if netAddr.String() != cidr {
			return diag.Errorf("cidr %s doesn't match subnet address %s for vkcs_networking_subnet", cidr, netAddr.String())
		}
		createOpts.CIDR = cidr
	}

	// Set gateway options if provided.
	if v, ok := d.GetOk("gateway_ip"); ok {
		gatewayIP := v.(string)
		createOpts.GatewayIP = &gatewayIP
	}

	noGateway := d.Get("no_gateway").(bool)
	if noGateway {
		gatewayIP := ""
		createOpts.GatewayIP = &gatewayIP
	}

	// Validate and set prefix options.
	if v, ok := d.GetOk("prefix_length"); ok {
		if d.Get("subnetpool_id").(string) == "" {
			return diag.Errorf("'prefix_length' is only valid if 'subnetpool_id' is set for vkcs_networking_subnet")
		}
		prefixLength := v.(int)
		createOpts.Prefixlen = prefixLength
	}

	// Set DHCP options if provided.
	enableDHCP := d.Get("enable_dhcp").(bool)
	createOpts.EnableDHCP = &enableDHCP

	log.Printf("[DEBUG] vkcs_networking_subnet create options: %#v", createOpts)
	s, err := isubnets.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_subnet: %s", err)
	}

	d.SetId(s.ID)

	log.Printf("[DEBUG] Waiting for vkcs_networking_subnet %s to become available", s.ID)
	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingSubnetStateRefreshFunc(networkingClient, s.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_subnet %s to become available: %s", s.ID, err)
	}

	tags := NetworkingAttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := iattributestags.ReplaceAll(networkingClient, "subnets", s.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating tags on vkcs_networking_subnet %s: %s", s.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_subnet %s", tags, s.ID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_subnet %s: %#v", s.ID, s)
	return resourceNetworkingSubnetRead(ctx, d, meta)
}

func resourceNetworkingSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var s subnetExtended
	err = isubnets.ExtractSubnetInto(isubnets.Get(networkingClient, d.Id()), &s)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_networking_subnet"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_subnet %s: %#v", d.Id(), s)

	d.Set("network_id", s.NetworkID)
	d.Set("cidr", s.CIDR)
	d.Set("name", s.Name)
	d.Set("description", s.Description)
	d.Set("dns_nameservers", s.DNSNameservers)
	d.Set("enable_dhcp", s.EnableDHCP)
	d.Set("network_id", s.NetworkID)
	d.Set("subnetpool_id", s.SubnetPoolID)

	NetworkingReadAttributesTags(d, s.Tags)

	// Set the allocation_pool attributes.
	allocationPools := flattenNetworkingSubnetAllocationPools(s.AllocationPools)
	if err := d.Set("allocation_pool", allocationPools); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_networking_subnet allocation_pool: %s", err)
	}

	// Set the subnet's "gateway_ip" and "no_gateway" attributes.
	d.Set("gateway_ip", s.GatewayIP)
	d.Set("no_gateway", false)
	if s.GatewayIP != "" {
		d.Set("no_gateway", false)
	} else {
		d.Set("no_gateway", true)
	}

	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", s.SDN)

	return nil
}

func resourceNetworkingSubnetUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var hasChange bool
	var updateOpts subnets.UpdateOpts

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

	if d.HasChange("gateway_ip") {
		hasChange = true
		updateOpts.GatewayIP = nil
		if v, ok := d.GetOk("gateway_ip"); ok {
			gatewayIP := v.(string)
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("no_gateway") {
		if d.Get("no_gateway").(bool) {
			hasChange = true
			gatewayIP := ""
			updateOpts.GatewayIP = &gatewayIP
		}
	}

	if d.HasChange("dns_nameservers") {
		if err := networkingSubnetDNSNameserverAreUnique(d.Get("dns_nameservers").([]interface{})); err != nil {
			return diag.Errorf("vkcs_networking_subnet dns_nameservers argument is invalid: %s", err)
		}
		hasChange = true
		nameservers := util.ExpandToStringSlice(d.Get("dns_nameservers").([]interface{}))
		updateOpts.DNSNameservers = &nameservers
	}

	if d.HasChange("enable_dhcp") {
		hasChange = true
		v := d.Get("enable_dhcp").(bool)
		updateOpts.EnableDHCP = &v
	}

	if d.HasChange("allocation_pool") {
		hasChange = true
		updateOpts.AllocationPools = expandNetworkingSubnetAllocationPools(d.Get("allocation_pool").(*schema.Set).List())
	}

	if hasChange {
		log.Printf("[DEBUG] Updating vkcs_networking_subnet %s with options: %#v", d.Id(), updateOpts)
		_, err = isubnets.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VKCS networking vkcs_networking_subnet %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		tags := NetworkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := iattributestags.ReplaceAll(networkingClient, "subnets", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating tags on vkcs_networking_subnet %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Updated tags %s on vkcs_networking_subnet %s", tags, d.Id())
	}

	return resourceNetworkingSubnetRead(ctx, d, meta)
}

func resourceNetworkingSubnetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingSubnetStateRefreshFuncDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_subnet %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
