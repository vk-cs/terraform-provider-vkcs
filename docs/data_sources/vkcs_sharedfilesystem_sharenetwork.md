---
layout: "vkcs"
page_title: "vkcs: sharedfilesystem_sharenetwork"
description: |-
  Get information on an Shared File System share network.
---

# vkcs\_sharedfilesystem\_sharenetwork

Use this data source to get the ID of an available Shared File System share network.

## Example Usage

```hcl
data "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
  name = "sharenetwork_1"
}
```

## Argument Reference

* `region` - (Optional) The region in which to obtain the Shared File System client.
    A Shared File System client is needed to read a share network. If omitted, the
    `region` argument of the provider is used.

* `name` - (Optional) The name of the share network.

* `description` - (Optional) The human-readable description of the share network.

* `project_id` - (Optional) The owner of the share network.

* `neutron_net_id` - (Optional) The neutron network UUID of the share network.

* `neutron_subnet_id` - (Optional) The neutron subnet UUID of the share network.

* `security_service_id` - (Optional) The security service IDs associated with the share network.

* `cidr` - (Optional) The share network CIDR.

* `ip_version` - (Optional) The IP version of the share network. Can either be 4 or 6.

## Attributes Reference

`id` is set to the ID of the found share network . In addition, the following
attributes are exported:

* `region` - See Argument Reference above.
* `project_id` - The owner of the Share Network.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `neutron_net_id` - See Argument Reference above.
* `neutron_subnet_id` - See Argument Reference above.
* `security_service_id` - See Argument Reference above.
* `cidr` - See Argument Reference above.
* `ip_version` - See Argument Reference above.
* `security_service_ids` - The list of security service IDs associated with the share network.
