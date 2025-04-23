package blockstorage

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	ivolumeactions "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumeactions"
	ivolumes "github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/services/blockstorage/v3/volumes"

	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"
)

const (
	bsVolumeCreateTimeout = 30 * time.Minute
	bsVolumeDelay         = 10 * time.Second
	bsVolumeMinTimeout    = 3 * time.Second
)

var (
	bsVolumeStatusBuild       = "creating"
	BSVolumeStatusActive      = "available"
	BSVolumeStatusInUse       = "in-use"
	BSVolumeStatusDetaching   = "detaching"
	bsVolumeStatusRetype      = "retyping"
	bsVolumeStatusExtending   = "extending"
	bsVolumeStatusAttaching   = "attaching"
	bsVolumeStatusShutdown    = "deleting"
	BSVolumeStatusDeleted     = "deleted"
	bsVolumeStatusDownloading = "downloading"
	bsVolumeMigrationPolicy   = "on-demand"
)

func ResourceBlockStorageVolume() *schema.Resource {
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
				Description: "Key-value map to configure metadata of the volume. _note_ Changes to keys that are not in scope, i.e. not configured here, will not be reflected in planned changes, if any, so those keys can be `silently` removed during an update.",
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

			"all_metadata": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed:    true,
				Description: "Map of key-value metadata of the volume.",
			},
		},
		Description: "Provides a blockstorage volume resource. This can be used to create, modify and delete blockstorage volume.",
	}
}

func resourceBlockStorageVolumeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
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
		Metadata:         util.ExpandToMapStringString(metadata),
	}

	log.Printf("[DEBUG] vkcs_blockstorage_volume create options: %#v", createOpts)

	v, err := ivolumes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating vkcs_blockstorage_volume: %s", err)
	}

	// Store the ID now
	d.SetId(v.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{bsVolumeStatusBuild, bsVolumeStatusDownloading},
		Target:     []string{BSVolumeStatusActive},
		Refresh:    BlockStorageVolumeStateRefreshFunc(blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      bsVolumeDelay,
		MinTimeout: bsVolumeMinTimeout,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for vkcs_blockstorage_volume %s to become ready: %s", v.ID, err)
	}

	return resourceBlockStorageVolumeRead(ctx, d, meta)
}

func resourceBlockStorageVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	v, err := ivolumes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "error retrieving vkcs_blockstorage_volume"))
	}

	log.Printf("[DEBUG] Retrieved vkcs_blockstorage_volume %s: %#v", d.Id(), v)

	d.Set("name", v.Name)
	d.Set("size", v.Size)
	d.Set("volume_type", v.VolumeType)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("description", v.Description)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("region", util.GetRegion(d, config))
	d.Set("all_metadata", v.Metadata)

	configMetadata := d.Get("metadata").(map[string]any)
	intersectionMetadata := make(map[string]string, len(configMetadata))
	for key := range configMetadata {
		if _, exist := v.Metadata[key]; exist {
			intersectionMetadata[key] = v.Metadata[key]
		}
	}
	d.Set("metadata", intersectionMetadata)

	return nil
}

func resourceBlockStorageVolumeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	metadata := util.ExpandToMapStringString(d.Get("metadata").(map[string]any))
	updateOpts := ivolumes.UpdateOpts{
		Name:        &name,
		Description: &description,
		Metadata:    metadata,
	}

	_, err = ivolumes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating vkcs_blockstorage_volume")
	}
	d.Set("metadata", metadata)

	if d.HasChange("size") {
		extendOpts := volumeactions.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}

		err = ivolumeactions.ExtendSize(blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error extending vkcs_blockstorage_volume %s size: %s", d.Id(), err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{bsVolumeStatusBuild, bsVolumeStatusDownloading, bsVolumeStatusExtending},
			Target:     []string{BSVolumeStatusActive, BSVolumeStatusInUse},
			Refresh:    BlockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
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
		err = ivolumeactions.ChangeType(blockStorageClient, d.Id(), changeTypeOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("error changing type of vkcs_blockstorage_volume %s", d.Id())
		}
		stateConf := &retry.StateChangeConf{
			Pending:    []string{bsVolumeStatusBuild, bsVolumeStatusRetype, bsVolumeStatusAttaching},
			Target:     []string{BSVolumeStatusActive, BSVolumeStatusInUse},
			Refresh:    BlockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
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
	blockStorageClient, err := config.BlockStorageV3Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS block storage client: %s", err)
	}

	err = ivolumes.Delete(blockStorageClient, d.Id(), nil).ExtractErr()
	if err != nil {
		return diag.FromErr(util.CheckDeleted(d, err, "Error deleting vkcs_blockstorage_volume"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{BSVolumeStatusActive, bsVolumeStatusShutdown},
		Target:     []string{BSVolumeStatusDeleted},
		Refresh:    BlockStorageVolumeStateRefreshFunc(blockStorageClient, d.Id()),
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
