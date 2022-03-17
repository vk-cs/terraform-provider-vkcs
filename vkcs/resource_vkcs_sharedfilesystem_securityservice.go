package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
)

func resourceSharedFilesystemSecurityService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemSecurityServiceCreate,
		ReadContext:   resourceSharedFilesystemSecurityServiceRead,
		UpdateContext: resourceSharedFilesystemSecurityServiceUpdate,
		DeleteContext: resourceSharedFilesystemSecurityServiceDelete,

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

			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active_directory", "kerberos", "ldap",
				}, true),
			},

			"dns_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ou": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"user": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"domain": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"server": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSharedFilesystemSecurityServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemMinMicroversion

	createOpts := securityservices.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        securityservices.SecurityServiceType(d.Get("type").(string)),
		DNSIP:       d.Get("dns_ip").(string),
		User:        d.Get("user").(string),
		Domain:      d.Get("domain").(string),
		Server:      d.Get("server").(string),
	}

	if v, ok := d.GetOkExists("ou"); ok {
		createOpts.OU = v.(string)

		sfsClient.Microversion = sharedFilesystemSecurityServiceOUMicroversion
	}

	log.Printf("[DEBUG] vkcs_sharedfilesystem_securityservice create options: %#v", createOpts)
	createOpts.Password = d.Get("password").(string)
	securityservice, err := securityservices.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_sharedfilesystem_securityservice: %s", err)
	}

	d.SetId(securityservice.ID)

	return resourceSharedFilesystemSecurityServiceRead(ctx, d, meta)
}

func resourceSharedFilesystemSecurityServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	// Select microversion to use.
	sfsClient.Microversion = sharedFilesystemMinMicroversion
	if _, ok := d.GetOkExists("ou"); ok {
		sfsClient.Microversion = sharedFilesystemSecurityServiceOUMicroversion
	}

	securityservice, err := securityservices.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error getting vkcs_sharedfilesystem_securityservice"))
	}

	// Workaround for resource import.
	if securityservice.OU == "" {
		sfsClient.Microversion = sharedFilesystemSecurityServiceOUMicroversion
		securityserviceOU, err := securityservices.Get(sfsClient, d.Id()).Extract()
		if err == nil {
			d.Set("ou", securityserviceOU.OU)
		}
	}

	nopassword := securityservice
	nopassword.Password = ""
	log.Printf("[DEBUG] Retrieved vkcs_sharedfilesystem_securityservice %s: %#v", d.Id(), nopassword)

	d.Set("name", securityservice.Name)
	d.Set("description", securityservice.Description)
	d.Set("type", securityservice.Type)
	d.Set("domain", securityservice.Domain)
	d.Set("dns_ip", securityservice.DNSIP)
	d.Set("user", securityservice.User)
	d.Set("server", securityservice.Server)

	// Computed.
	d.Set("project_id", securityservice.ProjectID)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceSharedFilesystemSecurityServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemMinMicroversion

	var updateOpts securityservices.UpdateOpts

	// Name should always be sent, otherwise it is vanished by manila backend.
	name := d.Get("name").(string)
	updateOpts.Name = &name

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("type") {
		updateOpts.Type = d.Get("type").(string)
	}

	if d.HasChange("dns_ip") {
		dnsIP := d.Get("dns_ip").(string)
		updateOpts.DNSIP = &dnsIP
	}

	if d.HasChange("ou") {
		ou := d.Get("ou").(string)
		updateOpts.OU = &ou

		sfsClient.Microversion = sharedFilesystemSecurityServiceOUMicroversion
	}

	if d.HasChange("user") {
		user := d.Get("user").(string)
		updateOpts.User = &user
	}

	if d.HasChange("domain") {
		domain := d.Get("domain").(string)
		updateOpts.Domain = &domain
	}

	if d.HasChange("server") {
		server := d.Get("server").(string)
		updateOpts.Server = &server
	}

	log.Printf("[DEBUG] vkcs_sharedfilesystem_securityservice %s update options: %#v", d.Id(), updateOpts)

	if d.HasChange("password") {
		password := d.Get("password").(string)
		updateOpts.Password = &password
	}

	_, err = securityservices.Update(sfsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_sharedfilesystem_securityservice %s: %s", d.Id(), err)
	}

	return resourceSharedFilesystemSecurityServiceRead(ctx, d, meta)
}

func resourceSharedFilesystemSecurityServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	if err := securityservices.Delete(sfsClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_sharedfilesystem_securityservice"))
	}

	return nil
}
