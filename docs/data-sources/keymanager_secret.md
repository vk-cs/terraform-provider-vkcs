---
subcategory: "Key Manager"
layout: "vkcs"
page_title: "vkcs: vkcs_keymanager_secret"
description: |-
  Get information on a Key secret resource within VKCS.
---

# vkcs_keymanager_secret

Use this data source to get the ID and the payload of an available Key secret

~> **Important Security Notice** The payload of this data source will be stored *unencrypted* in your Terraform state file. **Use of this resource for production deployments is *not* recommended**. [Read more about sensitive data in state](https://www.terraform.io/docs/language/state/sensitive-data.html).

## Example Usage

```terraform
data "vkcs_keymanager_secret" "example" {
  mode        = "cbc"
  secret_type = "passphrase"
}
```

## Argument Reference
- `acl_only` optional *boolean* &rarr;  Select the Secret with an ACL that contains the user. Project scope is ignored. Defaults to `false`.

- `algorithm` optional *string* &rarr;  The Secret algorithm.

- `bit_length` optional *number* &rarr;  The Secret bit length.

- `created_at_filter` optional *string* &rarr;  Date filter to select the Secret with created matching the specified criteria. See Date Filters below for more detail.

- `expiration_filter` optional *string* &rarr;  Date filter to select the Secret with expiration matching the specified criteria. See Date Filters below for more detail.

- `mode` optional *string* &rarr;  The Secret mode.

- `name` optional *string* &rarr;  The Secret name.

- `region` optional *string* &rarr;  The region in which to obtain the KeyManager client. A KeyManager client is needed to fetch a secret. If omitted, the `region` argument of the provider is used.

- `secret_type` optional *string* &rarr;  The Secret type. For more information see [Secret types](https://docs.openstack.org/barbican/latest/api/reference/secret_types.html).

- `updated_at_filter` optional *string* &rarr;  Date filter to select the Secret with updated matching the specified criteria. See Date Filters below for more detail.


## Attributes Reference
In addition to all arguments above, the following attributes are exported:
- `acl` *object* &rarr;  The list of ACLs assigned to a secret.

- `content_types` *map of* *string* &rarr;  The map of the content types, assigned on the secret.

- `created_at` *string* &rarr;  The date the secret was created.

- `creator_id` *string* &rarr;  The creator of the secret.

- `expiration` *string* &rarr;  The date the secret will expire.

- `id` *string* &rarr;  ID of the resource.

- `metadata` *map of* *string* &rarr;  The map of metadata, assigned on the secret, which has been explicitly and implicitly added.

- `payload` *string* &rarr;  The secret payload.

- `payload_content_encoding` *string* &rarr;  The Secret encoding.

- `payload_content_type` *string* &rarr;  The Secret content type.

- `secret_ref` *string* &rarr;  The secret reference / where to find the secret.

- `status` *string* &rarr;  The status of the secret.

- `updated_at` *string* &rarr;  The date the secret was last updated.



## Date Filters

The values for the `expiration_filter`, `created_at_filter`, and
`updated_at_filter` parameters are comma-separated lists of time stamps in
RFC3339 format. The time stamps can be prefixed with any of these comparison
operators: *gt:* (greater-than), *gte:* (greater-than-or-equal), *lt:*
(less-than), *lte:* (less-than-or-equal).

For example, to get a passphrase a Secret with CBC moda, that will expire in
January of 2020:

```hcl
data "vkcs_keymanager_secret" "date_filter_example" {
  mode              = "cbc"
  secret_type       = "passphrase"
  expiration_filter = "gt:2020-01-01T00:00:00Z"
}
```
