package networking

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func ResourceNetworkingRouter() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterCreate,
		ReadContext:   resourceNetworkingRouterRead,
		UpdateContext: resourceNetworkingRouterUpdate,
		DeleteContext: resourceNetworkingRouterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The region in which to obtain the networking client. A networking client is needed to create a router. If omitted, the `region` argument of the provider is used. Changing this creates a new router.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "A unique name for the router. Changing this updates the `name` of an existing router.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Human-readable description for the router.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "Administrative up/down status for the router (must be \"true\" or \"false\" if provided). Changing this updates the `admin_state_up` of an existing router.",
			},

			"external_network_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Computed:    true,
				Description: "The network UUID of an external gateway for the router. A router with an external gateway is required if any compute instances or load balancers will be using floating IPs. Changing this updates the external gateway of the router.",
			},

			"value_specs": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of additional driver-specific options.",
			},

			"vendor_options": {
				Type:     schema.TypeSet,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"set_router_gateway_after_create": {
							Type:        schema.TypeBool,
							Default:     false,
							Optional:    true,
							Description: "Boolean to control whether the Router gateway is assigned during creation or updated after creation.",
						},
					},
				},
				Description: "Map of additional vendor-specific options. Supported options are described below.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags for the router.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of tags assigned on the router, which have been explicitly and implicitly added.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},
		},
		Description: "Manages a router resource within VKCS.",
	}
}

func resourceNetworkingRouterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	createOpts := RouterCreateOpts{
		CreateOpts: routers.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
		ValueSpecs: util.MapValueSpecs(d),
	}

	if asuRaw, ok := d.GetOk("admin_state_up"); ok {
		asu := asuRaw.(bool)
		createOpts.AdminStateUp = &asu
	}

	// Get Vendor_options
	vendorOptionsRaw := d.Get("vendor_options").(*schema.Set)
	var vendorUpdateGateway bool
	if vendorOptionsRaw.Len() > 0 {
		vendorOptions := util.ExpandVendorOptions(vendorOptionsRaw.List())
		vendorUpdateGateway = vendorOptions["set_router_gateway_after_create"].(bool)
	}

	// Gateway settings
	var externalNetworkID string
	var gatewayInfo routers.GatewayInfo

	if v := d.Get("external_network_id").(string); v != "" {
		externalNetworkID = v
		gatewayInfo.NetworkID = externalNetworkID
	}

	// vendorUpdateGateway is a flag for certain vendor-specific virtual routers
	// which do not allow gateway settings to be set during router creation.
	// If this flag was not enabled, then we can safely set the gateway
	// information during create.
	if !vendorUpdateGateway && externalNetworkID != "" {
		createOpts.GatewayInfo = &gatewayInfo
	}

	var r *routers.Router
	log.Printf("[DEBUG] vkcs_networking_router create options: %#v", createOpts)

	// router creation without a retry
	r, err = routers.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_router: %s", err)
	}

	d.SetId(r.ID)

	log.Printf("[DEBUG] Waiting for vkcs_networking_router %s to become available.", r.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    resourceNetworkingRouterStateRefreshFunc(networkingClient, r.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_router %s to become available: %s", r.ID, err)
	}

	// If the vendorUpdateGateway flag was specified and if an external network
	// was specified, then set the gateway information after router creation.
	if vendorUpdateGateway && externalNetworkID != "" {
		log.Printf("[DEBUG] Adding external_network %s to vkcs_networking_router %s", externalNetworkID, r.ID)

		var updateOpts routers.UpdateOpts
		updateOpts.GatewayInfo = &gatewayInfo

		log.Printf("[DEBUG] Assigning external_gateway to vkcs_networking_router %s with options: %#v", r.ID, updateOpts)
		_, err = routers.Update(networkingClient, r.ID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating vkcs_networking_router: %s", err)
		}
	}

	tags := NetworkingAttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "routers", r.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_router %s: %s", r.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_router %s", tags, r.ID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_router %s: %#v", r.ID, r)
	return resourceNetworkingRouterRead(ctx, d, meta)
}

func resourceNetworkingRouterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	r, err := routers.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving vkcs_networking_router: %s", err)
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_router %s: %#v", d.Id(), r)

	// Basic settings.
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("admin_state_up", r.AdminStateUp)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", GetSDN(d))

	NetworkingReadAttributesTags(d, r.Tags)

	// Gateway settings.
	d.Set("external_network_id", r.GatewayInfo.NetworkID)

	return nil
}

func resourceNetworkingRouterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	routerID := d.Id()
	mutex := config.GetMutex()
	mutex.Lock(routerID)
	defer mutex.Unlock(routerID)

	var hasChange bool
	var updateOpts routers.UpdateOpts
	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		hasChange = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Gateway settings.
	var updateGatewaySettings bool
	var externalNetworkID string
	gatewayInfo := routers.GatewayInfo{}

	if v := d.Get("external_network_id").(string); v != "" {
		externalNetworkID = v
	}

	if externalNetworkID != "" {
		gatewayInfo.NetworkID = externalNetworkID
	}

	if d.HasChange("external_network_id") {
		updateGatewaySettings = true
	}

	if updateGatewaySettings {
		hasChange = true
		updateOpts.GatewayInfo = &gatewayInfo
	}

	if hasChange {
		log.Printf("[DEBUG] vkcs_networking_router %s update options: %#v", d.Id(), updateOpts)
		_, err = routers.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating vkcs_networking_router: %s", err)
		}
	}

	// Next, perform any required updates to the tags.
	if d.HasChange("tags") {
		tags := NetworkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "routers", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_router %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_router %s", tags, d.Id())
	}

	return resourceNetworkingRouterRead(ctx, d, meta)
}

func resourceNetworkingRouterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	if err := routers.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_networking_router"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingRouterStateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting vkcs_networking_router: %s", err)
	}

	d.SetId("")
	return nil
}
