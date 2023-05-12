package images

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/clients"
	"github.com/vk-cs/terraform-provider-vkcs/vkcs/internal/util"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

func DataSourceImagesImage() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImagesImageRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region in which to obtain the Image client. An Image client is needed to create an Image that can be used with a compute instance. If omitted, the `region` argument of the provider is used.",
			},

			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the image.",
			},

			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.ImageVisibilityPublic),
					string(images.ImageVisibilityPrivate),
					string(images.ImageVisibilityShared),
					string(images.ImageVisibilityCommunity),
				}, false),
				Description: "The visibility of the image. Must be one of \"private\", \"community\", or \"shared\". Defaults to \"private\".",
			},

			"member_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(images.ImageMemberStatusAccepted),
					string(images.ImageMemberStatusPending),
					string(images.ImageMemberStatusRejected),
					string(images.ImageMemberStatusAll),
				}, false),
				Description: "Status for adding a new member (tenant) to an image member list.",
			},

			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The owner (UUID) of the image.",
			},

			"size_min": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The minimum size (in bytes) of the image to return.",
			},

			"size_max": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The maximum size (in bytes) of the image to return.",
			},

			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Search for images with a specific tag.",
			},

			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				ForceNew:    true,
				Description: "If more than one result is returned, use the most recent image.",
			},

			"properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Description: "A map of key/value pairs to match an image with. All specified properties must be matched. Unlike other options filtering by `properties` does by client on the result of search query. Filtering is applied if server response contains at least 2 images. In case there is only one image the `properties` ignores.",
			},

			// Computed values
			"container_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The format of the image's container.",
			},

			"disk_format": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The format of the image's disk.",
			},

			"min_disk_gb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum amount of disk space required to use the image.",
			},

			"min_ram_mb": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum amount of ram required to use the image.",
			},

			"protected": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether or not the image is protected.",
			},

			"checksum": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The checksum of the data associated with the image.",
			},

			"size_bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the image (in bytes).",
			},

			"metadata": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The metadata associated with the image. Image metadata allow for meaningfully define the image properties and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.",
			},

			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the image was created.",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date the image was last updated.",
			},

			"file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The trailing path after the endpoint that represent the location of the image or the path to retrieve it.",
			},

			"schema": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The path to the JSON-schema that represent the image or image",
			},

			"tags": {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The tags list of the image.",
			},
		},
		Description: "Use this data source to get the ID of an available VKCS image.",
	}
}

// dataSourceImagesImageRead performs the image lookup.
func dataSourceImagesImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(clients.Config)
	imageClient, err := config.ImageV2Client(util.GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating VKCS image client: %s", err)
	}

	visibility := resourceImagesImageVisibilityFromString(d.Get("visibility").(string))
	memberStatus := resourceImagesImageMemberStatusFromString(d.Get("member_status").(string))

	var tags []string
	tag := d.Get("tag").(string)
	if tag != "" {
		tags = append(tags, tag)
	}

	listOpts := images.ListOpts{
		Name:         d.Get("name").(string),
		Visibility:   visibility,
		Owner:        d.Get("owner").(string),
		Status:       images.ImageStatusActive,
		SizeMin:      int64(d.Get("size_min").(int)),
		SizeMax:      int64(d.Get("size_max").(int)),
		Tags:         tags,
		MemberStatus: memberStatus,
	}

	log.Printf("[DEBUG] List Options: %#v", listOpts)

	var image images.Image
	allPages, err := images.List(imageClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to query images: %s", err)
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve images: %s", err)
	}

	properties := resourceImagesImageExpandProperties(
		d.Get("properties").(map[string]interface{}))

	if len(allImages) > 1 {
		allImages = imagesFilterByProperties(allImages, properties)

		log.Printf("[DEBUG] Image list filtered by properties: %#v", properties)
	}

	if len(allImages) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allImages) > 1 {
		recent := d.Get("most_recent").(bool)
		log.Printf("[DEBUG] Multiple results found and `most_recent` is set to: %t", recent)
		if recent {
			image = mostRecentImage(allImages)
		} else {
			log.Printf("[DEBUG] Multiple results found: %#v", allImages)
			return diag.Errorf("Your query returned more than one result. Please try a more " +
				"specific search criteria, or set `most_recent` attribute to true.")
		}
	} else {
		image = allImages[0]
	}

	log.Printf("[DEBUG] Single Image found: %s", image.ID)

	log.Printf("[DEBUG] vkcs_images_image details: %#v", image)

	d.SetId(image.ID)
	d.Set("name", image.Name)
	d.Set("region", util.GetRegion(d, config))
	d.Set("tags", image.Tags)
	d.Set("container_format", image.ContainerFormat)
	d.Set("disk_format", image.DiskFormat)
	d.Set("min_disk_gb", image.MinDiskGigabytes)
	d.Set("min_ram_mb", image.MinRAMMegabytes)
	d.Set("owner", image.Owner)
	d.Set("protected", image.Protected)
	d.Set("visibility", image.Visibility)
	d.Set("checksum", image.Checksum)
	d.Set("size_bytes", image.SizeBytes)
	d.Set("metadata", image.Metadata)
	d.Set("created_at", image.CreatedAt.Format(time.RFC3339))
	d.Set("updated_at", image.UpdatedAt.Format(time.RFC3339))
	d.Set("file", image.File)
	d.Set("schema", image.Schema)

	return nil
}

type imageSort []images.Image

func (a imageSort) Len() int      { return len(a) }
func (a imageSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a imageSort) Less(i, j int) bool {
	itime := a[i].CreatedAt
	jtime := a[j].CreatedAt
	return itime.Unix() < jtime.Unix()
}

// Returns the most recent Image out of a slice of images.
func mostRecentImage(images []images.Image) images.Image {
	sortedImages := images
	sort.Sort(imageSort(sortedImages))
	return sortedImages[len(sortedImages)-1]
}
