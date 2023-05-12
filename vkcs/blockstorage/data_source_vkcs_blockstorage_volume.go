package blockstorage

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func DataSourceBlockStorageVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBlockStorageVolumeRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Block Storage client. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the volume.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The status of the volume.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Computed:    true,
				Description: "Metadata key/value pairs associated with the volume.",
			},

			// Computed values
			"bootable": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Indicates if the volume is bootable.",
			},

			"volume_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The type of the volume.",
			},

			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the volume in GBs.",
			},

			"source_volume_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the volume from which the current volume was created.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the availability zone of the volume.",
			},
		},
		Description: "Use this data source to get information about an existing volume.",
	}
}

func dataSourceBlockStorageVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	listOpts := volumes.ListOpts{
		Metadata: util.ExpandToMapStringString(d.Get("metadata").(map[string]interface{})),
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
