package vkcs

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroupRule() *schema.Resource {
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
				ForceNew: true,
			},

			"direction": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"ethertype": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port_range_min": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"port_range_max": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"remote_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"remote_ip_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				StateFunc: func(v interface{}) string {
					return strings.ToLower(v.(string))
				},
			},

			"security_group_id": {
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

func resourceNetworkingSecGroupRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
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
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	sgRule, err := rules.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_networking_secgroup_rule"))
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
	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}

func resourceNetworkingSecGroupRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	securityGroupID := d.Get("security_group_id").(string)
	mutex := config.GetMutex()
	mutex.Lock(securityGroupID)
	defer mutex.Unlock(securityGroupID)

	if err := rules.Delete(networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_networking_secgroup_rule"))
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
