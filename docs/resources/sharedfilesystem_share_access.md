---
subcategory: "File Share (NFS)"
layout: "vkcs"
page_title: "vkcs: vkcs_sharedfilesystem_share_access"
description: |-
  Configure a Shared File System share access list.
---

# vkcs_sharedfilesystem_share_access

Use this resource to control the share access lists.

~> **Important Security Notice** The access key assigned by this resource will be stored *unencrypted* in your Terraform state file. If you use this resource in production, please make sure your state file is sufficiently protected. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage
### NFS
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
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = "${vkcs_sharedfilesystem_share.share_1.id}"
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "rw"
}
```

### CIFS
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
  name              = "test_sharenetwork_secure"
  description       = "share the secure love"
  neutron_net_id    = "${vkcs_networking_network.network_1.id}"
  neutron_subnet_id = "${vkcs_networking_subnet.subnet_1.id}"
  security_service_ids = [
	"${vkcs_sharedfilesystem_securityservice.securityservice_1.id}",
  ]
}

resource "vkcs_sharedfilesystem_share" "share_1" {
  name             = "cifs_share"
  share_proto      = "CIFS"
  size             = 1
  share_network_id = "${vkcs_sharedfilesystem_sharenetwork.sharenetwork_1.id}"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_1" {
  share_id     = "${vkcs_sharedfilesystem_share.share_1.id}"
  access_type  = "user"
  access_to    = "windows"
  access_level = "ro"
}

resource "vkcs_sharedfilesystem_share_access" "share_access_2" {
  share_id     = "${vkcs_sharedfilesystem_share.share_1.id}"
  access_type  = "user"
  access_to    = "linux"
  access_level = "rw"
}
```

## Argument Reference
- `access_level` **required** *string* &rarr;  The access level to the share. Can either be `rw` or `ro`.

- `access_to` **required** *string* &rarr;  The value that defines the access. Can either be an IP address or a username verified by configured Security Service of the Share Network.

- `access_type` **required** *string* &rarr;  The access rule type. Can either be an ip, user, cert, or cephx.

- `share_id` **required** *string* &rarr;  The UUID of the share to which you are granted access.

- `region` optional *string* &rarr;  The region in which to obtain the Shared File System client. A Shared File System client is needed to create a share access. Changing this creates a new share access.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.



## Import

This resource can be imported by specifying the ID of the share and the ID of the share access, separated by a slash, e.g.:

```shell
terraform import vkcs_sharedfilesystem_share_access.share_access_1 1c68f8cb-20b5-4f91-b761-6c612b4aae53/c8207c63-6a6d-4a7b-872f-047049582172
```
