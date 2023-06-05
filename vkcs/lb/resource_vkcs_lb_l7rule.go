package lb

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
)

func ResourceL7Rule() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7RuleCreate,
		ReadContext:   resourceL7RuleRead,
		UpdateContext: resourceL7RuleUpdate,
		DeleteContext: resourceL7RuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceL7RuleImport,
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
				Description: "The region in which to obtain the Loadbalancer client. If omitted, the `region` argument of the provider is used. Changing this creates a new L7 Rule.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"COOKIE", "FILE_TYPE", "HEADER", "HOST_NAME", "PATH",
				}, true),
				Description: "The L7 Rule type - can either be COOKIE, FILE\\_TYPE, HEADER, HOST\\_NAME or PATH.",
			},

			"compare_type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"CONTAINS", "STARTS_WITH", "ENDS_WITH", "EQUAL_TO", "REGEX",
				}, true),
				Description: "The comparison type for the L7 rule - can either be CONTAINS, STARTS\\_WITH, ENDS_WITH, EQUAL_TO or REGEX",
			},

			"l7policy_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the L7 Policy to query. Changing this creates a new L7 Rule.",
			},

			"listener_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the Listener owning this resource.",
			},

			"value": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					if len(v.(string)) == 0 {
						errors = append(errors, fmt.Errorf("'value' field should not be empty"))
					}
					return
				},
				Description: "The value to use for the comparison. For example, the file type to compare.",
			},

			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key to use for the comparison. For example, the name of the cookie to evaluate. Valid when `type` is set to COOKIE or HEADER.",
			},

			"invert": {
				Type:        schema.TypeBool,
				Default:     false,
				Optional:    true,
				Description: "When true the logic of the rule is inverted. For example, with invert true, equal to would become not equal to. Default is false.",
			},

			"admin_state_up": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "The administrative state of the L7 Rule. A valid value is true (UP) or false (DOWN).",
			},
		},
		Description: "Manages a L7 Rule resource within VKCS.",
	}
}

func resourceL7RuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Assign some required variables for use in creation.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := ""
	ruleType := d.Get("type").(string)
	key := d.Get("key").(string)
	compareType := d.Get("compare_type").(string)
	adminStateUp := d.Get("admin_state_up").(bool)

	// Ensure the right combination of options have been specified.
	err = checkL7RuleType(ruleType, key)
	if err != nil {
		return diag.Errorf("Unable to create L7 Rule: %s", err)
	}

	createOpts := l7policies.CreateRuleOpts{
		RuleType:     l7policies.RuleType(ruleType),
		CompareType:  l7policies.CompareType(compareType),
		Value:        d.Get("value").(string),
		Key:          key,
		Invert:       d.Get("invert").(bool),
		AdminStateUp: &adminStateUp,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return diag.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaS extension
		listenerID, err = getListenerIDForL7Policy(lbClient, l7policyID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBL7Policy(ctx, lbClient, parentListener, parentL7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create L7 Rule")
	var l7Rule *l7policies.Rule
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		l7Rule, err = l7policies.CreateRule(lbClient, l7policyID, createOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating L7 Rule: %s", err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBL7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Rule.ID)
	d.Set("listener_id", listenerID)

	return resourceL7RuleRead(ctx, d, meta)
}

func resourceL7RuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	l7policyID := d.Get("l7policy_id").(string)

	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "L7 Rule"))
	}

	log.Printf("[DEBUG] Retrieved L7 Rule %s: %#v", d.Id(), l7Rule)

	d.Set("l7policy_id", l7policyID)
	d.Set("type", l7Rule.RuleType)
	d.Set("compare_type", l7Rule.CompareType)
	d.Set("value", l7Rule.Value)
	d.Set("key", l7Rule.Key)
	d.Set("invert", l7Rule.Invert)
	d.Set("admin_state_up", l7Rule.AdminStateUp)

	return nil
}

func resourceL7RuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	// Assign some required variables for use in updating.
	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)
	ruleType := d.Get("type").(string)
	key := d.Get("key").(string)

	// Key should always be set
	updateOpts := l7policies.UpdateRuleOpts{
		Key: &key,
	}

	if d.HasChange("type") {
		updateOpts.RuleType = l7policies.RuleType(ruleType)
	}
	if d.HasChange("compare_type") {
		updateOpts.CompareType = l7policies.CompareType(d.Get("compare_type").(string))
	}
	if d.HasChange("value") {
		updateOpts.Value = d.Get("value").(string)
	}
	if d.HasChange("invert") {
		invert := d.Get("invert").(bool)
		updateOpts.Invert = &invert
	}

	// Ensure the right combination of options have been specified.
	err = checkL7RuleType(ruleType, key)
	if err != nil {
		return diag.FromErr(err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return diag.Errorf("Unable to get parent L7 Policy: %s", err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to get L7 Rule: %s", err)
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBL7Policy(ctx, lbClient, parentListener, parentL7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBL7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating L7 Rule %s with options: %#v", d.Id(), updateOpts)
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, err := l7policies.UpdateRule(lbClient, l7policyID, d.Id(), updateOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to update L7 Rule %s: %s", d.Id(), err)
	}

	// Wait for L7 Rule to become active before continuing
	err = waitForLBL7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7RuleRead(ctx, d, meta)
}

func resourceL7RuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	l7policyID := d.Get("l7policy_id").(string)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent listener (%s) for the L7 Rule: %s", listenerID, err)
	}

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent L7 Policy (%s) for the L7 Rule: %s", l7policyID, err)
	}

	// Get a clean copy of the L7 Rule.
	l7Rule, err := l7policies.GetRule(lbClient, l7policyID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Unable to retrieve L7 Rule"))
	}

	// Wait for parent L7 Policy to become active before continuing
	err = waitForLBL7Policy(ctx, lbClient, parentListener, parentL7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Rule %s", d.Id())
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = l7policies.DeleteRule(lbClient, l7policyID, d.Id()).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting L7 Rule"))
	}

	err = waitForLBL7Rule(ctx, lbClient, parentListener, parentL7Policy, l7Rule, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7RuleImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for L7 Rule. Format must be <policy id>/<rule id>")
		return nil, err
	}

	config := meta.(clients.Config)
	lbClient, err := config.LoadBalancerV2Client(util.GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS loadbalancer client: %s", err)
	}

	listenerID := ""
	l7policyID := parts[0]
	l7ruleID := parts[1]

	// Get a clean copy of the parent L7 Policy.
	parentL7Policy, err := l7policies.Get(lbClient, l7policyID).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to get parent L7 Policy: %s", err)
	}

	if parentL7Policy.ListenerID != "" {
		listenerID = parentL7Policy.ListenerID
	} else {
		// Fallback for the Neutron LBaaS extension
		listenerID, err = getListenerIDForL7Policy(lbClient, l7policyID)
		if err != nil {
			return nil, err
		}
	}

	d.SetId(l7ruleID)
	d.Set("l7policy_id", l7policyID)
	d.Set("listener_id", listenerID)

	return []*schema.ResourceData{d}, nil
}

func checkL7RuleType(ruleType, key string) error {
	keyRequired := []string{"COOKIE", "HEADER"}
	if util.StrSliceContains(keyRequired, ruleType) && key == "" {
		return fmt.Errorf("key attribute is required, when the L7 Rule type is %s", strings.Join(keyRequired, " or "))
	} else if !util.StrSliceContains(keyRequired, ruleType) && key != "" {
		return fmt.Errorf("key attribute must not be used, when the L7 Rule type is not %s", strings.Join(keyRequired, " or "))
	}
	return nil
}
