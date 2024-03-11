package networking

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	iports "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking/v2/ports"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	inetworking "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/networking"
)

func ResourceNetworkingPortSecGroupAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingPortSecGroupAssociateCreate,
		ReadContext:   resourceNetworkingPortSecGroupAssociateRead,
		UpdateContext: resourceNetworkingPortSecGroupAssociateUpdate,
		DeleteContext: resourceNetworkingPortSecGroupAssociateDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the networking client. A networking client is needed to manage a port. If omitted, the `region` argument of the provider is used. Changing this creates a new resource.",
			},

			"port_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "An UUID of the port to apply security groups to.",
			},

			"security_group_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "A list of security group IDs to apply to the port. The security groups must be specified by ID and not name (as opposed to how they are configured with the Compute Instance).",
			},

			"enforce": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to replace or append the list of security groups, specified in the `security_group_ids`. Defaults to `false`.",
			},

			"all_security_group_ids": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The collection of Security Group IDs on the port which have been explicitly and implicitly added.",
			},

			"sdn": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				Computed:         true,
				ValidateDiagFunc: ValidateSDN(),
				Description:      "SDN to use for this resource. Must be one of following: \"neutron\", \"sprut\". Default value is project's default SDN.",
			},
		},
		Description: "Manages a port's security groups within VKCS. Useful, when the port was created not by Terraform. It should not be used, when the port was created directly within Terraform.\n\n" +
			"When the resource is deleted, Terraform doesn't delete the port, but unsets the list of user defined security group IDs.  However, if `enforce` is set to `true` and the resource is deleted, Terraform will remove all assigned security group IDs.",
	}
}

func resourceNetworkingPortSecGroupAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	securityGroups := util.ExpandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	portID := d.Get("port_id").(string)

	port, err := iports.Get(networkingClient, portID).Extract()
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
		sg := util.SliceUnion(port.SecurityGroups, securityGroups)
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
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var port portExtended
	err = iports.Get(networkingClient, d.Id()).ExtractInto(&port)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error fetching port security groups"))
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
			d.Set("security_group_ids", util.ExpandToStringSlice(actualSet.List()))
		}
	}

	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", port.SDN)

	return nil
}

func resourceNetworkingPortSecGroupAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var updateOpts ports.UpdateOpts
	var enforce bool
	if v, ok := d.GetOk("enforce"); ok {
		enforce = v.(bool)
	}

	if enforce {
		securityGroups := util.ExpandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
		updateOpts.SecurityGroups = &securityGroups
	} else {
		allSet := d.Get("all_security_group_ids").(*schema.Set)
		oldIDs, newIDs := d.GetChange("security_group_ids")
		oldSet, newSet := oldIDs.(*schema.Set), newIDs.(*schema.Set)

		allWithoutOld := allSet.Difference(oldSet)

		newSecurityGroups := util.ExpandToStringSlice(allWithoutOld.Union(newSet).List())

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
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), inetworking.SearchInAllSDNs)
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

		newSecurityGroups := util.ExpandToStringSlice(allWithoutOld.List())

		updateOpts.SecurityGroups = &newSecurityGroups
	}

	log.Printf("[DEBUG] Port security groups disassociation options: %#v", updateOpts.SecurityGroups)

	_, err = ports.Update(networkingClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error disassociating port security groups"))
	}

	return nil
}
