---
subcategory: "Disks (block storage)"
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
- `availability_zone` **required** *string* &rarr;  The name of the availability zone of the volume.

- `size` **required** *number* &rarr;  The size of the volume.

- `volume_type` **required** *string* &rarr;  The type of the volume.

- `description` optional *string* &rarr;  The description of the volume.

- `image_id` optional *string* &rarr;  ID of the image to create volume with. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.

- `metadata` optional *map of* *string* &rarr;  Map of key-value metadata of the volume.

- `name` optional *string* &rarr;  The name of the volume.

- `region` optional *string* &rarr;  Region to create resource in.

- `snapshot_id` optional *string* &rarr;  ID of the snapshot of volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.

- `source_vol_id` optional *string* &rarr;  ID of the source volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Volumes can be imported using the `id`, e.g.

```shell
terraform import vkcs_blockstorage_volume.myvolume 64f3cfc5-226e-4388-a9b8-365b1441b94f
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
