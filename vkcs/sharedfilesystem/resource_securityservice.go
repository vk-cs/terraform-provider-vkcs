package sharedfilesystem

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	isecurityservices "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/sharedfilesystem/v2/securityservices"
)

func ResourceSharedFilesystemSecurityService() *schema.Resource {
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Shared File System client. A Shared File System client is needed to create a security service. If omitted, the `region` argument of the provider is used. Changing this creates a new security service.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the Security Service.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The name of the security service. Changing this updates the name of the existing security service.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the security service. Changing this updates the description of the existing security service.",
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"active_directory", "kerberos", "ldap",
				}, true),
				Description: "The security service type - can either be active\\_directory, kerberos or ldap.  Changing this updates the existing security service.",
			},

			"dns_ip": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The security service DNS IP address that is used inside the tenant network.",
			},

			"user": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The security service user or group name that is used by the tenant.",
			},

			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The user password, if you specify a user.",
			},

			"domain": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The security service domain.",
			},

			"server": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The security service host name or IP address.",
			},
		},
		Description: "Use this resource to configure a security service._note_ All arguments including the security service password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).\n\n" +
			"A security service stores configuration information for clients for authentication and authorization (AuthN/AuthZ). For example, a share server will be the client for an existing service such as LDAP, Kerberos, or Microsoft Active Directory.",
	}
}

func resourceSharedFilesystemSecurityServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = SharedFilesystemMinMicroversion

	createOpts := securityservices.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Type:        securityservices.SecurityServiceType(d.Get("type").(string)),
		DNSIP:       d.Get("dns_ip").(string),
		User:        d.Get("user").(string),
		Domain:      d.Get("domain").(string),
		Server:      d.Get("server").(string),
	}

	log.Printf("[DEBUG] vkcs_sharedfilesystem_securityservice create options: %#v", createOpts)
	createOpts.Password = d.Get("password").(string)
	securityservice, err := isecurityservices.Create(sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating vkcs_sharedfilesystem_securityservice: %s", err)
	}

	d.SetId(securityservice.ID)

	return resourceSharedFilesystemSecurityServiceRead(ctx, d, meta)
}

func resourceSharedFilesystemSecurityServiceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	// Select microversion to use.
	sfsClient.Microversion = SharedFilesystemMinMicroversion

	securityservice, err := isecurityservices.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error getting vkcs_sharedfilesystem_securityservice"))
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
	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourceSharedFilesystemSecurityServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = SharedFilesystemMinMicroversion

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

	_, err = isecurityservices.Update(sfsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating vkcs_sharedfilesystem_securityservice %s: %s", d.Id(), err)
	}

	return resourceSharedFilesystemSecurityServiceRead(ctx, d, meta)
}

func resourceSharedFilesystemSecurityServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	if err := isecurityservices.Delete(sfsClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_sharedfilesystem_securityservice"))
	}

	return nil
}
