---
layout: "vkcs"
page_title: "vkcs: sharedfilesystem_share"
description: |-
  Configure a Shared File System share.
---

# vkcs\_sharedfilesystem\_share

Use this resource to configure a share.

## Example Usage

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
  description       = "test share network with security services"
  neutron_net_id    = "${vkcs_networking_network.network_1.id}"
  neutron_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  size             = 1
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) The share size, in GBs. The requested share size cannot be greater
	than the allowed GB quota. Changing this resizes the existing share.

* `share_proto` - (Required) The share protocol - can either be NFS, CIFS,
	CEPHFS, GLUSTERFS, HDFS or MAPRFS. Changing this creates a new share.

* `availability_zone` - (Optional) The share availability zone. Changing this creates a new share.

* `description` - (Optional) The human-readable description for the share.
	Changing this updates the description of the existing share.

* `name` - (Optional) The name of the share. Changing this updates the name of the existing share.

* `region` - The region in which to obtain the Shared File System client.
	A Shared File System client is needed to create a share. Changing this creates a new share.

* `share_type` - (Optional) The share type name. If you omit this parameter, the default
	share type is used.

* `snapshot_id` - (Optional) The UUID of the share's base snapshot. Changing this creates a new share.

## Attributes Reference

* `id` - The unique ID for the Share.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Share.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `share_proto` - See Argument Reference above.
* `size` - See Argument Reference above.
* `share_type` - See Argument Reference above.
* `snapshot_id` - See Argument Reference above.
* `availability_zone` - See Argument Reference above.
* `share_server_id` - The UUID of the share server.
* `all_metadata` - The map of metadata, assigned on the share, which has been explicitly and implicitly added.

## Import

This resource can be imported by specifying the ID of the share:

```
$ terraform import vkcs_sharedfilesystem_share.share_1 <id>
```
