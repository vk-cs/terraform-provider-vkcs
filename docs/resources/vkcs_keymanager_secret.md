---
layout: "vkcs"
page_title: "vkcs: vkcs_keymanager_secret"
description: |-
  Manages a key secret resource within VKCS.
---

# vkcs_keymanager_secret

Manages a key secret resource within VKCS.

~> **Important Security Notice** The payload of this resource will be stored *unencrypted* in your Terraform state file. **Use of this resource for production deployments is *not* recommended**. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage
### Simple secret
```terraform
resource "vkcs_keymanager_secret" "secret_1" {
  algorithm            = "aes"
  bit_length           = 256
  mode                 = "cbc"
  name                 = "mysecret"
  payload              = "foobar"
  payload_content_type = "text/plain"
  secret_type          = "passphrase"

  metadata = {
    key = "foo"
  }
}
```

### Secret with whitespaces
```terraform
resource "vkcs_keymanager_secret" "secret_1" {
  name                     = "password"
  payload                  = "${base64encode("password with the whitespace at the end ")}"
  secret_type              = "passphrase"
  payload_content_type     = "application/octet-stream"
  payload_content_encoding = "base64"
}
```

### Secret with the expiration date
```terraform
resource "vkcs_keymanager_secret" "secret_1" {
  name                 = "certificate"
  payload              = "${file("certificate.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
  expiration           = "${timeadd(timestamp(), format("%dh", 8760))}" # one year in hours

  lifecycle {
    ignore_changes = [
      expiration
    ]
  }
}
```

### Secret with the ACL
~> **Note** Only read ACLs are supported
```terraform
resource "vkcs_keymanager_secret" "secret_1" {
  name                 = "certificate"
  payload              = "${file("certificate.pem")}"
  secret_type          = "certificate"
  payload_content_type = "text/plain"

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
- `acl` (*Optional*) Allows to control an access to a secret. Currently only the `read` operation is supported. If not specified, the secret is accessible project wide.
  - `read` (*Optional*) Block that describes read operation.
    - `project_access` **Boolean** (*Optional*) Whether the container is accessible project wide. Defaults to `true`.

    - `users` <strong>Set of </strong>**String** (*Optional*) The list of user IDs, which are allowed to access the container, when `project_access` is set to `false`.

- `algorithm` **String** (*Optional*) Metadata provided by a user or system for informational purposes.

- `bit_length` **Number** (*Optional*) Metadata provided by a user or system for informational purposes.

- `expiration` **String** (*Optional*) The expiration time of the secret in the RFC3339 timestamp format (e.g. `2019-03-09T12:58:49Z`). If omitted, a secret will never expire. Changing this creates a new secret.

- `metadata` <strong>Map of </strong>**String** (*Optional*) Additional Metadata for the secret.

- `mode` **String** (*Optional*) Metadata provided by a user or system for informational purposes.

- `name` **String** (*Optional*) Human-readable name for the Secret. Does not have to be unique.

- `payload` **String** (*Optional* Sensitive) The secret's data to be stored. **payload\_content\_type** must also be supplied if **payload** is included.

- `payload_content_encoding` **String** (*Optional*) (required if **payload** is encoded) The encoding used for the payload to be able to include it in the JSON request. Must be either `base64` or `binary`.

- `payload_content_type` **String** (*Optional*) (required if **payload** is included) The media type for the content of the payload. Must be one of `text/plain`, `text/plain;charset=utf-8`, `text/plain; charset=utf-8`, `application/octet-stream`, `application/pkcs8`.

- `region` **String** (*Optional*) The region in which to obtain the KeyManager client. A KeyManager client is needed to create a secret. If omitted, the `region` argument of the provider is used. Changing this creates a new V1 secret.

- `secret_type` **String** (*Optional*) Used to indicate the type of secret being stored. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).


## Attributes Reference
- `acl`  See Argument Reference above.
  - `read`  See Argument Reference above.
    - `project_access` **Boolean** See Argument Reference above.

    - `users` <strong>Set of </strong>**String** See Argument Reference above.

    - `created_at` **String** The date the container ACL was created.

    - `updated_at` **String** The date the container ACL was last updated.

- `algorithm` **String** See Argument Reference above.

- `bit_length` **Number** See Argument Reference above.

- `expiration` **String** See Argument Reference above.

- `metadata` <strong>Map of </strong>**String** See Argument Reference above.

- `mode` **String** See Argument Reference above.

- `name` **String** See Argument Reference above.

- `payload` **String** See Argument Reference above.

- `payload_content_encoding` **String** See Argument Reference above.

- `payload_content_type` **String** See Argument Reference above.

- `region` **String** See Argument Reference above.

- `secret_type` **String** See Argument Reference above.

- `all_metadata` <strong>Map of </strong>**String** The map of metadata, assigned on the secret, which has been explicitly and implicitly added.

- `content_types` <strong>Map of </strong>**String** The map of the content types, assigned on the secret.

- `created_at` **String** The date the secret ACL was created.

- `creator_id` **String** The creator of the secret.

- `id` **String** ID of the resource.

- `secret_ref` **String** The secret reference / where to find the secret.

- `status` **String** The status of the secret.

- `updated_at` **String** The date the secret ACL was last updated.



## Import

Secrets can be imported using the secret id (the last part of the secret reference), e.g.:

```shell
terraform import vkcs_keymanager_secret.secret_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
