package networking

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	"net/http"
	"strings"
)

func DataSourceNetworkingSDN() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSDNRead,

		Schema: map[string]*schema.Schema{
			"sdn": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Names of available VKCS SDN's in the current project.",
			},
		},
		Description: "Use this data source to get list of an available VKCS SDN in the current project.",
	}
}

func dataSourceNetworkingSDNRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(config.GetRegion(), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var sdn []string
	httpResp, err := networkingClient.Get(networkingClient.ServiceURL("available-sdn"), &sdn, nil)
	if err != nil {
		return diag.Errorf("Error getting avalible SDN's: %s", err)
	}
	if httpResp.StatusCode != http.StatusOK {
		return diag.Errorf("Error getting available SDN's: %s", httpResp.Status)
	}

	for i := 0; i < len(sdn); i++ {
		sdn[i] = strings.ToLower(sdn[i])
	}

	d.SetId(config.GetTenantID())
	d.Set("sdn", sdn)

	return nil
}
