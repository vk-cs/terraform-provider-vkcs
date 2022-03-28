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
)

const (
	bsVolumeCreateTimeout = 30 * time.Minute
	bsVolumeDelay         = 10 * time.Second
	bsVolumeMinTimeout    = 3 * time.Second
)

var (
	bsVolumeStatusBuild     = "creating"
	bsVolumeStatusActive    = "available"
	bsVolumeStatusInUse     = "in-use"
	bsVolumeStatusShutdown  = "deleting"
	bsVolumeStatusDeleted   = "deleted"
	bsVolumeMigrationPolicy = "on-demand"
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: false,
			},

			"volume_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_vol_id", "image_id"},
			},

			"source_vol_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "image_id"},
			},

			"image_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "source_vol_id"},
			},
		},
	}
}

func resourceBlockStorageVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
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
		Pending:    []string{bsVolumeStatusBuild},
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
	config := meta.(configer)
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
	config := meta.(configer)
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
			Pending:    []string{bsVolumeStatusBuild},
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
			Pending:    []string{bsVolumeStatusBuild},
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
	config := meta.(configer)
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
