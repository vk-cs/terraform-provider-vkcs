package networking

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	isdn "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
)

func DataSourceNetworkingSDN() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSDNRead,

		Schema: map[string]*schema.Schema{
			"sdn": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names of available VKCS SDNs in the current project.",
			},
		},
		Description: "Use this data source to get a list of available VKCS SDNs in the current project. The first SDN is default. You do not have to specify default sdn argument in resources and datasources. You may specify non default SDN only for root resources such as `vkcs_networking_router`, `vkcs_networking_network`, `vkcs_networking_secgroup` (they do not depend on any other resource/datasource with sdn argument).",
	}
}

func dataSourceNetworkingSDNRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(), isdn.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	sdn, err := isdn.GetAvailableSDNs(networkingClient)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	d.Set("sdn", sdn)

	return nil
}
