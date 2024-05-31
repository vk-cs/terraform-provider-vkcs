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
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/networking"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	iendpointgroups "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/vpnaas/v2/endpointgroups"
)

func ResourceEndpointGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointGroupCreate,
		ReadContext:   resourceEndpointGroupRead,
		UpdateContext: resourceEndpointGroupUpdate,
		DeleteContext: resourceEndpointGroupDelete,
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
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create an endpoint group. If omitted, the `region` argument of the provider is used. Changing this creates a new group.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the group. Changing this updates the name of the existing group.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the group. Changing this updates the description of the existing group.",
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ForceNew:    true,
				Description: "The type of the endpoints in the group. A valid value is subnet and cidr. For sprut SDN only cidr can be used, for neutron SDN - cidr for remote group, subnet for local. Changing this creates a new group.",
			},
			"endpoints": {
				Type:        schema.TypeSet,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "List of endpoints of the same type, for the endpoint group. The values will depend on the type. Changing this creates a new group.",
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
		Description: "Manages an Endpoint Group resource within VKCS.",
	}
}

func resourceEndpointGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var createOpts endpointgroups.CreateOptsBuilder

	endpointType := resourceEndpointGroupEndpointType(d.Get("type").(string))
	endpoints := util.ExpandToStringSlice(d.Get("endpoints").(*schema.Set).List())

	createOpts = EndpointGroupCreateOpts{
		CreateOpts: endpointgroups.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Endpoints:   endpoints,
			Type:        endpointType,
		},
	}

	log.Printf("[DEBUG] Create group: %#v", createOpts)

	group, err := iendpointgroups.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(group.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"PENDING_CREATE"},
		Target:     []string{"ACTIVE"},
		Refresh:    waitForEndpointGroupCreation(networkingClient, group.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] EndpointGroup created: %#v", group)

	return resourceEndpointGroupRead(ctx, d, meta)
}

func resourceEndpointGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about group: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var group groupExtended
	err = iendpointgroups.ExtractEndpointGroupInto(endpointgroups.Get(networkingClient, d.Id()), &group)
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "group"))
	}

	log.Printf("[DEBUG] Read VKCS Endpoint EndpointGroup %s: %#v", d.Id(), group)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("type", group.Type)
	d.Set("endpoints", group.Endpoints)
	d.Set("region", util.GetRegion(d, config))
	d.Set("sdn", group.SDN)

	return nil
}

func resourceEndpointGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	opts := endpointgroups.UpdateOpts{}

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

	var updateOpts endpointgroups.UpdateOptsBuilder = opts

	log.Printf("[DEBUG] Updating endpoint group with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		group, err := iendpointgroups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForEndpointGroupUpdate(networkingClient, group.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated group with id %s", d.Id())
	}

	return resourceEndpointGroupRead(ctx, d, meta)
}

func resourceEndpointGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy group: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	err = iendpointgroups.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForEndpointGroupDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	return diag.FromErr(err)
}

func waitForEndpointGroupDeletion(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := iendpointgroups.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got group %s => %#v", id, group)

		if err != nil {
			if errutil.IsNotFound(err) {
				log.Printf("[DEBUG] EndpointGroup %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] EndpointGroup %s deletion is pending", id)
		return group, "DELETING", nil
	}
}

func waitForEndpointGroupCreation(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := iendpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return group, "ACTIVE", nil
	}
}

func waitForEndpointGroupUpdate(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := iendpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return group, "UPDATED", nil
	}
}

func resourceEndpointGroupEndpointType(epType string) endpointgroups.EndpointType {
	var et endpointgroups.EndpointType
	switch epType {
	case "subnet":
		et = endpointgroups.TypeSubnet
	case "cidr":
		et = endpointgroups.TypeCIDR
	case "vlan":
		et = endpointgroups.TypeVLAN
	case "router":
		et = endpointgroups.TypeRouter
	case "network":
		et = endpointgroups.TypeNetwork
	}
	return et
}
