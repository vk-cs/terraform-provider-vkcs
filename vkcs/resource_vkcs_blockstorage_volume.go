package vkcs

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
)

const (
	bsVolumeCreateTimeout = 30 * time.Minute
	bsVolumeDelay         = 10 * time.Second
	bsVolumeMinTimeout    = 3 * time.Second
)

var (
	bsVolumeStatusBuild       = "creating"
	bsVolumeStatusActive      = "available"
	bsVolumeStatusInUse       = "in-use"
	bsVolumeStatusRetype      = "retyping"
	bsVolumeStatusShutdown    = "deleting"
	bsVolumeStatusDeleted     = "deleted"
	bsVolumeStatusDownloading = "downloading"
	bsVolumeMigrationPolicy   = "on-demand"
)

func resourceBlockStorageVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageVolumeCreate,
		ReadContext:   resourceBlockStorageVolumeRead,
		UpdateContext: resourceBlockStorageVolumeUpdate,
		DeleteContext: resourceBlockStorageVolumeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(bsVolumeCreateTimeout),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Region to create resource in.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    false,
				Description: "The name of the volume.",
			},

			"size": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    false,
				Description: "The size of the volume.",
			},

			"volume_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The type of the volume.",
			},

			"availability_zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The name of the availability zone of the volume.",
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
				Computed:    true,
				Description: "Map of key-value metadata of the volume.",
			},

			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_vol_id", "image_id"},
				Description:   "ID of the snapshot of volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.",
			},

			"source_vol_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "image_id"},
				Description:   "ID of the source volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.",
			},

			"image_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "source_vol_id"},
				Description:   "ID of the image to create volume with. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.",
			},
		},
		Description: "Provides a blockstorage volume resource. This can be used to create, modify and delete blockstorage volume.",
	}
}

func resourceBlockStorageVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	metadata := d.Get("metadata").(map[string]interface{})
	createOpts := &volumes.CreateOpts{
		AvailabilityZone: d.Get("availability_zone").(string),
		VolumeType:       d.Get("volume_type").(string),
		Size:             d.Get("size").(int),
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		SnapshotID:       d.Get("snapshot_id").(string),
		SourceVolID:      d.Get("source_vol_id").(string),
		ImageID:          d.Get("image_id").(string),
		Metadata:         expandToMapStringString(metadata),
	}

	log.Printf("[DEBUG] vkcs_blockstorage_volume create options: %#v", createOpts)

	v, err := volumes.Create(blockStorageClient, createOpts).Extract()

	if err != nil {
		return diag.Errorf("error creating vkcs_blockstorage_volume: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{bsVolumeStatusBuild, bsVolumeStatusDownloading},
		Target:     []string{bsVolumeStatusActive},
		Refresh:    blockStorageVolumeStateRefreshFunc(blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      bsVolumeDelay,
		MinTimeout: bsVolumeMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_blockstorage_volume %s to become ready: %s", v.ID, err)
	}

	// Store the ID now
	d.SetId(v.ID)

	return resourceBlockStorageVolumeRead(ctx, d, meta)
}

func resourceBlockStorageVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	v, err := volumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "error retrieving vkcs_blockstorage_volume"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_blockstorage_volume %s: %#v", d.Id(), v)

	d.Set("name", v.Name)
	d.Set("size", v.Size)
	d.Set("volume_type", v.VolumeType)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("description", v.Description)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("metadata", v.Metadata)
	d.Set("region", getRegion(d, config))

	return nil
}

func resourceBlockStorageVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateOpts := volumes.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	if d.HasChange("metadata") {
		metadata := d.Get("metadata").(map[string]interface{})
		updateOpts.Metadata = expandToMapStringString(metadata)
	}

	_, err = volumes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating vkcs_blockstorage_volume")
	}

	if d.HasChange("size") {
		extendOpts := volumeactions.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}

		err = volumeactions.ExtendSize(blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error extending vkcs_blockstorage_volume %s size: %s", d.Id(), err)
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{bsVolumeStatusBuild, bsVolumeStatusDownloading},
			Target:     []string{bsVolumeStatusActive, bsVolumeStatusInUse},
			Refresh:    blockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      bsVolumeDelay,
			MinTimeout: bsVolumeMinTimeout,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_blockstorage_volume %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("volume_type") || d.HasChange("availability_zone") {
		changeTypeOpts := volumeChangeTypeOpts{
			ChangeTypeOpts: volumeactions.ChangeTypeOpts{
				NewType:         d.Get("volume_type").(string),
				MigrationPolicy: volumeactions.MigrationPolicy(bsVolumeMigrationPolicy),
			},
		}
		if d.HasChange("availability_zone") {
			changeTypeOpts.AvailabilityZone = d.Get("availability_zone").(string)
		}
		err = volumeactions.ChangeType(blockStorageClient, d.Id(), changeTypeOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error changing type of vkcs_blockstorage_volume %s", d.Id())
		}
		stateConf := &resource.StateChangeConf{
			Pending:    []string{bsVolumeStatusBuild, bsVolumeStatusRetype},
			Target:     []string{bsVolumeStatusActive, bsVolumeStatusInUse},
			Refresh:    blockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      bsVolumeDelay,
			MinTimeout: bsVolumeMinTimeout,
		}
		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf("error waiting for vkcs_blockstorage_volume %s to become ready: %s", d.Id(), err)
		}
	}

	return resourceBlockStorageVolumeRead(ctx, d, meta)
}

func resourceBlockStorageVolumeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	err = volumes.Delete(blockStorageClient, d.Id(), nil).ExtractErr()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "Error deleting vkcs_blockstorage_volume"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{bsVolumeStatusActive, bsVolumeStatusShutdown, bsVolumeStatusInUse},
		Target:     []string{bsVolumeStatusDeleted},
		Refresh:    blockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      bsVolumeDelay,
		MinTimeout: bsVolumeMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_blockstorage_volume %s to delete : %s", d.Id(), err)
	}

	return nil
}
