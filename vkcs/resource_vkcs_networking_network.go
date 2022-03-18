package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
)

func resourceNetworkingNetwork() *schema.Resource {
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
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
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

			"port_security_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"private_dns_domain": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func resourceNetworkingNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := NetworkCreateOpts{
		networks.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		},
		MapValueSpecs(d),
		d.Get("private_dns_domain").(string),
	}

	if v, ok := d.GetOkExists("admin_state_up"); ok {
		asu := v.(bool)
		createOpts.AdminStateUp = &asu
	}

	// Declare a finalCreateOpts interface.
	var finalCreateOpts networks.CreateOptsBuilder
	finalCreateOpts = createOpts

	// Add the port security attribute if specified.
	if v, ok := d.GetOkExists("port_security_enabled"); ok {
		portSecurityEnabled := v.(bool)
		finalCreateOpts = portsecurity.NetworkCreateOptsExt{
			CreateOptsBuilder:   finalCreateOpts,
			PortSecurityEnabled: &portSecurityEnabled,
		}
	}

	log.Printf("[DEBUG] vkcs_networking_network create options: %#v", finalCreateOpts)
	n, err := networks.Create(networkingClient, finalCreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_network: %s", err)
	}

	log.Printf("[DEBUG] Waiting for vkcs_networking_network %s to become available.", n.ID)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"BUILD"},
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    resourceNetworkingNetworkStateRefreshFunc(networkingClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_network %s to become available: %s", n.ID, err)
	}

	d.SetId(n.ID)

	tags := networkingAttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "networks", n.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_network %s: %s", n.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_network %s", tags, n.ID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_network %s: %#v", n.ID, n)
	return resourceNetworkingNetworkRead(ctx, d, meta)
}

func resourceNetworkingNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var network networkExtended

	err = networks.Get(networkingClient, d.Id()).ExtractInto(&network)
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_network"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_network %s: %#v", d.Id(), network)

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("admin_state_up", network.AdminStateUp)
	d.Set("port_security_enabled", network.PortSecurityEnabled)
	d.Set("region", getRegion(d, config))
	d.Set("private_dns_domain", network.PrivateDNSDomain)
	d.Set("sdn", getSDN(d))

	networkingReadAttributesTags(d, network.Tags)

	return nil
}

func resourceNetworkingNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Declare finalUpdateOpts interface and basic updateOpts structure.
	var (
		finalUpdateOpts networks.UpdateOptsBuilder
		updateOpts      networks.UpdateOpts
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
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "networks", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_network %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_network %s", tags, d.Id())
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
	_, err = networks.Update(networkingClient, d.Id(), finalUpdateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_networking_network %s: %s", d.Id(), err)
	}

	return resourceNetworkingNetworkRead(ctx, d, meta)
}

func resourceNetworkingNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := networks.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_networking_network"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingNetworkStateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_network %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
