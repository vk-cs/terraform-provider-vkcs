---
subcategory: "Disks (block storage)"
layout: "vkcs"
page_title: "vkcs: vkcs_blockstorage_snapshot"
description: |-
  Manages a blockstorage snapshot.
---

# vkcs_blockstorage_snapshot

Provides a blockstorage snapshot resource. This can be used to create, modify and delete blockstorage snapshot.

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

resource "vkcs_blockstorage_snapshot" "snapshot" {
  volume_id = "${vkcs_blockstorage_volume.volume.id}"
  name = "snapshot"
  description = "test snapshot"
  metadata = {
    foo = "bar"
  }
}
```
## Argument Reference
- `volume_id` **required** *string* &rarr;  ID of the volume to create snapshot for. Changing this creates a new snapshot.

- `description` optional *string* &rarr;  The description of the volume.

- `force` optional *boolean* &rarr;  Allows or disallows snapshot of a volume when the volume is attached to an instance.

- `metadata` optional *map of* *string* &rarr;  Map of key-value metadata of the volume.

- `name` optional *string* &rarr;  The name of the snapshot.

- `region` optional *string*


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

Volume snapshots can be imported using the `id`, e.g.

```shell
terraform import vkcs_blockstorage_snapshot.myvolumesnapshot 0b4f5a9b-554e-4e80-b553-82aba6502315
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
