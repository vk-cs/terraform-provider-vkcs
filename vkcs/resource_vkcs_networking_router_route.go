package vkcs

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func resourceNetworkingRouterRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterRouteCreate,
		ReadContext:   resourceNetworkingRouterRouteRead,
		DeleteContext: resourceNetworkingRouterRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"destination_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"sdn": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validateSDN(),
			},
		},
	}
}

func resourceNetworkingRouterRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	mutex := config.GetMutex()
	mutex.Lock(routerID)
	defer mutex.Unlock(routerID)

	r, err := routers.Get(networkingClient, routerID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_router"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_router %s: %#v", routerID, r)

	routes := r.Routes
	dstCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)
	exists := false

	for _, route := range routes {
		if route.DestinationCIDR == dstCIDR && route.NextHop == nextHop {
			exists = true
			break
		}
	}

	if exists {
		log.Printf("[DEBUG] vkcs_networking_router %s already has route to %s via %s", routerID, dstCIDR, nextHop)
		return resourceNetworkingRouterRouteRead(ctx, d, meta)
	}

	routes = append(routes, routers.Route{
		DestinationCIDR: dstCIDR,
		NextHop:         nextHop,
	})
	updateOpts := routers.UpdateOpts{
		Routes: &routes,
	}
	log.Printf("[DEBUG] vkcs_networking_router %s update options: %#v", routerID, updateOpts)
	_, err = routers.Update(networkingClient, routerID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_networking_router: %s", err)
	}

	d.SetId(resourceNetworkingRouterRouteBuildID(routerID, dstCIDR, nextHop))

	return resourceNetworkingRouterRouteRead(ctx, d, meta)
}

func resourceNetworkingRouterRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	idFromResource, dstCIDR, nextHop, err := resourceNetworkingRouterRouteParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading vkcs_networking_router_route ID %s: %s", d.Id(), err)
	}

	routerID := d.Get("router_id").(string)
	if routerID == "" {
		routerID = idFromResource
	}
	d.Set("router_id", routerID)

	r, err := routers.Get(networkingClient, routerID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_router"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_router %s: %#v", routerID, r)

	for _, route := range r.Routes {
		if route.DestinationCIDR == dstCIDR && route.NextHop == nextHop {
			d.Set("destination_cidr", dstCIDR)
			d.Set("next_hop", nextHop)
			break
		}
	}

	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}

func resourceNetworkingRouterRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	mutex := config.GetMutex()
	mutex.Lock(routerID)
	defer mutex.Unlock(routerID)

	r, err := routers.Get(networkingClient, routerID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_router"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_router %s: %#v", routerID, r)

	dstCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)

	oldRoutes := r.Routes
	newRoute := []routers.Route{}

	for _, route := range oldRoutes {
		if route.DestinationCIDR != dstCIDR || route.NextHop != nextHop {
			newRoute = append(newRoute, route)
		}
	}

	if len(oldRoutes) == len(newRoute) {
		return diag.Errorf("Can't find route to %s via %s on vkcs_networking_router %s", dstCIDR, nextHop, routerID)
	}

	log.Printf("[DEBUG] Deleting vkcs_networking_router %s route to %s via %s", routerID, dstCIDR, nextHop)
	updateOpts := routers.UpdateOpts{
		Routes: &newRoute,
	}
	_, err = routers.Update(networkingClient, routerID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_networking_router: %s", err)
	}

	return nil
}
