package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupCreate,
		ReadContext:   resourceNetworkingSecGroupRead,
		UpdateContext: resourceNetworkingSecGroupUpdate,
		DeleteContext: resourceNetworkingSecGroupDelete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"delete_default_rules": {
				Type:     schema.TypeBool,
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
		},
	}
}

func resourceNetworkingSecGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[DEBUG] vkcs_networking_secgroup create options: %#v", opts)
	sg, err := groups.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_networking_secgroup: %s", err)
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		sgID := sg.ID
		sg, err := groups.Get(networkingClient, sgID).Extract()
		if err != nil {
			return diag.Errorf("Error retrieving the created vkcs_networking_secgroup %s: %s", sgID, err)
		}

		for _, rule := range sg.Rules {
			if err := rules.Delete(networkingClient, rule.ID).ExtractErr(); err != nil {
				return diag.Errorf("Error deleting a default rule for vkcs_networking_secgroup %s: %s", sgID, err)
			}
		}
	}

	d.SetId(sg.ID)

	tags := networkingAttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "security-groups", sg.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_secgroup %s: %s", sg.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_secgroup %s", tags, sg.ID)
	}

	log.Printf("[DEBUG] Created vkcs_networking_secgroup: %#v", sg)

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	sg, err := groups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_networking_secgroup"))
	}

	d.Set("description", sg.Description)
	d.Set("name", sg.Name)
	d.Set("region", getRegion(d, config))

	networkingReadAttributesTags(d, sg.Tags)

	return nil
}

func resourceNetworkingSecGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		updated    bool
		updateOpts groups.UpdateOpts
	)

	if d.HasChange("name") {
		updated = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updated = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if updated {
		log.Printf("[DEBUG] Updating vkcs_networking_secgroup %s with options: %#v", d.Id(), updateOpts)
		_, err = groups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating vkcs_networking_secgroup: %s", err)
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "security-groups", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on vkcs_networking_secgroup %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on vkcs_networking_secgroup %s", tags, d.Id())
	}

	return resourceNetworkingSecGroupRead(ctx, d, meta)
}

func resourceNetworkingSecGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingSecgroupStateRefreshFuncDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting vkcs_networking_secgroup: %s", err)
	}

	return diag.FromErr(err)
}
