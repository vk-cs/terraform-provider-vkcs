package vkcs

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func dataSourceNetworkingNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingNetworkRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"matching_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The ID of the Tenant (Identity v2) or Project (Identity v3)\n" +
					"to login with.",
			},

			"admin_state_up": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"shared": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"external": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"private_dns_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSDN(),
			},
		},
	}
}

func dataSourceNetworkingNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
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
	if v, ok := d.GetOkExists("external"); ok {
		isExternal := v.(bool)
		listOpts = external.ListOptsExt{
			ListOptsBuilder: listOpts,
			External:        &isExternal,
		}
	}

	tags := networkingAttributesTags(d)
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
	d.Set("region", getRegion(d, config))
	d.Set("private_dns_domain", network.PrivateDNSDomain)
	d.Set("sdn", getSDN(d))

	return nil
}
