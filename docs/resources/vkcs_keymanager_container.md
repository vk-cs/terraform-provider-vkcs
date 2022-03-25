---
layout: "vkcs"
page_title: "vkcs: keymanager_container"
description: |-
  Manages a key container resource within VKCS.
---

# vkcs\_keymanager\_container

Manages a key container resource within VKCS.

## Example Usage

### Simple secret

The container with the TLS certificates, which can be used by the loadbalancer HTTPS listener.

```hcl
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

```hcl
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

The following arguments are supported:

* `region` - (Optional) The region in which to obtain the KeyManager client.
    A KeyManager client is needed to create a container. If omitted, the
    `region` argument of the provider is used. Changing this creates a new
    container.

* `name` - (Optional) Human-readable name for the Container. Does not have
    to be unique.

* `type` - (Required) Used to indicate the type of container. Must be one of `generic`, `rsa` or `certificate`.

* `secret_refs` - (Optional) A set of dictionaries containing references to secrets. The structure is described
    below.

* `acl` - (Optional) Allows to control an access to a container. Currently only
  the `read` operation is supported. If not specified, the container is
  accessible project wide. The `read` structure is described below.

The `secret_refs` block supports:

* `name` - (Optional) The name of the secret reference. The reference names must correspond the container type, more details are available [here](https://docs.openstack.org/barbican/stein/api/reference/containers.html).

* `secret_ref` - (Required) The secret reference / where to find the secret, URL.

The `acl` `read` block supports:

* `project_access` - (Optional) Whether the container is accessible project wide.
  Defaults to `true`.

* `users` - (Optional) The list of user IDs, which are allowed to access the
  container, when `project_access` is set to `false`.

* `created_at` - (Computed) The date the container ACL was created.

* `updated_at` - (Computed) The date the container ACL was last updated.

## Attributes Reference

The following attributes are exported:

* `container_ref` - The container reference / where to find the container.
* `region` - See Argument Reference above.
* `name` - See Argument Reference above.
* `type` - See Argument Reference above.
* `secret_refs` - See Argument Reference above.
* `acl` - See Argument Reference above.
* `creator_id` - The creator of the container.
* `status` - The status of the container.
* `created_at` - The date the container was created.
* `updated_at` - The date the container was last updated.
* `consumers` - The list of the container consumers. The structure is described below.

The `consumers` block supports:

* `name` - The name of the consumer.

* `url` - The consumer URL.

## Import

Containers can be imported using the container id (the last part of the container reference), e.g.:

```
$ terraform import vkcs_keymanager_container.container_1 0c6cd26a-c012-4d7b-8034-057c0f1c2953
```
