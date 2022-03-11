---
layout: "vkcs"
page_title: "vkcs: blockstorage_snapshot"
subcategory: ""
description: |-
  Manages a blockstorage snapshot.
---

# vkcs\_blockstorage\_snapshot

Provides a blockstorage snapshot resource. This can be used to create, modify and delete blockstorage snapshot.

## Example Usage

```terraform
resource "vkcs_blockstorage_snapshot" "bs-snapshot" {
    name = "bs-volume-snapshot"
    volume_id = example_volume_id
}
```
## Argument Reference

The following arguments are supported:

* `volume_id` - (Required) ID of the volume to create snapshot for. Changing this creates a new snapshot.

* `name` - The name of the snapshot.

* `force` - Allows or disallows snapshot of a volume when the volume is attached to an instance.

* `description` - The description of the volume.

* `metadata` - Map of key-value metadata of the volume.
