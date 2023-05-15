---
subcategory: "File Share (NFS)"
layout: "vkcs"
page_title: "vkcs: vkcs_sharedfilesystem_share"
description: |-
  Get information on an Shared File System share.
---

# vkcs_sharedfilesystem_share

Use this data source to get the ID of an available Shared File System share.

## Example Usage

```terraform
data "vkcs_sharedfilesystem_share" "share_1" {
  name = "share_1"
}
```

## Argument Reference
- `share_network_id` **required** *string* &rarr;  The UUID of the share's share network.

- `description` optional *string* &rarr;  The human-readable description for the share.

- `export_location_path` optional *string* &rarr;  The export location path of the share.

- `name` optional *string* &rarr;  The name of the share.

- `region` optional *string* &rarr;  The region in which to obtain the Shared File System client.

- `snapshot_id` optional *string* &rarr;  The UUID of the share's base snapshot.

- `status` optional *string* &rarr;  A share status filter. A valid value is `creating`, `error`, `available`, `deleting`, `error_deleting`, `manage_starting`, `manage_error`, `unmanage_starting`, `unmanage_error`, `unmanaged`, `extending`, `extending_error`, `shrinking`, `shrinking_error`, or `shrinking_possible_data_loss_error`.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `availability_zone` *string* &rarr;  The share availability zone.

- `id` *string* &rarr;  ID of the resource.

- `project_id` *string* &rarr;  The owner of the share.

- `share_proto` *string* &rarr;  The share protocol.

- `size` *number* &rarr;  The share size, in GBs.


