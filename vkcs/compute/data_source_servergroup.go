package compute

import (
	"context"
	"sort"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iservergroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/servergroups"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourceComputeServerGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeServerGroupRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the server group.",
			},

			"policies": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The set of policies for the server group.",
			},

			"members": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The instances that are part of this server group.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The metadata of the server group.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS server group.",
	}
}

func dataSourceComputeServerGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	allPages, err := iservergroups.List(computeClient, servergroups.ListOpts{}).AllPages()
	if err != nil {
		return diag.Errorf("Error listing VKCS server groups: %s", err)
	}

	allServerGroups, err := servergroups.ExtractServerGroups(allPages)
	if err != nil {
		return diag.Errorf("Error extracting VKCS server groups: %s", err)
	}

	name := d.Get("name").(string)

	var matched []servergroups.ServerGroup
	for _, sg := range allServerGroups {
		if sg.Name == name {
			matched = append(matched, sg)
		}
	}

	if len(matched) < 1 {
		return diag.Errorf("Could not find any server group with this name: %s", name)
	}

	if len(matched) > 1 {
		return diag.Errorf("More than one server group found with this name: %s", name)
	}

	sg := matched[0]

	members := sg.Members
	sort.Strings(members)

	d.SetId(sg.ID)
	d.Set("name", sg.Name)
	d.Set("policies", sg.Policies)
	d.Set("members", members)
	d.Set("metadata", sg.Metadata)
	d.Set("region", util.GetRegion(d, config))

	return nil
}
