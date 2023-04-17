---
layout: "vkcs"
page_title: "vkcs: vkcs_blockstorage_snapshot"
description: |-
  Get information on an VKCS Volume Snapshot.
---

# vkcs_blockstorage_snapshot

Use this data source to get information about an existing snapshot.

## Example Usage

```terraform
data "vkcs_blockstorage_snapshot" "snapshot_1" {
  name        = "snapshot_1"
  most_recent = true
}
```

## Argument Reference
- `most_recent` **Boolean** (*Optional*) Pick the most recently created snapshot if there are multiple results.

- `name` **String** (*Optional*) The name of the snapshot.

- `region` **String** (*Optional*) The region in which to obtain the Block Storage client. If omitted, the `region` argument of the provider is used.

- `status` **String** (*Optional*) The status of the snapshot.

- `volume_id` **String** (*Optional*) The ID of the snapshot's volume.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `description` **String** The snapshot's description.

- `id` **String** ID of the resource.

- `metadata` <strong>Map of </strong>**String** The snapshot's metadata.

- `size` **Number** The size of the snapshot.


