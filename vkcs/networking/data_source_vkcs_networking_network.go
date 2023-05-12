package networking

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func DataSourceNetworkingNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingNetworkRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve networks ids. If omitted, the `region` argument of the provider is used.",
			},

			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the network.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the network.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description of the network.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The status of the network.",
			},

			"matching_subnet_cidr": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The CIDR of a subnet within the network.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The owner of the network.",
			},

			"admin_state_up": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The administrative state of the network.",
			},

			"shared": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Specifies whether the network resource can be accessed by any tenant or not.",
			},

			"external": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The external routing facility of the network.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of network tags to filter.",
			},

			"subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of subnet IDs belonging to the network.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The set of string tags applied on the network.",
			},

			"private_dns_domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Private dns domain name",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},

			"vkcs_services_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Specifies whether VKCS services access is enabled.",
			},

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the found network.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS network.",
	}
}

func dataSourceNetworkingNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	// Prepare basic listOpts.
	var listOpts networks.ListOptsBuilder

	var status string
	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}

	listOpts = networks.ListOpts{
		ID:          d.Get("network_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
		Status:      status,
	}

	// Add the external attribute if specified.
	if v, ok := d.GetOk("external"); ok {
		isExternal := v.(bool)
		listOpts = external.ListOptsExt{
			ListOptsBuilder: listOpts,
			External:        &isExternal,
		}
	}

	tags := NetworkingAttributesTags(d)
	if len(tags) > 0 {
		listOpts = networks.ListOpts{Tags: strings.Join(tags, ",")}
	}

	pages, err := networks.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	// First extract into a normal networks.Network in order to see if
	// there were any results at all.
	tmpAllNetworks, err := networks.ExtractNetworks(pages)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(tmpAllNetworks) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var allNetworks []networkExtended
	err = networks.ExtractNetworksInto(pages, &allNetworks)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_networking_networks: %s", err)
	}

	var refinedNetworks []networkExtended
	if cidr := d.Get("matching_subnet_cidr").(string); cidr != "" {
		for _, n := range allNetworks {
			for _, s := range n.Subnets {
				subnet, err := subnets.Get(networkingClient, s).Extract()
				if err != nil {
					if _, ok := err.(gophercloud.ErrDefault404); ok {
						continue
					}
					return diag.Errorf("Unable to retrieve vkcs_networking_network subnet: %s", err)
				}
				if cidr == subnet.CIDR {
					refinedNetworks = append(refinedNetworks, n)
				}
			}
		}
	} else {
		refinedNetworks = allNetworks
	}

	if len(refinedNetworks) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNetworks) > 1 {
		return diag.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	network := refinedNetworks[0]

	log.Printf("[DEBUG] Retrieved vkcs_networking_network %s: %+v", network.ID, network)
	d.SetId(network.ID)

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("admin_state_up", strconv.FormatBool(network.AdminStateUp))
	d.Set("shared", strconv.FormatBool(network.Shared))
	d.Set("external", network.External)
	d.Set("tenant_id", network.TenantID)
	d.Set("subnets", network.Subnets)
	d.Set("all_tags", network.Tags)
	d.Set("region", util.GetRegion(d, config))
	d.Set("private_dns_domain", network.PrivateDNSDomain)
	d.Set("sdn", GetSDN(d))
	d.Set("vkcs_services_access", network.ServicesAccess)

	return nil
}
