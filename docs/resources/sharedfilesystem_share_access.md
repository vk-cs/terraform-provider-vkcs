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
```terraform
resource "vkcs_sharedfilesystem_share_access" "opencloud" {
  share_id     = vkcs_sharedfilesystem_share.data.id
  access_type  = "ip"
  access_to    = "192.168.199.11"
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
