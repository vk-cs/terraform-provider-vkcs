package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/zones"
)

const (
	zoneDelay         = 10 * time.Second
	zoneMinTimeout    = 3 * time.Second
	zoneCreateTimeout = 10 * time.Minute
	zoneDeleteTimeout = 10 * time.Minute
)

const (
	zoneStatusActive  = "active"
	zoneStatusDeleted = "deleted"
)

func resourcePublicDNSZone() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicDNSZoneCreate,
		ReadContext:   resourcePublicDNSZoneRead,
		UpdateContext: resourcePublicDNSZoneUpdate,
		DeleteContext: resourcePublicDNSZoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(zoneCreateTimeout),
			Delete: schema.DefaultTimeout(zoneDeleteTimeout),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used. Changing this creates a new zone.",
			},

			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the zone. **Changes this creates a new zone**.",
			},

			"primary_dns": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The primary DNS of the zone SOA.",
			},

			"admin_email": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The admin email of the zone SOA.",
			},

			"serial": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The serial number of the zone SOA.",
			},

			"refresh": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The refresh time of the zone SOA.",
			},

			"retry": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The retry time of the zone SOA.",
			},

			"expire": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The expire time of the zone SOA.",
			},

			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The TTL (time to live) of the zone SOA.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the zone.",
			},
		},
		Description: "Manages a public DNS record resource within VKCS. **New since v.0.2.0**.",
	}
}

func resourcePublicDNSZoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	createOpts := zones.CreateOpts{
		SOAPrimaryDNS: d.Get("primary_dns").(string),
		SOAAdminEmail: d.Get("admin_email").(string),
		SOARefresh:    d.Get("refresh").(int),
		SOARetry:      d.Get("retry").(int),
		SOAExpire:     d.Get("expire").(int),
		SOATTL:        d.Get("ttl").(int),
		Zone:          d.Get("zone").(string),
	}

	log.Printf("[DEBUG] vkcs_publicdns_zone create options: %#v", createOpts)

	zone, err := zones.Create(client, createOpts).Extract()
	if err != nil {
		return diag.FromErr(checkAlreadyExists(err, "Error creating vkcs_publicdns_zone",
			"vkcs_publicdns_zone", fmt.Sprintf("\"zone\" = %s", createOpts.Zone)))
	}

	d.SetId(zone.ID)
	log.Printf("[DEBUG] Created vkcs_publicdns_zone %s: %#v", zone.ID, zone)

	return resourcePublicDNSZoneRead(ctx, d, meta)
}

func resourcePublicDNSZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	zone, err := zones.Get(client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_publicdns_zone"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_publicdns_zone %s: %#v", d.Id(), zone)

	d.Set("region", getRegion(d, config))
	d.Set("zone", zone.Zone)
	d.Set("primary_dns", zone.PrimaryDNS)
	d.Set("admin_email", zone.AdminEmail)
	d.Set("serial", zone.Serial)
	d.Set("refresh", zone.Refresh)
	d.Set("retry", zone.Retry)
	d.Set("expire", zone.Expire)
	d.Set("ttl", zone.TTL)
	d.Set("status", zone.Status)

	return nil
}

func resourcePublicDNSZoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	updateOpts := zones.UpdateOpts{
		SOAPrimaryDNS: d.Get("primary_dns").(string),
		SOAAdminEmail: d.Get("admin_email").(string),
		SOARefresh:    d.Get("refresh").(int),
		SOARetry:      d.Get("retry").(int),
		SOAExpire:     d.Get("expire").(int),
		SOATTL:        d.Get("ttl").(int),
	}

	log.Printf("[DEBUG] vkcs_publicdns_zone create options: %#v", updateOpts)

	_, err = zones.Update(client, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error updating vkcs_publicdns_zone"))
	}

	return resourcePublicDNSZoneRead(ctx, d, meta)
}

func resourcePublicDNSZoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	err = zones.Delete(client, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_publicdns_zone"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{zoneStatusActive},
		Target:     []string{zoneStatusDeleted},
		Refresh:    publicDNSZoneStateRefreshFunc(client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      zoneDelay,
		MinTimeout: zoneMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_publicdns_zone %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
