---
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
- `member_status` **String** (*Optional*) Status for adding a new member (tenant) to an image member list.

- `most_recent` **Boolean** (*Optional*) If more than one result is returned, use the most recent image.

- `name` **String** (*Optional*) The name of the image.

- `owner` **String** (*Optional*) The owner (UUID) of the image.

- `properties` <strong>Map of </strong>**String** (*Optional*) A map of key/value pairs to match an image with. All specified properties must be matched. Unlike other options filtering by `properties` does by client on the result of search query. Filtering is applied if server response contains at least 2 images. In case there is only one image the `properties` ignores.

- `region` **String** (*Optional*) The region in which to obtain the Image client. An Image client is needed to create an Image that can be used with a compute instance. If omitted, the `region` argument of the provider is used.

- `size_max` **Number** (*Optional*) The maximum size (in bytes) of the image to return.

- `size_min` **Number** (*Optional*) The minimum size (in bytes) of the image to return.

- `sort_direction` **String** (*Optional*) Order the results in either `asc` or `desc`.

- `sort_key` **String** (*Optional*) Sort images based on a certain key. Defaults to `name`.

- `tag` **String** (*Optional*) Search for images with a specific tag.

- `visibility` **String** (*Optional*) The visibility of the image. Must be one of "private", "community", or "shared". Defaults to "private".


## Attributes Reference
- `member_status` **String** See Argument Reference above.

- `most_recent` **Boolean** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `owner` **String** See Argument Reference above.

- `properties` <strong>Map of </strong>**String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `size_max` **Number** See Argument Reference above.

- `size_min` **Number** See Argument Reference above.

- `sort_direction` **String** See Argument Reference above.

- `sort_key` **String** See Argument Reference above.

- `tag` **String** See Argument Reference above.

- `visibility` **String** See Argument Reference above.

- `checksum` **String** The checksum of the data associated with the image.

- `container_format` **String** The format of the image's container.

- `created_at` **String** The date the image was created.

- `disk_format` **String** The format of the image's disk.

- `file` **String** The trailing path after the endpoint that represent the location of the image or the path to retrieve it.

- `id` **String** ID of the resource.

- `metadata` <strong>Map of </strong>**String** The metadata associated with the image. Image metadata allow for meaningfully define the image properties and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.

- `min_disk_gb` **Number** The minimum amount of disk space required to use the image.

- `min_ram_mb` **Number** The minimum amount of ram required to use the image.

- `protected` **Boolean** Whether or not the image is protected.

- `schema` **String** The path to the JSON-schema that represent the image or image

- `size_bytes` **Number** The size of the image (in bytes).

- `tags` <strong>Set of </strong>**String** The tags list of the image.

- `updated_at` **String** The date the image was last updated.


