package vkcs

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/l7policies"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/listeners"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func resourceL7Policy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyCreate,
		ReadContext:   resourceL7PolicyRead,
		UpdateContext: resourceL7PolicyUpdate,
		DeleteContext: resourceL7PolicyDelete,
		Importer: &schema.ResourceImporter{
			State: resourceL7PolicyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"action": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"REDIRECT_TO_POOL", "REDIRECT_TO_URL", "REJECT",
				}, true),
			},

			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"position": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"redirect_pool_id": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_url"},
				Optional:      true,
			},

			"redirect_url": {
				Type:          schema.TypeString,
				ConflictsWith: []string{"redirect_pool_id"},
				Optional:      true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					_, err := url.ParseRequestURI(value)
					if err != nil {
						errors = append(errors, fmt.Errorf("URL is not valid: %s", err))
					}
					return
				},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
		},
	}
}

func resourceL7PolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Assign some required variables for use in creation.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectURL := d.Get("redirect_url").(string)

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectURL, redirectPoolID)
	if err != nil {
		return diag.Errorf("Unable to create L7 Policy: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts := l7policies.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Action:         l7policies.Action(action),
		ListenerID:     listenerID,
		RedirectPoolID: redirectPoolID,
		RedirectURL:    redirectURL,
		AdminStateUp:   &adminStateUp,
	}

	if v, ok := d.GetOk("position"); ok {
		createOpts.Position = int32(v.(int))
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	// Make sure the associated pool is active before proceeding.
	if redirectPoolID != "" {
		pool, err := pools.Get(lbClient, redirectPoolID).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBPool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve listener %s: %s", listenerID, err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBListener(ctx, lbClient, parentListener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to create L7 Policy")
	var l7Policy *l7policies.L7Policy
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		l7Policy, err = l7policies.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating L7 Policy: %s", err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBL7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyRead(ctx, d, meta)
}

func resourceL7PolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "L7 Policy"))
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	d.Set("action", l7Policy.Action)
	d.Set("description", l7Policy.Description)
	d.Set("name", l7Policy.Name)
	d.Set("position", int(l7Policy.Position))
	d.Set("redirect_url", l7Policy.RedirectURL)
	d.Set("redirect_pool_id", l7Policy.RedirectPoolID)
	d.Set("region", getRegion(d, config))
	d.Set("admin_state_up", l7Policy.AdminStateUp)

	return nil
}

func resourceL7PolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	// Assign some required variables for use in updating.
	listenerID := d.Get("listener_id").(string)
	action := d.Get("action").(string)
	redirectPoolID := d.Get("redirect_pool_id").(string)
	redirectURL := d.Get("redirect_url").(string)

	var updateOpts l7policies.UpdateOpts

	if d.HasChange("action") {
		updateOpts.Action = l7policies.Action(action)
	}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("redirect_pool_id") {
		redirectPoolID = d.Get("redirect_pool_id").(string)

		updateOpts.RedirectPoolID = &redirectPoolID
	}
	if d.HasChange("redirect_url") {
		redirectURL = d.Get("redirect_url").(string)
		updateOpts.RedirectURL = &redirectURL
	}
	if d.HasChange("position") {
		updateOpts.Position = int32(d.Get("position").(int))
	}
	if d.HasChange("admin_state_up") {
		adminStateUp := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &adminStateUp
	}

	// Ensure the right combination of options have been specified.
	err = checkL7PolicyAction(action, redirectURL, redirectPoolID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Make sure the pool is active before continuing.
	timeout := d.Timeout(schema.TimeoutUpdate)
	if redirectPoolID != "" {
		pool, err := pools.Get(lbClient, redirectPoolID).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve %s: %s", redirectPoolID, err)
		}

		err = waitForLBPool(ctx, lbClient, pool, "ACTIVE", getLbPendingStatuses(), timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Get a clean copy of the parent listener.
	parentListener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent listener %s: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve L7 Policy: %s: %s", d.Id(), err)
	}

	// Wait for parent Listener to become active before continuing.
	err = waitForLBListener(ctx, lbClient, parentListener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBL7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, err = l7policies.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.Errorf("Unable to update L7 Policy %s: %s", d.Id(), err)
	}

	// Wait for L7 Policy to become active before continuing
	err = waitForLBL7Policy(ctx, lbClient, parentListener, l7Policy, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyRead(ctx, d, meta)
}

func resourceL7PolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS loadbalancer client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	listenerID := d.Get("listener_id").(string)

	// Get a clean copy of the listener.
	listener, err := listeners.Get(lbClient, listenerID).Extract()
	if err != nil {
		return diag.Errorf("Unable to retrieve parent listener (%s) for the L7 Policy: %s", listenerID, err)
	}

	// Get a clean copy of the L7 Policy.
	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Unable to retrieve L7 Policy"))
	}

	// Wait for Listener to become active before continuing.
	err = waitForLBListener(ctx, lbClient, listener, "ACTIVE", getLbPendingStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = l7policies.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting L7 Policy"))
	}

	err = waitForLBL7Policy(ctx, lbClient, listener, l7Policy, "DELETED", getLbPendingDeleteStatuses(), timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceL7PolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*config)
	lbClient, err := config.LoadBalancerV2Client(getRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS loadbalancer client: %s", err)
	}

	l7Policy, err := l7policies.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return nil, checkDeleted(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s during the import: %#v", d.Id(), l7Policy)

	if l7Policy.ListenerID != "" {
		d.Set("listener_id", l7Policy.ListenerID)
	} else {
		listenerID, err := getListenerIDForL7Policy(lbClient, d.Id())
		if err != nil {
			return nil, err
		}
		d.Set("listener_id", listenerID)
	}

	return []*schema.ResourceData{d}, nil
}

func checkL7PolicyAction(action, redirectURL, redirectPoolID string) error {
	if action == "REJECT" {
		if redirectURL != "" || redirectPoolID != "" {
			return fmt.Errorf(
				"redirect_url and redirect_pool_id must be empty when action is set to %s", action)
		}
	}

	if action == "REDIRECT_TO_POOL" && redirectURL != "" {
		return fmt.Errorf("redirect_url must be empty when action is set to %s", action)
	}

	if action == "REDIRECT_TO_URL" && redirectPoolID != "" {
		return fmt.Errorf("redirect_pool_id must be empty when action is set to %s", action)
	}

	return nil
}
