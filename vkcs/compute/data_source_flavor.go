package compute

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
	iflavors "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/flavors"
)

func DataSourceComputeFlavor() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeFlavorRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.",
			},

			"flavor_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name", "min_ram", "min_disk"},
				Description:   "The ID of the flavor. Conflicts with the `name`, `min_ram` and `min_disk`",
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"flavor_id"},
				Description:   "The name of the flavor. Conflicts with the `flavor_id`.",
			},

			"min_ram": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"flavor_id"},
				Description:   "The minimum amount of RAM (in megabytes). Conflicts with the `flavor_id`.",
			},

			"ram": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The exact amount of RAM (in megabytes).",
			},

			"vcpus": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The amount of VCPUs.",
			},

			"min_disk": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"flavor_id"},
				Description:   "The minimum amount of disk (in gigabytes). Conflicts with the `flavor_id`.",
			},

			"disk": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The exact amount of disk (in gigabytes).",
			},

			"swap": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The amount of swap (in gigabytes).",
			},

			"rx_tx_factor": {
				Type:        schema.TypeFloat,
				Optional:    true,
				ForceNew:    true,
				Description: "The `rx_tx_factor` of the flavor.",
			},

			"is_public": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Description: "The flavor visibility.",
			},

			// Computed values
			"extra_specs": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Key/Value pairs of metadata for the flavor.",
			},

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the found flavor.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS flavor.",
	}
}

// dataSourceComputeFlavorRead performs the flavor lookup.
func dataSourceComputeFlavorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	var allFlavors []flavors.Flavor
	if v := d.Get("flavor_id").(string); v != "" {
		flavor, err := iflavors.Get(computeClient, v).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return diag.Errorf("No Flavor found")
			}
			return diag.Errorf("Unable to retrieve VKCS %s flavor: %s", v, err)
		}

		allFlavors = append(allFlavors, *flavor)
	} else {
		accessType := flavors.AllAccess
		if v, ok := d.GetOk("is_public"); ok {
			if v, ok := v.(bool); ok {
				if v {
					accessType = flavors.PublicAccess
				} else {
					accessType = flavors.PrivateAccess
				}
			}
		}
		listOpts := flavors.ListOpts{
			MinDisk:    d.Get("min_disk").(int),
			MinRAM:     d.Get("min_ram").(int),
			AccessType: accessType,
		}

		log.Printf("[DEBUG] vkcs_compute_flavor ListOpts: %#v", listOpts)

		allPages, err := flavors.ListDetail(computeClient, listOpts).AllPages()
		if err != nil {
			return diag.Errorf("Unable to query VKCS flavors: %s", err)
		}

		allFlavors, err = flavors.ExtractFlavors(allPages)
		if err != nil {
			return diag.Errorf("Unable to retrieve VKCS flavors: %s", err)
		}
	}

	// Loop through all flavors to find a more specific one.
	if len(allFlavors) > 0 {
		var filteredFlavors []flavors.Flavor
		for _, flavor := range allFlavors {
			if v := d.Get("name").(string); v != "" {
				if flavor.Name != v {
					continue
				}
			}

			// d.GetOk is used because 0 might be a valid choice.
			if v, ok := d.GetOk("ram"); ok {
				if flavor.RAM != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("vcpus"); ok {
				if flavor.VCPUs != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("disk"); ok {
				if flavor.Disk != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("swap"); ok {
				if flavor.Swap != v.(int) {
					continue
				}
			}

			if v, ok := d.GetOk("rx_tx_factor"); ok {
				if flavor.RxTxFactor != v.(float64) {
					continue
				}
			}

			filteredFlavors = append(filteredFlavors, flavor)
		}

		allFlavors = filteredFlavors
	}

	if len(allFlavors) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allFlavors) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allFlavors)
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	return diag.FromErr(dataSourceComputeFlavorAttributes(d, computeClient, &allFlavors[0]))
}

// dataSourceComputeFlavorAttributes populates the fields of a Flavor resource.
func dataSourceComputeFlavorAttributes(d *schema.ResourceData, computeClient *gophercloud.ServiceClient, flavor *flavors.Flavor) error {
	log.Printf("[DEBUG] Retrieved vkcs_compute_flavor %s: %#v", flavor.ID, flavor)

	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("flavor_id", flavor.ID)
	d.Set("disk", flavor.Disk)
	d.Set("ram", flavor.RAM)
	d.Set("rx_tx_factor", flavor.RxTxFactor)
	d.Set("swap", flavor.Swap)
	d.Set("vcpus", flavor.VCPUs)
	d.Set("is_public", flavor.IsPublic)

	es, err := iflavors.ListExtraSpecs(computeClient, d.Id()).Extract()
	if err != nil {
		return err
	}

	if err := d.Set("extra_specs", es); err != nil {
		log.Printf("[WARN] Unable to set extra_specs for vkcs_compute_flavor %s: %s", d.Id(), err)
	}

	return nil
}
