---
layout: "vkcs"
page_title: "vkcs: vkcs_images_images"
description: |-
  Get information on available VKCS images
---

# vkcs_images_images



## Example Usage

```terraform
data "vkcs_images_images" "images" {
  visibility = "public"
  default = true
  properties = {
    mcs_os_distro = "debian"
  }
}
```

## Argument Reference
- `created_at` optional *string* &rarr;  Date filter to select images with created_at matching the specified criteria. Value should be either RFC3339 formatted time or time filter in format `filter:time`, where `filter` is one of [eq, neq, gt, gte, lt, lte] and `time` is RFC3339 formatted time.

- `default` optional *boolean* &rarr;  The flag used to filter images based on whether they are available for virtual machine creation.

- `owner` optional *string* &rarr;  The ID of the owner of images.

- `properties` optional *map of* *string* &rarr;  Search for images with specific properties.

- `region` optional *string* &rarr;  The region in which to obtain the Images client. If omitted, the `region` argument of the provider is used.

- `size_max` optional *number* &rarr;  The maximum size (in bytes) of images to return.

- `size_min` optional *number* &rarr;  The minimum size (in bytes) of images to return.

- `tags` optional *string* &rarr;  Search for images with specific tags.

- `updated_at` optional *string* &rarr;  Date filter to select images with updated_at matching the specified criteria. Value should be either RFC3339 formatted time or time filter in format `filter:time`, where `filter` is one of [eq, neq, gt, gte, lt, lte] and `time` is RFC3339 formatted time.

- `visibility` optional *string* &rarr;  The visibility of images. Must be one of "public", "private", "community", or "shared".


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  The ID of the data source

- `images`  *list* &rarr;  Images matching specified criteria.
  - `id` *string* &rarr;  ID of an image.

  - `name` *string* &rarr;  Name of an image.

  - `properties` *map of* *string* &rarr;  Properties associated with an image.



