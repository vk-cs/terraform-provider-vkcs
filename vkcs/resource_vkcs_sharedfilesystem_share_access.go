package vkcs

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func resourceSharedFilesystemShareAccess() *schema.Resource {
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ip", "user", "cert", "cephx",
				}, false),
			},

			"access_to": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"access_level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"rw", "ro",
				}, false),
			},
		},
	}
}

func resourceSharedFilesystemShareAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemMinMicroversion
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
	err = resource.Retry(timeout, func() *resource.RetryError {
		access, err = shares.GrantAccess(sfsClient, shareID, grantOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		detailedErr := errors.ErrorDetails{}
		e := errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error creating vkcs_sharedfilesystem_share_access: %s: %s", err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error creating vkcs_sharedfilesystem_share_access: %s (%d): %s", k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for vkcs_sharedfilesystem_share_access %s to become available.", access.ID)
	stateConf := &resource.StateChangeConf{
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

	d.SetId(access.ID)

	return resourceSharedFilesystemShareAccessRead(ctx, d, meta)
}

func resourceSharedFilesystemShareAccessRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	// Set the client to the minimum supported microversion.
	sfsClient.Microversion = sharedFilesystemMinMicroversion

	shareID := d.Get("share_id").(string)
	access, err := shares.ListAccessRights(sfsClient, shareID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_sharedfilesystem_share_access"))
	}

	for _, v := range access {
		if v.ID == d.Id() {
			log.Printf("[DEBUG] Retrieved vkcs_sharedfilesystem_share_access %s: %#v", d.Id(), v)
			d.Set("access_type", v.AccessType)
			d.Set("access_to", v.AccessTo)
			d.Set("access_level", v.AccessLevel)
			d.Set("region", getRegion(d, config))
			return nil
		}
	}

	log.Printf("[DEBUG] Unable to find vkcs_sharedfilesystem_share_access %s", d.Id())
	d.SetId("")

	return nil
}

func resourceSharedFilesystemShareAccessDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemMinMicroversion

	shareID := d.Get("share_id").(string)

	revokeOpts := shares.RevokeAccessOpts{AccessID: d.Id()}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete vkcs_sharedfilesystem_share_access %s", d.Id())
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = shares.RevokeAccess(sfsClient, shareID, revokeOpts).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		e := checkDeleted(d, err, "Error deleting vkcs_sharedfilesystem_share_access")
		if e == nil {
			return nil
		}
		detailedErr := errors.ErrorDetails{}
		e = errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error waiting for vkcs_sharedfilesystem_share_access on %s to be removed: %s: %s", shareID, err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error waiting for vkcs_sharedfilesystem_share_access on %s to be removed: %s (%d): %s", shareID, k, msg.Code, msg.Message)
		}
	}

	log.Printf("[DEBUG] Waiting for vkcs_sharedfilesystem_share_access %s to become denied.", d.Id())
	stateConf := &resource.StateChangeConf{
		Target:     []string{"denied"},
		Pending:    []string{"active", "new", "queued_to_deny", "denying"},
		Refresh:    sharedFilesystemShareAccessStateRefreshFunc(sfsClient, shareID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return nil
		}
		return diag.Errorf("Error waiting for vkcs_sharedfilesystem_share_access %s to become denied: %s", d.Id(), err)
	}

	return nil
}

func resourceSharedFilesystemShareAccessImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		err := fmt.Errorf("Invalid format specified for vkcs_sharedfilesystem_share_access. Format must be <share id>/<ACL id>")
		return nil, err
	}

	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = sharedFilesystemMinMicroversion

	shareID := parts[0]
	accessID := parts[1]

	access, err := shares.ListAccessRights(sfsClient, shareID).Extract()
	if err != nil {
		return nil, fmt.Errorf("Unable to get %s vkcs_sharedfilesystem_share: %s", shareID, err)
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
