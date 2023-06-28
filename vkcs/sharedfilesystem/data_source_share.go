package sharedfilesystem

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

const (
	// export_location_path filter appeared in 2.35.
	minManilaShareListExportLocationPath = "2.35"
)

func DataSourceSharedFilesystemShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedFilesystemShareRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Shared File System client.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the share.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the share.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The human-readable description for the share.",
			},

			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the share's base snapshot.",
			},

			"share_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the share's share network.",
			},

			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A share status filter. A valid value is `creating`, `error`, `available`, `deleting`, `error_deleting`, `manage_starting`, `manage_error`, `unmanage_starting`, `unmanage_error`, `unmanaged`, `extending`, `extending_error`, `shrinking`, `shrinking_error`, or `shrinking_possible_data_loss_error`.",
			},

			"export_location_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The export location path of the share.",
			},

			"share_proto": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The share protocol.",
			},

			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The share size, in GBs.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The share availability zone.",
			},
		},
		Description: "Use this data source to get the ID of an available Shared File System share.",
	}
}

func dataSourceSharedFilesystemShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem sfsClient: %s", err)
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

	if v, ok := d.GetOk("export_location_path"); ok {
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

	exportLocationPath, err := getShareExportLocationPath(sfsClient, share.ID)
	if err != nil {
		return diag.Errorf("Unable to retrieve export location path: %s", err)
	}

	d.SetId(share.ID)
	d.Set("name", share.Name)
	d.Set("region", util.GetRegion(d, config))
	d.Set("project_id", share.ProjectID)
	d.Set("snapshot_id", share.SnapshotID)
	d.Set("share_network_id", share.ShareNetworkID)
	d.Set("availability_zone", share.AvailabilityZone)
	d.Set("description", share.Description)
	d.Set("size", share.Size)
	d.Set("status", share.Status)
	d.Set("export_location_path", exportLocationPath)
	d.Set("share_proto", share.ShareProto)

	return nil
}
