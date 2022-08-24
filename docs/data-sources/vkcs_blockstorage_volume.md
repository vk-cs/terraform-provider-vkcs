---
layout: "vkcs"
page_title: "vkcs: vkcs_blockstorage_volume"
description: |-
  Get information on an VKCS Volume.
---

# vkcs_blockstorage_volume

Use this data source to get information about an existing volume.

## Example Usage

```terraform
data "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
}
```

## Argument Reference
- `bootable` **String** (*Optional*) Indicates if the volume is bootable.

- `metadata` <strong>Map of </strong>**String** (*Optional*) Metadata key/value pairs associated with the volume.

- `name` **String** (*Optional*) The name of the volume.

- `region` **String** (*Optional*) The region in which to obtain the Block Storage client. If omitted, the `region` argument of the provider is used.

- `status` **String** (*Optional*) The status of the volume.

- `volume_type` **String** (*Optional*) The type of the volume.


## Attributes Reference
- `bootable` **String** See Argument Reference above.

- `metadata` <strong>Map of </strong>**String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `volume_type` **String** See Argument Reference above.

- `availability_zone` **String** The name of the availability zone of the volume.

- `id` **String** ID of the resource.

- `size` **Number** The size of the volume in GBs.

- `source_volume_id` **String** The ID of the volume from which the current volume was created.


