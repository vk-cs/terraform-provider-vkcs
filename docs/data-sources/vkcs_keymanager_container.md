---
layout: "vkcs"
page_title: "vkcs: keymanager_container"
description: |-
  Get information on a Key container resource within VKCS.
---

# vkcs\_keymanager\_container

Use this data source to get the ID of an available Key container.

## Example Usage

```hcl
data "vkcs_keymanager_container" "example" {
  name = "my_container"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the KeyManager client.
  A KeyManager client is needed to fetch a container. If omitted, the `region`
  argument of the provider is used.

* `name` - (Optional) The Container name.

## Attributes Reference

The following attributes are exported:

* `container_ref` - The container reference / where to find the container.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - The container type.
* `secret_refs` - A set of dictionaries containing references to secrets. The
  structure is described below.
* `creator_id` - The creator of the container.
* `status` - The status of the container.
* `created_at` - The date the container was created.
* `updated_at` - The date the container was last updated.
* `consumers` - The list of the container consumers. The structure is described
  below.
* `acl` - The list of ACLs assigned to a container. The `read` structure is
  described below.

The `secret_refs` block supports:

* `name` - The name of the secret reference. The reference names must correspond
  the container type.

* `secret_ref` - The secret reference / where to find the secret, URL.

The `consumers` block supports:

* `name` - The name of the consumer.

* `url` - The consumer URL.

The `acl` `read` attribute supports:

* `project_access` - Whether the container is accessible project wide.

* `users` - The list of user IDs, which are allowed to access the container,
  when `project_access` is set to `false`.

* `created_at` - The date the container ACL was created.

* `updated_at` - The date the container ACL was last updated.
