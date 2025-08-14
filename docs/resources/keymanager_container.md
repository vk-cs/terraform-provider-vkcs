---
subcategory: "Key Manager"
layout: "vkcs"
page_title: "vkcs: vkcs_keymanager_container"
description: |-
  Manages a key container resource within VKCS.
---

# vkcs_keymanager_container

Manages a key container resource within VKCS.

## Example Usage
The container with the TLS certificate and private key which can be used by the loadbalancer HTTPS listener.
```terraform
resource "vkcs_keymanager_container" "lb_cert" {
  name = "container-tf-example"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = vkcs_keymanager_secret.certificate.secret_ref
  }

  secret_refs {
    name       = "private_key"
    secret_ref = vkcs_keymanager_secret.priv_key.secret_ref
  }
}
```

## Argument Reference
- `type` **required** *string* &rarr;  Used to indicate the type of container. Must be one of `generic`, `rsa` or `certificate`.

- `acl` optional &rarr;  Allows to control an access to a container. Currently only the `read` operation is supported. If not specified, the container is accessible project wide. The `read` structure is described below.
    - `read` optional &rarr;  Block that describes read operation.
        - `project_access` optional *boolean* &rarr;  Whether the container is accessible project wide. Defaults to `true`.

        - `users` optional *set of* *string* &rarr;  The list of user IDs, which are allowed to access the container, when `project_access` is set to `false`.

- `name` optional *string* &rarr;  Human-readable name for the Container. Does not have to be unique.

- `region` optional *string* &rarr;  The region in which to obtain the KeyManager client. A KeyManager client is needed to create a container. If omitted, the `region` argument of the provider is used. Changing this creates a new container.

- `secret_refs` optional &rarr;  A set of dictionaries containing references to secrets. The structure is described below.
    - `secret_ref` **required** *string* &rarr;  The secret reference / where to find the secret, URL.

    - `name` optional *string* &rarr;  The name of the secret reference. The reference names must correspond the container type, more details are available [here](https://docs.openstack.org/barbican/stein/api/reference/containers.html).


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `acl` 
    - `read` 
        - `created_at` *string* &rarr;  The date the container ACL was created.

        - `updated_at` *string* &rarr;  The date the container ACL was last updated.

- `consumers` *object* &rarr;  The list of the container consumers. The structure is described below.

- `container_ref` *string* &rarr;  The container reference / where to find the container.

- `created_at` *string* &rarr;  The date the container ACL was created.

- `creator_id` *string* &rarr;  The creator of the container.

- `id` *string* &rarr;  ID of the resource.

- `status` *string* &rarr;  The status of the container.

- `updated_at` *string* &rarr;  The date the container ACL was last updated.



## Import

Containers can be imported using the container id (the last part of the container reference), e.g.:

```shell
terraform import vkcs_keymanager_container.container_1 0c6cd26a-c012-4d7b-8034-057c0f1c2953
```
