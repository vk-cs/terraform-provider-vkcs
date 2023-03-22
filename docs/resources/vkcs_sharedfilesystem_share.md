---
layout: "vkcs"
page_title: "vkcs: vkcs_sharedfilesystem_share"
description: |-
  Configure a Shared File System share.
---

# vkcs_sharedfilesystem_share

Use this resource to configure a share.

## Example Usage
```terraform
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
- `name` **String** (***Required***) The name of the share. Changing this updates the name of the existing share.

- `share_network_id` **String** (***Required***) The UUID of the share network.

- `share_proto` **String** (***Required***) The share protocol - can either be NFS, CIFS, CEPHFS, GLUSTERFS, HDFS or MAPRFS. Changing this creates a new share.

- `size` **Number** (***Required***) The share size, in GBs. The requested share size cannot be greater than the allowed GB quota. Changing this resizes the existing share.

- `availability_zone` **String** (*Optional*) The share availability zone. Changing this creates a new share.

- `description` **String** (*Optional*) The human-readable description for the share. Changing this updates the description of the existing share.

- `region` **String** (*Optional*) The region in which to obtain the Shared File System client. A Shared File System client is needed to create a share. Changing this creates a new share.

- `share_type` **String** (*Optional*) The share type name. If you omit this parameter, the default share type is used.

- `snapshot_id` **String** (*Optional*) The UUID of the share's base snapshot. Changing this creates a new share.


## Attributes Reference
- `name` **String** See Argument Reference above.

- `share_network_id` **String** See Argument Reference above.

- `share_proto` **String** See Argument Reference above.

- `size` **Number** See Argument Reference above.

- `availability_zone` **String** See Argument Reference above.

- `description` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `share_type` **String** See Argument Reference above.

- `snapshot_id` **String** See Argument Reference above.

- `all_metadata` <strong>Map of </strong>**String** The map of metadata, assigned on the share, which has been explicitly and implicitly added.

- `export_location_path` **String** The export location path of the share

- `id` **String** ID of the resource.

- `project_id` **String** The owner of the Share.

- `share_server_id` **String** The UUID of the share server.



## Import

This resource can be imported by specifying the ID of the share:

```shell
terraform import vkcs_sharedfilesystem_share.share_1 829b7299-eae0-4860-88d4-13d03f0e9e2c
```
