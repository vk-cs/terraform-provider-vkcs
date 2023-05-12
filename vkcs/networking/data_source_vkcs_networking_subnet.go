package networking

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func DataSourceNetworkingSubnet() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSubnetRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve subnet ids. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The name of the subnet.",
			},

			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "Human-readable description of the subnet.",
			},

			"dhcp_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "If the subnet has DHCP enabled.",
			},

			"network_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The ID of the network the subnet belongs to.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The owner of the subnet.",
			},

			"gateway_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The IP of the subnet's gateway.",
			},

			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The CIDR of the subnet.",
			},

			"subnet_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The ID of the subnet.",
			},

			// Computed values
			"allocation_pools": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"start": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The starting address.",
						},
						"end": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ending address.",
						},
					},
				},
				Description: "Allocation pools of the subnet.",
			},

			"enable_dhcp": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the subnet has DHCP enabled or not.",
			},

			"dns_nameservers": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "DNS Nameservers of the subnet.",
			},

			"host_routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "Host Routes of the subnet.",
			},

			"subnetpool_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The ID of the subnetpool associated with the subnet.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of subnet tags to filter.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags applied on the subnet.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the found subnet.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS subnet.",
	}
}

func dataSourceNetworkingSubnetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	listOpts := subnets.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("dhcp_enabled"); ok {
		enableDHCP := v.(bool)
		listOpts.EnableDHCP = &enableDHCP
	}

	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("gateway_ip"); ok {
		listOpts.GatewayIP = v.(string)
	}

	if v, ok := d.GetOk("cidr"); ok {
		listOpts.CIDR = v.(string)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("subnetpool_id"); ok {
		listOpts.SubnetPoolID = v.(string)
	}

	tags := NetworkingAttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	pages, err := subnets.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_networking_subnet: %s", err)
	}

	allSubnets, err := subnets.ExtractSubnets(pages)
	if err != nil {
		return diag.Errorf("Unable to extract vkcs_networking_subnet: %s", err)
	}

	if len(allSubnets) < 1 {
		return diag.Errorf("Your query returned no vkcs_networking_subnet. " +
			"Please change your search criteria and try again.")
	}

	if len(allSubnets) > 1 {
		return diag.Errorf("Your query returned more than one vkcs_networking_subnet." +
			" Please try a more specific search criteria")
	}

	subnet := allSubnets[0]

	log.Printf("[DEBUG] Retrieved vkcs_networking_subnet %s: %+v", subnet.ID, subnet)
	d.SetId(subnet.ID)

	d.Set("name", subnet.Name)
	d.Set("description", subnet.Description)
	d.Set("tenant_id", subnet.TenantID)
	d.Set("network_id", subnet.NetworkID)
	d.Set("cidr", subnet.CIDR)
	d.Set("gateway_ip", subnet.GatewayIP)
	d.Set("enable_dhcp", subnet.EnableDHCP)
	d.Set("subnetpool_id", subnet.SubnetPoolID)
	d.Set("all_tags", subnet.Tags)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", GetSDN(d))

	if err := d.Set("dns_nameservers", subnet.DNSNameservers); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_networking_subnet dns_nameservers: %s", err)
	}

	if err := d.Set("host_routes", subnet.HostRoutes); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_networking_subnet host_routes: %s", err)
	}

	allocationPools := flattenNetworkingSubnetAllocationPools(subnet.AllocationPools)
	if err := d.Set("allocation_pools", allocationPools); err != nil {
		log.Printf("[DEBUG] Unable to set vkcs_networking_subnet allocation_pools: %s", err)
	}

	return nil
}
