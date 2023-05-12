package publicdns

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/publicdns/v2/zones"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func DataSourcePublicDNSZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourcePublicDNSZoneRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the V2 Public DNS client. If omitted, the `region` argument of the provider is used.",
			},

			"id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The UUID of the DNS zone.",
			},

			"zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The name of the zone.",
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
				Computed:    true,
				Description: "The admin email of the zone SOA.",
			},

			"serial": {
				Type:        schema.TypeInt,
				Optional:    true,
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
				Optional:    true,
				Computed:    true,
				Description: "The status of the zone.",
			},
		},
		Description: "Use this data source to get the ID of a VKCS public DNS zone. **New since v.0.2.0**.",
	}
}

func dataSourcePublicDNSZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	client, err := config.PublicDNSV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS public DNS client: %s", err)
	}

	opts := zones.ListOpts{}

	if v, ok := d.GetOk("id"); ok {
		opts.ID = v.(string)
	}

	if v, ok := d.GetOk("zone"); ok {
		opts.Zone = v.(string)
	}

	if v, ok := d.GetOk("primary_dns"); ok {
		opts.SOAPrimaryDNS = v.(string)
	}

	if v, ok := d.GetOk("admin_email"); ok {
		opts.SOAAdminEmail = v.(string)
	}

	if v, ok := d.GetOk("serial"); ok {
		opts.SOASerial = v.(int)
	}

	if v, ok := d.GetOk("refresh"); ok {
		opts.SOARefresh = v.(int)
	}

	if v, ok := d.GetOk("retry"); ok {
		opts.SOARetry = v.(int)
	}

	if v, ok := d.GetOk("expire"); ok {
		opts.SOAExpire = v.(int)
	}

	if v, ok := d.GetOk("ttl"); ok {
		opts.SOATTL = v.(int)
	}

	if v, ok := d.GetOk("status"); ok {
		opts.Status = v.(string)
	}

	log.Printf("[DEBUG] vkcs_publicdns_zone list options: %#v", opts)

	zones, err := zones.List(client, opts).Extract()
	if err != nil {
		return diag.Errorf("Error retrieving vkcs_publicdns_zone: %s", err)
	}

	if len(zones) < 1 {
		return diag.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	if len(zones) > 1 {
		return diag.Errorf("Your query returned more than one result. Please try a more specific search criteria.")
	}

	zone := zones[0]

	log.Printf("[DEBUG] Retrieved vkcs_publicdns_zone %s: %#v", zone.ID, zone)

	d.SetId(zone.ID)
	d.Set("region", util.GetRegion(d, config))
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
