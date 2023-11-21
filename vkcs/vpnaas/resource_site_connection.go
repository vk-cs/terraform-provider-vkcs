package vpnaas

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
)

func ResourceSiteConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSiteConnectionCreate,
		ReadContext:   resourceSiteConnectionRead,
		UpdateContext: resourceSiteConnectionUpdate,
		DeleteContext: resourceSiteConnectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create an IPSec site connection. If omitted, the `region` argument of the provider is used. Changing this creates a new site connection.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the connection. Changing this updates the name of the existing connection.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the connection. Changing this updates the description of the existing connection.",
			},
			"ikepolicy_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The ID of the IKE policy. Changing this creates a new connection.",
			},
			"peer_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The peer router identity for authentication. A valid value is an IPv4 address, IPv6 address, e-mail address, key ID, or FQDN. Typically, this value matches the peer_address value. Changing this updates the existing policy.",
			},
			"peer_address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The peer gateway public IPv4 or IPv6 address or FQDN.",
			},
			"peer_ep_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID for the endpoint group that contains private CIDRs in the form < net_address > / < prefix > for the peer side of the connection. You must specify this parameter with the local_ep_group_id parameter unless in backward-compatible mode where peer_cidrs is provided with a subnet_id for the VPN service.",
			},
			"local_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An ID to be used instead of the external IP address for a virtual router used in traffic between instances on different networks in east-west traffic. Most often, local ID would be domain name, email address, etc. If this is not configured then the external IP address will be used as the ID.",
			},
			"vpnservice_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The ID of the VPN service. Changing this creates a new connection.",
			},
			"local_ep_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The ID for the endpoint group that contains private subnets for the local side of the connection. You must specify this parameter with the peer_ep_group_id parameter unless in backward- compatible mode where peer_cidrs is provided with a subnet_id for the VPN service. Changing this updates the existing connection.",
			},
			"ipsecpolicy_id": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The ID of the IPsec policy. Changing this creates a new connection.",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The administrative state of the resource. Can either be up(true) or down(false). Changing this updates the administrative state of the existing connection.",
			},
			"psk": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The pre-shared key. A valid value is any string.",
			},
			"initiator": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "A valid value is response-only or bi-directional. Default is bi-directional.",
			},
			"mtu": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The maximum transmission unit (MTU) value to address fragmentation. Minimum value is 68 for IPv4, and 1280 for IPv6.",
			},
			"peer_cidrs": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Unique list of valid peer private CIDRs in the form < net_address > / < prefix >.",
			},
			"dpd": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
							Description: "The dead peer detection (DPD) action. A valid value is clear, hold, restart, disabled, or restart-by-peer. Default value is hold.",
						},
						"timeout": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The dead peer detection (DPD) timeout in seconds. A valid value is a positive integer that is greater than the DPD interval value. Default is 120.",
						},
						"interval": {
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
							Description: "The dead peer detection (DPD) interval, in seconds. A valid value is a positive integer. Default is 30.",
						},
					},
				},
				Description: "A dictionary with dead peer detection (DPD) protocol controls.",
			},
			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: networking.ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
			},
		},
		Description: "Manages a IPSec site connection resource within VKCS.",
	}
}

func resourceSiteConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var createOpts siteconnections.CreateOptsBuilder

	dpd := resourceSiteConnectionDPDCreateOpts(d.Get("dpd").(*schema.Set))

	v := d.Get("peer_cidrs").([]interface{})
	peerCidrs := make([]string, len(v))
	for i, v := range v {
		peerCidrs[i] = v.(string)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	initiator := resourceSiteConnectionInitiator(d.Get("initiator").(string))

	createOpts = SiteConnectionCreateOpts{
		CreateOpts: siteconnections.CreateOpts{
			Name:           d.Get("name").(string),
			Description:    d.Get("description").(string),
			AdminStateUp:   &adminStateUp,
			Initiator:      initiator,
			IKEPolicyID:    d.Get("ikepolicy_id").(string),
			PeerID:         d.Get("peer_id").(string),
			PeerAddress:    d.Get("peer_address").(string),
			PeerEPGroupID:  d.Get("peer_ep_group_id").(string),
			LocalID:        d.Get("local_id").(string),
			VPNServiceID:   d.Get("vpnservice_id").(string),
			LocalEPGroupID: d.Get("local_ep_group_id").(string),
			IPSecPolicyID:  d.Get("ipsecpolicy_id").(string),
			PSK:            d.Get("psk").(string),
			MTU:            d.Get("mtu").(int),
			PeerCIDRs:      peerCidrs,
			DPD:            &dpd,
		},
	}

	log.Printf("[DEBUG] Create site connection: %#v", createOpts)

	conn, err := siteconnections.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(conn.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"PENDING_CREATE"},
		Refresh:    waitForSiteConnectionCreation(networkingClient, conn.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] SiteConnection created: %#v", conn)

	return resourceSiteConnectionRead(ctx, d, meta)
}

func resourceSiteConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about site connection: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	conn, err := siteconnections.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "site_connection"))
	}

	log.Printf("[DEBUG] Read VKCS VPN SiteConnection %s: %#v", d.Id(), conn)

	d.Set("name", conn.Name)
	d.Set("description", conn.Description)
	d.Set("admin_state_up", conn.AdminStateUp)
	d.Set("initiator", conn.Initiator)
	d.Set("ikepolicy_id", conn.IKEPolicyID)
	d.Set("peer_id", conn.PeerID)
	d.Set("peer_address", conn.PeerAddress)
	d.Set("local_id", conn.LocalID)
	d.Set("peer_ep_group_id", conn.PeerEPGroupID)
	d.Set("vpnservice_id", conn.VPNServiceID)
	d.Set("local_ep_group_id", conn.LocalEPGroupID)
	d.Set("ipsecpolicy_id", conn.IPSecPolicyID)
	d.Set("psk", conn.PSK)
	d.Set("mtu", conn.MTU)
	d.Set("peer_cidrs", conn.PeerCIDRs)

	// Set the dpd
	dpdMap := make(map[string]interface{})
	dpdMap["action"] = conn.DPD.Action
	dpdMap["interval"] = conn.DPD.Interval
	dpdMap["timeout"] = conn.DPD.Timeout

	var dpd []map[string]interface{}
	dpd = append(dpd, dpdMap)
	if err := d.Set("dpd", &dpd); err != nil {
		log.Printf("[WARN] unable to set Site connection DPD")
	}

	return nil
}

func resourceSiteConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	opts := siteconnections.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		opts.AdminStateUp = &adminStateUp
		hasChange = true
	}

	if d.HasChange("local_id") {
		opts.LocalID = d.Get("local_id").(string)
		hasChange = true
	}

	if d.HasChange("peer_address") {
		opts.PeerAddress = d.Get("peer_address").(string)
		hasChange = true
	}

	if d.HasChange("peer_id") {
		opts.PeerID = d.Get("peer_id").(string)
		hasChange = true
	}

	if d.HasChange("local_ep_group_id") {
		opts.LocalEPGroupID = d.Get("local_ep_group_id").(string)
		hasChange = true
	}

	if d.HasChange("peer_ep_group_id") {
		opts.PeerEPGroupID = d.Get("peer_ep_group_id").(string)
		hasChange = true
	}

	if d.HasChange("psk") {
		opts.PSK = d.Get("psk").(string)
		hasChange = true
	}

	if d.HasChange("mtu") {
		opts.MTU = d.Get("mtu").(int)
		hasChange = true
	}

	if d.HasChange("initiator") {
		initiator := resourceSiteConnectionInitiator(d.Get("initiator").(string))
		opts.Initiator = initiator
		hasChange = true
	}

	if d.HasChange("peer_cidrs") {
		opts.PeerCIDRs = d.Get("peer_cidrs").([]string)
		hasChange = true
	}

	if d.HasChange("dpd") {
		dpdUpdateOpts := resourceSiteConnectionDPDUpdateOpts(d.Get("dpd").(*schema.Set))
		opts.DPD = &dpdUpdateOpts
		hasChange = true
	}

	var updateOpts siteconnections.UpdateOptsBuilder = opts

	log.Printf("[DEBUG] Updating site connection with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		conn, err := siteconnections.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForSiteConnectionUpdate(networkingClient, conn.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated connection with id %s", d.Id())
	}

	return resourceSiteConnectionRead(ctx, d, meta)
}

func resourceSiteConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy service: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	err = siteconnections.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForSiteConnectionDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	return diag.FromErr(err)
}

func waitForSiteConnectionDeletion(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := siteconnections.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got site connection %s => %#v", id, conn)

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] SiteConnection %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] SiteConnection %s deletion is pending", id)
		return conn, "DELETING", nil
	}
}

func waitForSiteConnectionCreation(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := siteconnections.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return service, "PENDING_CREATE", nil
	}
}

func waitForSiteConnectionUpdate(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		conn, err := siteconnections.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return conn, "UPDATED", nil
	}
}

func resourceSiteConnectionInitiator(initatorString string) siteconnections.Initiator {
	var ini siteconnections.Initiator
	switch initatorString {
	case "bi-directional":
		ini = siteconnections.InitiatorBiDirectional
	case "response-only":
		ini = siteconnections.InitiatorResponseOnly
	}
	return ini
}

func resourceSiteConnectionDPDCreateOpts(d *schema.Set) siteconnections.DPDCreateOpts {
	dpd := siteconnections.DPDCreateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		dpd.Action = resourceSiteConnectionAction(rawMap["action"].(string))

		timeout := rawMap["timeout"].(int)
		dpd.Timeout = timeout

		interval := rawMap["interval"].(int)
		dpd.Interval = interval
	}
	return dpd
}
func resourceSiteConnectionAction(actionString string) siteconnections.Action {
	var act siteconnections.Action
	switch actionString {
	case "hold":
		act = siteconnections.ActionHold
	case "restart":
		act = siteconnections.ActionRestart
	case "disabled":
		act = siteconnections.ActionDisabled
	case "restart-by-peer":
		act = siteconnections.ActionRestartByPeer
	case "clear":
		act = siteconnections.ActionClear
	}
	return act
}

func resourceSiteConnectionDPDUpdateOpts(d *schema.Set) siteconnections.DPDUpdateOpts {
	dpd := siteconnections.DPDUpdateOpts{}

	rawPairs := d.List()
	for _, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		dpd.Action = resourceSiteConnectionAction(rawMap["action"].(string))

		timeout := rawMap["timeout"].(int)
		dpd.Timeout = timeout

		interval := rawMap["interval"].(int)
		dpd.Interval = interval
	}
	return dpd
}
