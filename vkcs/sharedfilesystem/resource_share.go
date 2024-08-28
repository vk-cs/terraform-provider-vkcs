package sharedfilesystem

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud"
	sfserrors "github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/messages"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	ishares "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/sharedfilesystem/v2/shares"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util/errutil"
)

const (
	// Major share functionality appeared in 2.14.
	minManilaShareMicroversion  = "2.14"
	shareOperationCreateTimeout = 40
	shareOperationUpdateTimeout = 40
	shareOperationDeleteTimeout = 20
)

func ResourceSharedFilesystemShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareCreate,
		ReadContext:   resourceSharedFilesystemShareRead,
		UpdateContext: resourceSharedFilesystemShareUpdate,
		DeleteContext: resourceSharedFilesystemShareDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(shareOperationCreateTimeout * time.Minute),
			Update: schema.DefaultTimeout(shareOperationUpdateTimeout * time.Minute),
			Delete: schema.DefaultTimeout(shareOperationDeleteTimeout * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Shared File System client. A Shared File System client is needed to create a share. Changing this creates a new share.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The owner of the Share.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the share. Changing this updates the name of the existing share.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The human-readable description for the share. Changing this updates the description of the existing share.",
			},

			"export_location_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The export location path of the share.",
			},

			"share_proto": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NFS", "CIFS", "CEPHFS", "GLUSTERFS", "HDFS", "MAPRFS",
				}, true),
				Description: "The share protocol - can either be NFS, CIFS, CEPHFS, GLUSTERFS, HDFS or MAPRFS. Changing this creates a new share.",
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
				Description:  "The share size, in GBs. The requested share size cannot be greater than the allowed GB quota. Changing this resizes the existing share.",
			},

			"share_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The share type name. If you omit this parameter, the default share type is used.",
			},

			"snapshot_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The UUID of the share's base snapshot. Changing this creates a new share.",
			},

			"share_network_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the share network.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The share availability zone. Changing this creates a new share.",
			},

			"share_server_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the share server.",
			},

			"all_metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The map of metadata, assigned on the share, which has been explicitly and implicitly added.",
			},
		},
		Description: "Use this resource to configure a share.",
	}
}

func resourceSharedFilesystemShareCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
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
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		share, err = ishares.Create(sfsClient, createOpts).Extract()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		detailedErr := sfserrors.ErrorDetails{}
		e := sfserrors.ExtractErrorInto(err, &detailedErr)
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
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	share, err := ishares.Get(sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "share"))
	}

	log.Printf("[DEBUG] Retrieved share %s: %#v", d.Id(), share)

	exportLocationPath, err := getShareExportLocationPath(sfsClient, share.ID)
	if err != nil {
		return diag.Errorf("Error retrieving share export location path: %s", err)
	}

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
	d.Set("region", util.GetRegion(d, config))
	d.Set("project_id", share.ProjectID)
	d.Set("export_location_path", exportLocationPath)
	d.Set("share_server_id", share.ShareServerID)

	return nil
}

func resourceSharedFilesystemShareUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
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
		err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, err := ishares.Update(sfsClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return util.CheckForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			detailedErr := sfserrors.ErrorDetails{}
			e := sfserrors.ExtractErrorInto(err, &detailedErr)
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
			err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
				err := ishares.Extend(sfsClient, d.Id(), resizeOpts).Err
				log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
				if err != nil {
					return util.CheckForRetryableError(err)
				}
				return nil
			})
		} else if newSize.(int) < oldSize.(int) {
			pending = append(pending, "shrinking")
			resizeOpts := shares.ShrinkOpts{NewSize: newSize.(int)}
			log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
			err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
				err := ishares.Shrink(sfsClient, d.Id(), resizeOpts).Err
				log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)
				if err != nil {
					return util.CheckForRetryableError(err)
				}
				return nil
			})
		}

		if err != nil {
			detailedErr := sfserrors.ErrorDetails{}
			e := sfserrors.ExtractErrorInto(err, &detailedErr)
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
	config := meta.(clients.Config)
	sfsClient, err := config.SharedfilesystemV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS sharedfilesystem client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete share %s", d.Id())
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = ishares.Delete(sfsClient, d.Id()).ExtractErr()
		if err != nil {
			return util.CheckForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		e := util.CheckDeleted(d, err, "")
		if e == nil {
			return nil
		}
		detailedErr := sfserrors.ErrorDetails{}
		e = sfserrors.ExtractErrorInto(err, &detailedErr)
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

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceSFShareRefreshFunc(sfsClient, id),
		Timeout:    timeout,
		Delay:      1 * time.Second,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if errutil.IsNotFound(err) {
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

func resourceSFShareRefreshFunc(sfsClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		share, err := ishares.Get(sfsClient, id).Extract()
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
