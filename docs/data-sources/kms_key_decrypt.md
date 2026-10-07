---
subcategory: "Key Management Service (KMS)"
layout: "vkcs"
page_title: "vkcs: vkcs_kms_key_decrypt"
description: |-
  Decrypt data using a KMS key within VKCS.
---

# vkcs_kms_key_decrypt

Data source containing decrypt result by KMS key

The returned `plaintext` is marked as sensitive and is stored in Terraform state.

## Example Usage

```terraform
resource "vkcs_kms_key" "encrypt_decrypt_key" {
  name             = "encrypt_decrypt_key"
  type             = "aes128-gcm96"
  deletion_allowed = true
}

data "vkcs_kms_key_encrypt" "encrypted" {
  key       = vkcs_kms_key.encrypt_decrypt_key.id
  plaintext = base64encode("a secret message")
}

output "encrypted" {
  value = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}

data "vkcs_kms_key_decrypt" "decrypted" {
  key        = vkcs_kms_key.encrypt_decrypt_key.id
  ciphertext = data.vkcs_kms_key_encrypt.encrypted.ciphertext
}

output "decrypted" {
  value     = base64decode(data.vkcs_kms_key_decrypt.decrypted.plaintext)
  sensitive = true
}
```

## Argument Reference

- `key` **required** *string* &rarr; Name of the KMS key used to decrypt the ciphertext.

- `ciphertext` **required** *string* &rarr; Encrypted data returned by KMS encryption.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `plaintext` sensitive *string* &rarr; Resulting decrypted data encoded in Base64. This attribute is computed. Use `base64decode` to recover the original text, as shown in the example.

- `id` *string* &rarr; Hexadecimal SHA-256 hash of the returned plaintext value.
