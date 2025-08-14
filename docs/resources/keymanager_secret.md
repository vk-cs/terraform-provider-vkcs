---
subcategory: "Key Manager"
layout: "vkcs"
page_title: "vkcs: vkcs_keymanager_secret"
description: |-
  Manages a key secret resource within VKCS.
---

# vkcs_keymanager_secret

Manages a key secret resource within VKCS.

~> **Important Security Notice** The payload of this resource will be stored *unencrypted* in your Terraform state file. **Use of this resource for production deployments is *not* recommended**. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage
### Private key for TERMINATED_HTTPS loadbalancer listener
```terraform
resource "vkcs_keymanager_secret" "priv_key" {
  name                 = "priv-key-tf-example"
  secret_type          = "private"
  payload_content_type = "text/plain"
  payload              = file("${path.module}/private-key.key")
}
```

### Certificate for TERMINATED_HTTPS loadbalancer listener
```terraform
resource "vkcs_keymanager_secret" "certificate" {
  name                 = "certificate-tf-example"
  secret_type          = "certificate"
  payload_content_type = "text/plain"
  payload              = file("${path.module}/certificate.pem")
}
```

## Argument Reference
- `acl` optional &rarr;  Allows to control an access to a secret. Currently only the `read` operation is supported. If not specified, the secret is accessible project wide.
    - `read` optional &rarr;  Block that describes read operation.
        - `project_access` optional *boolean* &rarr;  Whether the container is accessible project wide. Defaults to `true`.

        - `users` optional *set of* *string* &rarr;  The list of user IDs, which are allowed to access the container, when `project_access` is set to `false`.

- `algorithm` optional *string* &rarr;  Metadata provided by a user or system for informational purposes.

- `bit_length` optional *number* &rarr;  Metadata provided by a user or system for informational purposes.

- `expiration` optional *string* &rarr;  The expiration time of the secret in the RFC3339 timestamp format (e.g. `2019-03-09T12:58:49Z`). If omitted, a secret will never expire. Changing this creates a new secret.

- `metadata` optional *map of* *string* &rarr;  Additional Metadata for the secret.

- `mode` optional *string* &rarr;  Metadata provided by a user or system for informational purposes.

- `name` optional *string* &rarr;  Human-readable name for the Secret. Does not have to be unique.

- `payload` optional sensitive *string* &rarr;  The secret's data to be stored. **payload\_content\_type** must also be supplied if **payload** is included.

- `payload_content_encoding` optional *string* &rarr;  (required if **payload** is encoded) The encoding used for the payload to be able to include it in the JSON request. Must be either `base64` or `binary`.

- `payload_content_type` optional *string* &rarr;  (required if **payload** is included) The media type for the content of the payload. Must be one of `text/plain`, `text/plain;charset=utf-8`, `text/plain; charset=utf-8`, `application/octet-stream`, `application/pkcs8`.

- `region` optional *string* &rarr;  The region in which to obtain the KeyManager client. A KeyManager client is needed to create a secret. If omitted, the `region` argument of the provider is used. Changing this creates a new V1 secret.

- `secret_type` optional *string* &rarr;  Used to indicate the type of secret being stored. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `acl` 
    - `read` 
        - `created_at` *string* &rarr;  The date the container ACL was created.

        - `updated_at` *string* &rarr;  The date the container ACL was last updated.

- `all_metadata` *map of* *string* &rarr;  The map of metadata, assigned on the secret, which has been explicitly and implicitly added.

- `content_types` *map of* *string* &rarr;  The map of the content types, assigned on the secret.

- `created_at` *string* &rarr;  The date the secret ACL was created.

- `creator_id` *string* &rarr;  The creator of the secret.

- `id` *string* &rarr;  ID of the resource.

- `secret_ref` *string* &rarr;  The secret reference / where to find the secret.

- `status` *string* &rarr;  The status of the secret.

- `updated_at` *string* &rarr;  The date the secret ACL was last updated.



## Import

Secrets can be imported using the secret id (the last part of the secret reference), e.g.:

```shell
terraform import vkcs_keymanager_secret.secret_1 8a7a79c2-cf17-4e65-b2ae-ddc8bfcf6c74
```
