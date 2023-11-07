package sharedfilesystem

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	sfserrors "github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	ishares "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/sharedfilesystem/v2/shares"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

func ResourceSharedFilesystemShareAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareAccessCreate,
		ReadContext:   resourceSharedFilesystemShareAccessRead,
		DeleteContext: resourceSharedFilesystemShareAccessDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSharedFilesystemShareAccessImport,
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
				Description: "The region in which to obtain the Shared File System client. A Shared File System client is needed to create a share access. Changing this creates a new share access.",
			},

			"share_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the share to which you are granted access.",
			},

			"access_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "user", "cert", "cephx",
				}, false),
				Description: "The access rule type. Can either be an ip, user, cert, or cephx.",
			},

			"access_to": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The value that defines the access. Can either be an IP address or a username verified by configured Security Service of the Share Network.",
			},

			"access_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"rw", "ro",
				}, false),
				Description: "The access level to the share. Can either be `rw` or `ro`.",
			},
		},
		Description: "Use this resource to control the share access lists.\n\n" +
			"~> **Important Security Notice** The access key assigned by this resource will be stored *unencrypted* in your Terraform state file. If you use this resource in production, please make sure your state file is sufficiently protected. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).",
	}
}

func resourceSharedFilesystemShareAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = SharedFilesystemMinMicroversion
	accessType := d.Get("access_type").(string)
	if accessType == "cephx" {
		sfsClient.Microversion = sharedFilesystemSharedAccessCephXMicroversion
	}

	shareID := d.Get("share_id").(string)

	grantOpts := shares.GrantAccessOpts{
		AccessType:  accessType,
		AccessTo:    d.Get("access_to").(string),
		AccessLevel: d.Get("access_level").(string),
	}

	log.Printf("[DEBUG] vkcs_sharedfilesystem_share_access create options: %#v", grantOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	var access *shares.AccessRight
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		access, err = ishares.GrantAccess(sfsClient, shareID, grantOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		detailedErr := sfserrors.ErrorDetails{}
		e := sfserrors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error creating vkcs_sharedfilesystem_share_access: %s: %s", err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error creating vkcs_sharedfilesystem_share_access: %s (%d): %s", k, msg.Code, msg.Message)
		}
	}

	d.SetId(access.ID)

	log.Printf("[DEBUG] Waiting for vkcs_sharedfilesystem_share_access %s to become available.", access.ID)
	stateConf := &retry.StateChangeConf{
		Target:     []string{"active"},
		Pending:    []string{"new", "queued_to_apply", "applying"},
		Refresh:    sharedFilesystemShareAccessStateRefreshFunc(sfsClient, shareID, access.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for vkcs_sharedfilesystem_share_access %s to become available: %s", access.ID, err)
	}

	return resourceSharedFilesystemShareAccessRead(ctx, d, meta)
}

func resourceSharedFilesystemShareAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	// Set the client to the minimum supported microversion.
	sfsClient.Microversion = SharedFilesystemMinMicroversion

	shareID := d.Get("share_id").(string)
	access, err := ishares.ListAccessRights(sfsClient, shareID).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_sharedfilesystem_share_access"))
	}

	for _, v := range access {
		if v.ID == d.Id() {
			log.Printf("[DEBUG] Retrieved vkcs_sharedfilesystem_share_access %s: %#v", d.Id(), v)
			d.Set("access_type", v.AccessType)
			d.Set("access_to", v.AccessTo)
			d.Set("access_level", v.AccessLevel)
			d.Set("region", util.GetRegion(d, config))
			return nil
		}
	}

	log.Printf("[DEBUG] Unable to find vkcs_sharedfilesystem_share_access %s", d.Id())
	d.SetId("")

	return nil
}

func resourceSharedFilesystemShareAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = SharedFilesystemMinMicroversion

	shareID := d.Get("share_id").(string)

	revokeOpts := shares.RevokeAccessOpts{AccessID: d.Id()}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete vkcs_sharedfilesystem_share_access %s", d.Id())
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = ishares.RevokeAccess(sfsClient, shareID, revokeOpts).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		e := util.CheckDeleted(d, err, "Error deleting vkcs_sharedfilesystem_share_access")
		if e == nil {
			return nil
		}
		detailedErr := sfserrors.ErrorDetails{}
		e = sfserrors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error waiting for vkcs_sharedfilesystem_share_access on %s to be removed: %s: %s", shareID, err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error waiting for vkcs_sharedfilesystem_share_access on %s to be removed: %s (%d): %s", shareID, k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for vkcs_sharedfilesystem_share_access %s to become denied.", d.Id())
	stateConf := &retry.StateChangeConf{
		Target:     []string{"denied"},
		Pending:    []string{"active", "new", "queued_to_deny", "denying"},
		Refresh:    sharedFilesystemShareAccessStateRefreshFunc(sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if errutil.IsNotFound(err) {
			return nil
		}
		return diag.Errorf("error waiting for vkcs_sharedfilesystem_share_access %s to become denied: %s", d.Id(), err)
	}

	return nil
}

func resourceSharedFilesystemShareAccessImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("invalid format specified for vkcs_sharedfilesystem_share_access. Format must be <share id>/<ACL id>")
		return nil, err
	}

	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = SharedFilesystemMinMicroversion

	shareID := parts[0]
	accessID := parts[1]

	access, err := ishares.ListAccessRights(sfsClient, shareID).Extract()
	if err != nil {
		return nil, fmt.Errorf("unable to get %s vkcs_sharedfilesystem_share: %s", shareID, err)
	}

	for _, v := range access {
		if v.ID == accessID {
			log.Printf("[DEBUG] Retrieved vkcs_sharedfilesystem_share_access %s: %#v", accessID, v)

			d.SetId(accessID)
			d.Set("share_id", shareID)
			return []*schema.ResourceData{d}, nil
		}
	}

	return nil, fmt.Errorf("[DEBUG] Unable to find vkcs_sharedfilesystem_share_access %s", accessID)
}
