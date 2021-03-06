---
layout: "vkcs"
page_title: "vkcs: sharedfilesystem_securityservice"
description: |-
  Configure a Shared File System security service.
---

# vkcs\_sharedfilesystem\_securityservice

Use this resource to configure a security service.

~> **Note:** All arguments including the security service password will be
stored in the raw state as plain-text. [Read more about sensitive data in
state](/docs/state/sensitive-data.html).

A security service stores configuration information for clients for
authentication and authorization (AuthN/AuthZ). For example, a share server
will be the client for an existing service such as LDAP, Kerberos, or
Microsoft Active Directory.

## Example Usage

```hcl
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

The following arguments are supported:

* `type` - (Required) The security service type - can either be active\_directory,
	kerberos or ldap.  Changing this updates the existing security service.

* `description` - (Optional) The human-readable description for the security service.
	Changing this updates the description of the existing security service.

* `dns_ip` - (Optional) The security service DNS IP address that is used inside the
	tenant network.

* `domain` - (Optional) The security service domain.

* `name` - (Optional) The name of the security service. Changing this updates the name
	of the existing security service.

* `password` - (Optional) The user password, if you specify a user.

* `region` - (Optional) The region in which to obtain the Shared File System client.
	A Shared File System client is needed to create a security service. If omitted, the
	`region` argument of the provider is used. Changing this creates a new
	security service.

* `server` - (Optional) The security service host name or IP address.

* `user` - (Optional) The security service user or group name that is used by the
	tenant.

## Attributes Reference

* `id` - The unique ID for the Security Service.
* `region` - See Argument Reference above.
* `project_id` - The owner of the Security Service.
* `name` - See Argument Reference above.
* `description` - See Argument Reference above.
* `type` - See Argument Reference above.
* `dns_ip` - See Argument Reference above.
* `user` - See Argument Reference above.
* `password` - See Argument Reference above.
* `domain` - See Argument Reference above.
* `server` - See Argument Reference above.

## Import

This resource can be imported by specifying the ID of the security service:

```
$ terraform import vkcs_sharedfilesystem_securityservice.securityservice_1 <id>
```
