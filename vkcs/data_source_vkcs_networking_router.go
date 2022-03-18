package vkcs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func dataSourceNetworkingRouter() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingRouterRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_snat": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"sdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSDN(),
			},
		},
	}
}

func dataSourceNetworkingRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := routers.ListOpts{}

	if v, ok := d.GetOk("router_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOkExists("admin_state_up"); ok {
		asu := v.(bool)
		listOpts.AdminStateUp = &asu
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	tags := networkingAttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	pages, err := routers.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list Routers: %s", err)
	}

	allRouters, err := routers.ExtractRouters(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Routers: %s", err)
	}

	if len(allRouters) < 1 {
		return diag.Errorf("No Router found")
	}

	if len(allRouters) > 1 {
		return diag.Errorf("More than one Router found")
	}

	router := allRouters[0]

	log.Printf("[DEBUG] Retrieved Router %s: %+v", router.ID, router)
	d.SetId(router.ID)

	d.Set("name", router.Name)
	d.Set("description", router.Description)
	d.Set("admin_state_up", router.AdminStateUp)
	d.Set("status", router.Status)
	d.Set("tenant_id", router.TenantID)
	d.Set("external_network_id", router.GatewayInfo.NetworkID)
	d.Set("enable_snat", router.GatewayInfo.EnableSNAT)
	d.Set("all_tags", router.Tags)
	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	externalFixedIPs := make([]map[string]string, 0, len(router.GatewayInfo.ExternalFixedIPs))
	for _, v := range router.GatewayInfo.ExternalFixedIPs {
		externalFixedIPs = append(externalFixedIPs, map[string]string{
			"subnet_id":  v.SubnetID,
			"ip_address": v.IPAddress,
		})
	}
	return nil
}
