---
subcategory: "File Share (NFS)"
layout: "vkcs"
page_title: "vkcs: vkcs_sharedfilesystem_sharenetwork"
description: |-
  Get information on an Shared File System share network.
---

# vkcs_sharedfilesystem_sharenetwork

Use this data source to get the ID of an available Shared File System share network.

## Example Usage

```terraform
data "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
  name = "sharenetwork_1"
}
```

## Argument Reference
- `description` optional *string* &rarr;  The human-readable description of the share network.

- `name` optional *string* &rarr;  The name of the share network.

- `neutron_net_id` optional *string* &rarr;  The neutron network UUID of the share network.

- `neutron_subnet_id` optional *string* &rarr;  The neutron subnet UUID of the share network.

- `region` optional *string* &rarr;  The region in which to obtain the Shared File System client. A Shared File System client is needed to read a share network. If omitted, the `region` argument of the provider is used.

- `security_service_id` optional *string* &rarr;  The security service IDs associated with the share network.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `cidr` *string* &rarr;  The share network CIDR.

- `id` *string* &rarr;  ID of the resource.

- `project_id` *string* &rarr;  The owner of the share network.

- `security_service_ids` *set of* *string* &rarr;  The list of security service IDs associated with the share network.


