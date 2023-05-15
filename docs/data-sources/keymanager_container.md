---
subcategory: "Key Manager"
layout: "vkcs"
page_title: "vkcs: vkcs_keymanager_container"
description: |-
  Get information on a Key container resource within VKCS.
---

# vkcs_keymanager_container

Use this data source to get the ID of an available Key container.

## Example Usage

```terraform
data "vkcs_keymanager_container" "example" {
  name = "my_container"
}
```

## Argument Reference
- `name` optional *string* &rarr;  The Container name.

- `region` optional *string* &rarr;  The region in which to obtain the KeyManager client. A KeyManager client is needed to fetch a container. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `acl` *object* &rarr;  The list of ACLs assigned to a container.

- `consumers` *object* &rarr;  The list of the container consumers.

- `container_ref` *string* &rarr;  The container reference / where to find the container.

- `created_at` *string* &rarr;  The date the container was created.

- `creator_id` *string* &rarr;  The creator of the container.

- `id` *string* &rarr;  ID of the resource.

- `secret_refs` *set of* *object* &rarr;  A set of dictionaries containing references to secrets.

- `status` *string* &rarr;  The status of the container.

- `type` *string* &rarr;  The container type.

- `updated_at` *string* &rarr;  The date the container was last updated.


