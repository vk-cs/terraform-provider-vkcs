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
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
)

func resourceComputeVolumeAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeVolumeAttachCreate,
		ReadContext:   resourceComputeVolumeAttachRead,
		DeleteContext: resourceComputeVolumeAttachDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceComputeVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	var (
		blockStorageClient *gophercloud.ServiceClient
	)

	blockStorageClient, err = config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	volumeID := d.Get("volume_id").(string)

	attachOpts := volumeattach.CreateOpts{
		VolumeID: volumeID,
	}

	log.Printf("[DEBUG] vkcs_compute_volume_attach attach options %s: %#v", instanceID, attachOpts)

	var attachment *volumeattach.VolumeAttachment
	timeout := d.Timeout(schema.TimeoutCreate)
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		attachment, err = volumeattach.Create(computeClient, instanceID, attachOpts).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault400); ok {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating vkcs_compute_volume_attach %s: %s", instanceID, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ATTACHING"},
		Target:     []string{"ATTACHED"},
		Refresh:    computeVolumeAttachAttachFunc(computeClient, blockStorageClient, instanceID, attachment.ID, volumeID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error attaching vkcs_compute_volume_attach %s: %s", instanceID, err)
	}

	// Use the instance ID and attachment ID as the resource ID.
	// This is because an attachment cannot be retrieved just by its ID alone.
	id := fmt.Sprintf("%s/%s", instanceID, attachment.ID)

	d.SetId(id)

	return resourceComputeVolumeAttachRead(ctx, d, meta)
}

func resourceComputeVolumeAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID, attachmentID, err := computeVolumeAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := volumeattach.Get(computeClient, instanceID, attachmentID).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error retrieving vkcs_compute_volume_attach"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_volume_attach %s: %#v", d.Id(), attachment)

	d.Set("instance_id", attachment.ServerID)
	d.Set("volume_id", attachment.VolumeID)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceComputeVolumeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	computeClient, err := config.ComputeV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID, attachmentID, err := computeVolumeAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeVolumeAttachDetachFunc(computeClient, instanceID, attachmentID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error detaching vkcs_compute_volume_attach"))
	}

	return nil
}
