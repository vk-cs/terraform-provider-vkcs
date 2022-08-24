---
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
- `volume_id` **String** (***Required***) ID of the volume to create snapshot for. Changing this creates a new snapshot.

- `description` **String** (*Optional*) The description of the volume.

- `force` **Boolean** (*Optional*) Allows or disallows snapshot of a volume when the volume is attached to an instance.

- `metadata` <strong>Map of </strong>**String** (*Optional*) Map of key-value metadata of the volume.

- `name` **String** (*Optional*) The name of the snapshot.

- `region` **String** (*Optional*)


## Attributes Reference
- `volume_id` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `force` **Boolean** See Argument Reference above.

- `metadata` <strong>Map of </strong>**String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String**

- `id` **String** ID of the resource.



## Import

Volume snapshots can be imported using the `id`, e.g.

```shell
terraform import vkcs_blockstorage_snapshot.myvolumesnapshot 0b4f5a9b-554e-4e80-b553-82aba6502315
```

After the import you can use ```terraform show``` to view imported fields and write their values to your .tf file.
