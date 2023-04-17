---
layout: "vkcs"
page_title: "vkcs: vkcs_blockstorage_volume"
description: |-
  Manages a blockstorage volume.
---

# vkcs_blockstorage_volume

Provides a blockstorage volume resource. This can be used to create, modify and delete blockstorage volume.

## Example Usage

```terraform
resource "vkcs_blockstorage_volume" "volume" {
  name = "volume"
  description = "test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  availability_zone = "GZ1"
  volume_type = "ceph-ssd"
}
```
## Argument Reference
- `availability_zone` **String** (***Required***) The name of the availability zone of the volume.

- `size` **Number** (***Required***) The size of the volume.

- `volume_type` **String** (***Required***) The type of the volume.

- `description` **String** (*Optional*) The description of the volume.

- `image_id` **String** (*Optional*) ID of the image to create volume with. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.

- `metadata` <strong>Map of </strong>**String** (*Optional*) Map of key-value metadata of the volume.

- `name` **String** (*Optional*) The name of the volume.

- `region` **String** (*Optional*) Region to create resource in.

- `snapshot_id` **String** (*Optional*) ID of the snapshot of volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.

- `source_vol_id` **String** (*Optional*) ID of the source volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` **String** ID of the resource.



## Import

Volumes can be imported using the `id`, e.g.

```shell
terraform import vkcs_blockstorage_volume.myvolume 64f3cfc5-226e-4388-a9b8-365b1441b94f
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
