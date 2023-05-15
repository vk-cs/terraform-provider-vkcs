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
### Simple secret
The container with the TLS certificates, which can be used by the loadbalancer HTTPS listener.
```terraform
resource "vkcs_keymanager_secret" "certificate_1" {
  name                 = "certificate"
  payload              = "${file("cert.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "private_key_1" {
  name                 = "private_key"
  payload              = "${file("cert-key.pem")}"
  secret_type          = "private"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_secret" "intermediate_1" {
  name                 = "intermediate"
  payload              = "${file("intermediate-ca.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
}

resource "vkcs_keymanager_container" "tls_1" {
  name = "tls"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = "${vkcs_keymanager_secret.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${vkcs_keymanager_secret.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = "${vkcs_keymanager_secret.intermediate_1.secret_ref}"
  }
}

data "vkcs_networking_subnet" "subnet_1" {
  name = "my-subnet"
}

resource "vkcs_lb_loadbalancer" "lb_1" {
  name          = "loadbalancer"
  vip_subnet_id = "${data.vkcs_networking_subnet.subnet_1.id}"
}

resource "vkcs_lb_listener" "listener_1" {
  name                      = "https"
  protocol                  = "TERMINATED_HTTPS"
  protocol_port             = 443
  loadbalancer_id           = "${vkcs_lb_loadbalancer.lb_1.id}"
  default_tls_container_ref = "${vkcs_keymanager_container.tls_1.container_ref}"
}
```

### Container with the ACL
~> **Note** Only read ACLs are supported
```terraform
resource "vkcs_keymanager_container" "tls_1" {
  name = "tls"
  type = "certificate"

  secret_refs {
    name       = "certificate"
    secret_ref = "${vkcs_keymanager_secret.certificate_1.secret_ref}"
  }

  secret_refs {
    name       = "private_key"
    secret_ref = "${vkcs_keymanager_secret.private_key_1.secret_ref}"
  }

  secret_refs {
    name       = "intermediates"
    secret_ref = "${vkcs_keymanager_secret.intermediate_1.secret_ref}"
  }

  acl {
    read {
      project_access = false
      users = [
        "userid1",
        "userid2",
      ]
    }
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
