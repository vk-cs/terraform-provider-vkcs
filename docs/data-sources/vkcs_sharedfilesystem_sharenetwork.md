---
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
- `description` **String** (*Optional*) The human-readable description of the share network.

- `name` **String** (*Optional*) The name of the share network.

- `neutron_net_id` **String** (*Optional*) The neutron network UUID of the share network.

- `neutron_subnet_id` **String** (*Optional*) The neutron subnet UUID of the share network.

- `region` **String** (*Optional*) The region in which to obtain the Shared File System client. A Shared File System client is needed to read a share network. If omitted, the `region` argument of the provider is used.

- `security_service_id` **String** (*Optional*) The security service IDs associated with the share network.


## Attributes Reference
- `description` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `neutron_net_id` **String** See Argument Reference above.

- `neutron_subnet_id` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `security_service_id` **String** See Argument Reference above.

- `cidr` **String** The share network CIDR.

- `id` **String** ID of the resource.

- `project_id` **String** The owner of the share network.

- `security_service_ids` <strong>Set of </strong>**String** The list of security service IDs associated with the share network.


