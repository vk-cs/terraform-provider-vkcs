---
subcategory: "Disks (block storage)"
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
- `bootable` optional *string* &rarr;  Indicates if the volume is bootable.

- `metadata` optional *map of* *string* &rarr;  Metadata key/value pairs associated with the volume.

- `name` optional *string* &rarr;  The name of the volume.

- `region` optional *string* &rarr;  The region in which to obtain the Block Storage client. If omitted, the `region` argument of the provider is used.

- `status` optional *string* &rarr;  The status of the volume.

- `volume_type` optional *string* &rarr;  The type of the volume.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `availability_zone` *string* &rarr;  The name of the availability zone of the volume.

- `id` *string* &rarr;  ID of the resource.

- `size` *number* &rarr;  The size of the volume in GBs.

- `source_volume_id` *string* &rarr;  The ID of the volume from which the current volume was created.


