package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

const (
	// export_location_path filter appeared in 2.35.
	minManilaShareListExportLocationPath = "2.35"
)

func dataSourceSharedFilesystemShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedFilesystemShareRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"share_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"export_location_path": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"share_proto": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSharedFilesystemShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem sfsClient: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	listOpts := shares.ListOpts{
		Name:               d.Get("name").(string),
		DisplayDescription: d.Get("description").(string),
		ProjectID:          d.Get("project_id").(string),
		SnapshotID:         d.Get("snapshot_id").(string),
		ShareNetworkID:     d.Get("share_network_id").(string),
		Status:             d.Get("status").(string),
	}

	if v, ok := d.GetOkExists("export_location_path"); ok {
		listOpts.ExportLocationPath = v.(string)
		sfsClient.Microversion = minManilaShareListExportLocationPath
	}

	allPages, err := shares.ListDetail(sfsClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query shares: %s", err)
	}

	allShares, err := shares.ExtractShares(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve shares: %s", err)
	}

	if len(allShares) < 1 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again")
	}

	var share shares.Share
	if len(allShares) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allShares)
		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}
	share = allShares[0]

	d.SetId(share.ID)
	d.Set("name", share.Name)
	d.Set("region", getRegion(d, config))
	d.Set("project_id", share.ProjectID)
	d.Set("snapshot_id", share.SnapshotID)
	d.Set("share_network_id", share.ShareNetworkID)
	d.Set("availability_zone", share.AvailabilityZone)
	d.Set("description", share.Description)
	d.Set("size", share.Size)
	d.Set("status", share.Status)
	d.Set("share_proto", share.ShareProto)

	return nil
}
