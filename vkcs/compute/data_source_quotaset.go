package compute

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/quotasets"
)

func DataSourceComputeQuotaset() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeQuotasetRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The id of the project to retrieve the quotaset.",
			},

			"injected_file_content_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed bytes of content for each injected file.",
			},

			"injected_file_path_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed bytes for each injected file path.",
			},

			"injected_files": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed injected files.",
			},

			"key_pairs": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed key pairs for each user.",
			},

			"metadata_items": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed metadata items for each server.",
			},

			"ram": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The amount of allowed server RAM, in MiB.",
			},

			"cores": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed server cores.",
			},

			"instances": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed servers.",
			},

			"server_groups": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed server groups.",
			},

			"server_group_members": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of allowed members for each server group.",
			},
		},
		Description: "Use this data source to get the compute quotaset of an VKCS project.",
	}
}

func dataSourceComputeQuotasetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	region := util.GetRegion(d, config)
	computeClient, err := config.ComputeV2Client(region)
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	projectID := d.Get("project_id").(string)

	q, err := quotasets.Get(computeClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_compute_quotaset"))
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
