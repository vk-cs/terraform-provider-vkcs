package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

func dataSourceSharedFilesystemShareNetwork() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedFilesystemShareNetworkRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Shared File System client. A Shared File System client is needed to read a share network. If omitted, the `region` argument of the provider is used.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the share network.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the share network.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The human-readable description of the share network.",
			},

			"neutron_net_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The neutron network UUID of the share network.",
			},

			"neutron_subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The neutron subnet UUID of the share network.",
			},

			"security_service_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The security service IDs associated with the share network.",
			},

			"security_service_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The list of security service IDs associated with the share network.",
			},

			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The share network CIDR.",
			},
		},
		Description: "Use this data source to get the ID of an available Shared File System share network.",
	}
}

func dataSourceSharedFilesystemShareNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem sfsClient: %s", err)
	}

	listOpts := sharenetworks.ListOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		ProjectID:       d.Get("project_id").(string),
		NeutronNetID:    d.Get("neutron_net_id").(string),
		NeutronSubnetID: d.Get("neutron_subnet_id").(string),
	}

	listOpts.IPVersion = gophercloud.IPVersion(4)

	allPages, err := sharenetworks.ListDetail(sfsClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query share networks: %s", err)
	}

	allShareNetworks, err := sharenetworks.ExtractShareNetworks(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve share networks: %s", err)
	}

	if len(allShareNetworks) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var securityServiceID string
	var securityServiceIDs []string
	if v, ok := d.GetOk("security_service_id"); ok {
		// filtering by "security_service_id"
		securityServiceID = v.(string)
		var filteredShareNetworks []sharenetworks.ShareNetwork

		log.Printf("[DEBUG] Filtering share networks by a %s security service ID", securityServiceID)
		for _, shareNetwork := range allShareNetworks {
			tmp, err := resourceSharedFilesystemShareNetworkGetSvcByShareNetID(sfsClient, shareNetwork.ID)
			if err != nil {
				return diag.FromErr(err)
			}
			if strSliceContains(tmp, securityServiceID) {
				filteredShareNetworks = append(filteredShareNetworks, shareNetwork)
				securityServiceIDs = tmp
			}
		}

		if len(filteredShareNetworks) == 0 {
			return diag.Errorf("Your query returned no results after the security service ID filter. " +
				"Please change your search criteria and try again")
		}
		allShareNetworks = filteredShareNetworks
	}

	var shareNetwork sharenetworks.ShareNetwork
	if len(allShareNetworks) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allShareNetworks)
		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria")
	}
	shareNetwork = allShareNetworks[0]

	// skip extra calls if "security_service_id" filter was already used
	if securityServiceID == "" {
		securityServiceIDs, err = resourceSharedFilesystemShareNetworkGetSvcByShareNetID(sfsClient, shareNetwork.ID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(shareNetwork.ID)
	d.Set("name", shareNetwork.Name)
	d.Set("description", shareNetwork.Description)
	d.Set("project_id", shareNetwork.ProjectID)
	d.Set("neutron_net_id", shareNetwork.NeutronNetID)
	d.Set("neutron_subnet_id", shareNetwork.NeutronSubnetID)
	d.Set("security_service_ids", securityServiceIDs)
	d.Set("cidr", shareNetwork.CIDR)
	d.Set("region", getRegion(d, config))

	return nil
}
