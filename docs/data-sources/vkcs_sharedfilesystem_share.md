---
layout: "vkcs"
page_title: "vkcs: sharedfilesystem_share"
description: |-
  Get information on an Shared File System share.
---

# vkcs\_sharedfilesystem\_share

Use this data source to get the ID of an available Shared File System share.

## Example Usage

```hcl
data "vkcs_sharedfilesystem_share" "share_1" {
  name = "share_1"
}
```

## Argument Reference

* `share_network_id` - (Required) The UUID of the share's share network.

* `name` - (Optional) The name of the share.

* `description` - (Optional) The human-readable description for the share.

* `project_id` - (Optional) The owner of the share.

* `snapshot_id` - (Optional) The UUID of the share's base snapshot.

* `export_location_path` - (Optional) The export location path of the share.

* `status` - (Optional) A share status filter. A valid value is `creating`,
    `error`, `available`, `deleting`, `error_deleting`, `manage_starting`,
    `manage_error`, `unmanage_starting`, `unmanage_error`, `unmanaged`,
    `extending`, `extending_error`, `shrinking`, `shrinking_error`, or
    `shrinking_possible_data_loss_error`.

## Attributes Reference

`id` is set to the ID of the found share. In addition, the following attributes
are exported:

* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `share_network_id` - See Argument Reference above.
* `export_location_path` - See Argument Reference above.
* `status` - See Argument Reference above.
* `region` - The region in which to obtain the Shared File System client.
* `availability_zone` - The share availability zone.
* `share_proto` - The share protocol.
* `size` - The share size, in GBs.
