---
layout: "vkcs"
page_title: "vkcs: compute_quotaset"
description: |-
  Get information on a Compute Quotaset of a project.
---

# vkcs\_compute\_quotaset

Use this data source to get the compute quotaset of an VKCS project.

## Example Usage

```hcl
data "vkcs_compute_quotaset" "quota" {
  project_id = "2e367a3d29f94fd988e6ec54e305ec9d"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Compute client.
    If omitted, the `region` argument of the provider is used.

* `project_id` - (Required) The id of the project to retrieve the quotaset.


## Attributes Reference

The following attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - See Argument Reference above.
* `cores` -  The number of allowed server cores.
* `instances` - The number of allowed servers.
* `key_pairs` - The number of allowed key pairs for each user.
* `metadata_items` - The number of allowed metadata items for each server.
* `ram` - The amount of allowed server RAM, in MiB.
* `server_groups` - The number of allowed server groups.
* `server_group_members` - The number of allowed members for each server group.
* `injected_file_content_bytes` - The number of allowed bytes of content for each injected file.
* `injected_file_path_bytes` - The number of allowed bytes for each injected file path.
* `injected_files` - The number of allowed injected files.
