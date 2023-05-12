package firewall

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func ResourceNetworkingSecGroupRule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupRuleCreate,
		ReadContext:   resourceNetworkingSecGroupRuleRead,
		DeleteContext: resourceNetworkingSecGroupRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the networking client. A networking client is needed to create a port. If omitted, the `region` argument of the provider is used. Changing this creates a new security group rule.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "A description of the rule. Changing this creates a new security group rule.",
			},

			"direction": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The direction of the rule, valid values are __ingress__ or __egress__. Changing this creates a new security group rule.",
			},

			"ethertype": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The layer 3 protocol type, the only valid value is __IPv4__. Changing this creates a new security group rule. **Note** This argument is deprecated, please do not use it.",
				Default:      "IPv4",
				Deprecated:   "Only IPv4 can be used as ethertype. This argument is deprecated, please do not use it.",
				ValidateFunc: validation.StringInSlice([]string{"IPv4"}, false),
			},

			"port_range_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The lower part of the allowed port range, valid integer value needs to be between 1 and 65535. Changing this creates a new security group rule.",
			},

			"port_range_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The higher part of the allowed port range, valid integer value needs to be between 1 and 65535. Changing this creates a new security group rule.",
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Description: "The layer 4 protocol type, valid values are following. Changing this creates a new security group rule. This is required if you want to specify a port range.\n" +
					"  * __tcp__\n" +
					"  * __udp__\n" +
					"  * __icmp__\n" +
					"  * __ah__\n" +
					"  * __dccp__\n" +
					"  * __egp__\n" +
					"  * __esp__\n" +
					"  * __gre__\n" +
					"  * __igmp__\n" +
					"  * __ospf__\n" +
					"  * __pgm__\n" +
					"  * __rsvp__\n" +
					"  * __sctp__\n" +
					"  * __udplite__\n" +
					"  * __vrrp__",
			},

			"remote_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"remote_ip_prefix"},
				Description:   "The remote group id, the value needs to be an ID of a security group in the same tenant. Changing this creates a new security group rule. **Note**: Only one of `remote_group_id` or `remote_ip_prefix` may be set.",
			},

			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
				ConflictsWith: []string{"remote_group_id"},
				Description:   "The remote CIDR, the value needs to be a valid CIDR (i.e. 192.168.0.0/16). Changing this creates a new security group rule. **Note**: Only one of `remote_group_id` or `remote_ip_prefix` may be set.",
			},

			"security_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The security group id the rule should belong to, the value needs to be an ID of a security group in the same tenant. Changing this creates a new security group rule.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: networking.ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is \"neutron\".",
			},
		},
		Description: "Manages a security group rule resource within VKCS.",
	}
}

func resourceNetworkingSecGroupRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	securityGroupID := d.Get("security_group_id").(string)
	mutex := config.GetMutex()
	mutex.Lock(securityGroupID)
	defer mutex.Unlock(securityGroupID)

	portRangeMin := d.Get("port_range_min").(int)
	portRangeMax := d.Get("port_range_max").(int)
	protocol := d.Get("protocol").(string)

	if protocol == "" {
		if portRangeMin != 0 || portRangeMax != 0 {
			return diag.Errorf("A protocol must be specified when using port_range_min and port_range_max for vkcs_networking_secgroup_rule")
		}
	}

	opts := rules.CreateOpts{
		Description:    d.Get("description").(string),
		SecGroupID:     d.Get("security_group_id").(string),
		PortRangeMin:   d.Get("port_range_min").(int),
		PortRangeMax:   d.Get("port_range_max").(int),
		RemoteGroupID:  d.Get("remote_group_id").(string),
		RemoteIPPrefix: d.Get("remote_ip_prefix").(string),
	}

	if v, ok := d.GetOk("direction"); ok {
		direction, err := resourceNetworkingSecGroupRuleDirection(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Direction = direction
	}

	if v, ok := d.GetOk("ethertype"); ok {
		ethertype, err := resourceNetworkingSecGroupRuleEtherType(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.EtherType = ethertype
	}

	if v, ok := d.GetOk("protocol"); ok {
		protocol, err := resourceNetworkingSecGroupRuleProtocol(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		opts.Protocol = protocol
	}

	log.Printf("[DEBUG] vkcs_networking_secgroup_rule create options: %#v", opts)

	sgRule, err := rules.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_secgroup_rule: %s", err)
	}

	d.SetId(sgRule.ID)

	log.Printf("[DEBUG] Created vkcs_networking_secgroup_rule %s: %#v", sgRule.ID, sgRule)
	return resourceNetworkingSecGroupRuleRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	sgRule, err := rules.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_networking_secgroup_rule"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_networking_secgroup_rule %s: %#v", d.Id(), sgRule)

	d.Set("description", sgRule.Description)
	d.Set("direction", sgRule.Direction)
	d.Set("ethertype", sgRule.EtherType)
	d.Set("protocol", sgRule.Protocol)
	d.Set("port_range_min", sgRule.PortRangeMin)
	d.Set("port_range_max", sgRule.PortRangeMax)
	d.Set("remote_group_id", sgRule.RemoteGroupID)
	d.Set("remote_ip_prefix", sgRule.RemoteIPPrefix)
	d.Set("security_group_id", sgRule.SecGroupID)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", networking.GetSDN(d))

	return nil
}

func resourceNetworkingSecGroupRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	securityGroupID := d.Get("security_group_id").(string)
	mutex := config.GetMutex()
	mutex.Lock(securityGroupID)
	defer mutex.Unlock(securityGroupID)

	if err := rules.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_networking_secgroup_rule"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingSecGroupRuleStateRefreshFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_networking_secgroup_rule %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}
