package compute

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/volumeattach"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/blockstorage"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ivolumeattach "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/compute/v2/volumeattach"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

func ResourceComputeVolumeAttach() *schema.Resource {
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
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Compute client. A Compute client is needed to create a volume attachment. If omitted, the `region` argument of the provider is used. Changing this creates a new volume attachment.",
			},

			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Instance to attach the Volume to.",
			},

			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the Volume to attach to an Instance.",
			},
		},
		Description: "Attaches a Block Storage Volume to an Instance using the VKCS Compute API.",
	}
}

func resourceComputeVolumeAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	var (
		blockStorageClient *gophercloud.ServiceClient
	)

	blockStorageClient, err = config.BlockStorageV3Client(util.GetRegion(d, config))
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
	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		attachment, err = ivolumeattach.Create(computeClient, instanceID, attachOpts).Extract()
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return diag.Errorf("Error creating vkcs_compute_volume_attach %s: %s", instanceID, err)
	}

	// Use the instance ID and attachment ID as the resource ID.
	// This is because an attachment cannot be retrieved just by its ID alone.
	id := fmt.Sprintf("%s/%s", instanceID, attachment.ID)

	d.SetId(id)

	stateConf := &retry.StateChangeConf{
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

	return resourceComputeVolumeAttachRead(ctx, d, meta)
}

func resourceComputeVolumeAttachRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID, attachmentID, err := ComputeVolumeAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	attachment, err := ivolumeattach.Get(computeClient, instanceID, attachmentID).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error retrieving vkcs_compute_volume_attach"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_compute_volume_attach %s: %#v", d.Id(), attachment)

	d.Set("instance_id", attachment.ServerID)
	d.Set("volume_id", attachment.VolumeID)
	d.Set("region", util.GetRegion(d, config))

	return nil
}

func resourceComputeVolumeAttachDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	computeClient, err := config.ComputeV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS compute client: %s", err)
	}

	instanceID, attachmentID, err := ComputeVolumeAttachParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{""},
		Target:     []string{"DETACHED"},
		Refresh:    computeVolumeAttachDetachFunc(computeClient, instanceID, attachmentID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error detaching vkcs_compute_volume_attach"))
	}

	// Volume may be still in detaching status after detach resource is deleted
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	volumeID := d.Get("volume_id").(string)

	volumeStateConf := &retry.StateChangeConf{
		Pending:    []string{blockstorage.BSVolumeStatusDetaching, blockstorage.BSVolumeStatusInUse},
		Target:     []string{blockstorage.BSVolumeStatusActive},
		Refresh:    blockstorage.BlockStorageVolumeStateRefreshFunc(blockStorageClient, volumeID),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = volumeStateConf.WaitForStateContext(ctx); err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error detaching vkcs_compute_volume_attach"))
	}

	return nil
}
