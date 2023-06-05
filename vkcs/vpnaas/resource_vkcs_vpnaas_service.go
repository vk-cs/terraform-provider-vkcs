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
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
)

func ResourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,
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
				Description: "The region in which to obtain the Networking client. A Networking client is needed to create a VPN service. If omitted, the `region` argument of the provider is used. Changing this creates a new service.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the service. Changing this updates the name of the existing service.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the service. Changing this updates the description of the existing service.",
			},
			"admin_state_up": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "The administrative state of the resource. Can either be up(true) or down(false). Changing this updates the administrative state of the existing service.",
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "SubnetID is the ID of the subnet. Default is null.",
			},
			"router_id": {
				Type:        schema.TypeString,
				Required:    true,
				Computed:    false,
				ForceNew:    true,
				Description: "The ID of the router. Changing this creates a new service.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Indicates whether IPsec VPN service is currently operational. Values are ACTIVE, DOWN, BUILD, ERROR, PENDING_CREATE, PENDING_UPDATE, or PENDING_DELETE.",
			},
			"external_v6_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The read-only external (public) IPv6 address that is used for the VPN service.",
			},
			"external_v4_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The read-only external (public) IPv4 address that is used for the VPN service.",
			},
		},
		Description: "Manages a VPN service resource within VKCS.",
	}
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	var createOpts services.CreateOptsBuilder

	adminStateUp := d.Get("admin_state_up").(bool)
	createOpts = ServiceCreateOpts{
		CreateOpts: services.CreateOpts{
			Name:         d.Get("name").(string),
			Description:  d.Get("description").(string),
			AdminStateUp: &adminStateUp,
			SubnetID:     d.Get("subnet_id").(string),
			RouterID:     d.Get("router_id").(string),
		},
	}

	log.Printf("[DEBUG] Create service: %#v", createOpts)

	service, err := services.Create(networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"PENDING_CREATE"},
		Refresh:    waitForServiceCreation(networkingClient, service.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)

	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Service created: %#v", service)

	d.SetId(service.ID)

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about service: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	service, err := services.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "service"))
	}

	log.Printf("[DEBUG] Read VKCS VPN Service %s: %#v", d.Id(), service)

	d.Set("name", service.Name)
	d.Set("description", service.Description)
	d.Set("subnet_id", service.SubnetID)
	d.Set("admin_state_up", service.AdminStateUp)
	d.Set("router_id", service.RouterID)
	d.Set("status", service.Status)
	d.Set("external_v6_ip", service.ExternalV6IP)
	d.Set("external_v4_ip", service.ExternalV4IP)
	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	opts := services.UpdateOpts{}

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

	var updateOpts services.UpdateOptsBuilder = opts

	log.Printf("[DEBUG] Updating service with id %s: %#v", d.Id(), updateOpts)

	if hasChange {
		service, err := services.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{"PENDING_UPDATE"},
			Target:     []string{"UPDATED"},
			Refresh:    waitForServiceUpdate(networkingClient, service.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}
		_, err = stateConf.WaitForStateContext(ctx)

		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated service with id %s", d.Id())
	}

	return resourceServiceRead(ctx, d, meta)
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy service: %s", d.Id())

	config := meta.(clients.Config)
	networkingClient, err := config.NetworkingV2Client(util.GetRegion(d, config), networking.GetSDN(d))
	if err != nil {
		return diag.Errorf("Error creating VKCS networking client: %s", err)
	}

	err = services.Delete(networkingClient, d.Id()).Err

	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"DELETING"},
		Target:     []string{"DELETED"},
		Refresh:    waitForServiceDeletion(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)

	return diag.FromErr(err)
}

func waitForServiceDeletion(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		serv, err := services.Get(networkingClient, id).Extract()
		log.Printf("[DEBUG] Got service %s => %#v", id, serv)

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Service %s is actually deleted", id)
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("unexpected error: %s", err)
		}

		log.Printf("[DEBUG] Service %s deletion is pending", id)
		return serv, "DELETING", nil
	}
}

func waitForServiceCreation(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := services.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "NOT_CREATED", nil
		}
		return service, "PENDING_CREATE", nil
	}
}

func waitForServiceUpdate(networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		service, err := services.Get(networkingClient, id).Extract()
		if err != nil {
			return "", "PENDING_UPDATE", nil
		}
		return service, "UPDATED", nil
	}
}
