---
subcategory: "Images"
layout: "vkcs"
page_title: "vkcs: vkcs_images_image"
description: |-
  Manages an Image resource within VKCS.
---

# vkcs_images_image

Manages an Image resource within VKCS.

~> **Note:** All arguments including the source image URL password will be stored in the raw state as plain-text. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

```terraform
resource "vkcs_images_image" "rancheros" {
  name             = "RancherOS"
  image_source_url = "https://releases.rancher.com/os/latest/rancheros-openstack.img"
  container_format = "bare"
  disk_format      = "qcow2"

  properties = {
    key = "value"
  }
}
```
## Argument Reference
- `container_format` **required** *string* &rarr;  The container format. Must be one of "bare".

- `disk_format` **required** *string* &rarr;  The disk format. Must be one of "raw", "iso".

- `name` **required** *string* &rarr;  The name of the image.

- `image_cache_path` optional *string* &rarr;  This is the directory where the images will be downloaded. Images will be stored with a filename corresponding to the url's md5 hash. Defaults to "$HOME/.terraform/image_cache"

- `image_source_password` optional sensitive *string* &rarr;  The password of basic auth to download `image_source_url`.

- `image_source_url` optional *string* &rarr;  This is the url of the raw image. The image will be downloaded in the `image_cache_path` before being uploaded to VKCS. Conflicts with `local_file_path`.

- `image_source_username` optional *string* &rarr;  The username of basic auth to download `image_source_url`.

- `local_file_path` optional *string* &rarr;  This is the filepath of the raw image file that will be uploaded to VKCS. Conflicts with `image_source_url`

- `min_disk_gb` optional *number* &rarr;  Amount of disk space (in GB) required to boot image. Defaults to 0.

- `min_ram_mb` optional *number* &rarr;  Amount of ram (in MB) required to boot image. Defauts to 0.

- `properties` optional *map of* *string* &rarr;  A map of key/value pairs to set freeform information about an image. See the "Notes" section for further information about properties.

- `protected` optional *boolean* &rarr;  If true, image will not be deletable. Defaults to false.

- `region` optional *string* &rarr;  The region in which to obtain the Image client. An Image client is needed to create an Image that can be used with a compute instance. If omitted, the `region` argument of the provider is used. Changing this creates a new Image.

- `tags` optional *set of* *string* &rarr;  The tags of the image. It must be a list of strings. At this time, it is not possible to delete all tags of an image.

- `verify_checksum` optional *boolean* &rarr;  If false, the checksum will not be verified once the image is finished uploading.

- `visibility` optional *string* &rarr;  The visibility of the image. Must be one of "private", "community", or "shared". The ability to set the visibility depends upon the configuration of the VKCS cloud.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `checksum` *string* &rarr;  The checksum of the data associated with the image.

- `created_at` *string* &rarr;  The date the image was created.

- `file` *string* &rarr;  The trailing path after the image endpoint that represent the location of the image or the path to retrieve it.

- `id` *string* &rarr;  ID of the resource.

- `metadata` *map of* *string* &rarr;  The metadata associated with the image. Image metadata allow for meaningfully define the image properties and tags. See https://docs.openstack.org/glance/latest/user/metadefs-concepts.html.

- `owner` *string* &rarr;  The id of the vkcs user who owns the image.

- `schema` *string* &rarr;  The path to the JSON-schema that represent the image or image

- `size_bytes` *number* &rarr;  The size in bytes of the data associated with the image.

- `status` *string* &rarr;  The status of the image. It can be "queued", "active" or "saving".

- `updated_at` *string* &rarr;  The date the image was last updated.



## Notes
### Properties

This resource supports the ability to add properties to a resource during creation as well as add, update, and delete properties during an update of this resource.

## Import

Images can be imported using the `id`, e.g.

```shell
terraform import vkcs_images_image.rancheros 89c60255-9bd6-460c-822a-e2b959ede9d2
```
