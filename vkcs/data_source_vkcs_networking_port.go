package vkcs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func dataSourceNetworkingPort() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingPortRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve port ids. If omitted, the `region` argument of the provider is used.",
			},

			"port_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the port.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the port.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description of the port.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "The administrative state of the port.",
			},

			"network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the network the port belongs to.",
			},

			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The tenant_id of the owner of the port.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The project_id of the owner of the port.",
			},

			"device_owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The device owner of the port.",
			},

			"mac_address": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The MAC address of the port.",
			},

			"device_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID of the device the port belongs to.",
			},

			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsIPAddress,
				Description:  "The port IP address filter.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The status of the port.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of port security group IDs to filter.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of port tags to filter.",
			},

			"allowed_address_pairs": {
				Type:     schema.TypeSet,
				Computed: true,
				Set:      resourceNetworkingPortAllowedAddressPairsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The additional IP address.",
						},
						"mac_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The additional MAC address.",
						},
					},
				},
				Description: "An IP/MAC Address pair of additional IP addresses that can be active on this port. The structure is described below.",
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
				Description: "The set of security group IDs applied on the port.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The set of string tags applied on the port.",
			},

			"extra_dhcp_option": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the DHCP option.",
						},
						"value": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Value of the DHCP option.",
						},
					},
				},
				Description: "An extra DHCP option configured on the port. The structure is described below.",
			},

			"dns_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The port DNS name to filter.",
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
				ValidateDiagFunc: validateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the found port.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS port.",
	}
}

func dataSourceNetworkingPortRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	listOpts := ports.ListOpts{}
	var listOptsBuilder ports.ListOptsBuilder

	if v, ok := d.GetOk("port_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("admin_state_up"); ok {
		asu := v.(bool)
		listOpts.AdminStateUp = &asu
	}

	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("device_owner"); ok {
		listOpts.DeviceOwner = v.(string)
	}

	if v, ok := d.GetOk("mac_address"); ok {
		listOpts.MACAddress = v.(string)
	}

	if v, ok := d.GetOk("device_id"); ok {
		listOpts.DeviceID = v.(string)
	}

	tags := networkingAttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	listOptsBuilder = listOpts

	if v, ok := d.GetOk("dns_name"); ok {
		listOptsBuilder = dns.PortListOptsExt{
			ListOptsBuilder: listOptsBuilder,
			DNSName:         v.(string),
		}
	}

	allPages, err := ports.List(networkingClient, listOptsBuilder).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list vkcs_networking_ports: %s", err)
	}

	var allPorts []portExtended

	err = ports.ExtractPortsInto(allPages, &allPorts)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_networking_ports: %s", err)
	}

	if len(allPorts) == 0 {
		return diag.Errorf("No vkcs_networking_port found")
	}

	var portsList []portExtended

	// Filter returned Fixed IPs by a "fixed_ip".
	if v, ok := d.GetOk("fixed_ip"); ok {
		for _, p := range allPorts {
			for _, ipObject := range p.FixedIPs {
				if v.(string) == ipObject.IPAddress {
					portsList = append(portsList, p)
				}
			}
		}
		if len(portsList) == 0 {
			log.Printf("No vkcs_networking_port found after the 'fixed_ip' filter")
			return diag.Errorf("No vkcs_networking_port found")
		}
	} else {
		portsList = allPorts
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	if len(securityGroups) > 0 {
		var sgPorts []portExtended
		for _, p := range portsList {
			for _, sg := range p.SecurityGroups {
				if strSliceContains(securityGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}
		if len(sgPorts) == 0 {
			log.Printf("[DEBUG] No vkcs_networking_port found after the 'security_group_ids' filter")
			return diag.Errorf("No vkcs_networking_port found")
		}
		portsList = sgPorts
	}

	if len(portsList) > 1 {
		return diag.Errorf("More than one vkcs_networking_port found (%d)", len(portsList))
	}

	port := portsList[0]

	log.Printf("[DEBUG] Retrieved vkcs_networking_port %s: %+v", port.ID, port)
	d.SetId(port.ID)

	d.Set("port_id", port.ID)
	d.Set("name", port.Name)
	d.Set("description", port.Description)
	d.Set("admin_state_up", port.AdminStateUp)
	d.Set("network_id", port.NetworkID)
	d.Set("tenant_id", port.TenantID)
	d.Set("project_id", port.ProjectID)
	d.Set("device_owner", port.DeviceOwner)
	d.Set("mac_address", port.MACAddress)
	d.Set("device_id", port.DeviceID)
	d.Set("region", getRegion(d, config))
	d.Set("all_tags", port.Tags)
	d.Set("all_security_group_ids", port.SecurityGroups)
	d.Set("all_fixed_ips", expandNetworkingPortFixedIPToStringSlice(port.FixedIPs))
	d.Set("allowed_address_pairs", flattenNetworkingPortAllowedAddressPairs(port.MACAddress, port.AllowedAddressPairs))
	d.Set("extra_dhcp_option", flattenNetworkingPortDHCPOpts(port.ExtraDHCPOptsExt))
	d.Set("dns_name", port.DNSName)
	d.Set("dns_assignment", port.DNSAssignment)
	d.Set("sdn", getSDN(d))

	return nil
}
