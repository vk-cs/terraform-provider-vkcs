package vkcs

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
)

func dataSourceComputeQuotaset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeQuotasetRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"injected_file_content_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"injected_file_path_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"injected_files": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"key_pairs": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"metadata_items": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"ram": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"cores": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"instances": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"server_groups": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"server_group_members": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceComputeQuotasetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	region := getRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	projectID := d.Get("project_id").(string)

	q, err := quotasets.Get(computeClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_compute_quotaset"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_quotaset %s: %#v", d.Id(), q)

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)
	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("injected_file_content_bytes", q.InjectedFileContentBytes)
	d.Set("injected_file_path_bytes", q.InjectedFilePathBytes)
	d.Set("injected_files", q.InjectedFiles)
	d.Set("key_pairs", q.KeyPairs)
	d.Set("metadata_items", q.MetadataItems)
	d.Set("ram", q.RAM)
	d.Set("cores", q.Cores)
	d.Set("instances", q.Instances)
	d.Set("server_groups", q.ServerGroups)
	d.Set("server_group_members", q.ServerGroupMembers)

	return nil
}
