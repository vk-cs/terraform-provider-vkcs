package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func dataSourceBlockStorageVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBlockStorageVolumeRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			// Computed values
			"bootable": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"source_volume_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	client, err := config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	listOpts := volumes.ListOpts{
		Metadata: expandToMapStringString(d.Get("metadata").(map[string]interface{})),
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
	}

	allPages, err := volumes.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query vkcs_blockstorage_volume: %s", err)
	}

	allVolumes, err := volumes.ExtractVolumes(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_blockstorage_volume: %s", err)
	}

	if len(allVolumes) > 1 {
		return diag.Errorf("Your vkcs_blockstorage_volume query returned multiple results")
	}

	if len(allVolumes) < 1 {
		return diag.Errorf("Your vkcs_blockstorage_volume query returned no results")
	}

	dataSourceBlockStorageVolumeAttributes(d, allVolumes[0])

	return nil
}

func dataSourceBlockStorageVolumeAttributes(d *schema.ResourceData, volume volumes.Volume) {
	d.SetId(volume.ID)
	d.Set("name", volume.Name)
	d.Set("status", volume.Status)
	d.Set("bootable", volume.Bootable)
	d.Set("volume_type", volume.VolumeType)
	d.Set("size", volume.Size)
	d.Set("source_volume_id", volume.SourceVolID)
	d.Set("availability_zone", volume.AvailabilityZone)

	if err := d.Set("metadata", volume.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for vkcs_blockstorage_volume %s: %s", volume.ID, err)
	}
}
