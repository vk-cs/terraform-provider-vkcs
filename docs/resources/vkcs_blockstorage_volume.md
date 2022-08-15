---
layout: "vkcs"
page_title: "vkcs: blockstorage_volume"
subcategory: ""
description: |-
  Manages a blockstorage volume.
---

# vkcs\_blockstorage\_volume

Provides a blockstorage volume resource. This can be used to create, modify and delete blockstorage volume.

## Example Usage

```terraform

resource "vkcs_blockstorage_volume" "bs-volume" {
  name = "bs-volume"
  size = 8
  volume_type = "ceph-hdd"
  availability_zone = "GZ1"
}
```
## Argument Reference

The following arguments are supported:

* `size` - (Required) The size of the volume.

* `volume_type` - (Required) The type of the volume.

* `availability_zone` - (Required) The name of the availability zone of the volume.

* `name` - The name of the volume.

* `description` - The description of the volume.

* `metadata` - Map of key-value metadata of the volume.

* `snapshot_id` - ID of the snapshot of volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set. 

* `source_volume_id` - ID of the source volume. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set. 

* `image_id` - ID of the image to create volume with. Changing this creates a new volume. Only one of snapshot_id, source_volume_id, image_id fields may be set. 
