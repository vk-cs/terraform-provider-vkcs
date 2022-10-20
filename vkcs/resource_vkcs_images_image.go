package vkcs

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

const (
	storeS3 = "s3"
)

func resourceImagesImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceImagesImageCreate,
		ReadContext:   resourceImagesImageRead,
		UpdateContext: resourceImagesImageUpdate,
		DeleteContext: resourceImagesImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: resourceImagesImageUpdateComputedAttributes,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Image client. An Image client is needed to create an Image that can be used with a compute instance. If omitted, the `region` argument of the provider is used. Changing this creates a new Image.",
			},

			"container_format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"bare"}, false),
				Description:  "The container format. Must be one of \"bare\".",
			},

			"disk_format": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"raw", "iso"}, false),
				Description:  "The disk format. Must be one of \"raw\", \"iso\".",
			},

			"file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The trailing path after the image endpoint that represent the location of the image or the path to retrieve it.",
			},

			"image_cache_path": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     fmt.Sprintf("%s/.terraform/image_cache", os.Getenv("HOME")),
				Description: "This is the directory where the images will be downloaded. Images will be stored with a filename corresponding to the url's md5 hash. Defaults to \"$HOME/.terraform/image_cache\"",
			},

			"image_source_url": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"local_file_path"},
				Description:   "This is the url of the raw image. The image will be downloaded in the `image_cache_path` before being uploaded to VKCS. Conflicts with `local_file_path`.",
			},

			"image_source_username": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"local_file_path"},
				Description:   "The username of basic auth to download `image_source_url`.",
			},

			"image_source_password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"local_file_path"},
				Description:   "The password of basic auth to download `image_source_url`.",
			},

			"local_file_path": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"image_source_url"},
				Description:   "This is the filepath of the raw image file that will be uploaded to VKCS. Conflicts with `image_source_url`",
			},

			"min_disk_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
				Description:  "Amount of disk space (in GB) required to boot image. Defaults to 0.",
			},

			"min_ram_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
				Default:      0,
				Description:  "Amount of ram (in MB) required to boot image. Defauts to 0.",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "The name of the image.",
			},

			"protected": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, image will not be deletable. Defaults to false.",
			},

			"tags": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The tags of the image. It must be a list of strings. At this time, it is not possible to delete all tags of an image.",
			},

			"verify_checksum": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    false,
				Description: "If false, the checksum will not be verified once the image is finished uploading.",
			},

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				ValidateFunc: validation.StringInSlice([]string{
					"private", "shared", "community",
				}, false),
				Default:     "private",
				Description: "The visibility of the image. Must be one of \"private\", \"community\", or \"shared\". The ability to set the visibility depends upon the configuration of the VKCS cloud.",
			},

			"properties": {
				Type:         schema.TypeMap,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateStoreInProperties,
				Description:  "A map of key/value pairs to set freeform information about an image. See the \"Notes\" section for further information about properties.",
			},

			// Computed-only
			"checksum": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The checksum of the data associated with the image.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the image was created.",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The metadata associated with the image. Image metadata allow for meaningfully define the image properties and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.",
			},

			"owner": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The id of the vkcs user who owns the image.",
			},

			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path to the JSON-schema that represent the image or image",
			},

			"size_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size in bytes of the data associated with the image.",
			},

			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the image. It can be \"queued\", \"active\" or \"saving\".",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the image was last updated.",
			},
		},
		Description: "Manages an Image resource within VKCS.\n\n" +
			"~> **Note:** All arguments including the source image URL password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).",
	}
}

func resourceImagesImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	imageClient, err := config.ImageV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS image client: %s", err)
	}

	protected := d.Get("protected").(bool)
	visibility := resourceImagesImageVisibilityFromString(d.Get("visibility").(string))

	properties := d.Get("properties").(map[string]interface{})
	imageProperties := resourceImagesImageExpandProperties(properties)
	if !resourceImagesImageNeedsDefaultStore(imageClient.Endpoint) {
		imageProperties["store"] = storeS3
	}

	createOpts := &images.CreateOpts{
		Name:            d.Get("name").(string),
		ContainerFormat: d.Get("container_format").(string),
		DiskFormat:      d.Get("disk_format").(string),
		MinDisk:         d.Get("min_disk_gb").(int),
		MinRAM:          d.Get("min_ram_mb").(int),
		Protected:       &protected,
		Visibility:      &visibility,
		Properties:      imageProperties,
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = resourceImagesImageBuildTags(tags)
	}

	d.Partial(true)

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	newImg, err := images.Create(imageClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating Image: %s", err)
	}

	d.SetId(newImg.ID)

	var fileChecksum string

	// variable declaration
	var imgFilePath string
	var fileSize int64
	var imgFile *os.File

	// downloading/getting image file props
	imgFilePath, err = resourceImagesImageFile(imageClient, d)
	if err != nil {
		return diag.Errorf("Error opening file for Image: %s", err)
	}
	fileSize, fileChecksum, err = resourceImagesImageFileProps(imgFilePath)
	if err != nil {
		return diag.Errorf("Error getting file props: %s", err)
	}

	// upload
	imgFile, err = os.Open(imgFilePath)
	if err != nil {
		return diag.Errorf("Error opening file %q: %s", imgFilePath, err)
	}
	defer imgFile.Close()
	log.Printf("[WARN] Uploading image %s (%d bytes). This can be pretty long.", d.Id(), fileSize)

	res := imagedata.Upload(imageClient, d.Id(), imgFile)
	if res.Err != nil {
		return diag.Errorf("Error while uploading file %q: %s", imgFilePath, res.Err)
	}

	// wait for active
	stateConf := &resource.StateChangeConf{
		Pending:    []string{string(images.ImageStatusQueued), string(images.ImageStatusSaving), string(images.ImageStatusImporting)},
		Target:     []string{string(images.ImageStatusActive)},
		Refresh:    resourceImagesImageRefreshFunc(imageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Image: %s", err)
	}

	img, err := images.Get(imageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "image"))
	}

	if v, ok := d.GetOk("verify_checksum"); !ok || (ok && v.(bool)) {
		if img.Checksum != fileChecksum {
			return diag.Errorf("Error wrong checksum: got %q, expected %q", img.Checksum, fileChecksum)
		}
	}

	d.Partial(false)

	return resourceImagesImageRead(ctx, d, meta)
}

func resourceImagesImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	imageClient, err := config.ImageV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS image client: %s", err)
	}

	img, err := images.Get(imageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(checkDeleted(d, err, "image"))
	}

	log.Printf("[DEBUG] Retrieved Image %s: %#v", d.Id(), img)

	d.Set("owner", img.Owner)
	d.Set("status", img.Status)
	d.Set("file", img.File)
	d.Set("schema", img.Schema)
	d.Set("checksum", img.Checksum)
	d.Set("size_bytes", img.SizeBytes)
	d.Set("metadata", img.Metadata)
	d.Set("created_at", img.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", img.UpdatedAt.Format(time.RFC3339))
	d.Set("container_format", img.ContainerFormat)
	d.Set("disk_format", img.DiskFormat)
	d.Set("min_disk_gb", img.MinDiskGigabytes)
	d.Set("min_ram_mb", img.MinRAMMegabytes)
	d.Set("file", img.File)
	d.Set("name", img.Name)
	d.Set("protected", img.Protected)
	d.Set("size_bytes", img.SizeBytes)
	d.Set("tags", img.Tags)
	d.Set("visibility", img.Visibility)
	d.Set("region", getRegion(d, config))

	properties := resourceImagesImageExpandProperties(img.Properties)
	if err := d.Set("properties", properties); err != nil {
		log.Printf("[WARN] unable to set properties for image %s: %s", img.ID, err)
	}

	return nil
}

func resourceImagesImageUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	imageClient, err := config.ImageV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS image client: %s", err)
	}

	updateOpts := make(images.UpdateOpts, 0)

	if d.HasChange("visibility") {
		visibility := resourceImagesImageVisibilityFromString(d.Get("visibility").(string))
		v := images.UpdateVisibility{Visibility: visibility}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("protected") {
		protected := d.Get("protected").(bool)
		v := images.ReplaceImageProtected{NewProtected: protected}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("min_disk_gb") {
		minDiskGb := d.Get("min_disk_gb").(int)
		v := images.ReplaceImageMinDisk{NewMinDisk: minDiskGb}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("min_ram_mb") {
		minRAMMb := d.Get("min_ram_mb").(int)
		v := images.ReplaceImageMinRam{NewMinRam: minRAMMb}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("name") {
		v := images.ReplaceImageName{NewName: d.Get("name").(string)}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("tags") {
		tags := d.Get("tags").(*schema.Set).List()
		v := images.ReplaceImageTags{
			NewTags: resourceImagesImageBuildTags(tags),
		}
		updateOpts = append(updateOpts, v)
	}

	if d.HasChange("properties") {
		o, n := d.GetChange("properties")
		oldProperties := resourceImagesImageExpandProperties(o.(map[string]interface{}))
		newProperties := resourceImagesImageExpandProperties(n.(map[string]interface{}))

		// Check for new and changed properties
		for newKey, newValue := range newProperties {
			var changed bool

			oldValue, found := oldProperties[newKey]
			if found && (newValue != oldValue) {
				changed = true
			}

			// os_ keys are provided by the VKCS Image service.
			// These are read-only properties that cannot be modified.
			// Ignore them here and let CustomizeDiff handle them.
			if strings.HasPrefix(newKey, "os_") {
				found = true
				changed = false
			}
			// This is a read-only property that cannot be modified.
			// Ignore it here and let CustomizeDiff handle it.
			if newKey == "direct_url" {
				found = true
				changed = false
			}

			if newKey == "locations" {
				found = true
				changed = false
			}

			if !found {
				v := images.UpdateImageProperty{
					Op:    images.AddOp,
					Name:  newKey,
					Value: newValue,
				}

				updateOpts = append(updateOpts, v)
			}

			if found && changed {
				v := images.UpdateImageProperty{
					Op:    images.ReplaceOp,
					Name:  newKey,
					Value: newValue,
				}

				updateOpts = append(updateOpts, v)
			}
		}

		// Check for removed properties
		for oldKey := range oldProperties {
			_, found := newProperties[oldKey]

			if !found {
				v := images.UpdateImageProperty{
					Op:   images.RemoveOp,
					Name: oldKey,
				}

				updateOpts = append(updateOpts, v)
			}
		}
	}

	log.Printf("[DEBUG] Update Options: %#v", updateOpts)

	_, err = images.Update(imageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating image: %s", err)
	}

	return resourceImagesImageRead(ctx, d, meta)
}

func resourceImagesImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(configer)
	imageClient, err := config.ImageV2Client(getRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating VKCS image client: %s", err)
	}

	log.Printf("[DEBUG] Deleting Image %s", d.Id())
	if err := images.Delete(imageClient, d.Id()).Err; err != nil {
		return diag.Errorf("error deleting Image: %s", err)
	}

	d.SetId("")
	return nil
}

func validateStoreInProperties(v interface{}, k string) (ws []string, errors []error) {
	if _, ok := v.(map[string]interface{})["store"]; ok {
		errors = append(errors, fmt.Errorf("error creating Image: set up store disabled"))
	}
	return
}
