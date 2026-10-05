---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_key_encrypt"
description: |-
  Encrypt data using a KMS key within VKCS.
---

# vkcs_kms_key_encrypt

Data source containing encrypt result by KMS key

The input `plaintext` is marked as sensitive and is stored in Terraform state.

## Example Usage

```terraform
data "vkcs_kms_key_encrypt" "encrypted" {
  key       = "key-tf-example"
  plaintext = base64encode("example message")
}

output "ciphertext" {
  value = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}
```

## Argument Reference

- `key` **required** *string* &rarr; Name of the KMS key used for encryption.

- `plaintext` **required** sensitive *string* &rarr; Data to encrypt, encoded in Base64. The provider does not encode this value automatically.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `ciphertext` *string* &rarr; Resulting encrypted data returned by KMS. This attribute is computed.

- `id` *string* &rarr; Hexadecimal SHA-256 hash of the returned ciphertext.
