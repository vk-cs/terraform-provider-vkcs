---
subcategory: "Disks (block storage)"
layout: "vkcs"
page_title: "vkcs: vkcs_blockstorage_snapshot"
description: |-
  Get information on an VKCS Volume Snapshot.
---

# vkcs_blockstorage_snapshot

Use this data source to get information about an existing snapshot.

## Example Usage

```terraform
data "vkcs_blockstorage_snapshot" "recent_snapshot" {
  name        = "snapshot-tf-example"
  most_recent = true
  # This is unnecessary in real life.
  # This is required here to let the example work with snapshot resource example.
  depends_on = [vkcs_blockstorage_snapshot.recent_snapshot]
}
```

## Argument Reference
- `most_recent` optional *boolean* &rarr;  Pick the most recently created snapshot if there are multiple results.

- `name` optional *string* &rarr;  The name of the snapshot.

- `region` optional *string* &rarr;  The region in which to obtain the Block Storage client. If omitted, the `region` argument of the provider is used.

- `status` optional *string* &rarr;  The status of the snapshot.

- `volume_id` optional *string* &rarr;  The ID of the snapshot's volume.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `description` *string* &rarr;  The snapshot's description.

- `id` *string* &rarr;  ID of the resource.

- `metadata` *map of* *string* &rarr;  The snapshot's metadata.

- `size` *number* &rarr;  The size of the snapshot.


