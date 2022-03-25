---
layout: "vkcs"
page_title: "vkcs: blockstorage_volume"
description: |-
  Get information on an VKCS Volume.
---

# vkcs\_blockstorage\_volume

Use this data source to get information about an existing volume.

## Example Usage

```hcl
data "vkcs_blockstorage_volume" "volume_1" {
  name = "volume_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Block Storage
    client. If omitted, the `region` argument of the provider is used.

* `name` - (Optional) The name of the volume.

* `status` - (Optional) The status of the volume.

* `metadata` - (Optional) Metadata key/value pairs associated with the volume.

## Attributes Reference

`id` is set to the ID of the found volume. In addition, the following attributes
are exported:

* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `status` - See Argument Reference above.
* `metadata` - See Argument Reference above.
* `volume_type` - The type of the volume.
* `bootable` - Indicates if the volume is bootable.
* `size` - The size of the volume in GBs.
* `source_volume_id` - The ID of the volume from which the current volume was created.
