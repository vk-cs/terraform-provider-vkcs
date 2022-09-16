---
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
- `name` **String** (*Optional*) The Container name.

- `region` **String** (*Optional*) The region in which to obtain the KeyManager client. A KeyManager client is needed to fetch a container. If omitted, the `region` argument of the provider is used.


## Attributes Reference
- `name` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `acl` **Object** The list of ACLs assigned to a container.

- `consumers` **Object** The list of the container consumers.

- `container_ref` **String** The container reference / where to find the container.

- `created_at` **String** The date the container was created.

- `creator_id` **String** The creator of the container.

- `id` **String** ID of the resource.

- `secret_refs` <strong>Set of </strong>**Object** A set of dictionaries containing references to secrets.

- `status` **String** The status of the container.

- `type` **String** The container type.

- `updated_at` **String** The date the container was last updated.


