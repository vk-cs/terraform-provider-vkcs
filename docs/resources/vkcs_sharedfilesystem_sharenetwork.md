---
layout: "vkcs"
page_title: "VKCS: sharedfilesystem_sharenetwork"
description: |-
  Configure a Shared File System share network.
---

# vkcs\_sharedfilesystem\_sharenetwork

Use this resource to configure a share network.

A share network stores network information that share servers can use when shares are created.

## Example Usage

### Basic share network

```hcl
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
  name              = "test_sharenetwork"
  description       = "test share network"
  neutron_net_id    = "${vkcs_networking_network.network_1.id}"
  neutron_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}
```

### Share network with associated security services

```hcl
resource "vkcs_networking_network" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "vkcs_networking_subnet" "subnet_1" {
  name       = "subnet_1"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${vkcs_networking_network.network_1.id}"
}

resource "vkcs_sharedfilesystem_securityservice" "securityservice_1" {
  name        = "security"
  description = "created by terraform"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
  domain      = "example.com"
  user        = "joinDomainUser"
  password    = "s8cret"
}

resource "vkcs_sharedfilesystem_sharenetwork" "sharenetwork_1" {
  name              = "test_sharenetwork"
  description       = "test share network with security services"
  neutron_net_id    = "${vkcs_networking_network.network_1.id}"
  neutron_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
  security_service_ids = [
	"${vkcs_sharedfilesystem_securityservice.securityservice_1.id}",
  ]
}
```

## Argument Reference

The following arguments are supported:

* `neutron_net_id` - (Required) The UUID of a neutron network when setting up or updating
	a share network. Changing this updates the existing share network if it's not used by shares.

* `neutron_subnet_id` - (Required) The UUID of the neutron subnet when setting up or
	updating a share network. Changing this updates the existing share network if it's not used by shares.

* `description` - (Optional) The human-readable description for the share network.
	Changing this updates the description of the existing share network.

* `name` - (Optional) The name for the share network. Changing this updates the name of the existing share network.

* `region` - (Optional) The region in which to obtain the V2 Shared File System client.
	A Shared File System client is needed to create a share network. If omitted, the
	`region` argument of the provider is used. Changing this creates a new share network.

* `security_service_ids` - (Optional) The list of security service IDs to associate with
	the share network. The security service must be specified by ID and not name.

## Attributes Reference

* `id` - The unique ID for the Share Network.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Share Network.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `neutron_net_id` - See Argument Reference above.
* `neutron_subnet_id` - See Argument Reference above.
* `security_service_ids` - See Argument Reference above.
* `cidr` - The share network CIDR.
* `ip_version` - The IP version of the share network. Can either be 4 or 6.

## Import

This resource can be imported by specifying the ID of the share network:

```
$ terraform import vkcs_sharedfilesystem_sharenetwork.sharenetwork_1 <id>
```
