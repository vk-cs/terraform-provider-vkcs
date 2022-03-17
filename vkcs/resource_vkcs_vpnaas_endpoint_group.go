package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
)

func resourceEndpointGroup() *schema.Resource {
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ForceNew: true,
			},
			"endpoints": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
		},
	}
}

func resourceEndpointGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var createOpts endpointgroups.CreateOptsBuilder

	endpointType := resourceEndpointGroupEndpointType(d.Get("type").(string))
	endpoints := expandToStringSlice(d.Get("endpoints").(*schema.Set).List())

	createOpts = EndpointGroupCreateOpts{
		endpointgroups.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Endpoints:   endpoints,
			Type:        endpointType,
		},
	}

	log.Printf("[DEBUG] Create group: %#v", createOpts)

	group, err := endpointgroups.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
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

	d.SetId(group.ID)

	return resourceEndpointGroupRead(ctx, d, meta)
}

func resourceEndpointGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about group: %s", d.Id())

	config := meta.(*config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	group, err := endpointgroups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "group"))
	}

	log.Printf("[DEBUG] Read VKCS Endpoint EndpointGroup %s: %#v", d.Id(), group)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("type", group.Type)
	d.Set("endpoints", group.Endpoints)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceEndpointGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
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

	var updateOpts endpointgroups.UpdateOptsBuilder
	updateOpts = opts

	log.Printf("[DEBUG] Updating endpoint group with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		group, err := endpointgroups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &resource.StateChangeConf{
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

	config := meta.(*config)
	networkingClient, err := config.NetworkingV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	err = endpointgroups.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
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

func waitForEndpointGroupDeletion(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got group %s => %#v", id, group)

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] EndpointGroup %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("Unexpected error: %s", err)
		}

		log.Printf("[DEBUG] EndpointGroup %s deletion is pending", id)
		return group, "DELETING", nil
	}
}

func waitForEndpointGroupCreation(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_CREATE", nil
		}
		return group, "ACTIVE", nil
	}
}

func waitForEndpointGroupUpdate(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := endpointgroups.Get(networkingClient, id).Extract()
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
