package vkcs

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func resourceNetworkingPortSecGroupAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingPortSecGroupAssociateCreate,
		ReadContext:   resourceNetworkingPortSecGroupAssociateRead,
		UpdateContext: resourceNetworkingPortSecGroupAssociateUpdate,
		DeleteContext: resourceNetworkingPortSecGroupAssociateDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"enforce": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"all_security_group_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: validateSDN(),
			},
		},
	}
}

func resourceNetworkingPortSecGroupAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	portID := d.Get("port_id").(string)

	port, err := ports.Get(networkingClient, portID).Extract()
	if err != nil {
		return diag.Errorf("Unable to get %s Port: %s", portID, err)
	}

	log.Printf("[DEBUG] Retrieved Port %s: %+v", portID, port)

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOk("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		updateOpts.SecurityGroups = &securityGroups
	} else {
		// append security groups
		sg := sliceUnion(port.SecurityGroups, securityGroups)
		updateOpts.SecurityGroups = &sg
	}

	log.Printf("[DEBUG] Port Security Group Associate Options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(networkingClient, portID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error associating %s port with '%s' security groups: %s", portID, strings.Join(securityGroups, ","), err)
	}

	d.SetId(portID)

	return resourceNetworkingPortSecGroupAssociateRead(ctx, d, meta)
}

func resourceNetworkingPortSecGroupAssociateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	port, err := ports.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error fetching port security groups"))
	}

	var enforce bool
	if v, ok := d.GetOk("enforce"); ok {
		enforce = v.(bool)
	}

	d.Set("all_security_group_ids", port.SecurityGroups)

	if enforce {
		d.Set("security_group_ids", port.SecurityGroups)
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		desiredSet := d.Get("security_group_ids").(*schema.Set)
		actualSet := allSet.Intersection(desiredSet)
		if !actualSet.Equal(desiredSet) {
			d.Set("security_group_ids", expandToStringSlice(actualSet.List()))
		}
	}

	d.Set("region", getRegion(d, config))
	d.Set("sdn", getSDN(d))

	return nil
}

func resourceNetworkingPortSecGroupAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOk("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
		updateOpts.SecurityGroups = &securityGroups
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldIDs, newIDs := d.GetChange("security_group_ids")
		oldSet, newSet := oldIDs.(*schema.Set), newIDs.(*schema.Set)

		allWithoutOld := allSet.Difference(oldSet)

		newSecurityGroups := expandToStringSlice(allWithoutOld.Union(newSet).List())

		updateOpts.SecurityGroups = &newSecurityGroups
	}

	if d.HasChange("security_group_ids") || d.HasChange("enforce") {
		log.Printf("[DEBUG] Port Security Group Update Options: %#v", updateOpts.SecurityGroups)
		_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating VKCS networking Port: %s", err)
		}
	}

	return resourceNetworkingPortSecGroupAssociateRead(ctx, d, meta)
}

func resourceNetworkingPortSecGroupAssociateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config), getSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOk("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		updateOpts.SecurityGroups = &[]string{}
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldSet := d.Get("security_group_ids").(*schema.Set)

		allWithoutOld := allSet.Difference(oldSet)

		newSecurityGroups := expandToStringSlice(allWithoutOld.List())

		updateOpts.SecurityGroups = &newSecurityGroups
	}

	log.Printf("[DEBUG] Port security groups disassociation options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error disassociating port security groups"))
	}

	return nil
}
