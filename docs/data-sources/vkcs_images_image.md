---
layout: "vkcs"
page_title: "vkcs: images_image"
description: |-
  Get information on an VKCS Image.
---

# vkcs\_images\_image

Use this data source to get the ID of an available VKCS image.

## Example Usage

```hcl
data "vkcs_images_image" "ubuntu" {
  name        = "Ubuntu 16.04"
  most_recent = true

  properties = {
    key = "value"
  }
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Image client.
    An Image client is needed to create an Image that can be used with
    a compute instance. If omitted, the `region` argument of the provider
    is used.

* `most_recent` - (Optional) If more than one result is returned, use the most
  recent image.

* `name` - (Optional) The name of the image.

* `owner` - (Optional) The owner (UUID) of the image.

* `properties` - (Optional) a map of key/value pairs to match an image with.
    All specified properties must be matched. Unlike other options filtering
    by `properties` does by client on the result of search query.
    Filtering is applied if server responce contains at least 2 images. In
    case there is only one image the `properties` ignores.

* `size_min` - (Optional) The minimum size (in bytes) of the image to return.

* `size_max` - (Optional) The maximum size (in bytes) of the image to return.

* `sort_direction` - (Optional) Order the results in either `asc` or `desc`.

* `sort_key` - (Optional) Sort images based on a certain key. Defaults to `name`.

* `tag` - (Optional) Search for images with a specific tag.

* `visibility` - (Optional) The visibility of the image. Must be one of 
    "private", "community", or "shared". Defaults to "private".

* `hidden` - (Optional) Whether or not the image is hidden from public list.

## Attributes Reference

`id` is set to the ID of the found image. In addition, the following attributes
are exported:

* `checksum` - The checksum of the data associated with the image.
* `created_at` - The date the image was created.
* `container_format`: The format of the image's container.
* `disk_format`: The format of the image's disk.
* `file` - the trailing path after the endpoint that represent the
location of the image or the path to retrieve it.
* `metadata` - The metadata associated with the image.
   Image metadata allow for meaningfully define the image properties
   and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.
* `min_disk_gb` - The minimum amount of disk space required to use the image.
* `min_ram_mb` - The minimum amount of ram required to use the image.
* `properties` - Freeform information about the image.
* `protected` - Whether or not the image is protected.
* `schema` - The path to the JSON-schema that represent
   the image or image
* `size_bytes` - The size of the image (in bytes).
* `tags` - The tags list of the image.
* `updated_at` - The date the image was last updated.
