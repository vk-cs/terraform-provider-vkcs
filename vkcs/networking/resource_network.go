package networking

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
	inetworks "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/networks"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	iattributestags "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/attributestags"
)

func ResourceNetworkingNetwork() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingNetworkCreate,
		ReadContext:   resourceNetworkingNetworkRead,
		UpdateContext: resourceNetworkingNetworkUpdate,
		DeleteContext: resourceNetworkingNetworkDelete,
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
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a network. If omitted, the `region` argument of the provider is used. Changing this creates a new network.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The name of the network. Changing this updates the name of the existing network.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "Human-readable description of the network. Changing this updates the name of the existing network.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Default:     true,
				Description: "The administrative state of the network. Acceptable values are \"true\" and \"false\". Changing this value updates the state of the existing network.",
			},

			"value_specs": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "Map of additional options.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A set of string tags for the network.",
			},

			"all_tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The collection of tags assigned on the network, which have been explicitly and implicitly added.",
			},

			"port_security_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether to explicitly enable or disable port security on the network. Port Security is usually enabled by default, so omitting this argument will usually result in a value of \"true\". Setting this explicitly to `false` will disable port security. Valid values are `true` and `false`.",
			},

			"private_dns_domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Private dns domain name",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
			},

			"vkcs_services_access": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether VKCS services access is enabled. This feature should be enabled globally for your project. Access can be enabled for new or existing networks, but cannot be disabled for existing networks. Valid values are `true` and `false`.",
			},
		},
		Description: "Manages a network resource within VKCS.",
	}
}

func resourceNetworkingNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	createOpts := inetworks.NetworkCreateOpts{
		CreateOpts: networks.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
		ValueSpecs:       util.MapValueSpecs(d),
		PrivateDNSDomain: d.Get("private_dns_domain").(string),
		ServicesAccess:   d.Get("vkcs_services_access").(bool),
	}

	v := d.Get("admin_state_up")
	asu := v.(bool)
	createOpts.AdminStateUp = &asu

	// Declare a finalCreateOpts interface.
	var finalCreateOpts networks.CreateOptsBuilder
	finalCreateOpts = createOpts

	v = d.Get("port_security_enabled")
	pse := v.(bool)
	finalCreateOpts = portsecurity.NetworkCreateOptsExt{
		CreateOptsBuilder:   finalCreateOpts,
		PortSecurityEnabled: &pse,
	}

	log.Printf("[DEBUG] vkcs_networking_network create options: %#v", finalCreateOpts)
	n, err := inetworks.Create(networkingClient, finalCreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_network: %s", err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Waiting for vkcs_networking_network %s to become available.", n.ID)

	var createErrDetails error
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    resourceNetworkingNetworkStateRefreshFunc(networkingClient, n.ID, &createErrDetails),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if createErrDetails != nil {
			var timeoutErr *retry.TimeoutError
			if errors.As(err, &timeoutErr) {
				timeoutErr.LastError = createErrDetails
				return diag.Errorf("Error waiting for vkcs_networking_network %s to become available: %s", d.Id(), timeoutErr)
			}
		}

		return diag.Errorf("Error waiting for vkcs_networking_network %s to become available: %s", n.ID, err)
	}

	tags := NetworkingAttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := iattributestags.ReplaceAll(networkingClient, "networks", n.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_network %s: %s", n.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_network %s", tags, n.ID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_network %s: %#v", n.ID, n)
	return resourceNetworkingNetworkRead(ctx, d, meta)
}

func resourceNetworkingNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var network networkExtended

	err = inetworks.Get(networkingClient, d.Id()).ExtractInto(&network)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_networking_network"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_network %s: %#v", d.Id(), network)

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("admin_state_up", network.AdminStateUp)
	d.Set("port_security_enabled", network.PortSecurityEnabled)
	d.Set("region", util.GetRegion(d, config))
	d.Set("private_dns_domain", network.PrivateDNSDomain)
	d.Set("sdn", network.SDN)
	d.Set("vkcs_services_access", network.ServicesAccess)

	NetworkingReadAttributesTags(d, network.Tags)

	return nil
}

func resourceNetworkingNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	// Declare finalUpdateOpts interface and basic updateOpts structure.
	var (
		finalUpdateOpts networks.UpdateOptsBuilder
		updateOpts      inetworks.NetworkUpdateOpts
	)

	// Populate basic updateOpts.
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("admin_state_up") {
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	// Change tags if needed.
	if d.HasChange("tags") {
		tags := NetworkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := iattributestags.ReplaceAll(networkingClient, "networks", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_network %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_network %s", tags, d.Id())
	}

	if d.HasChange("vkcs_services_access") {
		servicesAccess := d.Get("vkcs_services_access").(bool)
		if !servicesAccess {
			return diag.Errorf("services_access cannot be disabled")
		}
		updateOpts.ServicesAccess = &servicesAccess
	}

	// Save basic updateOpts into finalUpdateOpts.
	finalUpdateOpts = updateOpts

	// Populate port security options.
	if d.HasChange("port_security_enabled") {
		portSecurityEnabled := d.Get("port_security_enabled").(bool)
		finalUpdateOpts = portsecurity.NetworkUpdateOptsExt{
			UpdateOptsBuilder:   finalUpdateOpts,
			PortSecurityEnabled: &portSecurityEnabled,
		}
	}

	log.Printf("[DEBUG] vkcs_networking_network %s update options: %#v", d.Id(), finalUpdateOpts)
	_, err = inetworks.Update(networkingClient, d.Id(), finalUpdateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_networking_network %s: %s", d.Id(), err)
	}

	return resourceNetworkingNetworkRead(ctx, d, meta)
}

func resourceNetworkingNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	if err := inetworks.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_networking_network"))
	}

	var deleteErrDetails error
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingNetworkStateRefreshFunc(networkingClient, d.Id(), &deleteErrDetails),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if deleteErrDetails != nil {
			var timeoutErr *retry.TimeoutError
			if errors.As(err, &timeoutErr) {
				timeoutErr.LastError = deleteErrDetails
				return diag.Errorf("Error waiting for vkcs_networking_network %s to become deleted: %s", d.Id(), timeoutErr)
			}
		}

		return diag.Errorf("Error waiting for vkcs_networking_network %s to become deleted:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
