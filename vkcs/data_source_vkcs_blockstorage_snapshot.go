package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
)

func dataSourceBlockStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBlockStorageSnapshotRead,

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
				Description: "The name of the snapshot.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The status of the snapshot.",
			},

			"volume_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the snapshot's volume.",
			},

			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Pick the most recently created snapshot if there are multiple results.",
			},

			// Computed values
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The snapshot's description.",
			},

			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the snapshot.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The snapshot's metadata.",
			},
		},
		Description: "Use this data source to get information about an existing snapshot.",
	}
}

func dataSourceBlockStorageSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	listOpts := snapshots.ListOpts{
		Name:     d.Get("name").(string),
		Status:   d.Get("status").(string),
		VolumeID: d.Get("volume_id").(string),
	}

	allPages, err := snapshots.List(client, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query vkcs_blockstorage_snapshots: %s", err)
	}

	allSnapshots, err := snapshots.ExtractSnapshots(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve vkcs_blockstorage_snapshots: %s", err)
	}

	if len(allSnapshots) < 1 {
		return diag.Errorf("Your vkcs_blockstorage_snapshot query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var snapshot snapshots.Snapshot
	if len(allSnapshots) > 1 {
		recent := d.Get("most_recent").(bool)

		if recent {
			snapshot = dataSourceBlockStorageMostRecentSnapshot(allSnapshots)
		} else {
			log.Printf("[DEBUG] Multiple vkcs_blockstorage_snapshot results found: %#v", allSnapshots)

			return diag.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		snapshot = allSnapshots[0]
	}

	dataSourceBlockStorageSnapshotAttributes(d, snapshot)

	return nil
}

func dataSourceBlockStorageSnapshotAttributes(d *schema.ResourceData, snapshot snapshots.Snapshot) {
	d.SetId(snapshot.ID)
	d.Set("name", snapshot.Name)
	d.Set("description", snapshot.Description)
	d.Set("size", snapshot.Size)
	d.Set("status", snapshot.Status)
	d.Set("volume_id", snapshot.VolumeID)

	if err := d.Set("metadata", snapshot.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for snapshot %s: %s", snapshot.ID, err)
	}
}
