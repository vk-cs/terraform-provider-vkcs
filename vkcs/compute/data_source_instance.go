package compute

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	iflavors "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/flavors"
	iservers "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servers"
	itags "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/tags"
)

func DataSourceComputeInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeInstanceRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.",
			},
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the instance",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the server.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image ID used to create the server.",
			},
			"image_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The image name used to create the server.",
			},
			"flavor_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor ID used to create the server.",
			},
			"flavor_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The flavor name used to create the server.",
			},
			"user_data": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The user data added when the server was created.",
				// just stash the hash for state & diff comparisons
			},
			"security_groups": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "An array of security group names associated with this server.",
			},
			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability zone of this server.",
			},
			"network": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The UUID of the network",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the network",
						},
						"port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The port UUID for this network",
						},
						"fixed_ip_v4": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IPv4 address assigned to this network port.",
						},
						"mac": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The MAC address assigned to this network interface.",
						},
					},
				},
				Description: "An array of maps, detailed below.",
			},
			"access_ip_v4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first IPv4 address assigned to this server.",
			},
			"key_pair": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the key pair assigned to this server.",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "A set of key/value pairs made available to the server.",
			},
			"power_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "VM state",
			},
			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags for the instance.",
			},
		},
		Description: "Use this data source to get the details of a running server",
	}
}

func dataSourceComputeInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	log.Print("[DEBUG] Creating compute client")
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	id := d.Get("id").(string)
	log.Printf("[DEBUG] Attempting to retrieve server %s", id)
	server, err := iservers.Get(computeClient, id).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "server"))
	}

	log.Printf("[DEBUG] Retrieved Server %s: %+v", id, server)

	d.SetId(server.ID)

	d.Set("name", server.Name)
	d.Set("image_id", server.Image["ID"])

	// Get the instance network and address information
	networks, err := flattenInstanceNetworks(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	// Determine the best IPv4 addresses to access the instance with
	hostv4 := getInstanceAccessAddresses(d, networks)

	// AccessIPv4 isn't standard in OpenStack, but there have been reports
	// of them being used in some environments.
	if server.AccessIPv4 != "" && hostv4 == "" {
		hostv4 = server.AccessIPv4
	}

	log.Printf("[DEBUG] Setting networks: %+v", networks)

	d.Set("network", networks)
	d.Set("access_ip_v4", hostv4)

	d.Set("metadata", server.Metadata)

	secGrpNames := []string{}
	for _, sg := range server.SecurityGroups {
		secGrpNames = append(secGrpNames, sg["name"].(string))
	}

	log.Printf("[DEBUG] Setting security groups: %+v", secGrpNames)

	d.Set("security_groups", secGrpNames)

	flavorID, ok := server.Flavor["id"].(string)
	if !ok {
		return diag.Errorf("Error setting VKCS server's flavor: %v", server.Flavor)
	}
	d.Set("flavor_id", flavorID)

	d.Set("key_pair", server.KeyName)
	flavor, err := iflavors.Get(computeClient, flavorID).Extract()
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("flavor_name", flavor.Name)

	// Set the instance's image information appropriately
	if err := setImageInformation(computeClient, server, d); err != nil {
		return diag.FromErr(err)
	}

	// Build a custom struct for the availability zone extension
	var serverWithAZ struct {
		servers.Server
		availabilityzones.ServerAvailabilityZoneExt
	}

	// Do another Get so the above work is not disturbed.
	err = iservers.Get(computeClient, d.Id()).ExtractInto(&serverWithAZ)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "server"))
	}
	// Set the availability zone
	d.Set("availability_zone", serverWithAZ.AvailabilityZone)

	// Set the region
	d.Set("region", util.GetRegion(d, config))

	// Set the current power_state
	currentStatus := strings.ToLower(server.Status)
	switch currentStatus {
	case "active", "shutoff", "error", "migrating", "shelved_offloaded", "shelved":
		d.Set("power_state", currentStatus)
	default:
		return diag.Errorf("Invalid power_state for instance %s: %s", d.Id(), server.Status)
	}

	// Populate tags.
	computeClient.Microversion = computeAPIMicroVersion
	instanceTags, err := itags.List(computeClient, server.ID).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to get tags for vkcs_compute_instance: %s", err)
	} else {
		d.Set("tags", instanceTags)
	}

	return nil
}
