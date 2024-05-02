package compute

import (
	"cmp"
	"context"
	"encoding/json"
	"log"
	"slices"

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

			"cpu_generation": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The cpu generation of the flavor.",
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

type RequiredFlavor struct {
	// Disk is the amount of root disk, measured in GB.
	Disk    int  `json:"disk"`
	HasDisk bool `json:"has_disk"`

	// Disk is the amount of root disk, measured in GB.
	MinDisk    int  `json:"min_disk"`
	HasMinDisk bool `json:"has_min_disk"`

	// RAM is the amount of memory, measured in MB.
	RAM    int  `json:"ram"`
	HasRam bool `json:"has_ram"`

	// MinRAM is the amount of memory, measured in MB.
	MinRAM    int  `json:"min_ram"`
	HasMinRam bool `json:"has_min_ram"`

	// Name is the name of the flavor.
	Name    string `json:"name"`
	HasName bool   `json:"has_name"`

	// RxTxFactor describes bandwidth alterations of the flavor.
	RxTxFactor    float64 `json:"rxtx_factor"`
	HasRxTxFactor bool    `json:"has_rxtx_factor"`

	// Swap is the amount of swap space, measured in MB.
	Swap    int  `json:"-"`
	HasSwap bool `json:"has_swap"`

	// VCPUs indicates how many (virtual) CPUs are available for this flavor.
	VCPUs    int  `json:"vcpus"`
	HasVCPUs bool `json:"has_vcpus"`

	// CPUGeneration is the cpu generation of flavor
	CPUGeneration    string `json:"cpu_generation"`
	HasCPUGeneration bool   `json:"has_cpu_generation"`

	AccessType flavors.AccessType `json:"access_type"`
}

func NewRequiredFlavorFromResourceData(d *schema.ResourceData) *RequiredFlavor {
	name, hasName := d.GetOk("name")
	ram, hasRam := d.GetOk("ram")
	VCPUs, hasVCPUs := d.GetOk("vcpus")
	disk, hasDisk := d.GetOk("disk")
	minDisk, hasMinDisk := d.GetOk("min_disk")
	minRAM, hasMinRam := d.GetOk("min_ram")
	rxTxFactor, hasRxTxFactor := d.GetOk("rx_tx_factor")
	swap, hasSwap := d.GetOk("swap")
	cpuGeneration, hasCPUGeneration := d.GetOk("cpu_generation")

	if minRAM.(int) > ram.(int) {
		ram = minRAM
	}
	if minDisk.(int) > disk.(int) {
		disk = minDisk
	}

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

	return &RequiredFlavor{
		Disk:             disk.(int),
		HasDisk:          hasDisk,
		MinDisk:          minDisk.(int),
		HasMinDisk:       hasMinDisk,
		RAM:              ram.(int),
		HasRam:           hasRam,
		MinRAM:           minRAM.(int),
		HasMinRam:        hasMinRam,
		Name:             name.(string),
		HasName:          hasName,
		RxTxFactor:       rxTxFactor.(float64),
		HasRxTxFactor:    hasRxTxFactor,
		Swap:             swap.(int),
		HasSwap:          hasSwap,
		VCPUs:            VCPUs.(int),
		HasVCPUs:         hasVCPUs,
		CPUGeneration:    cpuGeneration.(string),
		HasCPUGeneration: hasCPUGeneration,
		AccessType:       accessType,
	}
}

// FlavorWithExtraInfo needs for extract ExtraSpecs info from flavors.FlavorPage
type FlavorWithExtraInfo struct {
	// ID is the flavor's unique ID.
	ID string `json:"id"`

	// Disk is the amount of root disk, measured in GB.
	Disk int `json:"disk"`

	// RAM is the amount of memory, measured in MB.
	RAM int `json:"ram"`

	// Name is the name of the flavor.
	Name string `json:"name"`

	// RxTxFactor describes bandwidth alterations of the flavor.
	RxTxFactor float64 `json:"rxtx_factor"`

	// Swap is the amount of swap space, measured in MB.
	Swap int `json:"-"`

	// VCPUs indicates how many (virtual) CPUs are available for this flavor.
	VCPUs int `json:"vcpus"`

	// IsPublic indicates whether the flavor is public.
	IsPublic bool `json:"os-flavor-access:is_public"`

	// Ephemeral is the amount of ephemeral disk space, measured in GB.
	Ephemeral int `json:"OS-FLV-EXT-DATA:ephemeral"`

	ExtraSpecs map[string]interface{} `json:"extra_specs"`
}

func (f *FlavorWithExtraInfo) toFlavor() *flavors.Flavor {
	return &flavors.Flavor{
		ID:         f.ID,
		Disk:       f.Disk,
		RAM:        f.RAM,
		Name:       f.Name,
		RxTxFactor: f.RxTxFactor,
		Swap:       f.Swap,
		VCPUs:      f.VCPUs,
		IsPublic:   f.IsPublic,
		Ephemeral:  f.Ephemeral,
	}
}

// dataSourceComputeFlavorRead performs the flavor lookup.
func dataSourceComputeFlavorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	// choose only one by flavor_id
	if v := d.Get("flavor_id").(string); v != "" {
		flavor, err := iflavors.Get(computeClient, v).Extract()
		if err != nil {
			if errutil.IsNotFound(err) {
				return diag.Errorf("No Flavor found")
			}
			return diag.Errorf("Unable to retrieve VKCS %s flavor: %s", v, err)
		}

		return diag.FromErr(dataSourceComputeFlavorAttributes(d, computeClient, flavor))
	}

	requiredFlavor := NewRequiredFlavorFromResourceData(d)
	listOpts := flavors.ListOpts{
		MinDisk:    requiredFlavor.MinDisk,
		MinRAM:     requiredFlavor.MinRAM,
		AccessType: requiredFlavor.AccessType,
	}

	log.Printf("[DEBUG] vkcs_compute_flavor ListOpts: %#v", listOpts)

	allPages, err := flavors.ListDetail(computeClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query VKCS flavors: %s", err)
	}

	var s struct {
		FlavorsExtra []FlavorWithExtraInfo `json:"flavors"`
	}
	err = (allPages.(flavors.FlavorPage)).ExtractInto(&s)

	if err != nil {
		return diag.Errorf("Unable to retrieve VKCS flavors: %s", err)
	}
	allFlavors := s.FlavorsExtra

	// Loop through all flavors to find a more specific one.
	if len(allFlavors) > 0 {
		var filteredFlavors []FlavorWithExtraInfo
		for _, flavor := range allFlavors {
			switch {
			case requiredFlavor.HasName && flavor.Name != requiredFlavor.Name:
				continue
			case requiredFlavor.HasRam && flavor.RAM != requiredFlavor.RAM:
				continue
			case requiredFlavor.HasVCPUs && flavor.VCPUs != requiredFlavor.VCPUs:
				continue
			case requiredFlavor.HasDisk && flavor.Disk != requiredFlavor.Disk:
				continue
			case requiredFlavor.HasSwap && flavor.Swap != requiredFlavor.Swap:
				continue
			case requiredFlavor.HasRxTxFactor && flavor.RxTxFactor != requiredFlavor.RxTxFactor:
				continue
			case requiredFlavor.HasCPUGeneration && flavor.ExtraSpecs == nil:
				continue
			}
			if !requiredFlavor.HasCPUGeneration {
				filteredFlavors = append(filteredFlavors, flavor)
				continue
			}

			if flavorCPU, ok := flavor.ExtraSpecs["mcs:cpu_generation"]; ok {
				if requiredFlavor.CPUGeneration == flavorCPU {
					filteredFlavors = append(filteredFlavors, flavor)
				}
			}
		}

		allFlavors = filteredFlavors
	}

	if len(allFlavors) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	// if we find many flavors and the user sets the min_ram or min_disk values
	// we give him the flavor with the minimum amount of RAM from the found flavors
	if len(allFlavors) > 1 && (requiredFlavor.HasMinRam || requiredFlavor.HasMinDisk) {
		minFlavor := slices.MinFunc(allFlavors, func(a, b FlavorWithExtraInfo) int {
			if a.RAM != b.RAM {
				return cmp.Compare(a.RAM, b.RAM)
			}
			return cmp.Compare(a.Disk, b.Disk)
		})

		return diag.FromErr(dataSourceComputeFlavorAttributes(d, computeClient, minFlavor.toFlavor()))
	}

	if len(allFlavors) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allFlavors)
		if len(allFlavors) > 2 {
			return diag.Errorf("Found %d avalaible flavors. Please try a more specific search criteria", len(allFlavors))
		}
		bytes, _ := json.MarshalIndent(allFlavors, "", "\t")

		return diag.Errorf("Found %d avalaible flavors. Avalaible flavors:\n%s\n"+
			"Please try a more specific search criteria", len(allFlavors), bytes)
	}

	return diag.FromErr(dataSourceComputeFlavorAttributes(d, computeClient, allFlavors[0].toFlavor()))
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
