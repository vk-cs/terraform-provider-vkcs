package firewall

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	igroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/firewall/v2/groups"
)

func DataSourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSecGroupRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Network client. A Network client is needed to retrieve security groups ids. If omitted, the `region` argument of the provider is used.",
			},
			"secgroup_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"id"},
				Deprecated:    "This argument is deprecated, please, use the `id` attribute instead.",
				Description:   "The ID of the security group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the security group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Human-readable description the the subnet.",
			},
			"tenant_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The owner of the security group.",
			},
			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of security group tags to filter.",
			},
			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The set of string tags applied on the security group.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: networking.ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
			},

			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				Description: "The ID of the security group.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS security group.",
	}
}

func dataSourceNetworkingSecGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	listOpts := groups.ListOpts{
		ID:          util.GetFirstNotEmpty(d.Get("id").(string), d.Get("secgroup_id").(string)),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
	}

	tags := networking.NetworkingAttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	pages, err := groups.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.FromErr(err)
	}

	var allSecGroups []securityGroupExtended
	err = igroups.ExtractSecurityGroupsInto(pages, &allSecGroups)
	if err != nil {
		return diag.Errorf("Unable to retrieve security groups: %s", err)
	}

	if len(allSecGroups) < 1 {
		return diag.Errorf("No Security Group found with name: %s", d.Get("name"))
	}

	if len(allSecGroups) > 1 {
		return diag.Errorf("More than one Security Group found with name: %s", d.Get("name"))
	}

	secGroup := allSecGroups[0]

	log.Printf("[DEBUG] Retrieved Security Group %s: %+v", secGroup.ID, secGroup)
	d.SetId(secGroup.ID)
	d.Set("secgroup_id", secGroup.ID)

	d.Set("name", secGroup.Name)
	d.Set("description", secGroup.Description)
	d.Set("tenant_id", secGroup.TenantID)
	d.Set("all_tags", secGroup.Tags)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", secGroup.SDN)

	return nil
}
