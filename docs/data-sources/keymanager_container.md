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
data "vkcs_keymanager_container" "lb_cert" {
  name = "container-tf-example"
  # This is unnecessary in real life.
  # This is required here to let the example work with container resource example. 
  depends_on = [vkcs_keymanager_container.lb_cert]
}
```

## Argument Reference
- `name` optional *string* &rarr;  The Container name.

- `region` optional *string* &rarr;  The region in which to obtain the VKCS Key Manager client. If omitted, the `region` argument of the provider is used.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `acl`  *list* &rarr;  ACLs assigned to a container.
    - `read`  *list* &rarr;  Object that describes read operation.
        - `created_at` *string* &rarr;  The date the container ACL was created.

        - `project_access` *boolean* &rarr;  Whether the container is accessible project wide.

        - `updated_at` *string* &rarr;  The date the container ACL was last updated.

        - `users` *set of* *string* &rarr;  The list of user IDs, which are allowed to access the container, when `project_access` is set to `false`.



- `consumers`  *list* &rarr;  The list of the container consumers.
    - `name` *string* &rarr;  The name of the consumer.

    - `url` *string* &rarr;  The consumer URL.


- `container_ref` *string* &rarr;  The container reference / where to find the container.

- `created_at` *string* &rarr;  The date the container was created.

- `creator_id` *string* &rarr;  The creator of the container.

- `id` *string* &rarr;  ID of the resource.

- `secret_refs`  *set* &rarr;  A set of dictionaries containing references to secrets.
    - `name` *string*

    - `secret_ref` *string*


- `status` *string* &rarr;  The status of the container.

- `type` *string* &rarr;  The container type.

- `updated_at` *string* &rarr;  The date the container was last updated.


