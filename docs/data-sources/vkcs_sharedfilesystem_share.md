---
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
- `share_network_id` **String** (***Required***) The UUID of the share's share network.

- `description` **String** (*Optional*) The human-readable description for the share.

- `export_location_path` **String** (*Optional*) The export location path of the share.

- `name` **String** (*Optional*) The name of the share.

- `region` **String** (*Optional*) The region in which to obtain the Shared File System client.

- `snapshot_id` **String** (*Optional*) The UUID of the share's base snapshot.

- `status` **String** (*Optional*) A share status filter. A valid value is `creating`, `error`, `available`, `deleting`, `error_deleting`, `manage_starting`, `manage_error`, `unmanage_starting`, `unmanage_error`, `unmanaged`, `extending`, `extending_error`, `shrinking`, `shrinking_error`, or `shrinking_possible_data_loss_error`.


## Attributes Reference
- `share_network_id` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `export_location_path` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `snapshot_id` **String** See Argument Reference above.

- `status` **String** See Argument Reference above.

- `availability_zone` **String** The share availability zone.

- `id` **String** ID of the resource.

- `project_id` **String** The owner of the share.

- `share_proto` **String** The share protocol.

- `size` **Number** The share size, in GBs.


