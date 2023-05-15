---
subcategory: "File Share (NFS)"
layout: "vkcs"
page_title: "vkcs: vkcs_sharedfilesystem_securityservice"
description: |-
  Configure a Shared File System security service.
---

# vkcs_sharedfilesystem_securityservice

Use this resource to configure a security service.

~> **Note:** All arguments including the security service password will be stored in the raw state as plain-text. [Read more about sensitive data in state](/docs/state/sensitive-data.html).

A security service stores configuration information for clients for authentication and authorization (AuthN/AuthZ). For example, a share server will be the client for an existing service such as LDAP, Kerberos, or Microsoft Active Directory.

## Example Usage
```terraform
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
```

## Argument Reference
- `type` **required** *string* &rarr;  The security service type - can either be active\_directory, kerberos or ldap.  Changing this updates the existing security service.

- `description` optional *string* &rarr;  The human-readable description for the security service. Changing this updates the description of the existing security service.

- `dns_ip` optional *string* &rarr;  The security service DNS IP address that is used inside the tenant network.

- `domain` optional *string* &rarr;  The security service domain.

- `name` optional *string* &rarr;  The name of the security service. Changing this updates the name of the existing security service.

- `password` optional sensitive *string* &rarr;  The user password, if you specify a user.

- `region` optional *string* &rarr;  The region in which to obtain the Shared File System client. A Shared File System client is needed to create a security service. If omitted, the `region` argument of the provider is used. Changing this creates a new security service.

- `server` optional *string* &rarr;  The security service host name or IP address.

- `user` optional *string* &rarr;  The security service user or group name that is used by the tenant.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `id` *string* &rarr;  ID of the resource.

- `project_id` *string* &rarr;  The owner of the Security Service.



## Import

This resource can be imported by specifying the ID of the security service:

```shell
terraform import vkcs_sharedfilesystem_securityservice.securityservice_1 048d7c1c-4187-4370-89ce-da638120d168
```
