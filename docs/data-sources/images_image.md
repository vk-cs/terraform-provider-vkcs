---
subcategory: "Images"
layout: "vkcs"
page_title: "vkcs: vkcs_images_image"
description: |-
  Get information on an VKCS Image.
---

# vkcs_images_image

Use this data source to get the ID of an available VKCS image.

## Example Usage

```terraform
data "vkcs_images_image" "ubuntu" {
  name        = "Ubuntu 16.04"
  most_recent = true

  properties = {
    key = "value"
  }
}
```

## Argument Reference
- `member_status` optional *string* &rarr;  Status for adding a new member (tenant) to an image member list.

- `most_recent` optional *boolean* &rarr;  If more than one result is returned, use the most recent image.

- `name` optional *string* &rarr;  The name of the image.

- `owner` optional *string* &rarr;  The owner (UUID) of the image.

- `properties` optional *map of* *string* &rarr;  A map of key/value pairs to match an image with. All specified properties must be matched. Unlike other options filtering by `properties` does by client on the result of search query. Filtering is applied if server response contains at least 2 images. In case there is only one image the `properties` ignores.

- `region` optional *string* &rarr;  The region in which to obtain the Image client. An Image client is needed to create an Image that can be used with a compute instance. If omitted, the `region` argument of the provider is used.

- `size_max` optional *number* &rarr;  The maximum size (in bytes) of the image to return.

- `size_min` optional *number* &rarr;  The minimum size (in bytes) of the image to return.

- `tag` optional *string* &rarr;  Search for images with a specific tag.

- `visibility` optional *string* &rarr;  The visibility of the image. Must be one of "private", "community", or "shared". Defaults to "private".


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `checksum` *string* &rarr;  The checksum of the data associated with the image.

- `container_format` *string* &rarr;  The format of the image's container.

- `created_at` *string* &rarr;  The date the image was created.

- `disk_format` *string* &rarr;  The format of the image's disk.

- `file` *string* &rarr;  The trailing path after the endpoint that represent the location of the image or the path to retrieve it.

- `id` *string* &rarr;  ID of the resource.

- `metadata` *map of* *string* &rarr;  The metadata associated with the image. Image metadata allow for meaningfully define the image properties and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.

- `min_disk_gb` *number* &rarr;  The minimum amount of disk space required to use the image.

- `min_ram_mb` *number* &rarr;  The minimum amount of ram required to use the image.

- `protected` *boolean* &rarr;  Whether or not the image is protected.

- `schema` *string* &rarr;  The path to the JSON-schema that represent the image or image

- `size_bytes` *number* &rarr;  The size of the image (in bytes).

- `tags` *set of* *string* &rarr;  The tags list of the image.

- `updated_at` *string* &rarr;  The date the image was last updated.


