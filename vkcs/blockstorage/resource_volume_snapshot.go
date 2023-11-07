package blockstorage

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/snapshots"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	isnapshots "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/snapshots"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	bsSnapshotCreateTimeout = 30 * time.Minute
	bsSnapshotDelay         = 10 * time.Second
	bsSnapshotMinTimeout    = 3 * time.Second
)

var (
	bsSnapshotStatusBuild    = "creating"
	bsSnapshotStatusActive   = "available"
	bsSnapshotStatusShutdown = "deleting"
	bsSnapshotStatusDeleted  = "deleted"
)

func ResourceBlockStorageSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageSnapshotCreate,
		ReadContext:   resourceBlockStorageSnapshotRead,
		UpdateContext: resourceBlockStorageSnapshotUpdate,
		DeleteContext: resourceBlockStorageSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(bsSnapshotCreateTimeout),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"volume_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the volume to create snapshot for. Changing this creates a new snapshot.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The name of the snapshot.",
			},

			"force": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "Allows or disallows snapshot of a volume when the volume is attached to an instance.",
			},

			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The description of the volume.",
			},

			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    false,
				Description: "Map of key-value metadata of the volume.",
			},
		},
		Description: "Provides a blockstorage snapshot resource. This can be used to create, modify and delete blockstorage snapshot.",
	}
}

func resourceBlockStorageSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS block storage client: %s", err)
	}

	metadata := d.Get("metadata").(map[string]interface{})
	createOpts := &snapshots.CreateOpts{
		VolumeID:    d.Get("volume_id").(string),
		Name:        d.Get("name").(string),
		Force:       d.Get("force").(bool),
		Description: d.Get("description").(string),
		Metadata:    util.ExpandToMapStringString(metadata),
	}

	log.Printf("[DEBUG] vkcs_blockstorage_snapshot create options: %#v", createOpts)

	snapshot, err := isnapshots.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_blockstorage_snapshot: %s", err)
	}

	// Store the ID now
	d.SetId(snapshot.ID)

	// Wait for the volume snapshot to become available.
	log.Printf("[DEBUG] Waiting for vkcs_blockstorage_volume_snapshot %s to become available", snapshot.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{bsSnapshotStatusBuild},
		Target:     []string{bsSnapshotStatusActive},
		Refresh:    blockStorageSnapshotStateRefreshFunc(blockStorageClient, snapshot.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      bsSnapshotDelay,
		MinTimeout: bsSnapshotMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_blockstorage_volume_snapshot %s to become ready: %s", snapshot.ID, err)
	}

	return resourceBlockStorageSnapshotRead(ctx, d, meta)
}

func resourceBlockStorageSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS block storage client: %s", err)
	}

	snapshot, err := isnapshots.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "error retrieving vkcs_blockstorage_snapshot"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_blockstorage_snapshot %s: %#v", d.Id(), snapshot)

	d.Set("name", snapshot.Name)
	d.Set("description", snapshot.Description)
	d.Set("volume_id", snapshot.VolumeID)
	d.Set("metadata", snapshot.Metadata)

	return nil
}

func resourceBlockStorageSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateOpts := snapshots.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	_, err = isnapshots.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating vkcs_blockstorage_snapshot")
	}
	stateConf := &retry.StateChangeConf{
		Pending:    []string{bsSnapshotStatusBuild},
		Target:     []string{bsSnapshotStatusActive},
		Refresh:    blockStorageSnapshotStateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      bsSnapshotDelay,
		MinTimeout: bsSnapshotMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_blockstorage_snapshot %s to become ready: %s", d.Id(), err)
	}

	if d.HasChange("metadata") {
		updateMetadataOpts := snapshots.UpdateMetadataOpts{
			Metadata: d.Get("metadata").(map[string]interface{}),
		}

		_, err = isnapshots.UpdateMetadata(blockStorageClient, d.Id(), updateMetadataOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating vkcs_blockstorage_snapshot metadata")
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{bsSnapshotStatusBuild},
			Target:     []string{bsSnapshotStatusActive},
			Refresh:    blockStorageSnapshotStateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      bsSnapshotDelay,
			MinTimeout: bsSnapshotMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_blockstorage_snapshot %s to become ready: %s", d.Id(), err)
		}
	}

	return resourceBlockStorageSnapshotRead(ctx, d, meta)
}

func resourceBlockStorageSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS block storage client: %s", err)
	}
	err = isnapshots.Delete(blockStorageClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_blockstorage_snapshot"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{bsSnapshotStatusActive, bsSnapshotStatusShutdown},
		Target:     []string{bsSnapshotStatusDeleted},
		Refresh:    BlockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      bsSnapshotDelay,
		MinTimeout: bsSnapshotMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_blockstorage_snapshot %s to delete : %s", d.Id(), err)
	}

	return nil
}
