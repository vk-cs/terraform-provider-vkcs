package vkcs

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/messages"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

const (
	// Major share functionality appeared in 2.14.
	minManilaShareMicroversion = "2.14"
)

func resourceSharedFilesystemShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareCreate,
		ReadContext:   resourceSharedFilesystemShareRead,
		UpdateContext: resourceSharedFilesystemShareUpdate,
		DeleteContext: resourceSharedFilesystemShareDelete,
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
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"share_proto": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NFS", "CIFS", "CEPHFS", "GLUSTERFS", "HDFS", "MAPRFS",
				}, true),
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"share_network_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"share_server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"all_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceSharedFilesystemShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	createOpts := shares.CreateOpts{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		ShareProto:       d.Get("share_proto").(string),
		Size:             d.Get("size").(int),
		SnapshotID:       d.Get("snapshot_id").(string),
		ShareNetworkID:   d.Get("share_network_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	if v, ok := d.GetOk("share_type"); ok {
		createOpts.ShareType = v.(string)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	log.Printf("[DEBUG] Attempting to create share")
	var share *shares.Share
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		share, err = shares.Create(sfsClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		detailedErr := errors.ErrorDetails{}
		e := errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error creating share: %s: %s", err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Error creating share: %s (%d): %s", k, msg.Code, msg.Message)
		}
	}

	d.SetId(share.ID)

	// Wait for share to become active before continuing
	err = waitForSFShare(ctx, sfsClient, share.ID, "available", []string{"creating", "manage_starting"}, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSharedFilesystemShareRead(ctx, d, meta)
}

func resourceSharedFilesystemShareRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	share, err := shares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "share"))
	}

	log.Printf("[DEBUG] Retrieved share %s: %#v", d.Id(), share)

	d.Set("name", share.Name)
	d.Set("description", share.Description)
	d.Set("share_proto", share.ShareProto)
	d.Set("size", share.Size)
	d.Set("share_type", share.ShareTypeName)
	d.Set("snapshot_id", share.SnapshotID)
	d.Set("all_metadata", share.Metadata)
	d.Set("share_network_id", share.ShareNetworkID)
	d.Set("availability_zone", share.AvailabilityZone)
	// Computed
	d.Set("region", getRegion(d, config))
	d.Set("project_id", share.ProjectID)
	d.Set("share_server_id", share.ShareServerID)

	return nil
}

func resourceSharedFilesystemShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	timeout := d.Timeout(schema.TimeoutUpdate)

	var updateOpts shares.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.DisplayName = &name
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.DisplayDescription = &description
	}

	if updateOpts != (shares.UpdateOpts{}) {
		// Wait for share to become active before continuing
		err = waitForSFShare(ctx, sfsClient, d.Id(), "available", []string{"creating", "manage_starting", "extending", "shrinking"}, timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Attempting to update share")
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err := shares.Update(sfsClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			detailedErr := errors.ErrorDetails{}
			e := errors.ExtractErrorInto(err, &detailedErr)
			if e != nil {
				return diag.Errorf("Error updating %s share: %s: %s", d.Id(), err, e)
			}
			for k, msg := range detailedErr {
				return diag.Errorf("Error updating %s share: %s (%d): %s", d.Id(), k, msg.Code, msg.Message)
			}
		}

		// Wait for share to become active before continuing
		err = waitForSFShare(ctx, sfsClient, d.Id(), "available", []string{"creating", "manage_starting", "extending", "shrinking"}, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("size") {
		var pending []string
		oldSize, newSize := d.GetChange("size")

		if newSize.(int) > oldSize.(int) {
			pending = append(pending, "extending")
			resizeOpts := shares.ExtendOpts{NewSize: newSize.(int)}
			log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
			err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
				err := shares.Extend(sfsClient, d.Id(), resizeOpts).Err
				log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
				if err != nil {
					return checkForRetryableError(err)
				}
				return nil
			})
		} else if newSize.(int) < oldSize.(int) {
			pending = append(pending, "shrinking")
			resizeOpts := shares.ShrinkOpts{NewSize: newSize.(int)}
			log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
			err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
				err := shares.Shrink(sfsClient, d.Id(), resizeOpts).Err
				log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
				if err != nil {
					return checkForRetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			detailedErr := errors.ErrorDetails{}
			e := errors.ExtractErrorInto(err, &detailedErr)
			if e != nil {
				return diag.Errorf("Unable to resize %s share: %s: %s", d.Id(), err, e)
			}
			for k, msg := range detailedErr {
				return diag.Errorf("Unable to resize %s share: %s (%d): %s", d.Id(), k, msg.Code, msg.Message)
			}
		}

		// Wait for share to become active before continuing
		err = waitForSFShare(ctx, sfsClient, d.Id(), "available", pending, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceSharedFilesystemShareRead(ctx, d, meta)
}

func resourceSharedFilesystemShareDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config)
	sfsClient, err := config.SharedfilesystemV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete share %s", d.Id())
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err = shares.Delete(sfsClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		e := checkDeleted(d, err, "")
		if e == nil {
			return nil
		}
		detailedErr := errors.ErrorDetails{}
		e = errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Unable to delete %s share: %s: %s", d.Id(), err, e)
		}
		for k, msg := range detailedErr {
			return diag.Errorf("Unable to delete %s share: %s (%d): %s", d.Id(), k, msg.Code, msg.Message)
		}
	}

	// Wait for share to become deleted before continuing
	pending := []string{"", "deleting", "available"}
	err = waitForSFShare(ctx, sfsClient, d.Id(), "deleted", pending, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Full list of the share statuses: https://developer.vkcs.org/api-ref/shared-file-system/#shares
func waitForSFShare(ctx context.Context, sfsClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for share %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceSFShareRefreshFunc(sfsClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			switch target {
			case "deleted":
				return nil
			default:
				return fmt.Errorf("error: share %s not found: %s", id, err)
			}
		}
		errorMessage := fmt.Sprintf("error waiting for share %s to become %s", id, target)
		msg := resourceSFSShareManilaMessage(sfsClient, id)
		if msg == nil {
			return fmt.Errorf("%s: %s", errorMessage, err)
		}
		return fmt.Errorf("%s: %s: the latest manila message (%s): %s", errorMessage, err, msg.CreatedAt, msg.UserMessage)
	}

	return nil
}

func resourceSFShareRefreshFunc(sfsClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		share, err := shares.Get(sfsClient, id).Extract()
		if err != nil {
			return nil, "", err
		}
		return share, share.Status, nil
	}
}

func resourceSFSShareManilaMessage(sfsClient *gophercloud.ServiceClient, id string) *messages.Message {
	// we can simply set this, because this function is called after the error occurred
	sfsClient.Microversion = "2.37"

	listOpts := messages.ListOpts{
		ResourceID: id,
		SortKey:    "created_at",
		SortDir:    "desc",
		Limit:      1,
	}
	allPages, err := messages.List(sfsClient, listOpts).AllPages()
	if err != nil {
		log.Printf("[DEBUG] Unable to retrieve messages: %v", err)
		return nil
	}

	allMessages, err := messages.ExtractMessages(allPages)
	if err != nil {
		log.Printf("[DEBUG] Unable to extract messages: %v", err)
		return nil
	}

	if len(allMessages) == 0 {
		log.Printf("[DEBUG] No messages found")
		return nil
	}

	return &allMessages[0]
}
