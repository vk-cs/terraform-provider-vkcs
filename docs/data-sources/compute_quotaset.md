---
subcategory: "Virtual Machines"
layout: "vkcs"
page_title: "vkcs: vkcs_compute_quotaset"
description: |-
  Get information on a Compute Quotaset of a project.
---

# vkcs_compute_quotaset

Use this data source to get the compute quotaset of an VKCS project.

## Example Usage

```terraform
data "vkcs_compute_quotaset" "quota" {
  project_id = "2e367a3d29f94fd988e6ec54e305ec9d"
}
```

## Argument Reference
- `project_id` **required** *string* &rarr;  The id of the project to retrieve the quotaset.

- `region` optional *string* &rarr;  The region in which to obtain the Compute client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `cores` *number* &rarr;  The number of allowed server cores.

- `id` *string* &rarr;  ID of the resource.

- `injected_file_content_bytes` *number* &rarr;  The number of allowed bytes of content for each injected file.

- `injected_file_path_bytes` *number* &rarr;  The number of allowed bytes for each injected file path.

- `injected_files` *number* &rarr;  The number of allowed injected files.

- `instances` *number* &rarr;  The number of allowed servers.

- `key_pairs` *number* &rarr;  The number of allowed key pairs for each user.

- `metadata_items` *number* &rarr;  The number of allowed metadata items for each server.

- `ram` *number* &rarr;  The amount of allowed server RAM, in MiB.

- `server_group_members` *number* &rarr;  The number of allowed members for each server group.

- `server_groups` *number* &rarr;  The number of allowed server groups.


